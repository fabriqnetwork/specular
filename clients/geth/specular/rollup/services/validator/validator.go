package validator

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

type challengeCtx struct {
	opponentAssertion      *rollupTypes.Assertion
	ourAssertion           *rollupTypes.Assertion
	lastValidatedAssertion *rollupTypes.Assertion
}

var errAssertionOverflowedLocalInbox = fmt.Errorf("assertion overflowed inbox")
var errValidationFailed = fmt.Errorf("validation failed")

type Validator struct {
	*services.BaseService

	newBatchCh            chan struct{}
	challengeCh           chan *challengeCtx
	challengeResolutionCh chan *rollupTypes.Assertion
}

// TODO: this shares a lot of code with sequencer
func New(eth services.Backend, proofBackend proof.Backend, l1Client client.L1BridgeClient, cfg *services.Config) (*Validator, error) {
	base, err := services.NewBaseService(eth, proofBackend, l1Client, cfg)
	if err != nil {
		return nil, err
	}
	v := &Validator{
		BaseService:           base,
		newBatchCh:            make(chan struct{}, 4096),
		challengeCh:           make(chan *challengeCtx),
		challengeResolutionCh: make(chan *rollupTypes.Assertion),
	}
	return v, nil
}

// This function will try to validate a pending assertion
func (v *Validator) tryValidateAssertion(lastValidatedAssertion, assertion *rollupTypes.Assertion) error {
	// Find asserted blocks in local blockchain
	inboxSizeDiff := new(big.Int).Sub(assertion.InboxSize, lastValidatedAssertion.InboxSize)
	currentBlockNum := assertion.StartBlock
	currentChainHeight := v.Chain().CurrentBlock().NumberU64()
	var block *types.Block
	targetGasUsed := new(big.Int).Set(lastValidatedAssertion.CumulativeGasUsed)
	for inboxSizeDiff.Cmp(common.Big0) > 0 {
		if currentBlockNum > currentChainHeight {
			return errAssertionOverflowedLocalInbox
		}
		block = v.Chain().GetBlockByNumber(currentBlockNum)
		if block == nil {
			return errAssertionOverflowedLocalInbox
		}
		numTxs := uint64(len(block.Transactions()))
		if numTxs > inboxSizeDiff.Uint64() {
			return fmt.Errorf("UNHANDLED: Assertion created in the middle of block, validator state corrupted!")
		}
		targetGasUsed.Add(targetGasUsed, new(big.Int).SetUint64(block.GasUsed()))
		inboxSizeDiff = new(big.Int).Sub(inboxSizeDiff, new(big.Int).SetUint64(numTxs))
		currentBlockNum++
	}
	assertion.EndBlock = currentBlockNum - 1
	targetVmHash := block.Root()
	if targetVmHash != assertion.VmHash || targetGasUsed.Cmp(assertion.CumulativeGasUsed) != 0 {
		// Validation failed
		ourAssertion := &rollupTypes.Assertion{
			VmHash:                targetVmHash,
			CumulativeGasUsed:     targetGasUsed,
			InboxSize:             assertion.InboxSize,
			StartBlock:            assertion.StartBlock,
			EndBlock:              assertion.EndBlock,
			PrevCumulativeGasUsed: new(big.Int).Set(lastValidatedAssertion.CumulativeGasUsed),
		}
		v.challengeCh <- &challengeCtx{assertion, ourAssertion, lastValidatedAssertion}
		return errValidationFailed
	}
	// Validation succeeded, confirm assertion and advance stake
	// if assertion.ID
	_, err := v.L1Client.AdvanceStake(assertion.ID)
	if errors.Is(err, core.ErrInsufficientFunds) {
		return fmt.Errorf("Insufficient Funds to send Tx, err: %w", err)
	}
	if err != nil {
		return fmt.Errorf("UNHANDLED: Can't advance stake, validator state corrupted, err: %w", err)
	}
	return nil
}

