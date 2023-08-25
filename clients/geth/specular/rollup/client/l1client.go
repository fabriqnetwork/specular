package client

import (
	"context"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
)

// TODO: delete this file
type L1BridgeClient interface {
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*types.Header, error)
	BlockNumber(ctx context.Context) (uint64, error)
	FilterTxBatchAppendedEvents(opts *bind.FilterOpts) (*bindings.ISequencerInboxTxBatchAppendedIterator, error)
}

type EthBridgeClient struct {
	*eth.EthClient
	inbox *bindings.ISequencerInboxSession
}

func NewEthBridgeClient(
	ctx context.Context,
	l1Endpoint string,
	sequencerInboxAddress common.Address,
	retryOpts []retry.Option,
) (*EthBridgeClient, error) {
	client, err := eth.DialWithRetry(ctx, l1Endpoint, retryOpts...)
	if err != nil {
		return nil, err
	}
	inbox, err := bindings.NewISequencerInbox(sequencerInboxAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := bind.CallOpts{Pending: true, Context: ctx}
	inboxSession := &bindings.ISequencerInboxSession{Contract: inbox, CallOpts: callOpts}
	return &EthBridgeClient{EthClient: client, inbox: inboxSession}, nil
}

func (c *EthBridgeClient) FilterTxBatchAppendedEvents(
	opts *bind.FilterOpts,
) (*bindings.ISequencerInboxTxBatchAppendedIterator, error) {
	return c.inbox.Contract.FilterTxBatchAppended(opts)
}
<<<<<<< HEAD
=======

func (c *EthBridgeClient) DecodeAppendTxBatchInput(tx *types.Transaction) ([]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inboxAbi.Methods["appendTxBatch"].Inputs.Unpack(tx.Data()[5:])
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

func (c *EthBridgeClient) VerifyOneStepProof(
	proof []byte,
	txInclusionProof []byte,
	verificationRawCtx bindings.VerificationContextLibRawContext,
	challengedStepIndex *big.Int,
	prevBisection [][32]byte,
	prevChallengedSegmentStart *big.Int,
	prevChallengedSegmentLength *big.Int,
) (*types.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f := func() (*types.Transaction, error) {
		return c.challenge.VerifyOneStepProof(
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

func retryTransactingFunction(f func() (*types.Transaction, error), retryOpts []retry.Option) (*types.Transaction, error) {
	var result *types.Transaction
	var err error
	err = retry.Do(func() error {
		result, err = f()
		return err
	}, retryOpts...)
	return result, err
}
>>>>>>> 64d7b5b (clean up)
