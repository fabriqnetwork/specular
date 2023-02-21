package testutils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From      common.Address  // the sender of the 'transaction'
	To        *common.Address // the destination contract (nil for contract creation)
	Gas       uint64          // if 0, the call executes with near-infinite gas
	GasPrice  *big.Int        // wei <-> gas exchange ratio
	GasFeeCap *big.Int        // EIP-1559 fee cap per gas.
	GasTipCap *big.Int        // EIP-1559 tip per gas.
	Value     *big.Int        // amount of wei sent along with the call
	Data      []byte          // input data, usually an ABI-encoded contract method invocation

	AccessList types.AccessList // EIP-2930 access list.
}

type EthClient interface {
	Close()
	ChainID(ctx context.Context) (*big.Int, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	PeerCount(ctx context.Context) (uint64, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	SyncProgress(ctx context.Context) (*interface{}, error)
	// Adding
	CallContract(ctx context.Context, call CallMsg, blockNumber *big.Int) ([]byte, error)
}
