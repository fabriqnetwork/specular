package geth_prover

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	my_state "github.com/specularl2/specular/clients/geth/specular/prover/state"
)

type GethEVM struct {
	EVM *vm.EVM
}

// implements L2ELClientEVMInterface
func (e *GethEVM) Context() my_state.L2ELClientBlockContextInterface {
	return GethBlockContext{e.EVM.Context}
}

func (e *GethEVM) StateDB() my_state.L2ELClientStateInterface {
	// assert EVM.StateDB is a *state.StateDB, as the prover only works with state.StateDB
	return GethState{StateDB: e.EVM.StateDB.(*state.StateDB)}
}
