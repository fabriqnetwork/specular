package validator

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type challengeCtx struct {
	opponentAssertion      *rollupTypes.Assertion
	ourAssertion           *rollupTypes.Assertion
	lastValidatedAssertion *rollupTypes.Assertion
}

var errAssertionOverflowedLocalInbox = fmt.Errorf("[Validator] assertion overflowed inbox")
var errValidationFailed = fmt.Errorf("[Validator] validation failed")

type Validator struct {
	*services.BaseService

	newBatchCh           chan struct{}
	challengeCh          chan *challengeCtx
	challengeResoutionCh chan *rollupTypes.Assertion
}

// TODO: this shares a lot of code with sequencer
func New(eth services.Backend, proofBackend proof.Backend, cfg *services.Config, auth *bind.TransactOpts) (*Validator, error) {
	base, err := services.NewBaseService(eth, proofBackend, cfg, auth)
	if err != nil {
		return nil, err
	}
	v := &Validator{
		BaseService:          base,
		newBatchCh:           make(chan struct{}, 4096),
		challengeCh:          make(chan *challengeCtx),
		challengeResoutionCh: make(chan *rollupTypes.Assertion),
	}
	return v, nil
}

// This function will try to validate a pending assertion
func (v *Validator) tryValidateAssertion(lastValidatedAssertion, assertion *rollupTypes.Assertion) error {
	// Find asserted blocks in local blockchain
	inboxSizeDiff := new(big.Int).Sub(assertion.InboxSize, lastValidatedAssertion.InboxSize)
	currentBlockNum := assertion.StartBlock
	currentChainHeight := v.Chain.CurrentBlock().NumberU64()
	var block *types.Block
	targetGasUsed := new(big.Int).Set(lastValidatedAssertion.CumulativeGasUsed)
	for inboxSizeDiff.Cmp(common.Big0) > 0 {
		if currentBlockNum > currentChainHeight {
			return errAssertionOverflowedLocalInbox
		}
		block = v.Chain.GetBlockByNumber(currentBlockNum)
		if block == nil {
			return errAssertionOverflowedLocalInbox
		}
		numTxs := uint64(len(block.Transactions()))
		if numTxs > inboxSizeDiff.Uint64() {
			return fmt.Errorf("[Validator: tryValidateAssertion] UNHANDLED: Assertion created in the middle of block, validator state corrupted!")
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
	_, err := v.Rollup.AdvanceStake(assertion.ID)
	if errors.Is(err, core.ErrInsufficientFunds) {
		return fmt.Errorf("[Validator: tryValidateAssertion] Insufficient Funds to send Tx, err: %w", err)
	}
	if err != nil {
		return fmt.Errorf("[Validator: tryValidateAssertion] UNHANDLED: Can't advance stake, validator state corrupted, err: %w", err)
	}
	return nil
}

// This function runs as a goroutine. It listens for and validates assertions posted to the L1 Rollup contract,
// advances its stake if validated, and challenges if not.
func (v *Validator) validationLoop(lastValidatedAssertion *rollupTypes.Assertion) {
	defer v.Wg.Done()

	// Listen to AssertionCreated event
	assertionEventCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	assertionEventSub, err := v.Rollup.Contract.WatchAssertionCreated(&bind.WatchOpts{Context: v.Ctx}, assertionEventCh)
	if err != nil {
		log.Crit("[Validator: validationLoop] Failed to watch rollup event", "err", err)
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
				log.Warn("[Validator: validationLoop] Assertion overflowed local inbox, wait for next batch event", "expected size", currentAssertion.InboxSize)
				return nil
			default:
				return err
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
			case ourAssertion := <-v.challengeResoutionCh:
				log.Info("challenge finished")
				isInChallenge = false
				lastValidatedAssertion = ourAssertion
				currentAssertion = nil
			case <-v.Ctx.Done():
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
						log.Crit("[Validator: validationLoop] UNHANDLED: Can't validate assertion, validator state corrupted", "err", err)
					}
				}
			case ev := <-assertionEventCh:
				if ev.AsserterAddr == common.Address(v.Config.Coinbase) {
					// Create by our own for challenge
					continue
				}
				// New assertion created on Rollup
				assertion := &rollupTypes.Assertion{
					ID:                    ev.AssertionID,
					VmHash:                ev.VmHash,
					CumulativeGasUsed:     ev.L2GasUsed,
					InboxSize:             ev.InboxSize,
					StartBlock:            lastValidatedAssertion.EndBlock + 1,
					PrevCumulativeGasUsed: new(big.Int).Set(lastValidatedAssertion.CumulativeGasUsed),
				}
				if currentAssertion != nil {
					// TODO: handle concurrent assertions
					log.Crit("[Validator: validationLoop] UNHANDLED: concurrent assertion")
					continue
				}
				currentAssertion = assertion
				err := validateCurrentAssertion()
				if err != nil {
					// TODO: error handling instead of panic
					log.Crit("[Validator: validationLoop] UNHANDLED: Can't validate assertion, validator state corrupted", "err", err)
				}
			case <-v.Ctx.Done():
				return
			}
		}
	}
}

