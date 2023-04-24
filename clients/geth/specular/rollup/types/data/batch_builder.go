package data

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
)

type BatchBuilder struct {
	maxBatchSize uint64

	pendingBlocks []DerivationBlock
	builtBatch    []byte

	lastAppendedBlockNumber uint64
	lastAppendedBlockHash   common.Hash
}

type Header interface {
	Hash() common.Hash
	ParentHash() common.Hash
}

// maxBatchSize is a soft cap on the size of the batch (# of bytes).
func NewBatchBuilder(maxBatchSize uint64) (*BatchBuilder, error) {
	return &BatchBuilder{maxBatchSize: maxBatchSize}, nil
}

func (b *BatchBuilder) LastAppended() uint64 {
	return b.lastAppendedBlockNumber
}

// Appends a block, to be processed and batched.
// Returns a `ReorgDetectedError` if the block is not a child of the last appended block.
func (b *BatchBuilder) Append(block DerivationBlock, header Header) error {
	if (header.ParentHash() != b.lastAppendedBlockHash) && (b.lastAppendedBlockHash != common.Hash{}) {
		return &utils.L2ReorgDetectedError{Msg: "Appended block is not a child of the last appended block"}
	}
	b.pendingBlocks = append(b.pendingBlocks, block)
	b.lastAppendedBlockNumber = block.blockNumber
	b.lastAppendedBlockHash = header.Hash()
	return nil
}

func (b *BatchBuilder) Reset(lastAppendedBlockNumber uint64, lastAppendedBlockHash common.Hash) {
	b.pendingBlocks = []DerivationBlock{}
	b.builtBatch = []byte{}
	b.lastAppendedBlockNumber = lastAppendedBlockNumber
	b.lastAppendedBlockHash = lastAppendedBlockHash
}

// This short-circuits the build process if a batch is
// already built and `Advance` hasn't been called.
func (b *BatchBuilder) Build() ([]byte, error) {
	if len(b.builtBatch) > 0 {
		return b.builtBatch, nil
	}
	if len(b.pendingBlocks) == 0 {
		return nil, io.EOF
	}
	batch, err := b.serializeToBytes()
	if err != nil {
		return nil, fmt.Errorf("Failed to build batch, err: %w", err)
	}
	b.builtBatch = batch
	return batch, nil
}

func (b *BatchBuilder) Advance() {
	b.builtBatch = []byte{}
}

func (b *BatchBuilder) serializeToBytes() ([]byte, error) {
	contexts, txLengths, txs, err := b.serializeToArgs()
	if err != nil {
		return nil, err
	}
	return client.PackAppendTxBatchInput(contexts, txLengths, txs)
}

func (b *BatchBuilder) serializeToArgs() ([]*big.Int, []*big.Int, []byte, error) {
	var contexts, txLengths []*big.Int
	buf := new(bytes.Buffer)
	var numBytes uint64
	var block DerivationBlock
	var idx int
	for idx, block = range b.pendingBlocks {
		// Construct context (`contexts` is a flat array of 3-tuples)
		contexts = append(contexts, big.NewInt(0).SetUint64(block.numTxs))
		contexts = append(contexts, big.NewInt(0).SetUint64(block.blockNumber))
		contexts = append(contexts, big.NewInt(0).SetUint64(block.timestamp))
		numBytes += 3 * 8
		// Construct txData.
		for _, tx := range block.txs {
			curLen := buf.Len()
			if _, err := buf.Write(tx); err != nil {
				return nil, nil, nil, err
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
	return contexts, txLengths, buf.Bytes(), nil
}

// TODO: unused
func writeContext(w *bytes.Buffer, ctx *DerivationContext) error {
	if err := writePrimitive(w, ctx.numTxs); err != nil {
		return err
	}
	if err := writePrimitive(w, ctx.blockNumber); err != nil {
		return err
	}
	return writePrimitive(w, ctx.timestamp)
}

func writePrimitive(w *bytes.Buffer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}
