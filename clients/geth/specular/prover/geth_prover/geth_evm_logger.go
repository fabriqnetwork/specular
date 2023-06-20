package geth_prover

import "github.com/ethereum/go-ethereum/core/vm"

type GethEVMLogger struct {
	Logger *vm.EVMLogger
}
