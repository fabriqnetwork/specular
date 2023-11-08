package derivation

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/bridge"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type BatchEncoderVersion = byte

const V0 BatchEncoderVersion = 0x0

type Config interface {
	GetL1OracleAddr() common.Address
	GetSeqWindowSize() uint64
	GetSubSafetyMargin() uint64
}

type VersionedDataEncoder interface {
	// Returns an encoded batch if one is ready (or if forced).
	// If not forced, an error is returned if one cannot yet be built.
	GetBatch(force bool) ([]byte, error)
	// Processes a block, adding its data to the current batch.
	// Returns an `errBatchFull` if the block would cause the batch to exceed the target size.
	// Returns an error if the block could not be processed.
	ProcessBlock(block *ethTypes.Block, isNewEpoch bool) error
	// Resets the encoder, discarding all buffered data.
	Reset()
}

type InvalidBlockError struct{ Msg string }

func (e InvalidBlockError) Error() string { return e.Msg }

type batchBuilder struct {
	cfg           Config
	encoder       VersionedDataEncoder
	pendingBlocks []*ethTypes.Block
	lastEnqueued  types.BlockID
	lastBuilt     []byte

	timeout uint64
}

func NewBatchBuilder(cfg Config, encoder VersionedDataEncoder) *batchBuilder {
	return &batchBuilder{cfg, encoder, nil, types.BlockID{}, nil, 0}
}

func (b *batchBuilder) LastEnqueued() types.BlockID { return b.lastEnqueued }

// Enqueues a block, to be processed and batched.
// Returns a `InvalidBlockError` if the block is not a child of the last enqueued block.
func (b *batchBuilder) Enqueue(block *ethTypes.Block) error {
	// Ensure block is a child of the last appended block. Not enforced when no prior blocks.
	if (b.lastEnqueued.GetHash() != common.Hash{}) && (block.ParentHash() != b.lastEnqueued.GetHash()) {
		return InvalidBlockError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, block)
	b.lastEnqueued = types.NewBlockID(block.NumberU64(), block.Hash())
	return nil
}

// Resets the builder, discarding all pending blocks.
func (b *batchBuilder) Reset(lastEnqueued types.BlockID) {
	b.encoder.Reset()
	b.pendingBlocks = []*ethTypes.Block{}
	b.lastEnqueued = lastEnqueued
	b.lastBuilt = nil
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
// An l1Head must be provided to allow the encoder to determine if the batch is ready.
func (b *batchBuilder) Build(l1Head types.BlockID) ([]byte, error) {
	if b.lastBuilt != nil {
		return b.lastBuilt, nil
	}
	if err := b.encodePending(); err != nil {
		return nil, fmt.Errorf("failed to encode pending blocks into a new batch: %w", err)
	}
	return b.getBatch(l1Head)
}

// Advances the builder, clearing the last built batch.
func (b *batchBuilder) Advance() {
	b.lastBuilt = nil
}

// Tries to get the current batch.
func (b *batchBuilder) getBatch(l1Head types.BlockID) ([]byte, error) {
	// Force-build batch if necessary (timeout exceeded).
	force := b.timeout != 0 && l1Head.GetNumber() >= b.timeout
	batch, err := b.encoder.GetBatch(force)
	if err != nil {
		if errors.Is(err, errBatchTooSmall) {
			log.Warn("Batch too small, waiting for more blocks")
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get batch: %w", err)
	}
	// Cache last built batch.
	b.lastBuilt = batch
	return batch, nil
}

// Encodes pending blocks into a new batch, constrained by `maxBatchSize`.
// Returns an `io.EOF` error if there are no pending blocks.
func (b *batchBuilder) encodePending() error {
	if len(b.pendingBlocks) == 0 {
		return io.EOF
	}
	// Process all pending blocks (until the batch is full).
	numProcessed := 0
	for _, block := range b.pendingBlocks {
		if err := b.processBlock(block); err != nil {
			if errors.Is(err, errBatchFull) {
				log.Info("Batch is full, stopping processing")
				break
			}
			return fmt.Errorf("failed to process block: %w", err)
		}
		numProcessed += 1
	}
	log.Info("Encoded l2 blocks", "num_processed", numProcessed)
	// Advance queue.
	b.pendingBlocks = b.pendingBlocks[numProcessed:]
	log.Trace("Advanced pending blocks", "pending left", len(b.pendingBlocks))
	return nil
}

// Processes a block, adding its data to the current batch.
func (b *batchBuilder) processBlock(block *ethTypes.Block) (err error) {
	var epoch uint64
	// Decode oracle tx to update timeout.
	if block.Transactions().Len() > 0 {
		// Process oracle tx, if it exists.
		var firstTx = block.Transactions()[0]
		if *firstTx.To() == b.cfg.GetL1OracleAddr() {
			epoch, _, _, _, _, err = bridge.UnpackL1OracleInput(firstTx)
			if err != nil {
				return fmt.Errorf("could not unpack oracle tx: %w", err)
			}
			// Note: it may make sense to move this to after processing the block.
			b.updateTimeout(epoch)
		}
	}
	// Process block.
	if err := b.encoder.ProcessBlock(block, epoch != 0); err != nil {
		return err
	}
	return err
}

// Updates the batch timeout if the given L1 epoch is earlier than the current timeout.
// Note: the timeout won't be updated more than once assuming the L1 epoch is monotonically increasing.
func (b *batchBuilder) updateTimeout(l1Epoch uint64) {
	timeout := l1Epoch + b.cfg.GetSeqWindowSize() - b.cfg.GetSubSafetyMargin()
	if b.timeout == 0 || b.timeout > timeout {
		b.timeout = timeout
	}
}

type DecodeTxBatchError struct{ msg string }

func (e *DecodeTxBatchError) Error() string {
	return fmt.Sprintf("failed to decode batch: %s", e.msg)
}

// TODO: this is not currently called anywhere but will be useful for testing.
func DecodeBatch(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, &DecodeTxBatchError{"empty batch data"}
	}
	// TODO: use map.
	switch data[0] {
	case V0:
		return decodeV0(data[1:])
	default:
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid batch version: %d", data[0])}
	}
}
