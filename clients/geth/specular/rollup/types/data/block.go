package data

import (
	"fmt"
	"math/big"
)

// `DerivationBlock` represents a block from which the L2 chain can be derived.
type DerivationBlock struct {
	DerivationContext
	txs []EncodedTransaction
}

type EncodedTransaction []byte

// DerivationContext is the relevent context of each block sequenced to
// L1 SeqeuncerInbox to ensure deterministic recomputation.
type DerivationContext struct {
	numTxs      uint64
	blockNumber uint64
	timestamp   uint64
}

type DecodeTxBatchError struct{ msg string }

func NewDerivationBlock(blockNumber, timestamp uint64, txs []EncodedTransaction) DerivationBlock {
	return DerivationBlock{
		DerivationContext: DerivationContext{
			numTxs:      uint64(len(txs)),
			blockNumber: blockNumber,
			timestamp:   timestamp,
		},
		txs: txs,
	}
}

func (b *DerivationBlock) BlockNumber() uint64 {
	return b.blockNumber
}

func (b *DerivationBlock) Timestamp() uint64 {
	return b.timestamp
}

func (b *DerivationBlock) Txs() []EncodedTransaction {
	return b.txs
}

func (e *DecodeTxBatchError) Error() string {
	return fmt.Sprintf("Failed to create TxBatch from decoded tx data - %s", e.msg)
}

// Decodes the input of `SequencerInbox.appendTxBatch` call
func BlocksFromDecoded(decoded []interface{}) ([]*DerivationBlock, error) {
	if len(decoded) != 3 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid decoded array length %d", len(decoded))}
	}
	contexts := decoded[0].([]*big.Int)
	txLengths := decoded[1].([]*big.Int)
	txBatch := decoded[2].([]byte)
	if len(contexts)%3 != 0 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid contexts length %d", len(contexts))}
	}

	var batchOffset uint64
	var numTxs uint64

	blocks := make([]*DerivationBlock, 0, len(contexts)/3)
	for i := 0; i < len(contexts); i += 3 {
		// For each context, decode a block.
		var txs []EncodedTransaction
		ctx := DerivationContext{
			numTxs:      contexts[i].Uint64(),
			blockNumber: contexts[i+1].Uint64(),
			timestamp:   contexts[i+2].Uint64(),
		}
		for j := uint64(0); j < ctx.numTxs; j++ {
			encodedTx := txBatch[batchOffset : batchOffset+txLengths[numTxs].Uint64()]
			txs = append(txs, encodedTx)
			numTxs++
			batchOffset += txLengths[numTxs-1].Uint64()
		}
		blocks = append(blocks, &DerivationBlock{ctx, txs})
	}
	return blocks, nil
}
