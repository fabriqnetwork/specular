package derivation

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type batchBuilder struct {
	maxBatchSize uint64

	pendingBlocks  []DerivationBlock
	builtBatchData *[]byte

	lastAppended types.BlockID
}

type BlockAttributes struct {
	Timestamp *big.Int
	Txs       [][]byte // encoded batch of transactions.
}

type BatchAttributes struct {
	FirstL2BlockNumber *big.Int
	Blocks             []BlockAttributes
}

// TODO: find a better place to not hardcode this
func TxBatchVersion() byte { return byte(0x00) }

type HeaderRef interface {
	GetHash() common.Hash
	GetParentHash() common.Hash
}

type InvalidBlockError struct{ Msg string }

func (e InvalidBlockError) Error() string { return e.Msg }

// maxBatchSize is a soft cap on the size of the batch (# of bytes).
func NewBatchBuilder(maxBatchSize uint64) (*batchBuilder, error) {
	return &batchBuilder{maxBatchSize: maxBatchSize}, nil
}

func (b *batchBuilder) LastAppended() types.BlockID { return b.lastAppended }

// Appends a block, to be processed and batched.
// Returns a `InvalidBlockError` if the block is not a child of the last appended block.
func (b *batchBuilder) Append(block DerivationBlock, header HeaderRef) error {
	// Ensure block is a child of the last appended block. Not enforced when no prior blocks.
	if (header.GetParentHash() != b.lastAppended.GetHash()) && (b.lastAppended.GetHash() != common.Hash{}) {
		return InvalidBlockError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, block)
	b.lastAppended = types.NewBlockID(block.BlockNumber(), header.GetHash())
	return nil
}

func (b *batchBuilder) Reset(lastAppended types.BlockID) {
	b.pendingBlocks = []DerivationBlock{}
	b.builtBatchData = nil
	b.lastAppended = lastAppended
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
func (b *batchBuilder) Build() (*[]byte, error) {
	if b.builtBatchData != nil {
		return b.builtBatchData, nil
	}
	if len(b.pendingBlocks) == 0 {
		return nil, io.EOF
	}
	batchAttrs, err := b.serializeToAttrs()
	if err != nil {
		return nil, fmt.Errorf("failed to build batch: %w", err)
	}
	batchData, err := encodeBatch(batchAttrs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode batch: %w", err)
	}

	b.builtBatchData = &batchData
	return b.builtBatchData, nil
}

func (b *batchBuilder) Advance() {
	b.builtBatchData = nil
}

func (b *batchBuilder) serializeToAttrs() (*BatchAttributes, error) {
	var (
		block  DerivationBlock
		idx    int
		blocks []BlockAttributes
	)
	firstL2BlockNumber := big.NewInt(0).SetUint64(b.pendingBlocks[0].BlockNumber())
	numBytes := uint64(rlp.IntSize(firstL2BlockNumber.Uint64()))

	// Iterate block-by-block to enforce soft cap on batch size.
	for idx, block = range b.pendingBlocks {
		blockTimestamp := big.NewInt(0).SetUint64(block.Timestamp())
		// Calculate size of txData.
		txListSize := uint64(0)
		for _, tx := range block.Txs() {
			txListSize += uint64(len(tx))
		}
		blocks = append(blocks, BlockAttributes{blockTimestamp, block.Txs()})
		numBytes += rlp.ListSize(uint64(rlp.IntSize(blockTimestamp.Uint64())) + rlp.ListSize(txListSize))

		// Enforce soft cap on batch size.
		if rlp.ListSize(numBytes) > b.maxBatchSize {
			log.Info("Reached max batch size", "numBytes", numBytes, "maxBatchSize", b.maxBatchSize, "numBlocks", len(blocks))
			break
		}
	}
	// Construct batch attributes.
	attrs := &BatchAttributes{firstL2BlockNumber, blocks}
	log.Info("Serialized l2 blocks", "first", firstL2BlockNumber, "last", b.pendingBlocks[idx].BlockNumber())
	// Advance queue.
	b.pendingBlocks = b.pendingBlocks[idx+1:]
	log.Trace("Advanced pending blocks", "len", len(b.pendingBlocks))
	return attrs, nil
}

func encodeBatch(b *BatchAttributes) ([]byte, error) {
	var w bytes.Buffer
	// Batch starts with version byte
	if err := w.WriteByte(TxBatchVersion()); err != nil {
		return nil, err
	}

	buf := rlp.NewEncoderBuffer(&w)
	err := rlp.Encode(buf, b)
	return w.Bytes(), err
}
