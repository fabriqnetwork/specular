package geth

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Engine

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type PayloadStatus = beacon.PayloadStatusV1
type PayloadID = beacon.PayloadID

var InvalidForkChoiceState = beacon.InvalidForkChoiceState
var STATUS_INVALID = beacon.STATUS_INVALID

var (
	// VALID is returned by the engine API in the following calls:
	//   - newPayloadV1:       if the payload was already known or was just validated and executed
	//   - forkchoiceUpdateV1: if the chain accepted the reorg (might ignore if it's stale)
	VALID = "VALID"

	// INVALID is returned by the engine API in the following calls:
	//   - newPayloadV1:       if the payload failed to execute on top of the local chain
	//   - forkchoiceUpdateV1: if the new head is unknown, pre-merge, or reorg to it fails
	INVALID = "INVALID"
)

// Header

type Header struct{ h *types.Header }

func NewHeader(header *types.Header) *Header { return &Header{h: header} }
func (h *Header) Hash() common.Hash          { return h.h.Hash() }
func (h *Header) ParentHash() common.Hash    { return h.h.ParentHash }

// Transaction

func EncodeRLP(txs types.Transactions) ([][]byte, error) {
	var encodedTxs [][]byte
	for _, tx := range txs {
		var txBuf bytes.Buffer
		if err := tx.EncodeRLP(&txBuf); err != nil {
			return nil, err
		}
		encodedTxs = append(encodedTxs, txBuf.Bytes())
	}
	return encodedTxs, nil
}

func DecodeRLP(txs [][]byte) ([]*types.Transaction, error) {
	var decodedTxs []*types.Transaction
	for _, tx := range txs {
		// TODO: use tx.DecodeRLP instead?
		var decodedTx *types.Transaction
		err := rlp.DecodeBytes(tx, decodedTx)
		if err != nil {
			return nil, err
		}
		decodedTxs = append(decodedTxs, decodedTx)
	}
	return decodedTxs, nil
}
