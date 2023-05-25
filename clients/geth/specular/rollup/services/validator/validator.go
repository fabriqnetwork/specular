package validator

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types/assertion"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

const interval = 10 * time.Second

type errAssertionOverflowedLocalInbox struct {
	Msg string
}

func (e *errAssertionOverflowedLocalInbox) Error() string {
	return fmt.Sprint("assertion overflowed local inbox with msg:", e.Msg)
}

type Validator struct {
	*services.BaseService

	cfg               ValidatorServiceConfig
	l2ClientCreatorFn l2ClientCreatorFn
	l2Client          L2Client
	l1TxMgr           TxManager
	rollupState       RollupState
	proofBackend      proof.Backend
}

type l2ClientCreatorFn func(ctx context.Context) (L2Client, error)

func NewValidator(
	cfg ValidatorServiceConfig,
	l2ClientCreatorFn l2ClientCreatorFn,
	l1TxMgr TxManager,
	proofBackend proof.Backend,
	rollupState RollupState,
) *Validator {
	return &Validator{
		BaseService:  &services.BaseService{},
		cfg:          cfg,
		l2Client:     nil, // Initialized in `Start()`
		l1TxMgr:      l1TxMgr,
		proofBackend: proofBackend,
		rollupState:  rollupState,
	}
}

func (v *Validator) Start() error {
	log.Info("Starting validator...")
	ctx := v.BaseService.Start()
	// Connect to L2 client.
	l2Client, err := v.l2ClientCreatorFn(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create L2 client: %w", err)
	}
	v.l2Client = l2Client
	if v.cfg.Validator().IsActiveStaker() {
		if err := v.Stake(ctx); err != nil {
			return fmt.Errorf("Failed to stake: %w", err)
		}
	}
	// end, err := v.SyncL2ChainToL1Head(ctx, v.Config.L1RollupGenesisBlock)
	// if err != nil {
	// 	return fmt.Errorf("Failed to sync L2 chain to head: %w", err)
	// }
	// TODO: handle synchronization between two parties modifying blockchain.
	// go v.SyncLoop(ctx, end+1, nil)
	if v.cfg.Validator().IsActiveStaker() {
		v.Eg.Go(func() error { return v.eventLoop(ctx) })
	}
	if v.cfg.Validator().IsActiveChallenger() {
		v.Eg.Go(func() error { return v.challengeLoop(ctx) })
	}
	return nil
}

func (v *Validator) APIs() []rpc.API {
	return []rpc.API{}
}

func (v *Validator) eventLoop(ctx context.Context) error {
	// createdCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionCreated](
	// 	ctx,
	// 	v.l1Syncer.LatestHeaderBroker,
	// 	v.l1Client.FilterAssertionCreated,
	// 	v.l1State.Head().Number.Uint64(),
	// )
	var createdCh = make(chan *bindings.IRollupAssertionCreated)

	// TODO: configure.
	var ticker = time.NewTicker(interval)
	defer ticker.Stop()
	// TODO: case for handling detected reorgs.
	for {
		select {
		case <-ticker.C:
			// Note: assertions are created/resolved relatively infrequently,
			// compared to batch sequencing. However, the queue may grow larger
			// particularly during disputes.
			if v.cfg.Validator().IsResolver() {
				err := v.tryResolve(ctx)
				if err != nil {
					return fmt.Errorf("Failed to resolve assertions: %w", err)
				}
			}
			if v.cfg.Validator().IsActiveCreator() {
				_, err := v.createAssertion(ctx)
				if errors.Is(err, core.ErrInsufficientFunds) {
					return fmt.Errorf("Insufficient funds to send tx: %w", err)
				} else if err != nil {
					log.Error("Failed to create DA", "error", err)
					continue // Try again...
				}
			}
		case ev := <-createdCh:
			log.Info("Received `AssertionCreated` event.", "assertion id", ev.AssertionID)
			// Validate. This blocks assertion creation, but that's fine.
			if common.Address(ev.AsserterAddr) == v.cfg.Validator().AccountAddr() {
				// No need to validate, advance stake or resolve (already done above).
				return nil
			}
			err := v.validate()
			if err != nil {
				log.Error("Failed to validate assertion", "error", err)
				if v.cfg.Validator().IsActiveChallenger() {
					// If incorrect, challenge (fork assertion chain if necessary).
					// Note: we can continue to resolve prior assertions concurrently.
				}
			} else {
				// If correct, advance stake.
				_, err := v.l1TxMgr.AdvanceStake(ctx, ev.AssertionID)
				if err != nil {
					return fmt.Errorf("Failed to advance stake: %w", err)
				}
			}
		case <-ctx.Done():
			log.Warn("Aborting")
			return nil
		}
	}
}