func (v *Validator) challengeLoop() {
	defer v.Wg.Done()

	abi, err := bindings.IChallengeMetaData.GetAbi()
	if err != nil {
		log.Crit("[Validator: challengeLoop] Failed to get IChallenge ABI", "err", err)
	}

	// Watch AssertionCreated event
	createdCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	createdSub, err := v.Rollup.Contract.WatchAssertionCreated(&bind.WatchOpts{Context: v.Ctx}, createdCh)
	if err != nil {
		log.Crit("[Validator: challengeLoop] Failed to watch rollup event", "err", err)
	}
	defer createdSub.Unsubscribe()

	challengedCh := make(chan *bindings.IRollupAssertionChallenged, 4096)
	challengedSub, err := v.Rollup.Contract.WatchAssertionChallenged(&bind.WatchOpts{Context: v.Ctx}, challengedCh)
	if err != nil {
		log.Crit("[Validator: challengeLoop] Failed to watch rollup event", "err", err)
	}
	defer challengedSub.Unsubscribe()

	// Watch L1 blockchain for challenge timeout
	headCh := make(chan *types.Header, 4096)
	headSub, err := v.L1.SubscribeNewHead(v.Ctx, headCh)
	if err != nil {
		log.Crit("[Validator: challengeLoop] Failed to watch l1 chain head", "err", err)
	}
	defer headSub.Unsubscribe()

	var challengeSession *bindings.IChallengeSession
	var states []*proof.ExecutionState

	var bisectedCh chan *bindings.IChallengeBisected
	var bisectedSub event.Subscription
	var challengeCompletedCh chan *bindings.IChallengeChallengeCompleted
	var challengeCompletedSub event.Subscription

	inChallenge := false
	var ctx *challengeCtx
	var opponentTimeoutBlock uint64

	for {
		if inChallenge {
			select {
			case ev := <-bisectedCh:
				// case get bisection, if is our turn
				//   if in single step, submit proof
				//   if multiple step, track current segment, update
				responder, err := challengeSession.CurrentResponder()
				if err != nil {
					// TODO: error handling
					log.Error("[Validator: challengeLoop] Can not get current responder", "error", err)
					continue
				}
				// If it's our turn
				if responder == common.Address(v.Config.Coinbase) {
					err := services.RespondBisection(v.BaseService, abi, challengeSession, ev, states, ctx.opponentAssertion.VmHash, false)
					if err != nil {
						// TODO: error handling
						log.Error("[Validator: challengeLoop] Can not respond to bisection", "error", err)
						continue
					}
				} else {
					opponentTimeLeft, err := challengeSession.CurrentResponderTimeLeft()
					if err != nil {
						// TODO: error handling
						log.Error("[Validator: challengeLoop] Can not get current responder left time", "error", err)
						continue
					}
					log.Info("[Validator: challengeLoop] Opponent time left", "time", opponentTimeLeft)
					opponentTimeoutBlock = ev.Raw.BlockNumber + opponentTimeLeft.Uint64()
				}
			case header := <-headCh:
				if opponentTimeoutBlock == 0 {
					continue
				}
				// TODO: can we use >= here?
				if header.Number.Uint64() > opponentTimeoutBlock {
					_, err := challengeSession.Timeout()
					if err != nil {
						log.Error("[Validator: challengeLoop] Can not timeout opponent", "error", err)
						continue
						// TODO: wait some time before retry
						// TODO: fix race condition
					}
				}
			case ev := <-challengeCompletedCh:
				// TODO: handle if we are not winner --> state corrupted
				log.Info("[Validator: challengeLoop] Challenge completed", "winner", ev.Winner)
				bisectedSub.Unsubscribe()
				challengeCompletedSub.Unsubscribe()
				states = []*proof.ExecutionState{}
				inChallenge = false
				v.challengeResoutionCh <- ctx.ourAssertion
			case <-v.Ctx.Done():
				bisectedSub.Unsubscribe()
				challengeCompletedSub.Unsubscribe()
				return
			}
		} else {
			select {
			case ctx = <-v.challengeCh:
				_, err = v.Rollup.CreateAssertion(
					ctx.ourAssertion.VmHash,
					ctx.ourAssertion.InboxSize,
					ctx.ourAssertion.CumulativeGasUsed,
					ctx.lastValidatedAssertion.VmHash,
					ctx.lastValidatedAssertion.CumulativeGasUsed,
				)
				if errors.Is(err, core.ErrInsufficientFunds) {
					log.Crit("[Validator: challengeLoop] Insufficient Funds to send Tx", "error", err)
				}
				if err != nil {
					log.Crit("[Validator: challengeLoop] UNHANDLED: Can't create assertion for challenge, validator state corrupted", "err", err)
				}
			case ev := <-createdCh:
				if common.Address(ev.AsserterAddr) == v.Config.Coinbase {
					if ev.VmHash == ctx.ourAssertion.VmHash {
						_, err := v.Rollup.ChallengeAssertion(
							[2]common.Address{
								common.Address(v.Config.SequencerAddr),
								common.Address(v.Config.Coinbase),
							},
							[2]*big.Int{
								ctx.opponentAssertion.ID,
								ev.AssertionID,
							},
						)
						if errors.Is(err, core.ErrInsufficientFunds) {
							log.Crit("[Validator: challengeLoop] Insufficient Funds to send Tx", "error", err)
						}
						if err != nil {
							log.Crit("[Validator: challengeLoop] UNHANDLED: Can't start challenge, validator state corrupted", "err", err)
						}
					}
				}
			case ev := <-challengedCh:
				if ctx == nil {
					continue
				}
				log.Info("[Validator: challengeLoop] validator saw challenge", "assertion id", ev.AssertionID, "expected id", ctx.opponentAssertion.ID, "block", ev.Raw.BlockNumber)
				if ev.AssertionID.Cmp(ctx.opponentAssertion.ID) == 0 {
					// start := ev.Raw.BlockNumber - 2
					challenge, err := bindings.NewIChallenge(ev.ChallengeAddr, v.L1)
					if err != nil {
						log.Crit("[Validator: challengeLoop] Failed to access ongoing challenge", "address", ev.ChallengeAddr, "err", err)
					}
					challengeSession = &bindings.IChallengeSession{
						Contract:     challenge,
						CallOpts:     bind.CallOpts{Pending: true, Context: v.Ctx},
						TransactOpts: *v.TransactOpts,
					}
					bisectedCh = make(chan *bindings.IChallengeBisected, 4096)
					bisectedSub, err = challenge.WatchBisected(&bind.WatchOpts{Context: v.Ctx}, bisectedCh)
					if err != nil {
						log.Crit("[Validator: challengeLoop] Failed to watch challenge event", "err", err)
					}
					challengeCompletedCh = make(chan *bindings.IChallengeChallengeCompleted, 4096)
					challengeCompletedSub, err = challenge.WatchChallengeCompleted(&bind.WatchOpts{Context: v.Ctx}, challengeCompletedCh)
					if err != nil {
						log.Crit("[Validator: challengeLoop] Failed to watch challenge event", "err", err)
					}
					states, err = proof.GenerateStates(
						v.ProofBackend,
						v.Ctx,
						ctx.opponentAssertion.PrevCumulativeGasUsed,
						ctx.opponentAssertion.StartBlock,
						ctx.opponentAssertion.EndBlock+1,
						nil,
					)
					if err != nil {
						log.Crit("[Validator: challengeLoop] Failed to generate states", "err", err)
					}
					inChallenge = true
				}
			case <-headCh:
				continue // consume channel values
			case <-v.Ctx.Done():
				return
			}
		}
	}
}

