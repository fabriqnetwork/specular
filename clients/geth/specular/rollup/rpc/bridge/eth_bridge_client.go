package bridge

import (
	"context"
	"math/big"
	"sync"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

// Basically a thread-safe shim for `ethclient.Client` and `bindings`.
// TODO: depricated; delete
type EthBridgeClient struct {
	client       *eth.EthClient
	transactOpts *bind.TransactOpts
	retryOpts    []retry.Option
	// Lock, conservatively on all functions.
	mu    sync.Mutex
	inbox *bindings.ISequencerInboxSession
	// IRollup.sol
	rollup *bindings.IRollupSession
	// IChallenge.sol
	// `challenge` initialized separately through `InitNewChallengeSession`
	challenge *bindings.ISymChallengeSession
}

func NewEthBridgeClient(
	ctx context.Context,
	l1Client *eth.EthClient,
	l1Endpoint string,
	genesisL1Block uint64,
	sequencerInboxAddress common.Address,
	rollupAddress common.Address,
	auth *bind.TransactOpts,
	retryOpts []retry.Option,
) (*EthBridgeClient, error) {
	if l1Client == nil {
		var err error
		l1Client, err = eth.DialWithRetry(ctx, l1Endpoint, retryOpts)
		if err != nil {
			return nil, err
		}
	}
	callOpts := bind.CallOpts{Pending: true, Context: ctx}
	transactOpts := bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: big.NewInt(800000000),
		Context:  ctx,
	}
	inbox, err := bindings.NewISequencerInbox(sequencerInboxAddress, l1Client)
	if err != nil {
		return nil, err
	}
	inboxSession := &bindings.ISequencerInboxSession{
		Contract:     inbox,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	rollup, err := bindings.NewIRollup(rollupAddress, l1Client)
	if err != nil {
		return nil, err
	}
	rollupSession := &bindings.IRollupSession{
		Contract:     rollup,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}

	return &EthBridgeClient{
		client:       l1Client,
		transactOpts: &transactOpts,
		inbox:        inboxSession,
		rollup:       rollupSession,
	}, nil
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

func (c *EthBridgeClient) GetStaker() (bindings.IRollupStaker, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rollup.GetStaker(c.transactOpts.From)
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
