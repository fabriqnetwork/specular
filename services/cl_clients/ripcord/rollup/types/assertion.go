package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Assertion represents disputable assertion for L1 rollup contract
type Assertion struct {
	ID         *big.Int
	VmHash     common.Hash
	InboxSize  *big.Int
	Deadline   *big.Int
	StartBlock uint64
	EndBlock   uint64
}

func (a *Assertion) Copy() *Assertion {
	return &Assertion{
		ID:         new(big.Int).Set(a.ID),
		VmHash:     a.VmHash,
		InboxSize:  new(big.Int).Set(a.InboxSize),
		Deadline:   new(big.Int).Set(a.Deadline),
		StartBlock: a.StartBlock,
		EndBlock:   a.EndBlock,
	}
}
