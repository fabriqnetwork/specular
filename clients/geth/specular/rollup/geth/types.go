package geth

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Header

type Header struct {
	header *types.Header
}

func NewHeader(header *types.Header) *Header {
	return &Header{header: header}
}

func (h *Header) Hash() common.Hash {
	return h.header.Hash()
}

func (h *Header) ParentHash() common.Hash {
	return h.header.ParentHash
}

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
