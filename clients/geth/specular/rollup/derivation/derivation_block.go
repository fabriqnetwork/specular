package derivation

import (
	"math/big"

	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// `DerivationBlock` represents a block from which the L2 chain can be derived.
type DerivationBlock struct {
	DerivationContext
	txs [][]byte
}

// DerivationContext is the relevent context of each block sequenced to
// L1 SeqeuncerInbox to ensure deterministic recomputation.
type DerivationContext struct {
	numTxs      uint64
	blockNumber uint64
	timestamp   uint64
}

func (c *DerivationContext) NumTxs() uint64      { return c.numTxs }
func (c *DerivationContext) BlockNumber() uint64 { return c.blockNumber }
func (c *DerivationContext) Timestamp() uint64   { return c.timestamp }

func NewDerivationBlock(blockNumber, timestamp uint64, txs [][]byte) DerivationBlock {
	return DerivationBlock{
		DerivationContext: DerivationContext{
			numTxs:      uint64(len(txs)),
			blockNumber: blockNumber,
			timestamp:   timestamp,
		},
		txs: txs,
	}
}

func (b *DerivationBlock) BlockNumber() uint64 { return b.blockNumber }
func (b *DerivationBlock) Timestamp() uint64   { return b.timestamp }
func (b *DerivationBlock) Txs() [][]byte       { return b.txs }

type DecodeTxBatchError struct{ msg string }

func (e *DecodeTxBatchError) Error() string {
	return fmt.Sprintf("failed to create TxBatch from decoded tx data - %s", e.msg)
}

// Decodes the input of `SequencerInbox.appendTxBatch` call
func BlocksFromData(calldata []any) ([]DerivationBlock, error) {
	if len(calldata) != 4 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid decoded array length %d", len(calldata))}
	}
	var (
		contexts           = calldata[0].([]*big.Int)
		txLengths          = calldata[1].([]*big.Int)
		firstL2BlockNumber = calldata[2].(*big.Int)
		// TODO: commented until multiple versions have been used
		//txBatchVersion     = calldata[3].(*big.Int)
		txBatch = calldata[4].([]byte)
	)
	if len(contexts)%2 != 0 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid contexts length %d", len(contexts))}
	}
	var (
		batchOffset, numTxs uint64
		blocks              = make([]DerivationBlock, 0, len(contexts)/2)
		blockNumber         = firstL2BlockNumber.Uint64()
	)
	for i := 0; i < len(contexts); i += 2 {
		// For each context, decode a block.
		ctx := DerivationContext{
			numTxs:      contexts[i].Uint64(),
			blockNumber: blockNumber + uint64(i),
			timestamp:   contexts[i+1].Uint64(),
		}
		var txs [][]byte
		for j := uint64(0); j < ctx.numTxs; j++ {
			encodedTx := txBatch[batchOffset : batchOffset+txLengths[numTxs].Uint64()]
			txs = append(txs, encodedTx)
			numTxs++
			batchOffset += txLengths[numTxs-1].Uint64()
		}
		log.Trace("Block decoded", "block#", ctx.blockNumber, "numTxs", ctx.numTxs)
		blocks = append(blocks, DerivationBlock{ctx, txs})
	}
	return blocks, nil
}