// This function runs as a goroutine. It listens for and validates assertions posted to the L1 Rollup contract,
// advances its stake if validated, and challenges if not.
func (v *Validator) validationLoop(ctx context.Context) {
	defer v.Wg.Done()

	lastValidatedAssertion, err := v.GetLastValidatedAssertion(ctx)
	if err != nil {
		log.Crit("Failed to start validating, couldn't get last validated assertion", "err", err)
	}
	// Listen to AssertionCreated event
	assertionEventCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	assertionEventSub, err := v.L1Client.WatchAssertionCreated(&bind.WatchOpts{Context: ctx}, assertionEventCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer assertionEventSub.Unsubscribe()

	// Current agreed assertion, initalize to genesis assertion
	// The next assertion to be validated
	var currentAssertion *rollupTypes.Assertion

	isInChallenge := false

	validateCurrentAssertion := func() error {
		// Validate current assertion
		err := v.tryValidateAssertion(lastValidatedAssertion, currentAssertion)
		if err != nil {
			switch {
			case errors.Is(err, errValidationFailed):
				// Validation failed, challenge
				isInChallenge = true
				return nil
			case errors.Is(err, errAssertionOverflowedLocalInbox):
				// Assertion overflowed local inbox, wait for next batch event
				log.Warn("Assertion overflowed local inbox, wait for next batch event", "expected size", currentAssertion.InboxSize)
				return nil
			default:
				return fmt.Errorf("tryValidateAssertion Failed",err)
			}
		}
		// Validation success, clean up
		lastValidatedAssertion = currentAssertion
		currentAssertion = nil
		return nil
	}

	for {
		if isInChallenge {
			// Wait for the challenge resolution
			select {
			case ourAssertion := <-v.challengeResolutionCh:
				log.Info("challenge finished")
				isInChallenge = false
				lastValidatedAssertion = ourAssertion
				currentAssertion = nil
			case <-ctx.Done():
				return
			}
		} else {
			select {
			case <-v.newBatchCh:
				// New block committed, try to validate all pending assertion
				if currentAssertion != nil {
					err := validateCurrentAssertion()
					if err != nil {
						// TODO: error handling instead of panic
						log.Crit("UNHANDLED: Can't validate assertion, validator state corrupted", "err", err)
					}
				}
			case ev := <-assertionEventCh:
				if ev.AsserterAddr == common.Address(v.Config.Coinbase) {
					// Create by our own for challenge
					continue
				}
				// New assertion created on Rollup
				assertionFromRollup, err := v.L1Client.GetAssertion(ev.AssertionID)
				if err != nil {
					log.Crit("Could not get DA", "error", err)
				}
				assertion := &rollupTypes.Assertion{
					ID:                    ev.AssertionID,
					VmHash:                ev.VmHash,
					CumulativeGasUsed:     ev.L2GasUsed,
					InboxSize:             assertionFromRollup.InboxSize,
					StartBlock:            lastValidatedAssertion.EndBlock + 1,
					PrevCumulativeGasUsed: new(big.Int).Set(lastValidatedAssertion.CumulativeGasUsed),
				}
				if currentAssertion != nil {
					// TODO: handle concurrent assertions
					log.Crit("UNHANDLED: concurrent assertion")
					continue
				}
				currentAssertion = assertion
				err = validateCurrentAssertion()
				if err != nil {
					// TODO: error handling instead of panic
					log.Crit("UNHANDLED: Can't validate assertion, validator state corrupted", "err", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (v *Validator) challengeLoop(ctx context.Context) {
	defer v.Wg.Done()

	// Watch AssertionCreated event
	createdCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	createdSub, err := v.L1Client.WatchAssertionCreated(&bind.WatchOpts{Context: ctx}, createdCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer createdSub.Unsubscribe()

	challengedCh := make(chan *bindings.IRollupAssertionChallenged, 4096)
	challengedSub, err := v.L1Client.WatchAssertionChallenged(&bind.WatchOpts{Context: ctx}, challengedCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer challengedSub.Unsubscribe()

	// Watch L1 blockchain for challenge timeout
	headCh := make(chan *types.Header, 4096)
	headSub, err := v.L1Client.ResubscribeErrNewHead(ctx, headCh)
	if err != nil {
		log.Crit("Failed to watch l1 chain head", "err", err)
	}
	defer headSub.Unsubscribe()

	var states []*proof.ExecutionState

	var bisectedCh chan *bindings.IChallengeBisected
	var bisectedSub event.Subscription
	var challengeCompletedCh chan *bindings.IChallengeChallengeCompleted
	var challengeCompletedSub event.Subscription

	inChallenge := false
	var chalCtx *challengeCtx
	var opponentTimeoutBlock uint64

	for {
		if inChallenge {
			select {
			case ev := <-bisectedCh:
				// case get bisection, if is our turn
				//   if in single step, submit proof
				//   if multiple step, track current segment, update
				responder, err := v.L1Client.CurrentChallengeResponder()
				if err != nil {
					// TODO: error handling
					log.Error("Can not get current responder", "error", err)
					continue
				}
				// If it's our turn
				if responder == common.Address(v.Config.Coinbase) {
					err := services.RespondBisection(ctx, v.ProofBackend, v.L1Client, ev, states, chalCtx.opponentAssertion.VmHash, false)
					if err != nil {
						// TODO: error handling
						log.Error("Can not respond to bisection", "error", err)
						continue
					}
				} else {
					opponentTimeLeft, err := v.L1Client.CurrentChallengeResponderTimeLeft()
					if err != nil {
						// TODO: error handling
						log.Error("Can not get current responder left time", "error", err)
						continue
					}
					log.Info("Opponent time left", "time", opponentTimeLeft)
					opponentTimeoutBlock = ev.Raw.BlockNumber + opponentTimeLeft.Uint64()
				}
			case header := <-headCh:
				if opponentTimeoutBlock == 0 {
					continue
				}
				// TODO: can we use >= here?
				if header.Number.Uint64() > opponentTimeoutBlock {
					_, err := v.L1Client.TimeoutChallenge()
					if err != nil {
						log.Error("Can not timeout opponent", "error", err)
						continue
						// TODO: wait some time before retry
						// TODO: fix race condition
					}
				}
			case ev := <-challengeCompletedCh:
				// TODO: handle if we are not winner --> state corrupted
				log.Info("Challenge completed", "winner", ev.Winner)
				bisectedSub.Unsubscribe()
				challengeCompletedSub.Unsubscribe()
				states = []*proof.ExecutionState{}
				inChallenge = false
				v.challengeResolutionCh <- chalCtx.ourAssertion
			case <-ctx.Done():
				bisectedSub.Unsubscribe()
				challengeCompletedSub.Unsubscribe()
				return
			}
		} else {
			select {
			case chalCtx = <-v.challengeCh:
				_, err = v.L1Client.CreateAssertion(
					chalCtx.ourAssertion.VmHash,
					chalCtx.ourAssertion.InboxSize,
					chalCtx.ourAssertion.CumulativeGasUsed,
					chalCtx.lastValidatedAssertion.VmHash,
					chalCtx.lastValidatedAssertion.CumulativeGasUsed,
				)
				if errors.Is(err, core.ErrInsufficientFunds) {
					log.Crit("Insufficient Funds to send Tx", "error", err)
				}
				if err != nil {
					log.Crit("UNHANDLED: Can't create assertion for challenge, validator state corrupted", "err", err)
				}
			case ev := <-createdCh:
				if common.Address(ev.AsserterAddr) == v.Config.Coinbase {
					if ev.VmHash == chalCtx.ourAssertion.VmHash {
						_, err := v.L1Client.ChallengeAssertion(
							[2]common.Address{
								common.Address(v.Config.SequencerAddr),
								common.Address(v.Config.Coinbase),
							},
							[2]*big.Int{
								chalCtx.opponentAssertion.ID,
								ev.AssertionID,
							},
						)
						if errors.Is(err, core.ErrInsufficientFunds) {
							log.Crit("Insufficient Funds to send Tx", "error", err)
						}
						if err != nil {
							log.Crit("UNHANDLED: Can't start challenge, validator state corrupted", "err", err)
						}
					}
				}
			case ev := <-challengedCh:
				if chalCtx == nil {
					continue
				}
				log.Info("validator saw challenge", "assertion id", ev.AssertionID, "expected id", chalCtx.opponentAssertion.ID, "block", ev.Raw.BlockNumber)
				if ev.AssertionID.Cmp(chalCtx.opponentAssertion.ID) == 0 {
					// start := ev.Raw.BlockNumber - 2
					err := v.L1Client.InitNewChallengeSession(context.Background(), ev.ChallengeAddr)
					if err != nil {
						log.Crit("Failed to access ongoing challenge", "address", ev.ChallengeAddr, "err", err)
					}
					bisectedCh = make(chan *bindings.IChallengeBisected, 4096)
					bisectedSub, err = v.L1Client.WatchBisected(&bind.WatchOpts{Context: ctx}, bisectedCh)
					if err != nil {
						log.Crit("Failed to watch challenge event", "err", err)
					}
					challengeCompletedCh = make(chan *bindings.IChallengeChallengeCompleted, 4096)
					challengeCompletedSub, err = v.L1Client.WatchChallengeCompleted(&bind.WatchOpts{Context: ctx}, challengeCompletedCh)
					if err != nil {
						log.Crit("Failed to watch challenge event", "err", err)
					}
					states, err = proof.GenerateStates(
						v.ProofBackend,
						ctx,
						chalCtx.opponentAssertion.PrevCumulativeGasUsed,
						chalCtx.opponentAssertion.StartBlock,
						chalCtx.opponentAssertion.EndBlock+1,
						nil,
					)
					if err != nil {
						log.Crit("Failed to generate states", "err", err)
					}
					inChallenge = true
				}
			case <-headCh:
				continue // consume channel values
			case <-ctx.Done():
				return
			}
		}
	}
}

func (v *Validator) Start() error {
	log.Info("Starting validator...")
	ctx, err := v.BaseService.Start()
	if err != nil {
		return err
	}
	if err := v.Stake(ctx); err != nil {
		return fmt.Errorf("Failed to stake, err: %w", err)
	}
	// This is necessary despite `SyncLoop`, since we need an up-to-date L2 chain
	// before we resolve all information about the last validated assertion.
	end, err := v.SyncL2ChainToL1Head(ctx, v.Config.L1RollupGenesisBlock)
	if err != nil {
		return fmt.Errorf("Failed to start sequencer: %w", err)
	}
	v.Wg.Add(3)
	go v.SyncLoop(ctx, end+1, v.newBatchCh)
	go v.validationLoop(ctx)
	go v.challengeLoop(ctx)
	log.Info("Validator started.")
	return nil
}

func (v *Validator) APIs() []rpc.API {
	// TODO: validator APIs
	return []rpc.API{}
}
