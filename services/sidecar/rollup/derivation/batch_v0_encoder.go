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

var errBatchFull = errors.New("batch full")
var errBatchTooSmall = errors.New("batch too small")

type V0Config interface {
	GetTargetBatchSize() uint64
}

type BatchV0Encoder struct {
	cfg      V0Config
	currBuf  *bytes.Buffer
	subBatch *subBatch
}

func NewBatchV0Encoder(cfg V0Config) *BatchV0Encoder {
	return &BatchV0Encoder{cfg, &bytes.Buffer{}, newSubBatch()}
}

func (e *BatchV0Encoder) GetBatch(force bool) ([]byte, error) {
	// Return error if the batch is too small and the timeout hasn't been reached.
	// If the timeout has been reached, the batch will be closed regardless of its current size.
	if !force && e.currBuf.Len() < int(e.cfg.GetTargetBatchSize()) {
		return nil, errBatchTooSmall
	}
	var data = e.currBuf.Bytes()
	e.currBuf.Reset()
	return data, nil
}

func (e *BatchV0Encoder) Len() int { return e.currBuf.Len() }
func (e *BatchV0Encoder) Reset() {
	e.currBuf.Reset()
	e.subBatch = newSubBatch()
}

// Processes a block. If the block is non-empty and fits, add it to the current sub-batch.
// If the block would cause the batch to exceed the target size, close the entire batch (by writing to `currBuf“).
// Returns an error if the block could not be processed.
func (e *BatchV0Encoder) ProcessBlock(block *types.Block, isNewEpoch bool) error {
	// Initialize buffer with version byte if it hasn't been already.
	if e.currBuf.Len() == 0 {
		if err := e.currBuf.WriteByte(e.getVersion()); err != nil {
			return fmt.Errorf("failed to encode version: %w", err)
		}
	}
	var (
		// Block is empty
		shouldSkipBlock = len(block.Transactions()) == 0
		// Batch would exceed the target size with the current sub-batch.
		shouldCloseBatch = e.shouldCloseBatch()
		// Should close batch entirely, or block is empty.
		shouldCloseSubBatch = shouldCloseBatch || (shouldSkipBlock && e.subBatch.size() > 0)
	)
	if shouldCloseSubBatch {
		if err := e.closeSubBatch(); err != nil {
			return fmt.Errorf("could not close sub-batch: %w", err)
		}
	}
	// Enforce soft cap on batch size.
	if shouldCloseBatch {
		return errBatchFull
	}
	// Skip intrinsically-derivable blocks.
	if shouldSkipBlock {
		log.Info("Skipping intrinsically-derivable block", "block#", block.NumberU64())
		return nil
	}
	// Append a block's txs to the current sub-batch.
	if err := e.subBatch.appendTxBlock(block.NumberU64(), block.Transactions()); err != nil {
		return fmt.Errorf("could not append block of txs: %w", err)
	}
	return nil
}

// Returns the data format version (v0).
func (e *BatchV0Encoder) getVersion() BatchEncoderVersion { return V0 }
func (e *BatchV0Encoder) shouldCloseBatch() bool {
	return uint64(e.currBuf.Len())+e.subBatch.size() > e.cfg.GetTargetBatchSize()
}

// Writes encoded sub-batch out to the buffer.
func (e *BatchV0Encoder) closeSubBatch() error {
	log.Info("Closing sub-batch...")
	if err := rlp.Encode(e.currBuf, e.subBatch); err != nil {
		return fmt.Errorf("could not encode sub-batch: %w", err)
	}
	e.subBatch = newSubBatch()
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
	// Set the first L2 block number if it hasn't been set yet.
	if len(s.txBlocks) == 0 {
		s.firstL2BlockNum = blockNum
		s.contentSize += uint64(rlp.IntSize(blockNum))
	}
	// Append the block of txs to the sub-batch.
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
