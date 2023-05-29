package validator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation"
	"golang.org/x/sync/errgroup"
)

type Config interface {
	AccountAddr() common.Address
	IsActiveStaker() bool
	IsActiveCreator() bool
	IsActiveChallenger() bool
	IsResolver() bool
	StakeAmount() uint64
}

type BaseService interface {
	Start() context.Context
	Stop() error
	Eg() *errgroup.Group
}

type L1Config interface {
	Endpoint() string
	ChainID() uint64
}

type RollupState interface {
	GetStaker(ctx context.Context, stakerAddr common.Address) (*derivation.Staker, error)
	GetAssertion(ctx context.Context, assertionID *big.Int) (*derivation.Assertion, error)
}

type TxManager interface {
	Stake(ctx context.Context, stakeAmount *big.Int) (*types.Receipt, error)
	AdvanceStake(ctx context.Context, assertionID *big.Int) (*types.Receipt, error)
	CreateAssertion(ctx context.Context, vmHash common.Hash, inboxSize *big.Int) (*types.Receipt, error)
	ConfirmFirstUnresolvedAssertion(ctx context.Context) (*types.Receipt, error)
	RejectFirstUnresolvedAssertion(ctx context.Context, stakerAddress common.Address) (*types.Receipt, error)
}

type L2Client interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type ChallengeClient interface {
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)

	InitNewChallengeSession(ctx context.Context, challengeAddress common.Address) error
	InitializeChallengeLength(numSteps *big.Int) (*types.Transaction, error)
	CurrentChallengeResponder() (common.Address, error)
	CurrentChallengeResponderTimeLeft() (*big.Int, error)
	TimeoutChallenge() (*types.Transaction, error)
	BisectExecution(
		bisection [][32]byte,
		challengedSegmentIndex *big.Int,
		prevBisection [][32]byte,
		prevChallengedSegmentStart *big.Int,
		prevChallengedSegmentLength *big.Int,
	) (*types.Transaction, error)
	VerifyOneStepProof(
		proof []byte,
		txInclusionProof []byte,
		verificationRawCtx bindings.VerificationContextLibRawContext,
		challengedStepIndex *big.Int,
		prevBisection [][32]byte,
		prevChallengedSegmentStart *big.Int,
		prevChallengedSegmentLength *big.Int,
	) (*types.Transaction, error)
	WatchChallengeCompleted(opts *bind.WatchOpts, sink chan<- *bindings.ISymChallengeCompleted) (event.Subscription, error)
	FilterChallengeCompleted(opts *bind.FilterOpts) (*bindings.ISymChallengeCompletedIterator, error)
}