func (v *Validator) challengeLoop(ctx context.Context) error {
	// challengeCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionChallenged](
	// 	ctx,
	// 	v.L1Syncer.LatestHeaderBroker,
	// 	v.l1Client.FilterAssertionCreated,
	// 	v.L1State.Head().Number.Uint64(),
	// )
	challengeCh := make(chan *bindings.IRollupAssertionChallenged)

	for {
		select {
		case ev := <-challengeCh:
			log.Info("Received `AssertionChallenged` event.", "assertion id", ev.AssertionID)
			continue
		case <-ctx.Done():
			log.Warn("Aborting")
			return nil
		}
	}
}

func (v *Validator) handleChallenge(ev *bindings.IRollupAssertionChallenged) {

}

func (v *Validator) validate() error {
	// TODO: refactor `tryValidateAssertion`.
	return nil
}

// Tries to resolve as many assertions as possible, starting from last resolved.
// TODO: timeout early if too many to resolve.
func (v *Validator) tryResolve(ctx context.Context) error {
	tail := big.NewInt(1)
	for id := new(big.Int); id.Cmp(tail) < 0; id.Add(id, common.Big1) {
		_, err := v.rollupState.GetAssertion(ctx, id)
		if err != nil {
			return err
		}
		// if v.SysState.L1State.Number.Uint64() < assertion.Deadline.Uint64() {
		// 	break
		// }
		// TODO: confirm OR reject.
		// Or possibly issue a challenge if necessary---although any
		// challenges should have been issued during validation.
		_, err = v.l1TxMgr.ConfirmFirstUnresolvedAssertion(ctx)
		if err != nil {
			return fmt.Errorf("Failed to confirm, err: %w", err)
		}
	}
	return nil
}

func (v *Validator) createAssertion(ctx context.Context) (*assertion.Assertion, error) {
	vmHash, inboxSize := v.getNextAssertion(ctx)
	// TODO: assertion mgr
	_, err := v.l1TxMgr.CreateAssertion(ctx, vmHash, inboxSize)
	if err != nil {
		return nil, err
	}
	staker, err := v.rollupState.GetStaker(ctx, v.cfg.Validator().AccountAddr())
	if err != nil {
		return nil, fmt.Errorf("Failed to get assertion ID (through staker), err: %w", err)
	}
	log.Info("Created assertion", "ID", staker.AssertionID, "vmHash", vmHash, "inboxSize", inboxSize)
	// Create assertion on L1 Rollup
	// pendingAssertion = queuedAssertion.Copy()
	// queuedAssertion.StartBlock = a.assertionMgr.GetAssertionAux(ctx, id).l2BlockEnd // queuedAssertion.EndBlock + 1
	// "start block", pendingAssertion.StartBlock,
	// "end block", pendingAssertion.EndBlock,
	return nil, nil
}

// TODO: Safe or finalized block.
func (v *Validator) getNextAssertion(ctx context.Context) (common.Hash, *big.Int) {
	// createCh <- struct{}{}
	// Update queued assertion to latest batch
	// vmHash := batch.LastBlockRoot()
	// inboxSize.Add(queuedAssertion.InboxSize, batch.Size())
	// queuedAssertion.EndBlock = batch.LastBlockNumber()
	return common.Hash{}, nil
}

// Gets the last validated assertion.
// func (v *Validator) GetLastValidatedAssertion(ctx context.Context) (*assertion.Assertion, error) {
// 	opts := bind.FilterOpts{Start: v.cfg.L1().RollupGenesisBlock(), Context: ctx}
// 	assertionID, err := v.rollupState.GetLastValidatedAssertionID(&opts)

