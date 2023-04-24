package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
)

// Client to both the L1 chain and the bridge contracts deployed on it.
type EthBridgeClient interface {
	EthClient
	BridgeClient
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	ResubscribeErrNewHead(
		ctx context.Context,
		sink chan<- *types.Header,
	) (event.Subscription, error)
	Close()
}

type BridgeClient interface {
	SequencerInboxClient
	RollupClient
	ChallengeClient
}

type SequencerInboxClient interface {
	WatchTxBatchAppended(opts *bind.WatchOpts, sink chan<- *bindings.ISequencerInboxTxBatchAppended) (event.Subscription, error)
	FilterTxBatchAppendedEvents(opts *bind.FilterOpts) (*bindings.ISequencerInboxTxBatchAppendedIterator, error)
}

type RollupClient interface {
	Stake(amount *big.Int) error
	GetStaker() (bindings.IRollupStaker, error)
	AdvanceStake(assertionID *big.Int) (*types.Transaction, error)
	CreateAssertion(vmHash [32]byte, inboxSize *big.Int) (*types.Transaction, error)
	ChallengeAssertion(players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error)
	ConfirmFirstUnresolvedAssertion() (*types.Transaction, error)
	RejectFirstUnresolvedAssertion(stakerAddress common.Address) (*types.Transaction, error)
	GetLastValidatedAssertionID(opts *bind.FilterOpts) (*big.Int, error)
	GetAssertion(assertionID *big.Int) (bindings.IRollupAssertion, error)
	WatchAssertionCreated(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionCreated) (event.Subscription, error)
	WatchAssertionChallenged(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionChallenged) (event.Subscription, error)
	WatchAssertionConfirmed(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionConfirmed) (event.Subscription, error)
	WatchAssertionRejected(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionRejected) (event.Subscription, error)
	FilterAssertionCreated(opts *bind.FilterOpts) (*bindings.IRollupAssertionCreatedIterator, error)
	FilterAssertionChallenged(opts *bind.FilterOpts) (*bindings.IRollupAssertionChallengedIterator, error)
	FilterAssertionConfirmed(opts *bind.FilterOpts) (*bindings.IRollupAssertionConfirmedIterator, error)
	FilterAssertionRejected(opts *bind.FilterOpts) (*bindings.IRollupAssertionRejectedIterator, error)
	GetGenesisAssertionCreated(opts *bind.FilterOpts) (*bindings.IRollupAssertionCreated, error)
}

type ChallengeClient interface {
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
	WatchBisected(opts *bind.WatchOpts, sink chan<- *bindings.ISymChallengeBisected) (event.Subscription, error)
	WatchChallengeCompleted(opts *bind.WatchOpts, sink chan<- *bindings.ISymChallengeCompleted) (event.Subscription, error)
	FilterBisected(opts *bind.FilterOpts) (*bindings.ISymChallengeBisectedIterator, error)
	FilterChallengeCompleted(opts *bind.FilterOpts) (*bindings.ISymChallengeCompletedIterator, error)
}