func (v *Validator) Start() error {
	log.Info("Starting validator...")
	err := v.BaseService.Start(true, true)
	if err != nil {
		return err
	}
	lastValidatedAssertion, err := v.getLastValidatedAssertion()
	if err != nil {
		return err
	}
	v.Wg.Add(3)
	go v.SyncLoop(v.newBatchCh)
	go v.validationLoop(lastValidatedAssertion)
	go v.challengeLoop()
	log.Info("[Validator] Validator started.")
	return nil
}

func (v *Validator) Stop() error {
	log.Info("[Validator] Stopping validator...")
	v.Cancel()
	v.Wg.Wait()
	log.Info("Validator stopped.")
	return nil
}

func (v *Validator) APIs() []rpc.API {
	// TODO: validator APIs
	return []rpc.API{}
}

// Gets the last validated assertion.
func (v *Validator) getLastValidatedAssertion() (*rollupTypes.Assertion, error) {
	// TODO: set FilterOpts.Start
	iter, err := v.Rollup.Contract.FilterStakerStaked(&bind.FilterOpts{Context: v.Ctx}, []common.Address{v.TransactOpts.From})
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
		return nil, iter.Error()
	}
	return v.getAssertion(lastValidatedAssertionID)
}

func (v *Validator) getAssertion(assertionID *big.Int) (*rollupTypes.Assertion, error) {
	if assertionID.Cmp(common.Big1) == 0 {
		return &rollupTypes.Assertion{
			ID:                    assertionID,
			VmHash:                v.Eth.BlockChain().CurrentBlock().Root(),
			CumulativeGasUsed:     new(big.Int),
			InboxSize:             new(big.Int),
			Deadline:              new(big.Int),
			PrevCumulativeGasUsed: new(big.Int),
		}, nil
	} else {
		assertion, err := v.AssertionMap.Assertions(assertionID)
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
}
