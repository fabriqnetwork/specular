package validator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
)

type Config interface{ GetAssertInterval() time.Duration }

type TxManager interface {
	Stake(ctx context.Context, stakeAmount *big.Int) (*types.Receipt, error)
	AdvanceStake(ctx context.Context, assertionID *big.Int) (*types.Receipt, error)
	CreateAssertion(ctx context.Context, vmHash common.Hash, inboxSize *big.Int) (*types.Receipt, error)
	ConfirmFirstUnresolvedAssertion(ctx context.Context) (*types.Receipt, error)
}

type L2Client interface {
	EnsureDialed(ctx context.Context) error
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}
