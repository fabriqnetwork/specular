package geth_prover

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

type GethTransaction struct {
	*types.Transaction
}

func (tx *GethTransaction) AsMessage(s interface{}, baseFee *big.Int, rules interface{}) (interface{}, error) {
	return tx.Transaction.AsMessage(s.(types.Signer), baseFee)
}

type GethBlock struct {
	*types.Block
}

func (b *GethBlock) Transaction(hash common.Hash) prover_types.Transaction {
	return &GethTransaction{b.Block.Transaction(hash)}
}

func (b *GethBlock) Transactions() prover_types.Transactions {
	geth_txs := b.Block.Transactions()
	txs := make(prover_types.Transactions, len(geth_txs))
	for i, tx := range geth_txs {
		txs[i] = &GethTransaction{tx}
	}
	return txs
}
