package sequencer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
)

type batchBuilder struct {
	maxBatchSize uint64

	pendingBlocks   []l2types.DerivationBlock
	builtBatchAttrs *batchAttributes

	lastAppended l2types.BlockID
}

type batchAttributes struct {
	contexts  []*big.Int
	txLengths []*big.Int
	txs       []byte
}

type Header interface {
	Hash() common.Hash
	ParentHash() common.Hash
}

// maxBatchSize is a soft cap on the size of the batch (# of bytes).
func NewBatchBuilder(maxBatchSize uint64) (*batchBuilder, error) {
	return &batchBuilder{maxBatchSize: maxBatchSize}, nil
}

func (b *batchBuilder) LastAppended() l2types.BlockID {
	return b.lastAppended
}

// Appends a block, to be processed and batched.
// Returns a `ReorgDetectedError` if the block is not a child of the last appended block.
func (b *batchBuilder) Append(block l2types.DerivationBlock, header Header) error {
	if (header.ParentHash() != b.lastAppended.Hash()) && (b.lastAppended.Hash() != common.Hash{}) {
		return &utils.L2ReorgDetectedError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, block)
	b.lastAppended = l2types.NewBlockID(block.BlockNumber(), header.Hash())
	return nil
}

func (b *batchBuilder) Reset(lastAppended l2types.BlockID) {
	b.pendingBlocks = []l2types.DerivationBlock{}
	b.builtBatchAttrs = nil
	b.lastAppended = lastAppended
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
func (b *batchBuilder) Build() (*batchAttributes, error) {
	if b.builtBatchAttrs != nil {
		return b.builtBatchAttrs, nil
	}
	if len(b.pendingBlocks) == 0 {
		return nil, io.EOF
	}
	batchAttrs, err := b.serializeToAttrs()
	if err != nil {
		return nil, fmt.Errorf("Failed to build batch, err: %w", err)
	}
	b.builtBatchAttrs = batchAttrs
	return batchAttrs, nil
}

func (b *batchBuilder) Advance() {
	b.builtBatchAttrs = nil
}

func (b *batchBuilder) serializeToAttrs() (*batchAttributes, error) {
	var contexts, txLengths []*big.Int
	buf := new(bytes.Buffer)
	var numBytes uint64
	var block l2types.DerivationBlock
	var idx int
	for idx, block = range b.pendingBlocks {
		// Construct context (`contexts` is a flat array of 3-tuples)
		contexts = append(contexts, big.NewInt(0).SetUint64(block.NumTxs()))
		contexts = append(contexts, big.NewInt(0).SetUint64(block.BlockNumber()))
		contexts = append(contexts, big.NewInt(0).SetUint64(block.Timestamp()))
		numBytes += 3 * 8
		// Construct txData.
		for _, tx := range block.Txs() {
			curLen := buf.Len()
			if _, err := buf.Write(tx); err != nil {
				return nil, err
			}
			txSize := buf.Len() - curLen
			txLengths = append(txLengths, big.NewInt(int64(txSize)))
			numBytes += uint64(txSize)
		}
		// Enforce soft cap on batch size.
		if numBytes > b.maxBatchSize {
			break
		}
	}
	// Advance queue.
	b.pendingBlocks = b.pendingBlocks[idx:]
	return &batchAttributes{contexts, txLengths, buf.Bytes()}, nil
}

// TODO: unused
func writeContext(w *bytes.Buffer, derivCtx *l2types.DerivationContext) error {
	if err := writePrimitive(w, derivCtx.NumTxs()); err != nil {
		return err
	}
	if err := writePrimitive(w, derivCtx.BlockNumber()); err != nil {
		return err
	}
	return writePrimitive(w, derivCtx.Timestamp())
}

func writePrimitive(w *bytes.Buffer, data any) error {
	return binary.Write(w, binary.BigEndian, data)
}
