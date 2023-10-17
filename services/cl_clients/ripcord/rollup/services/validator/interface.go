package validator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/services/cl_clients/ripcord/bindings"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/eth"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/types"
)

type Config interface {
	GetAccountAddr() common.Address
	GetValidationInterval() time.Duration
}

type TxManager interface {
	Stake(ctx context.Context, stakeAmount *big.Int) (*ethTypes.Receipt, error)
	AdvanceStake(ctx context.Context, assertionID *big.Int) (*ethTypes.Receipt, error)
	CreateAssertion(ctx context.Context, vmHash common.Hash, inboxSize *big.Int) (*ethTypes.Receipt, error)
	ConfirmFirstUnresolvedAssertion(ctx context.Context) (*ethTypes.Receipt, error)
}

type BridgeClient interface {
	GetRequiredStakeAmount(ctx context.Context) (*big.Int, error)
	GetStaker(ctx context.Context, addr common.Address) (bindings.IRollupStaker, error)
	GetAssertion(ctx context.Context, assertionID *big.Int) (bindings.IRollupAssertion, error)
	GetLastConfirmedAssertionID(ctx context.Context) (*big.Int, error)
	RequireFirstUnresolvedAssertionIsConfirmable(ctx context.Context) error
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
