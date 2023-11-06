package derivation

import (
	"bytes"
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
	WriteData(b *bytes.Buffer) (int, error)
}

type batchBuilder struct {
	maxBatchSize  uint64
	pendingBlocks []ethTypes.Block
	lastAppended  types.BlockID
	lastBuilt     []byte
}

// maxBatchSize is a soft cap on the size of the batch (# of bytes).
func NewBatchBuilder(maxBatchSize uint64) *batchBuilder {
	return &batchBuilder{maxBatchSize: maxBatchSize}
}

func (b *batchBuilder) LastAppended() types.BlockID { return b.lastAppended }

// Appends a block, to be processed and batched.
// Returns a `InvalidBlockError` if the block is not a child of the last appended block.
func (b *batchBuilder) Append(block *ethTypes.Block) error {
	// Ensure block is a child of the last appended block. Not enforced when no prior blocks.
	if (b.lastAppended.GetHash() != common.Hash{}) && (block.ParentHash() != b.lastAppended.GetHash()) {
		return InvalidBlockError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, *block)
	b.lastAppended = types.NewBlockID(block.NumberU64(), block.Hash())
	return nil
}

// Resets the builder, discarding all pending blocks.
func (b *batchBuilder) Reset(lastAppended types.BlockID) {
	b.pendingBlocks = []ethTypes.Block{}
	b.lastAppended = lastAppended
	b.lastBuilt = nil
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
func (b *batchBuilder) Build() ([]byte, error) {
	if b.lastBuilt != nil {
		return b.lastBuilt, nil
	}
	encodedData, err := b.encodePending()
	if err != nil {
		return nil, fmt.Errorf("failed to encode pending blocks into a new batch: %w", err)
	}
	b.lastBuilt = encodedData
	return b.lastBuilt, nil
}

// Advances the builder, clearing the last built batch.
func (b *batchBuilder) Advance() {
	b.lastBuilt = nil
}

// Encodes pending blocks into a new batch, constrained by `maxBatchSize`.
// Returns an `io.EOF` error if there are no pending blocks.
func (b *batchBuilder) encodePending() ([]byte, error) {
	if len(b.pendingBlocks) == 0 {
		return nil, io.EOF
	}
	// Encode data. TODO: abstract encoder creation out.
	batcherDataV0, numProcessed, err := encodeVersionedData(NewBatchV0Encoder(b.pendingBlocks, b.maxBatchSize))
	if err != nil {
		return nil, err
	}
	log.Info("Encoded l2 blocks", "num_processed", numProcessed)
	// Advance queue.
	b.pendingBlocks = b.pendingBlocks[numProcessed:]
	log.Trace("Advanced pending blocks", "len", len(b.pendingBlocks))
	return batcherDataV0, nil
}

// Encodes a versioned batch.
// Returns the encoded data and the number of blocks processed.
func encodeVersionedData(e VersionedDataEncoder) ([]byte, int, error) {
	var buf bytes.Buffer
	if err := buf.WriteByte(e.GetVersion()); err != nil {
		return nil, 0, fmt.Errorf("failed to encode version: %w", err)
	}
	numProcessed, err := e.WriteData(&buf)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to encode data: %w", err)
	}
	return buf.Bytes(), numProcessed, nil
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
