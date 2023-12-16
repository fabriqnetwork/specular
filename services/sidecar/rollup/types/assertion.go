package types

import (
	"math/big"
)

type Bytes32 [32]byte

// Assertion represents disputable assertion for L1 rollup contract
type Assertion struct {
	ID              *big.Int
	StateCommitment Bytes32
	BlockNum        *big.Int
	Deadline        *big.Int
	StartBlock      uint64
	EndBlock        uint64
}

func (a *Assertion) Copy() *Assertion {
	return &Assertion{
		ID:              new(big.Int).Set(a.ID),
		StateCommitment: a.StateCommitment,
		BlockNum:        new(big.Int).Set(a.BlockNum),
		Deadline:        new(big.Int).Set(a.Deadline),
		StartBlock:      a.StartBlock,
		EndBlock:        a.EndBlock,
	}
}
