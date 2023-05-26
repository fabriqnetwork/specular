package sequencer

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
)

type batchBuilder struct {
	maxBatchSize uint64

	pendingBlocks   []types.DerivationBlock
	builtBatchAttrs *batchAttributes

	lastAppended types.BlockID
}

type batchAttributes struct {
	contexts           []*big.Int
	txLengths          []*big.Int
	firstL2BlockNumber *big.Int
	txs                []byte
}

type Header interface {
	GetHash() common.Hash
	GetParentHash() common.Hash
}

// maxBatchSize is a soft cap on the size of the batch (# of bytes).
func NewBatchBuilder(maxBatchSize uint64) (*batchBuilder, error) {
	return &batchBuilder{maxBatchSize: maxBatchSize}, nil
}

func (b *batchBuilder) LastAppended() types.BlockID {
	return b.lastAppended
}

// Appends a block, to be processed and batched.
// Returns a `ReorgDetectedError` if the block is not a child of the last appended block.
func (b *batchBuilder) Append(block types.DerivationBlock, header Header) error {
	if (header.GetParentHash() != b.lastAppended.GetHash()) && (b.lastAppended.GetHash() != common.Hash{}) {
		return &utils.L2ReorgDetectedError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, block)
	b.lastAppended = types.NewBlockID(block.BlockNumber(), header.GetHash())
	return nil
}

func (b *batchBuilder) Reset(lastAppended types.BlockID) {
	b.pendingBlocks = []types.DerivationBlock{}
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
	var (
		contexts, txLengths []*big.Int
		buf                 = new(bytes.Buffer)
		numBytes            uint64
		block               types.DerivationBlock
		idx                 int
	)
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
	firstL2BlockNumber := big.NewInt(0).SetUint64(b.pendingBlocks[0].BlockNumber())
	return &batchAttributes{contexts, txLengths, firstL2BlockNumber, buf.Bytes()}, nil
}
