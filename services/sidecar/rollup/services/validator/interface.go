package validator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/specularL2/specular/services/sidecar/bindings"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/bridge"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
)

type Config interface {
	GetAccountAddr() common.Address
	GetValidationInterval() time.Duration
}

type TxManager interface {
	Stake(ctx context.Context, stakeAmount *big.Int) (*ethTypes.Receipt, error)
	CreateAssertion(
		ctx context.Context,
		stateCommitment common.Hash,
		blockNum *big.Int,
		l1BlockHash common.Hash,
		l1BlockNum *big.Int,
	) (*ethTypes.Receipt, error)
	ConfirmFirstUnresolvedAssertion(ctx context.Context) (*ethTypes.Receipt, error)
	RejectFirstUnresolvedAssertion(context.Context, common.Address) (*ethTypes.Receipt, error)
	RemoveStake(context.Context, common.Address) (*ethTypes.Receipt, error)
}

type BridgeClient interface {
	GetRequiredStakeAmount(context.Context) (*big.Int, error)
	GetStaker(context.Context, common.Address) (bindings.IRollupStaker, error)
	GetAssertion(context.Context, *big.Int) (bindings.IRollupAssertion, error)
	RequireFirstUnresolvedAssertionIsConfirmable(context.Context) (*bridge.UnsatisfiedCondition, error)
	RequireFirstUnresolvedAssertionIsRejectable(context.Context, common.Address) (*bridge.UnsatisfiedCondition, error)
}

type EthState interface {
	Head() types.BlockID
	Safe() types.BlockID
	Finalized() types.BlockID
}

type L2Client interface {
	EnsureDialed(ctx context.Context) error
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}

type ErrGroup interface{ Go(f func() error) }
