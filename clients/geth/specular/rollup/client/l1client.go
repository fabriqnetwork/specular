package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

const syncRange uint64 = 10000

type L1BridgeClient interface {
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	BlockNumber(ctx context.Context) (uint64, error)
	SubscribeNewHead(ctx context.Context, sink chan<- *types.Header) (ethereum.Subscription, error)
	Close()
	// ISequencerInbox.sol
	AppendTxBatch(contexts []*big.Int, txLengths []*big.Int, txBatch []byte) (*types.Transaction, error)
	WatchTxBatchAppended(opts *bind.WatchOpts, sink chan<- *bindings.ISequencerInboxTxBatchAppended) (event.Subscription, error)
	FilterTxBatchAppendedEvents(opts *bind.FilterOpts) (*bindings.ISequencerInboxTxBatchAppendedIterator, error)
	DecodeAppendTxBatchInput(tx *types.Transaction) ([]interface{}, error)
	// IRollup.sol
	IsStaked(context.Context) (bool, error)
	Stake(amount *big.Int) error
	AdvanceStake(assertionID *big.Int) (*types.Transaction, error)
	CreateAssertion(
		vmHash [32]byte,
		inboxSize *big.Int,
		cumulativeGasUsed *big.Int,
		prevVMHash common.Hash,
		prevL2GasUsed *big.Int,
	) (*types.Transaction, error)
	ChallengeAssertion(players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error)
	WatchAssertionCreated(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionCreated) (event.Subscription, error)
	WatchAssertionChallenged(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionChallenged) (event.Subscription, error)
	WatchAssertionConfirmed(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionConfirmed) (event.Subscription, error)
	WatchAssertionRejected(opts *bind.WatchOpts, sink chan<- *bindings.IRollupAssertionRejected) (event.Subscription, error)
	GetLastValidatedAssertionID(ctx context.Context, start uint64) (*big.Int, error)
	ConfirmFirstUnresolvedAssertion() (*types.Transaction, error)
	// AssertionMap.sol
	GetAssertion(assertionID *big.Int) (*rollupTypes.Assertion, error)
	GetDeadline(assertionID *big.Int) (*big.Int, error)
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
		proof []byte,
		challengedStepIndex *big.Int,
		prevBisection [][32]byte,
		prevChallengedSegmentStart *big.Int,
		prevChallengedSegmentLength *big.Int,
	) (*types.Transaction, error)
	WatchBisected(opts *bind.WatchOpts, sink chan<- *bindings.IChallengeBisected) (event.Subscription, error)
	WatchChallengeCompleted(opts *bind.WatchOpts, sink chan<- *bindings.IChallengeChallengeCompleted) (event.Subscription, error)
	DecodeBisectExecutionInput(tx *types.Transaction) ([]interface{}, error)
}

// Basically a shim for `ethclient.Client` and `bindings`.
type EthBridgeClient struct {
	client       *ethclient.Client
	transactOpts *bind.TransactOpts
	// ISequencerInbox.sol
	inboxAbi *abi.ABI
	inbox    *bindings.ISequencerInboxSession
	// IRollup.sol
	rollup *bindings.IRollupSession
	// AssertionMap.sol
	assertionMap *bindings.AssertionMapCallerSession
	// IChallenge.sol
	// `challenge` initialized separately through `InitNewChallengeSession`
	challengeAbi *abi.ABI
	challenge    *bindings.IChallengeSession
}

type ContractAddressBook struct {
	SequencerInboxAddress common.Address
	RollupAddress         common.Address
}

func NewEthBridgeClient(
	ctx context.Context,
	l1Endpoint string,
	genesisL1Block uint64,
	sequencerInboxAddress common.Address,
	rollupAddress common.Address,
	auth *bind.TransactOpts,
) (*EthBridgeClient, error) {
	client, err := ethclient.DialContext(ctx, l1Endpoint)
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
	assertionMapAddr, err := rollupSession.Assertions()
	if err != nil {
		return nil, err
	}
	assertionMap, err := bindings.NewAssertionMapCaller(assertionMapAddr, client)
	if err != nil {
		return nil, err
	}
	assertionMapSession := &bindings.AssertionMapCallerSession{Contract: assertionMap, CallOpts: callOpts}
	inboxAbi, err := bindings.ISequencerInboxMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("Failed to get ISequencerInbox ABI, err: %w", err)
	}

	challengeAbi, err := bindings.IChallengeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &EthBridgeClient{
		client:       client,
		transactOpts: &transactOpts,
		inboxAbi:     inboxAbi,
		inbox:        inboxSession,
		rollup:       rollupSession,
		assertionMap: assertionMapSession,
		challengeAbi: challengeAbi,
	}, nil
}

