package client

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

const syncRange uint64 = 10000

type L1BridgeClient interface {
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	BlockNumber(ctx context.Context) (uint64, error)
	ResubscribeErrNewHead(ctx context.Context, sink chan<- *types.Header) (event.Subscription, error)
	SubscribeNewHeadByPolling(
		ctx context.Context,
		headCh chan<- *types.Header,
		tag BlockTag,
		interval time.Duration,
		requestTimeout time.Duration,
	) event.Subscription
	Close()
	// ISequencerInbox.sol
	AppendTxBatch(contexts []*big.Int, txLengths []*big.Int, firstL2BlockNumber *big.Int, txBatch []byte) (*types.Transaction, error)
	WatchTxBatchAppended(opts *bind.WatchOpts, sink chan<- *bindings.ISequencerInboxTxBatchAppended) (event.Subscription, error)
	FilterTxBatchAppendedEvents(opts *bind.FilterOpts) (*bindings.ISequencerInboxTxBatchAppendedIterator, error)
	DecodeAppendTxBatchInput(tx *types.Transaction) ([]interface{}, error)
	// IRollup.sol
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
	// IChallenge.sol
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
		verifierSel uint8,
		proof []byte,
		txInclusionProof []byte,
		verificationRawCtx bindings.VerificationContextRawContext,
		challengedStepIndex *big.Int,
		prevBisection [][32]byte,
		prevChallengedSegmentStart *big.Int,
		prevChallengedSegmentLength *big.Int,
	) (*types.Transaction, error)
	WatchBisected(opts *bind.WatchOpts, sink chan<- *bindings.ISymChallengeBisected) (event.Subscription, error)
	WatchChallengeCompleted(opts *bind.WatchOpts, sink chan<- *bindings.ISymChallengeCompleted) (event.Subscription, error)
	FilterBisected(opts *bind.FilterOpts) (*bindings.ISymChallengeBisectedIterator, error)
	FilterChallengeCompleted(opts *bind.FilterOpts) (*bindings.ISymChallengeCompletedIterator, error)
	DecodeBisectExecutionInput(tx *types.Transaction) ([]interface{}, error)
}

// Basically a thread-safe shim for `ethclient.Client` and `bindings`.
// TODO: clean up retries
type EthBridgeClient struct {
	client       *EthClient
	transactOpts *bind.TransactOpts
	retryOpts    []retry.Option
	// Lock, conservatively on all functions.
	mu sync.Mutex
	// ISequencerInbox.sol
	inboxAbi *abi.ABI
	inbox    *bindings.ISequencerInboxSession
	// IRollup.sol
	rollup *bindings.IRollupSession
	// IChallenge.sol
	// `challenge` initialized separately through `InitNewChallengeSession`
	challengeAbi *abi.ABI
	challenge    *bindings.ISymChallengeSession
}

func NewEthBridgeClient(
	ctx context.Context,
	l1Endpoint string,
	genesisL1Block uint64,
	sequencerInboxAddress common.Address,
	rollupAddress common.Address,
	auth *bind.TransactOpts,
	retryOpts []retry.Option,
) (*EthBridgeClient, error) {
	client, err := DialWithRetry(ctx, l1Endpoint, retryOpts)
	if err != nil {
		return nil, err
	}
	callOpts := bind.CallOpts{Pending: true, Context: ctx}
	transactOpts := bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: big.NewInt(800000000),
		Context:  ctx,
	}
	inbox, err := bindings.NewISequencerInbox(sequencerInboxAddress, client)
	if err != nil {
		return nil, err
	}
	inboxSession := &bindings.ISequencerInboxSession{
		Contract:     inbox,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	rollup, err := bindings.NewIRollup(rollupAddress, client)
	if err != nil {
		return nil, err
	}
	rollupSession := &bindings.IRollupSession{
		Contract:     rollup,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	inboxAbi, err := bindings.ISequencerInboxMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("Failed to get ISequencerInbox ABI, err: %w", err)
	}

	challengeAbi, err := bindings.ISymChallengeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &EthBridgeClient{
		client:       client,
		transactOpts: &transactOpts,
		inboxAbi:     inboxAbi,
		inbox:        inboxSession,
		rollup:       rollupSession,
		challengeAbi: challengeAbi,
	}, nil
}

