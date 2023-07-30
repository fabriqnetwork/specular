package geth_prover

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

type GethEVM struct {
	EVM *vm.EVM
}

// implements L2ELClientEVMInterface
func (e *GethEVM) Context() prover_types.L2ELClientBlockContextInterface {
	return GethBlockContext{e.EVM.Context}
}

func (e *GethEVM) StateDB() prover_types.L2ELClientStateInterface {
	// assert EVM.StateDB is a *state.StateDB, as the prover only works with state.StateDB
	return GethState{StateDB: e.EVM.StateDB.(*state.StateDB)}
}

type GethContract struct {
	*vm.Contract
}

func (c *GethContract) Code() []byte {
	return c.Contract.Code
}

type GethScopeContext struct {
	ScopeContext *vm.ScopeContext
}

func (s *GethScopeContext) Memory() prover_types.L2ELClientMemoryInterface {
	return s.ScopeContext.Memory
}

func (s *GethScopeContext) Stack() prover_types.L2ELClientStackInterface {
	return s.ScopeContext.Stack
}

func (s *GethScopeContext) Contract() prover_types.L2ELClientContractInterface {
	return &GethContract{s.ScopeContext.Contract}
}