func (c *EthBridgeClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return c.client.TransactionByHash(ctx, hash)
}

func (c *EthBridgeClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.client.BlockNumber(ctx)
}

func (c *EthBridgeClient) SubscribeNewHead(ctx context.Context, headCh chan<- *types.Header) (ethereum.Subscription, error) {
	return c.client.SubscribeNewHead(ctx, headCh)
}

func (c *EthBridgeClient) Close() {
	c.client.Close()
}

func (c *EthBridgeClient) AppendTxBatch(contexts []*big.Int, txLengths []*big.Int, txBatch []byte) (*types.Transaction, error) {
	return c.inbox.AppendTxBatch(contexts, txLengths, txBatch)
}

func (c *EthBridgeClient) WatchTxBatchAppended(
	opts *bind.WatchOpts,
	sink chan<- *bindings.ISequencerInboxTxBatchAppended,
) (event.Subscription, error) {
	return c.inbox.Contract.WatchTxBatchAppended(opts, sink)
}

func (c *EthBridgeClient) FilterTxBatchAppendedEvents(
	opts *bind.FilterOpts,
) (*bindings.ISequencerInboxTxBatchAppendedIterator, error) {
	return c.inbox.Contract.FilterTxBatchAppended(opts)
}

func (c *EthBridgeClient) DecodeAppendTxBatchInput(tx *types.Transaction) ([]interface{}, error) {
	return c.inboxAbi.Methods["appendTxBatch"].Inputs.Unpack(tx.Data()[4:])
}

func (c *EthBridgeClient) IsStaked(ctx context.Context) (bool, error) {
	return c.rollup.Contract.IsStaked(&bind.CallOpts{Pending: true, Context: ctx}, c.transactOpts.From)
}

func (c *EthBridgeClient) Stake(amount *big.Int) error {
	log.Info("Staking...")
	stakeOpts := c.rollup.TransactOpts
	stakeOpts.Value = amount
	_, err := c.rollup.Contract.Stake(&stakeOpts)
	if err != nil {
		return fmt.Errorf("Failed to stake, err: %w", err)
	}
	log.Info("Staked %v ETH successfully.", amount)
	return nil
}

func (c *EthBridgeClient) AdvanceStake(assertionID *big.Int) (*types.Transaction, error) {
	return c.rollup.AdvanceStake(assertionID)
}

func (c *EthBridgeClient) CreateAssertion(
	vmHash [32]byte,
	inboxSize *big.Int,
	cumulativeGasUsed *big.Int,
	prevVMHash common.Hash,
	prevL2GasUsed *big.Int,
) (*types.Transaction, error) {
	return c.rollup.CreateAssertion(vmHash, inboxSize, cumulativeGasUsed, prevVMHash, prevL2GasUsed)
}

func (c *EthBridgeClient) ChallengeAssertion(players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error) {
	return c.rollup.ChallengeAssertion(players, assertionIDs)
}

func (c *EthBridgeClient) WatchAssertionCreated(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IRollupAssertionCreated,
) (event.Subscription, error) {
	return c.rollup.Contract.WatchAssertionCreated(opts, sink)
}

