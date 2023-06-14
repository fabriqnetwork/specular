package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Assertion represents disputable assertion for L1 rollup contract
type Assertion struct {
	ID                    *big.Int
	VmHash                common.Hash
	CumulativeGasUsed     *big.Int
	InboxSize             *big.Int
	Deadline              *big.Int
	StartBlock            uint64
	EndBlock              uint64
	PrevCumulativeGasUsed *big.Int
}

func (a *Assertion) Copy() *Assertion {
	copied := &Assertion{
		ID:                    new(big.Int).Set(a.ID),
		VmHash:                a.VmHash,
		CumulativeGasUsed:     new(big.Int),
		InboxSize:             new(big.Int).Set(a.InboxSize),
		Deadline:              new(big.Int).Set(a.Deadline),
		StartBlock:            a.StartBlock,
		EndBlock:              a.EndBlock,
		PrevCumulativeGasUsed: new(big.Int),
	}
	if a.CumulativeGasUsed != nil {
		copied.CumulativeGasUsed.Set(a.CumulativeGasUsed)
	}
	if a.PrevCumulativeGasUsed != nil {
		copied.PrevCumulativeGasUsed.Set(a.PrevCumulativeGasUsed)
	}
	return copied
}
