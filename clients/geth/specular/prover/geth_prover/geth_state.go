package geth_prover

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

type GethState struct {
	StateDB *state.StateDB
}

func (g GethState) Prepare(thash, bhash common.Hash, ti int) {
	g.StateDB.Prepare(thash, ti)
}

func (g GethState) GetRootForProof() common.Hash {
	return g.StateDB.GetRootForProof()
}

func (g GethState) GetRefund() uint64 {
	return g.StateDB.GetRefund()
}

func (g GethState) CommitForProof() {
	g.StateDB.CommitForProof()
}

func (g GethState) GetCurrentLogs() []*types.Log {
	return g.StateDB.GetCurrentLogs()
}

func (g GethState) GetCode(address common.Address) []byte {
	return g.StateDB.GetCode(address)
}

func (g GethState) GetProof(address common.Address) ([][]byte, error) {
	return g.StateDB.GetProof(address)
}

func (g GethState) GetStorageProof(address common.Address, hash common.Hash) ([][]byte, error) {
	return g.StateDB.GetStorageProof(address, hash)
}

func (g GethState) SubBalance(address common.Address, b *big.Int) {
	g.StateDB.SubBalance(address, b)
}

func (g GethState) SetNonce(address common.Address, u uint64) {
	g.StateDB.SetNonce(address, u)
}

func (g GethState) GetNonce(address common.Address) uint64 {
	return g.StateDB.GetNonce(address)
}

func (g GethState) AddBalance(address common.Address, b *big.Int) {
	g.StateDB.AddBalance(address, b)
}

func (g GethState) DeleteSuicidedAccountForProof(addr common.Address) {
	g.StateDB.DeleteSuicidedAccountForProof(addr)
}

func (g GethState) SetCode(address common.Address, bytes []byte) {
	g.StateDB.SetCode(address, bytes)
}

func (g GethState) GetBalance(address common.Address) *big.Int {
	return g.StateDB.GetBalance(address)
}

func (g GethState) GetCodeHash(address common.Address) common.Hash {
	return g.StateDB.GetCodeHash(address)
}

func (g GethState) Copy() prover_types.L2ELClientStateInterface {
	return &GethState{StateDB: g.StateDB.Copy()}
}
