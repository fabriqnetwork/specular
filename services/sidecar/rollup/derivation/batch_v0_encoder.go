package derivation

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

const v0 = 0x0

type BatchV0Encoder struct {
	minSize  uint64
	maxSize  uint64
	currBuf  *bytes.Buffer
	subBatch *subBatch
}

func NewBatchV0Encoder(minSize uint64, maxSize uint64) *BatchV0Encoder {
	return &BatchV0Encoder{minSize, maxSize, nil, newSubBatch()}
}

// Returns the data format version (v0, i.e. 0x0).
func (b *BatchV0Encoder) GetVersion() byte { return v0 }

func (b *BatchV0Encoder) GetBatch() ([]byte, error) {
	if b.currBuf.Len() < int(b.minSize) {
		return nil, errors.New("insufficient data to form batch")
	}
	data := b.currBuf.Bytes()
	b.currBuf.Reset()
	return data, nil
}

func (b *BatchV0Encoder) Reset() {
	b.currBuf.Reset()
	b.subBatch = newSubBatch()
}

// Returns an error if the block could not be processed.
func (b *BatchV0Encoder) ProcessBlock(block *types.Block) error {
	// Initialize buffer with version byte if it hasn't been already.
	if b.currBuf.Len() == 0 {
		if err := b.currBuf.WriteByte(b.GetVersion()); err != nil {
			return fmt.Errorf("failed to encode version: %w", err)
		}
	}
	var (
		shouldSkipBlock     = len(block.Transactions()) == 0
		shouldCloseBatch    = uint64(b.currBuf.Len())+b.subBatch.size() > b.maxSize
		shouldCloseSubBatch = shouldCloseBatch || (shouldSkipBlock && b.subBatch.size() > 0)
	)
	if shouldCloseSubBatch {
		if err := b.closeSubBatch(); err != nil {
			return fmt.Errorf("could not close sub-batch: %w", err)
		}
	}
	// Enforce soft cap on batch size.
	if shouldCloseBatch {
		log.Info("Closing batch")
		return errors.New("full batch")
	}
	// Skip intrinsically-derivable blocks.
	if shouldSkipBlock {
		log.Info("Skipping intrinsically-derivable block", "block#", block.NumberU64())
		return nil
	}
	// Append a block's txs to the current sub-batch.
	if err := b.subBatch.appendTxBlock(block.NumberU64(), block.Transactions()); err != nil {
		return fmt.Errorf("could not append block of txs: %w", err)
	}
	return nil
}

// Writes encoded sub-batch out to the buffer.
func (b *BatchV0Encoder) closeSubBatch() error {
	log.Info("Closing sub-batch...")
	if err := rlp.Encode(b.currBuf, b.subBatch); err != nil {
		return fmt.Errorf("could not encode sub-batch: %w", err)
	}
	b.Reset()
	return nil
}

type subBatch struct {
	firstL2BlockNum uint64
	txBlocks        []rawTxBlock
	contentSize     uint64 `rlp:"-"` // size of sub-batch content (# of bytes)
}

type rawTxBlock []hexutil.Bytes

func newSubBatch() *subBatch     { return &subBatch{contentSize: 0} }
func (s *subBatch) size() uint64 { return rlp.ListSize(s.contentSize) }

func (s *subBatch) appendTxBlock(blockNum uint64, txs types.Transactions) error {
	if len(s.txBlocks) == 0 {
		s.firstL2BlockNum = blockNum
		s.contentSize += uint64(rlp.IntSize(blockNum))
	}
	marshalled, numBytes, err := marshallTxs(txs)
	if err != nil {
		return fmt.Errorf("could not marshall txs: %w", err)
	}
	s.txBlocks = append(s.txBlocks, marshalled)
	s.contentSize += uint64(numBytes)
	return nil
}

func marshallTxs(txs types.Transactions) (rawTxs rawTxBlock, numBytes int, err error) {
	rawTxs = make([]hexutil.Bytes, 0, len(txs))
	for i, tx := range txs {
		rawTx, err := tx.MarshalBinary()
		if err != nil {
			return nil, 0, fmt.Errorf("could not marshall tx %v in block %v: %w", i, tx.Hash(), err)
		}
		rawTxs = append(rawTxs, rawTx)
		numBytes += len(rawTx)
	}
	return rawTxs, numBytes, err
}

func decodeV0(data []byte) ([]subBatch, error) {
	var decoded []subBatch
	if err := rlp.Decode(bytes.NewReader(data), &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}