func (c *EthBridgeClient) WatchAssertionChallenged(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IRollupAssertionChallenged,
) (event.Subscription, error) {
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

// Returns the last assertion ID that was validated by the staker.
func (c *EthBridgeClient) GetLastValidatedAssertionID(ctx context.Context, start uint64) (*big.Int, error) {
	filterOpts := bind.FilterOpts{Start: start, Context: ctx}
	iter, err := c.rollup.Contract.FilterStakerStaked(&filterOpts, []common.Address{c.transactOpts.From})
	if err != nil {
		return nil, err
	}
	lastValidatedAssertionID := common.Big1
	for iter.Next() {
		// Note: this should always be true if the iterator iterates in time order.
		if lastValidatedAssertionID.Cmp(iter.Event.AssertionID) == 1 {
			lastValidatedAssertionID = iter.Event.AssertionID
		}
	}
	if iter.Error() != nil {
		return nil, fmt.Errorf("Failed to iterate through validated assertion IDs, err: %w", iter.Error())
	}
	return lastValidatedAssertionID, nil
}

func (c *EthBridgeClient) ConfirmFirstUnresolvedAssertion() (*types.Transaction, error) {
	return c.rollup.ConfirmFirstUnresolvedAssertion()
}

func (c *EthBridgeClient) RejectFirstUnresolvedAssertion(stakerAddress common.Address) (*types.Transaction, error) {
	return c.rollup.RejectFirstUnresolvedAssertion(stakerAddress)
}

func (c *EthBridgeClient) GetAssertion(assertionID *big.Int) (*rollupTypes.Assertion, error) {
	assertion, err := c.assertionMap.Assertions(assertionID)
	if err != nil {
		return nil, err
	}
	return &rollupTypes.Assertion{
		ID: assertionID,
		// VmHash: assertion.vmHash(), // TODO change contract
		InboxSize: assertion.InboxSize,
		Deadline:  assertion.Deadline,
	}, nil
}

func (c *EthBridgeClient) GetDeadline(assertionID *big.Int) (*big.Int, error) {
	return c.assertionMap.GetDeadline(assertionID)
}

func (c *EthBridgeClient) InitNewChallengeSession(ctx context.Context, challengeAddress common.Address) error {
	challenge, err := bindings.NewIChallenge(challengeAddress, c.client)
	if err != nil {
		return fmt.Errorf("Failed to initialize challenge contract, err: %w", err)
	}
	c.challenge = &bindings.IChallengeSession{
		Contract:     challenge,
		CallOpts:     bind.CallOpts{Pending: true, Context: ctx},
		TransactOpts: *c.transactOpts,
	}
	return nil
}

func (c *EthBridgeClient) InitializeChallengeLength(numSteps *big.Int) (*types.Transaction, error) {
	return c.challenge.InitializeChallengeLength(numSteps)
}

func (c *EthBridgeClient) CurrentChallengeResponder() (common.Address, error) {
	return c.challenge.CurrentResponder()
}

func (c *EthBridgeClient) CurrentChallengeResponderTimeLeft() (*big.Int, error) {
	return c.challenge.CurrentResponderTimeLeft()
}

func (c *EthBridgeClient) TimeoutChallenge() (*types.Transaction, error) {
	return c.challenge.Timeout()
}

func (c *EthBridgeClient) BisectExecution(
	bisection [][32]byte,
	challengedSegmentIndex *big.Int,
	prevBisection [][32]byte,
	prevChallengedSegmentStart *big.Int,
	prevChallengedSegmentLength *big.Int,
) (*types.Transaction, error) {
	return c.challenge.BisectExecution(
		bisection,
		challengedSegmentIndex,
		prevBisection,
		prevChallengedSegmentStart,
		prevChallengedSegmentLength,
	)
}

func (c *EthBridgeClient) VerifyOneStepProof(
	proof []byte,
	challengedStepIndex *big.Int,
	prevBisection [][32]byte,
	prevChallengedSegmentStart *big.Int,
	prevChallengedSegmentLength *big.Int,
) (*types.Transaction, error) {
	return c.challenge.VerifyOneStepProof(
		proof,
		challengedStepIndex,
		prevBisection,
		prevChallengedSegmentStart,
		prevChallengedSegmentLength,
	)
}

func (c *EthBridgeClient) WatchBisected(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IChallengeBisected,
) (event.Subscription, error) {
	return c.challenge.Contract.WatchBisected(opts, sink)
}

func (c *EthBridgeClient) WatchChallengeCompleted(
	opts *bind.WatchOpts,
	sink chan<- *bindings.IChallengeChallengeCompleted,
) (event.Subscription, error) {
	return c.challenge.Contract.WatchChallengeCompleted(opts, sink)
}

func (c *EthBridgeClient) DecodeBisectExecutionInput(tx *types.Transaction) ([]interface{}, error) {
	return c.challengeAbi.Methods["bisectExecution"].Inputs.Unpack(tx.Data()[4:])
}