// 	var assertionCreatedEvent *bindings.IRollupAssertionCreated
// 	var lastValidatedAssertion bindings.IRollupAssertion
// 	if err != nil {
// 		// If no assertion was validated (or other errors encountered), try to use the genesis assertion.
// 		log.Warn("No validated assertions found, using genesis assertion", "err", err)
// 		assertionCreatedEvent, err = v.rollupState.GetGenesisAssertionCreated(&opts)
// 		if err != nil {
// 			return nil, fmt.Errorf("Failed to get `AssertionCreated` event for last validated assertion, err: %w", err)
// 		}
// 		// Check that the genesis assertion is correct.
// 		vmHash := common.BytesToHash(assertionCreatedEvent.VmHash[:])
// 		genesisBlock, err := v.l2Client.BlockByNumber(ctx, common.Big0)
// 		if err != nil {
// 			return nil, fmt.Errorf("Failed to get genesis root, err: %w", err)
// 		}
// 		genesisRoot := genesisBlock.Root()
// 		if vmHash != genesisRoot {
// 			return nil, fmt.Errorf("Mismatching genesis %s vs %s", vmHash, genesisRoot.String())
// 		}
// 		log.Info("Genesis assertion found", "assertionID", assertionCreatedEvent.AssertionID)
// 		// Get assertion.
// 		lastValidatedAssertion, err = v.rollupState.GetAssertion(ctx, assertionCreatedEvent.AssertionID)
// 	} else {
// 		// If an assertion was validated, use it.
// 		log.Info("Last validated assertion ID found", "assertionID", assertionID)
// 		lastValidatedAssertion, err = v.rollupState.GetAssertion(ctx, assertionID)
// 		if err != nil {
// 			return nil, fmt.Errorf("Failed to get last validated assertion, err: %w", err)
// 		}
// 		opts = bind.FilterOpts{Start: lastValidatedAssertion.ProposalTime.Uint64(), Context: ctx}
// 		assertionCreatedIter, err := v.rollupState.FilterAssertionCreated(&opts)
// 		if err != nil {
// 			return nil, fmt.Errorf("Failed to get `AssertionCreated` event for last validated assertion, err: %w", err)
// 		}
// 		assertionCreatedEvent, err = filterAssertionCreatedWithID(assertionCreatedIter, assertionID)
// 	}
// 	// Initialize assertion.
// 	assertion := NewAssertionFrom(&lastValidatedAssertion, assertionCreatedEvent)
// 	// Set its boundaries using parent. TODO: move this out. Use local caching.
// 	opts = bind.FilterOpts{Start: v.cfg.L1().RollupGenesisBlock(), Context: ctx}
// 	parentAssertionCreatedIter, err := v.rollupState.FilterAssertionCreated(&opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to get `AssertionCreated` event for parent assertion, err: %w", err)
// 	}
// 	parentAssertionCreatedEvent, err := filterAssertionCreatedWithID(parentAssertionCreatedIter, lastValidatedAssertion.Parent)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to get `AssertionCreated` event for parent assertion, err: %w", err)
// 	}
// 	err = v.setL2BlockBoundaries(ctx, assertion, parentAssertionCreatedEvent)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to set L2 block boundaries for last validated assertion, err: %w", err)
// 	}
// 	return assertion, nil
// }

func filterAssertionCreatedWithID(iter *bindings.IRollupAssertionCreatedIterator, assertionID *big.Int) (*bindings.IRollupAssertionCreated, error) {
	var assertionCreated *bindings.IRollupAssertionCreated
	for iter.Next() {
		// Assumes invariant: only one `AssertionCreated` event per assertion ID.
		if iter.Event.AssertionID.Cmp(assertionID) == 0 {
			assertionCreated = iter.Event
			break
		}
	}
	if iter.Error() != nil {
		return nil, fmt.Errorf("Failed to iterate through `AssertionCreated` events, err: %w", iter.Error())
	}
	if assertionCreated == nil {
		return nil, fmt.Errorf("No `AssertionCreated` event found for %v.", assertionID)
	}
	return assertionCreated, nil
}

// TODO: clean up.
func (v *Validator) setL2BlockBoundaries(
	ctx context.Context,
	assertion *assertion.Assertion,
	parentAssertionCreatedEvent *bindings.IRollupAssertionCreated,
) error {
	block, err := v.l2Client.BlockByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to get current L2 block, err: %w", err)
	}
	numBlocks := block.Number().Uint64()
	if numBlocks == 0 {
		log.Info("Zero-initializing assertion block boundaries.")
		assertion.StartBlock = 0
		assertion.EndBlock = 0
		return nil
	}
	startFound := false
	// Note: by convention defined in Rollup.sol, the parent VmHash is the
	// same as the child only when the assertion is the genesis assertion.
	// This is a hack to avoid mis-setting `assertion.StartBlock`.
	if assertion.ID == parentAssertionCreatedEvent.AssertionID {
		parentAssertionCreatedEvent.VmHash = common.Hash{}
		startFound = true
	}
	log.Info("Searching for start and end blocks for assertion.", "id", assertion.ID)
	// Find start and end blocks using L2 chain (assumes it's synced at least up to the assertion).
	for i := uint64(0); i <= numBlocks; i++ {
		// TODO: remove assumption of VM hash being the block root.
		block, err := v.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(i))
		if err != nil {
			return fmt.Errorf("Failed to get L2 block, err: %w", err)
		}
		root := block.Root()
		if root == parentAssertionCreatedEvent.VmHash {
			log.Info("Found start block", "l2 block#", i+1)
			assertion.StartBlock = i + 1
			startFound = true
		} else if root == assertion.VmHash {
			log.Info("Found end block", "l2 block#", i)
			assertion.EndBlock = i
			if !startFound {
				return fmt.Errorf("Found end block before start block for assertion with hash %d", assertion.VmHash)
			}
			return nil
		}
	}
	return fmt.Errorf("Could not find start or end block for assertion with hash %s", assertion.VmHash)
}

