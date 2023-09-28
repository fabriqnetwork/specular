package derivation

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
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
	if len(calldata) != 2 {
		return nil, &DecodeTxBatchError{fmt.Sprintf("invalid decoded array length %d", len(calldata))}
	}
	var (
		// TODO: commented until multiple versions have been used
		//txBatchVersion     = calldata[0].(*big.Int)
		txBatch = calldata[1].([]byte)
	)

	var decodedBatch BatchAttributes
	if err := rlp.Decode(bytes.NewReader(txBatch), &decodedBatch); err != nil {
		return nil, err
	}

	blocks := make([]DerivationBlock, 0, len(decodedBatch.Blocks))
	for i, block := range decodedBatch.Blocks {
		// For each context, decode a block.
		ctx := DerivationContext{
			numTxs:      uint64(len(block.Txs)),
			blockNumber: decodedBatch.FirstL2BlockNumber.Uint64() + uint64(i),
			timestamp:   block.Timestamp.Uint64(),
		}
		log.Trace("Block decoded", "block#", ctx.blockNumber, "numTxs", ctx.numTxs)
		blocks = append(blocks, DerivationBlock{ctx, block.Txs})
	}
	return blocks, nil
}
