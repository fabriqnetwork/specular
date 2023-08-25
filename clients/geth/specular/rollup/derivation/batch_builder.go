package derivation

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

type batchBuilder struct {
	maxBatchSize uint64

	pendingBlocks   []DerivationBlock
	builtBatchAttrs *BatchAttributes

	lastAppended types.BlockID
}

type BatchAttributes struct {
	contexts           []*big.Int
	txLengths          []*big.Int
	firstL2BlockNumber *big.Int
	txBatch            []byte // encoded batch of transactions.
}

func (a *BatchAttributes) Contexts() []*big.Int         { return a.contexts }
func (a *BatchAttributes) TxLengths() []*big.Int        { return a.txLengths }
func (a *BatchAttributes) FirstL2BlockNumber() *big.Int { return a.firstL2BlockNumber }
func (a *BatchAttributes) TxBatch() []byte              { return a.txBatch }

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
	if (header.GetParentHash() != b.lastAppended.GetHash()) && (b.lastAppended.GetHash() != common.Hash{}) {
		return InvalidBlockError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, block)
	b.lastAppended = types.NewBlockID(block.BlockNumber(), header.GetHash())
	return nil
}

func (b *batchBuilder) Reset(lastAppended types.BlockID) {
	b.pendingBlocks = []DerivationBlock{}
	b.builtBatchAttrs = nil
	b.lastAppended = lastAppended
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
func (b *batchBuilder) Build() (*BatchAttributes, error) {
	if b.builtBatchAttrs != nil {
		return b.builtBatchAttrs, nil
	}
	if len(b.pendingBlocks) == 0 {
		return nil, io.EOF
	}
	batchAttrs, err := b.serializeToAttrs()
	if err != nil {
		return nil, fmt.Errorf("failed to build batch: %w", err)
	}
	b.builtBatchAttrs = batchAttrs
	return batchAttrs, nil
}

func (b *batchBuilder) Advance() {
	b.builtBatchAttrs = nil
}

func (b *batchBuilder) serializeToAttrs() (*BatchAttributes, error) {
	var (
		contexts, txLengths []*big.Int
		buf                 = new(bytes.Buffer)
		numBytes            uint64
		block               DerivationBlock
		idx                 int
	)
	for idx, block = range b.pendingBlocks {
		// Construct context (`contexts` is a flat array of 2-tuples)
		contexts = append(contexts, big.NewInt(0).SetUint64(block.NumTxs()))
		contexts = append(contexts, big.NewInt(0).SetUint64(block.Timestamp()))
		numBytes += 2 * 8
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
			log.Info("Reached max batch size", "numBytes", numBytes, "maxBatchSize", b.maxBatchSize)
			break
		}
	}
	// Construct batch attributes.
	var (
		firstL2BlockNumber = big.NewInt(0).SetUint64(b.pendingBlocks[0].BlockNumber())
		attrs              = &BatchAttributes{contexts, txLengths, firstL2BlockNumber, buf.Bytes()}
	)
	log.Info("Serialized l2 blocks", "first", firstL2BlockNumber, "last", b.pendingBlocks[idx].BlockNumber())
	// Advance queue.
	b.pendingBlocks = b.pendingBlocks[idx+1:]
	log.Trace("Advanced pending blocks", "len", len(b.pendingBlocks))
	return attrs, nil
}