func (c *EthBridgeClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.client.TransactionByHash(ctx, hash)
}

func (c *EthBridgeClient) BlockNumber(ctx context.Context) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.client.BlockNumber(ctx)
}

func (c *EthBridgeClient) ResubscribeErrNewHead(ctx context.Context, headCh chan<- *types.Header) (event.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sub := event.ResubscribeErr(
		time.Second*10,
		func(ctx context.Context, err error) (event.Subscription, error) {
			if err != nil {
				log.Warn("Error in NewHead subscription, resubscribing", "err", err)
			}
			return c.client.SubscribeNewHead(ctx, headCh)
		},
	)
	return sub, nil
}

func (c *EthBridgeClient) SubscribeNewHeadByPolling(
	ctx context.Context,
	headCh chan<- *types.Header,
	tag BlockTag,
	interval time.Duration,
	requestTimeout time.Duration,
) event.Subscription {
	return event.NewSubscription(func(unsub <-chan struct{}) error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		lastHeaderNumber := big.NewInt(-1)
		for {
			select {
			case <-ticker.C:
				reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
				header, err := c.client.HeaderByTag(reqCtx, tag)
				cancel()
				if err != nil {
					log.Warn("Failed to poll for latest L1 block header", "err", err)
					continue
				}
				if header.Number.Cmp(lastHeaderNumber) <= 0 {
					log.Warn("Polled header is not new", "number", header.Number, "newest", lastHeaderNumber)
					continue
				}
				headCh <- header
				lastHeaderNumber = header.Number
			case <-ctx.Done():
				return ctx.Err()
			case <-unsub:
				return nil
			}
		}
	})
}

func (c *EthBridgeClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.client.Close()
}

func (c *EthBridgeClient) AppendTxBatch(contexts []*big.Int, txLengths []*big.Int, firstL2BlockNumber *big.Int, txBatch []byte) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) {
		return c.inbox.AppendTxBatch(contexts, txLengths, firstL2BlockNumber, txBatch)
	}
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) WatchTxBatchAppended(
	opts *bind.WatchOpts,
	sink chan<- *bindings.ISequencerInboxTxBatchAppended,
) (event.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inbox.Contract.WatchTxBatchAppended(opts, sink)
}

func (c *EthBridgeClient) FilterTxBatchAppendedEvents(
	opts *bind.FilterOpts,
) (*bindings.ISequencerInboxTxBatchAppendedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inbox.Contract.FilterTxBatchAppended(opts)
}

func (c *EthBridgeClient) DecodeAppendTxBatchInput(tx *types.Transaction) ([]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inboxAbi.Methods["appendTxBatch"].Inputs.Unpack(tx.Data()[4:])
}

func (c *EthBridgeClient) GetStaker() (bindings.IRollupStaker, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.GetStaker(c.transactOpts.From)
}

func (c *EthBridgeClient) Stake(amount *big.Int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	log.Info("Staking...")
	opts := c.transactOpts
	opts.Value = amount
	_, err := c.rollup.Contract.Stake(opts)
	if err != nil {
		return fmt.Errorf("Failed to stake, err: %w", err)
	}
	log.Info("Staked successfully.", "amount (ETH)", amount)
	return nil
}

