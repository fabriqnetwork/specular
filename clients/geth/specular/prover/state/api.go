package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type L2ELClientStateInterface interface {
	Prepare(thash common.Hash, ti int)
	Copy() L2ELClientStateInterface
	GetRootForProof() common.Hash
	GetRefund() uint64
	CommitForProof()
	GetCurrentLogs() []*types.Log
	GetCode(address common.Address) []byte
	GetProof(common.Address) ([][]byte, error)
	GetStorageProof(common.Address, common.Hash) ([][]byte, error)
	SubBalance(common.Address, *big.Int)
	SetNonce(common.Address, uint64)
	GetNonce(common.Address) uint64
	AddBalance(common.Address, *big.Int)
	DeleteSuicidedAccountForProof(addr common.Address)
	SetCode(common.Address, []byte)
	GetBalance(common.Address) *big.Int
	GetCodeHash(common.Address) common.Hash
}

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(L2ELClientStateInterface, common.Address, *big.Int) bool
	// GetHashFunc returns the n'th block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

type L2ELClientBlockContextInterface interface {
	CanTransfer() CanTransferFunc
	GetHash() GetHashFunc
	Coinbase() common.Address
	GasLimit() uint64
	BlockNumber() *big.Int
	Time() *big.Int
	Difficulty() *big.Int
	BaseFee() *big.Int
	Random() *common.Hash
}

type SpecularEVMLoggerInterface interface {
}

// Config are the configuration options for the Interpreter
type L2ELClientConfig struct {
	Debug  bool                       // Enables debugging
	Tracer SpecularEVMLoggerInterface // Opcode logger
}

type L2ELClientEVMInterface interface {
	// Context provides auxiliary blockchain related information
	Context() L2ELClientBlockContextInterface
	// StateDB gives access to the underlying state
	StateDB() L2ELClientStateInterface
}
