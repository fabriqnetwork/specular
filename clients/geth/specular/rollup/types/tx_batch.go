package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// TxBatch represents a transaction batch to be sequenced to L1 sequencer inbox
// It may contain multiple blocks
type TxBatch struct {
	Blocks				types.Blocks
	FirstL2BlockNumber  *big.Int
	Contexts			[]SequenceContext
	Txs					types.Transactions
	GasUsed				*big.Int
}

// SequenceBlock represents a block sequenced to L1 sequencer inbox
// Used by validators to reconstruct L2 chain from L1 data.
type SequenceBlock struct {
	SequenceContext
	Txs types.Transactions
}

// SequenceContext is the relavent context of each block sequenced to L1 sequncer inbox
type SequenceContext struct {
	NumTxs      uint64
	Timestamp   uint64
}

type DecodeTxBatchError struct{ msg string }

func (e *DecodeTxBatchError) Error() string {
	return fmt.Sprintf("Failed to create TxBatch from decoded tx data - %s", e.msg)
}

func NewTxBatch(blocks []*types.Block, maxBatchSize uint64) *TxBatch {
	// TODO: handle maxBatchSize constraint
	var contexts []SequenceContext
	var txs []*types.Transaction
	gasUsed := new(big.Int)

	firstL2BlockNumber := blocks[0].Number()

	for _, block := range blocks {
		blockTxs := block.Transactions()
		ctx := SequenceContext{
			NumTxs:      uint64(len(blockTxs)),
			Timestamp:   block.Time(),
		}
		contexts = append(contexts, ctx)
		txs = append(txs, blockTxs...)
		gasUsed.Add(gasUsed, new(big.Int).SetUint64(block.GasUsed()))
	}

	return &TxBatch{blocks, firstL2BlockNumber, contexts, txs, gasUsed}
}

func (b *TxBatch) LastBlockNumber() uint64 {
	return b.FirstL2BlockNumber.Uint64() + uint64(len(b.Contexts)) - 1
}

func (b *TxBatch) LastBlockRoot() common.Hash {
	return b.Blocks[len(b.Blocks)-1].Root()
}

func (b *TxBatch) InboxSize() *big.Int {
	return new(big.Int).SetUint64(uint64(b.Txs.Len()))
}

func (b *TxBatch) SerializeToArgs() ([]*big.Int, []*big.Int, *big.Int, []byte, error) {
	var contexts, txLengths []*big.Int
	for _, ctx := range b.Contexts {
		contexts = append(contexts, big.NewInt(0).SetUint64(ctx.NumTxs))
		contexts = append(contexts, big.NewInt(0).SetUint64(ctx.Timestamp))
	}

	buf := new(bytes.Buffer)
	for _, tx := range b.Txs {
		curLen := buf.Len()
		if err := writeTx(buf, tx); err != nil {
			return nil, nil, nil, nil, err
		}
		txLengths = append(txLengths, big.NewInt(int64(buf.Len()-curLen)))
	}

	firstL2BlockNumber := b.FirstL2BlockNumber

	return contexts, txLengths, firstL2BlockNumber, buf.Bytes(), nil
}

// Splits batch into blocks
func (b *TxBatch) SplitToBlocks() []*SequenceBlock {
	txNum := 0
	blocks := make([]*SequenceBlock, 0, len(b.Contexts))

	for _, ctx := range b.Contexts {
		block := &SequenceBlock{
			SequenceContext: ctx,
			Txs:             b.Txs[txNum : txNum+int(ctx.NumTxs)],
		}
		blocks = append(blocks, block)
		txNum += int(ctx.NumTxs)
	}
	return blocks
}

func writeContext(w *bytes.Buffer, ctx *SequenceContext) error {
	if err := writePrimitive(w, ctx.NumTxs); err != nil {
		return err
	}
	return writePrimitive(w, ctx.Timestamp)
}

func writeTx(w *bytes.Buffer, tx *types.Transaction) error {
	var txBuf bytes.Buffer
	if err := tx.EncodeRLP(&txBuf); err != nil {
		return err
	}
	txBytes := txBuf.Bytes()
	_, err := w.Write(txBytes)
	return err
}

func writePrimitive(w *bytes.Buffer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

// TxBatchFromDecoded decodes the input of SequencerInbox#appendTxBatch call
// It will only fill Contexts and Txs fields
func TxBatchFromDecoded(decoded []interface{}) (*TxBatch, error) {
	if len(decoded) != 4 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid decoded array length %d", len(decoded))}
	}
	contexts := decoded[0].([]*big.Int)
	txLengths := decoded[1].([]*big.Int)
	firstL2BlockNumber := decoded[2].(*big.Int)
	txBatch := decoded[3].([]byte)

	if len(contexts) % 2 != 0 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid contexts length %d", len(contexts))}
	}

	var txs []*types.Transaction
	var ctxs []SequenceContext
	var batchOffset uint64
	var numTxs uint64
	for i := 0; i < len(contexts); i += 2 {
		ctx := SequenceContext{
			NumTxs:      contexts[i].Uint64(),
			Timestamp:   contexts[i+1].Uint64(),
		}
		ctxs = append(ctxs, ctx)
		for j := uint64(0); j < ctx.NumTxs; j++ {
			raw := txBatch[batchOffset : batchOffset+txLengths[numTxs].Uint64()]
			var tx types.Transaction
			err := rlp.DecodeBytes(raw, &tx)
			if err != nil {
				return nil, &DecodeTxBatchError{err.Error()}
			}
			txs = append(txs, &tx)
			numTxs++
			batchOffset += txLengths[numTxs-1].Uint64()
		}
	}
	batch := &TxBatch{
		FirstL2BlockNumber: firstL2BlockNumber,
		Contexts: ctxs,
		Txs:      txs,
	}
	return batch, nil
}