func (c *EthBridgeClient) AdvanceStake(assertionID *big.Int) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) { return c.rollup.AdvanceStake(assertionID) }
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) CreateAssertion(vmHash [32]byte, inboxSize *big.Int) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) {
		return c.rollup.CreateAssertion(vmHash, inboxSize)
	}
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) ChallengeAssertion(players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) { return c.rollup.ChallengeAssertion(players, assertionIDs) }
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) ConfirmFirstUnresolvedAssertion() (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) { return c.rollup.ConfirmFirstUnresolvedAssertion() }
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) RejectFirstUnresolvedAssertion(stakerAddress common.Address) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) { return c.rollup.RejectFirstUnresolvedAssertion(stakerAddress) }
	return retryTransactingFunction(f, c.retryOpts)
}

// Returns the last assertion ID that was validated *by us*.
func (c *EthBridgeClient) GetLastValidatedAssertionID(opts *bind.FilterOpts) (*big.Int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	iter, err := c.rollup.Contract.FilterStakerStaked(opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to filter through `StakerStaked` events to get last validated assertion ID, err: %w", err)
	}
	lastValidatedAssertionID := common.Big0
	for iter.Next() {
		// Note: the second condition should always hold true if the iterator iterates in time order.
		if iter.Event.StakerAddr == c.transactOpts.From && iter.Event.AssertionID.Cmp(lastValidatedAssertionID) == 1 {
			log.Debug("StakerStaked event found", "staker", iter.Event.StakerAddr, "assertionID", iter.Event.AssertionID)
			lastValidatedAssertionID = iter.Event.AssertionID
		}
	}
	if iter.Error() != nil {
		return nil, fmt.Errorf("Failed to iterate through validated assertion IDs, err: %w", iter.Error())
	}
	if lastValidatedAssertionID.Cmp(common.Big0) == 0 {
		return nil, fmt.Errorf("No validated assertion IDs found")
	}
	return lastValidatedAssertionID, nil
}

func (c *EthBridgeClient) GetAssertion(assertionID *big.Int) (bindings.IRollupAssertion, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.GetAssertion(assertionID)
}

func (c *EthBridgeClient) WatchAssertionCreated(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IRollupAssertionCreated,
) (event.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.WatchAssertionCreated(opts, sink)
}

func (c *EthBridgeClient) WatchAssertionChallenged(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IRollupAssertionChallenged,
) (event.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.WatchAssertionChallenged(opts, sink)
}

func (c *EthBridgeClient) WatchAssertionConfirmed(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IRollupAssertionConfirmed,
) (event.Subscription, error) {
	return c.rollup.Contract.WatchAssertionConfirmed(opts, sink)
}

func (c *EthBridgeClient) WatchAssertionRejected(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IRollupAssertionRejected,
) (event.Subscription, error) {
	return c.rollup.Contract.WatchAssertionRejected(opts, sink)
}

func (c *EthBridgeClient) GetAllAssertionCreated(opts *bind.FilterOpts) (*bindings.IRollupAssertionCreatedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.FilterAssertionCreated(opts)
}

func (c *EthBridgeClient) FilterAssertionCreated(opts *bind.FilterOpts) (*bindings.IRollupAssertionCreatedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.FilterAssertionCreated(opts)
}

func (c *EthBridgeClient) FilterAssertionChallenged(
	opts *bind.FilterOpts,
) (*bindings.IRollupAssertionChallengedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.FilterAssertionChallenged(opts)
}

func (c *EthBridgeClient) FilterAssertionConfirmed(
	opts *bind.FilterOpts,
) (*bindings.IRollupAssertionConfirmedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.FilterAssertionConfirmed(opts)
}

func (c *EthBridgeClient) FilterAssertionRejected(
	opts *bind.FilterOpts,
) (*bindings.IRollupAssertionRejectedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.Contract.FilterAssertionRejected(opts)
}

