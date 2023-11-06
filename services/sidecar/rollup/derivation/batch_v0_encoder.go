package derivation

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

const v0 = 0x0

type BatchV0Encoder struct {
	blocks  []types.Block
	maxSize uint64
}

func NewBatchV0Encoder(blocks []types.Block, maxSize uint64) *BatchV0Encoder {
	return &BatchV0Encoder{blocks, maxSize}
}

// Returns the data format version (v0, i.e. 0x0).
func (b *BatchV0Encoder) GetVersion() byte { return v0 }

// Writes the batch data in the v0-specified format.
func (b *BatchV0Encoder) WriteData(buf *bytes.Buffer) (int, error) {
	var (
		idx      int
		block    types.Block
		numBytes uint64
		subBatch = newSubBatch()
	)
	// Iterate block-by-block, enforcing a soft cap on the total batch size.
	for idx, block = range b.blocks {
		var (
			shouldSkipBlock     = len(block.Transactions()) == 0
			shouldCloseBatch    = numBytes+subBatch.size() > b.maxSize
			shouldCloseSubBatch = shouldCloseBatch || (shouldSkipBlock && subBatch.size() > 0)
		)
		if shouldCloseSubBatch {
			log.Info("Closing sub-batch...")
			// Write encoded sub-batch.
			if err := rlp.Encode(buf, subBatch); err != nil {
				return 0, fmt.Errorf("could not encode sub-batch: %w", err)
			}
			numBytes += subBatch.size()
			// Initialize a new sub-batch.
			subBatch = newSubBatch()
		}
		// Enforce soft cap on batch size.
		if shouldCloseBatch {
			log.Info("Closing batch")
			break
		}
		// Skip intrinsically-derivable blocks.
		if shouldSkipBlock {
			log.Info("Skipping intrinsically-derivable block", "block#", block.NumberU64())
			continue
		}
		// Append a block's txs to the current sub-batch.
		if err := subBatch.appendTxBlock(block.NumberU64(), block.Transactions()); err != nil {
			return 0, fmt.Errorf("could not append block of txs: %w", err)
		}
	}
	// Handle the case where the last sub-batch is non-empty.
	if subBatch.size() > 0 {
		if err := rlp.Encode(buf, subBatch); err != nil {
			return 0, fmt.Errorf("could not encode last sub-batch: %w", err)
		}
		idx += 1
	}
	return idx, nil
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
