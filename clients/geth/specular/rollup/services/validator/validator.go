package validator

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

var errAssertionOverflowedLocalInbox = fmt.Errorf("assertion overflowed inbox")
var errValidationFailed = fmt.Errorf("validation failed")

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
			log.Crit("UNHANDELED: Assertion created in the middle of block, validator state corrupted!")
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
	_, err := v.Rollup.AdvanceStake(assertion.ID)
	if err != nil {
		log.Crit("UNHANDELED: Can't advance stake, validator state corrupted", "err", err)
	}
	return nil
}

// This goroutine validates the assertion posted to L1 Rollup, advances
// stake if validated, or challenges if not
func (v *Validator) validationLoop(genesisRoot common.Hash) {
	defer v.Wg.Done()

	// Listen to AssertionCreated event
	assertionEventCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	assertionEventSub, err := v.Rollup.Contract.WatchAssertionCreated(&bind.WatchOpts{Context: v.Ctx}, assertionEventCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer assertionEventSub.Unsubscribe()

	// Current agreed assertion, initalize to genesis assertion
	// TODO: sync from L1 when restart
	lastValidatedAssertion := &rollupTypes.Assertion{
		ID:                    new(big.Int),
		VmHash:                genesisRoot,
		CumulativeGasUsed:     new(big.Int),
		InboxSize:             new(big.Int),
		Deadline:              new(big.Int),
		PrevCumulativeGasUsed: new(big.Int),
	}
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
						log.Crit("UNHANDELED: Can't validate assertion, validator state corrupted", "err", err)
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
					log.Crit("UNHANDELED: concurrent assertion")
					continue
				}
				currentAssertion = assertion
				err := validateCurrentAssertion()
				if err != nil {
					// TODO: error handling instead of panic
					log.Crit("UNHANDELED: Can't validate assertion, validator state corrupted", "err", err)
				}
			case <-v.Ctx.Done():
				return
			}
		}
	}
}

func (v *Validator) Start() error {
	genesis := v.BaseService.Start(true, true)

	v.Wg.Add(3)
	go v.SyncLoop(v.newBatchCh)
	go v.validationLoop(genesis.Root())
	go v.BaseService.ChallengeLoop()
	log.Info("Validator started")
	return nil
}

func (v *Validator) Stop() error {
	log.Info("Validator stopped")
	v.Cancel()
	v.Wg.Wait()
	return nil
}

func (v *Validator) APIs() []rpc.API {
	// TODO: validator APIs
	return []rpc.API{}
}
