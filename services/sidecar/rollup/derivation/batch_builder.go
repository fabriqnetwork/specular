package derivation

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type HeaderRef interface {
	GetHash() common.Hash
	GetParentHash() common.Hash
}

type InvalidBlockError struct{ Msg string }

func (e InvalidBlockError) Error() string { return e.Msg }

type VersionedDataEncoder interface {
	GetVersion() byte
	GetBatch() ([]byte, error)
	Reset()
	ProcessBlock(block *ethTypes.Block) error
}

type ProtocolConfig interface {
	GetSeqWindowSize() uint64
}

type batchBuilder struct {
	encoder       VersionedDataEncoder
	pendingBlocks []ethTypes.Block
	lastEnqueued  types.BlockID
	lastBuilt     []byte
}

// maxBatchSize is a soft cap on the size of the batch (# of bytes).
func NewBatchBuilder(minBatchSize, maxBatchSize uint64) *batchBuilder {
	// TODO: abstract out encoder creation.
	return &batchBuilder{encoder: NewBatchV0Encoder(minBatchSize, maxBatchSize)}
}

func (b *batchBuilder) LastEnqueued() types.BlockID { return b.lastEnqueued }

// Enqueues a block, to be processed and batched.
// Returns a `InvalidBlockError` if the block is not a child of the last enqueued block.
func (b *batchBuilder) Enqueue(block *ethTypes.Block) error {
	// Ensure block is a child of the last appended block. Not enforced when no prior blocks.
	if (b.lastEnqueued.GetHash() != common.Hash{}) && (block.ParentHash() != b.lastEnqueued.GetHash()) {
		return InvalidBlockError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, *block)
	b.lastEnqueued = types.NewBlockID(block.NumberU64(), block.Hash())
	return nil
}

// Resets the builder, discarding all pending blocks.
func (b *batchBuilder) Reset(lastEnqueued types.BlockID) {
	b.encoder.Reset()
	b.pendingBlocks = []ethTypes.Block{}
	b.lastEnqueued = lastEnqueued
	b.lastBuilt = nil
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
func (b *batchBuilder) Build() ([]byte, error) {
	if b.lastBuilt != nil {
		return b.lastBuilt, nil
	}
	if err := b.encodePending(); err != nil {
		return nil, fmt.Errorf("failed to encode pending blocks into a new batch: %w", err)
	}
	batch, err := b.encoder.GetBatch()
	if err != nil {
		// process
	}
	b.lastBuilt = batch
	return b.lastBuilt, nil
}

// Advances the builder, clearing the last built batch.
func (b *batchBuilder) Advance() {
	b.lastBuilt = nil
}

// Encodes pending blocks into a new batch, constrained by `maxBatchSize`.
// Returns an `io.EOF` error if there are no pending blocks.
func (b *batchBuilder) encodePending() error {
	if len(b.pendingBlocks) == 0 {
		return io.EOF
	}
	// Process pending blocks.
	numProcessed := 0
	for _, block := range b.pendingBlocks {
		if err := b.encoder.ProcessBlock(&block); err != nil {
			if errors.Is(err, errors.New("full batch")) {
				break
			}
		}
		numProcessed += 1
	}
	log.Info("Encoded l2 blocks", "num_processed", numProcessed)
	// Advance queue.
	b.pendingBlocks = b.pendingBlocks[numProcessed:]
	log.Trace("Advanced pending blocks", "len", len(b.pendingBlocks))
	return nil
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
	version := data[0]
	switch version {
	case 0:
		return decodeV0(data[1:])
	default:
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid batch version: {%d}", version)}
	}
}