func (v *Validator) Stake(ctx context.Context) error {
	staker, err := v.rollupState.GetStaker(ctx, v.cfg.Validator().AccountAddr())
	if err != nil {
		return fmt.Errorf("Failed to get staker, to stake, err: %w", err)
	}
	if !staker.IsStaked {
		_, err = v.l1TxMgr.Stake(ctx, big.NewInt(int64(v.cfg.Validator().StakeAmount())))
		if err != nil {
			return fmt.Errorf("Failed to stake, err: %w", err)
		}
	}
	return nil
}

// This goroutine tries to confirm created assertions
// func (a *Asserter) resolutionLoop(ctx context.Context) {
// 	defer a.wg.Done()

// 	headCh := a.L1Syncer.LatestHeaderBroker.Subscribe()
// 	confirmedCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionConfirmed](
// 		ctx, a.L1Syncer.LatestHeaderBroker, a.L1Client.FilterAssertionConfirmed, a.L1State.Latest().Number.Uint64(),
// 	)
// 	challengedCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionChallenged](
// 		ctx, a.L1Syncer.LatestHeaderBroker, a.L1Client.FilterAssertionChallenged, a.L1State.Latest().Number.Uint64(),
// 	)

// 	// Current pending assertion from sequencing goroutine
// 	// TODO: watch multiple pending assertions
// 	var pendingAssertion *assertion.Assertion
// 	pendingConfirmationSent := true
// 	pendingConfirmed := true

// 	for {
// 		select {
// 		case header := <-headCh:
// 			// New block mined on L1
// 			log.Info("Received new header", "number", header.Number.Uint64())
// 			if !pendingConfirmationSent && !pendingConfirmed {
// 				if header.Number.Uint64() >= pendingAssertion.Deadline.Uint64() {
// 					log.Info("We can now confirm", "pending assertion", pendingAssertion.Deadline.Uint64())
// 					// Confirmation period has past, confirm it
// 					_, err := s.L1Client.ConfirmFirstUnresolvedAssertion()
// 					if errors.Is(err, core.ErrInsufficientFunds) {
// 						log.Crit("Insufficient Funds to send Tx", "error", err)
// 					}
// 					if err != nil {
// 						// log.Error("Failed to confirm DA", "error", err)
// 						log.Crit("Failed to confirm DA", "err", err)
// 						// TODO: wait some time before retry
// 					}
// 					pendingConfirmationSent = true
// 				}
// 			}
// 		case ev := <-confirmedCh:
// 			log.Info("Received `AssertionConfirmed` event ", "assertion id", ev.AssertionID)
// 			// New confirmed assertion
// 			if ev.AssertionID.Cmp(pendingAssertion.ID) == 0 {
// 				// Notify sequencing goroutine
// 				s.confirmedIDCh <- pendingAssertion.ID
// 				pendingConfirmed = true
// 			}
// 		case newPendingAssertion := <-s.resolveAssertionCh:
// 			log.Info("Received pending assertion")
// 			// New assertion created by sequencing goroutine
// 			if !pendingConfirmed {
// 				// TODO: support multiple pending assertion
// 				log.Error("Got another DA request before current is confirmed")
// 				continue
// 			}
// 			pendingAssertion = newPendingAssertion.Copy()
// 			pendingConfirmationSent = false
// 			pendingConfirmed = false
// 		case ev := <-challengedCh:
// 			// New challenge raised
// 			log.Info("Received `AssertionChallenged` event ", "assertion id", ev.AssertionID)
// 			// if ev.AssertionID.Cmp(pendingAssertion.ID) == 0 {
// 			// 	a.challengeCh <- &challengeCtx{
// 			// 		ev.ChallengeAddr,
// 			// 		pendingAssertion,
// 			// 	}
// 			// 	wait(ctx, s.challengeResoutionCh, "challenge resolution")
// 			// }
// 		case <-ctx.Done():
// 			log.Info("Aborting.")
// 			return
// 		}
// 	}
// }