func (c *EthBridgeClient) GetGenesisAssertionCreated(opts *bind.FilterOpts) (*bindings.IRollupAssertionCreated, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// We could probably do this from initialization calldata too.
	iter, err := c.rollup.Contract.FilterAssertionCreated(opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to filter through `AssertionCreated` events to get genesis assertion ID, err: %w", err)
	}
	if iter.Next() {
		return iter.Event, nil
	}
	if iter.Error() != nil {
		return nil, fmt.Errorf("No genesis `AssertionCreated` event found, err: %w", iter.Error())
	}
	return nil, fmt.Errorf("No genesis `AssertionCreated` event found")
}

func (c *EthBridgeClient) InitNewChallengeSession(ctx context.Context, challengeAddress common.Address) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	challenge, err := bindings.NewISymChallenge(challengeAddress, c.client)
	if err != nil {
		return fmt.Errorf("Failed to initialize challenge contract, err: %w", err)
	}
	c.challenge = &bindings.ISymChallengeSession{
		Contract:     challenge,
		CallOpts:     bind.CallOpts{Pending: true, Context: ctx},
		TransactOpts: *c.transactOpts,
	}
	return nil
}

func (c *EthBridgeClient) InitializeChallengeLength(numSteps *big.Int) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) { return c.challenge.InitializeChallengeLength(numSteps) }
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) CurrentChallengeResponder() (common.Address, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challenge.CurrentResponder()
}

func (c *EthBridgeClient) CurrentChallengeResponderTimeLeft() (*big.Int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challenge.CurrentResponderTimeLeft()
}

func (c *EthBridgeClient) TimeoutChallenge() (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) { return c.challenge.Timeout() }
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) BisectExecution(
	bisection [][32]byte,
	challengedSegmentIndex *big.Int,
	prevBisection [][32]byte,
	prevChallengedSegmentStart *big.Int,
	prevChallengedSegmentLength *big.Int,
) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) {
		return c.challenge.BisectExecution(
			bisection,
			challengedSegmentIndex,
			prevBisection,
			prevChallengedSegmentStart,
			prevChallengedSegmentLength,
		)
	}
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) VerifyOneStepProof(
	verifierSel uint8,
	proof []byte,
	txInclusionProof []byte,
	verificationRawCtx bindings.VerificationContextRawContext,
	challengedStepIndex *big.Int,
	prevBisection [][32]byte,
	prevChallengedSegmentStart *big.Int,
	prevChallengedSegmentLength *big.Int,
) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) {
		return c.challenge.VerifyOneStepProof(
			verifierSel,
			proof,
			txInclusionProof,
			verificationRawCtx,
			challengedStepIndex,
			prevBisection,
			prevChallengedSegmentStart,
			prevChallengedSegmentLength,
		)
	}
	return retryTransactingFunction(f, c.retryOpts)
}

func (c *EthBridgeClient) WatchBisected(
	opts *bind.WatchOpts,
	sink chan<- *bindings.ISymChallengeBisected,
) (event.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challenge.Contract.WatchBisected(opts, sink)
}

func (c *EthBridgeClient) WatchChallengeCompleted(
	opts *bind.WatchOpts,
	sink chan<- *bindings.ISymChallengeCompleted,
) (event.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challenge.Contract.WatchCompleted(opts, sink)
}

func (c *EthBridgeClient) FilterBisected(opts *bind.FilterOpts) (*bindings.ISymChallengeBisectedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challenge.Contract.FilterBisected(opts)
}

func (c *EthBridgeClient) FilterChallengeCompleted(
	opts *bind.FilterOpts,
) (*bindings.ISymChallengeCompletedIterator, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challenge.Contract.FilterCompleted(opts)
}

func (c *EthBridgeClient) DecodeBisectExecutionInput(tx *types.Transaction) ([]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.challengeAbi.Methods["bisectExecution"].Inputs.Unpack(tx.Data()[4:])
}

func retryTransactingFunction(f func() (*types.Transaction, error), retryOpts []retry.Option) (*types.Transaction, error) {
	var result *types.Transaction
	var err error
	err = retry.Do(func() error {
		result, err = f()
		return err
	}, retryOpts...)
	return result, err
}
