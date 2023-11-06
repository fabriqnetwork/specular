package validator

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/bridge"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/services/api"
	specularTypes "github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

var transactTimeout = 10 * time.Minute

type unexpectedSystemStateError struct{ msg string }

func (e unexpectedSystemStateError) Error() string {
	return fmt.Sprintf("service entered unexpected state: %s", e.msg)
}

type Validator struct {
	cfg            Config
	l1TxMgr        TxManager
	l1BridgeClient BridgeClient
	l1State        EthState
	l2Client       L2Client

	lastCreatedAssertionAttrs assertionAttributes
}

type assertionAttributes struct {
	l2BlockNum        uint64
	l2StateCommitment specularTypes.Bytes32
}

func NewValidator(
	cfg Config,
	l1TxMgr TxManager,
	l1BridgeClient BridgeClient,
	l1State EthState,
	l2Client L2Client,
) *Validator {
	return &Validator{cfg: cfg, l1TxMgr: l1TxMgr, l1BridgeClient: l1BridgeClient, l1State: l1State, l2Client: l2Client}
}

func (v *Validator) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting validator...")
	if err := v.l2Client.EnsureDialed(ctx); err != nil {
		return fmt.Errorf("failed to create L2 client: %w", err)
	}
	eg.Go(func() error { return v.start(ctx) })
	log.Info("Validator started")
	return nil
}

// Advances validator step-by-step.
func (v *Validator) start(ctx context.Context) error {
	var ticker = time.NewTicker(v.cfg.GetValidationInterval())
	defer ticker.Stop()
	// TODO: do this in the L2 consensus client, not here.
	if err := v.validateGenesis(ctx); err != nil {
		return fmt.Errorf("failed to validate genesis: %w", err)
	}
	if err := v.ensureStaked(ctx); err != nil {
		return fmt.Errorf("failed to ensure validator is staked: %w", err)
	}
	if err := v.rollback(ctx); err != nil {
		return fmt.Errorf("failed to initialize state: %w", err)
	}
	for {
		select {
		case <-ticker.C:
			if err := v.step(ctx); err != nil {
				log.Errorf("Failed to advance: %w", err)
				if errors.As(err, &unexpectedSystemStateError{}) {
					return fmt.Errorf("aborting: %w", err)
				}
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

// Attempts to create a new assertion and confirm an existing assertion.
func (v *Validator) step(ctx context.Context) error {
	// Try to create a new assertion.
	// TODO: do this only if configured to be an active validator.
	if err := v.createAssertion(ctx); err != nil {
		return fmt.Errorf("failed to create assertion: %w", err)
	}
	// TODO: validate assertions locally.
	// Resolve the first unresolved assertion.
	if err := v.resolveFirstUnresolvedAssertion(ctx); err != nil {
		return fmt.Errorf("failed to resolve assertion: %w", err)
	}
	return nil
}

// If enough time has passed and txs have been sequenced to L1, create a new assertion.
// Add it to the queue for confirmation.
func (v *Validator) createAssertion(ctx context.Context) error {
	assertionAttrs, err := v.getNextAssertionAttrs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next assertion attrs: %w", err)
	}
	// TODO fix assumptions: not reorg-resistant. Other validators may have inserted new assertions.
	if assertionAttrs.l2BlockNum <= v.lastCreatedAssertionAttrs.l2BlockNum {
		log.Info("No new blocks to create assertion for yet.")
		return nil
	}
	cCtx, cancel := context.WithTimeout(ctx, transactTimeout)
	defer cancel()
	// TOOD: GasLimit: 0 ...?
	receipt, err := v.l1TxMgr.CreateAssertion(cCtx, assertionAttrs.l2StateCommitment, big.NewInt(0).SetUint64(assertionAttrs.l2BlockNum))
	if err != nil {
		return err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		log.Error("Tx successfully published but reverted", "tx_hash", receipt.TxHash)
	} else {
		log.Info("Tx successfully published", "tx_hash", receipt.TxHash)
		log.Info("Created assertion", "l2Block#", assertionAttrs.l2BlockNum)
		v.lastCreatedAssertionAttrs = assertionAttrs
	}
	return nil
}

// If the first unresolved assertion is eligible for confirmation, trigger its confirmation. Otherwise, wait.
// TODO: reject or challenge, depending on circumstances.
func (v *Validator) resolveFirstUnresolvedAssertion(ctx context.Context) error {
	// TODO: so maybe us as a validator we are probably staked... shall we maybe
	//       run some checks first before we spend all this money on gas if we're properly staked?
	//		Or maybe we check this upfront when we run this? Just food for thought.

	// TODO: For security we should check what assertion are we dealing with and keep ID in memory
	//       same with its parent to assure we have a proper relation in the validation phase.
	//	     This is used later on
	// assertionID = v.lastCreatedAssertionAttrs.<ID> (? field not implemented, where to get ID from ?)
	// parentAssertionID, err := v.l1BridgeClient.GetLastConfirmedAssertionID(ctx)
	// FIXME: remove faking below
	assertionID := big.NewInt(0)
	parentAssertionID := big.NewInt(1)

	// FIXME: All IRollup methods should deal with the assertionID as an argument
	//        to assure we operate only on the right assertion

	// Simulate a confirmation attempt.
	err := v.l1BridgeClient.RequireFirstUnresolvedAssertionIsConfirmable(ctx)
	// err -> "No"
	if err != nil {
		// "No"
		errStr := err.Error()
		if errStr == bridge.NoUnresolvedAssertionErr {
			log.Trace("No unresolved assertion to resolve.")
		} else if errStr == bridge.ConfirmationPeriodPendingErr {
			log.Trace("Too early to confirm first unresolved assertion.")
		} else {
			return &unexpectedSystemStateError{"failed to validate assertion (breaks current assumptions): " + err.Error()}
		}

		// If not confirmable could still be rejectable
		// TODO: confirm that the stakerAddress should come from configuration
		err := v.l1BridgeClient.RequireFirstUnresolvedAssertionIsRejectable(ctx, v.cfg.GetAccountAddr())
		if err != nil {
			// It is not rejectable
			// TODO: confirm there is no need to do errStr := err.Error() and cmp with IRollup "requireFirstUnresolvedAssertionIsRejectable"
			log.Trace("No unresolved assertion to be rejected.")
			// TODO: confirm exit with nil below
			return nil
		}

		// It is not confirmable, but it is rejectable so let's reject
		_, err = v.l1BridgeClient.RejectFirstUnresolvedAssertion(ctx, v.cfg.GetAccountAddr())
		if err != nil {
			// It should be rejectable, but it failed to reject, so it must be a bigger problem
			// TODO: confirm solution
			return err
		}
	}

	// At this point we know it's confirmable
	// But first we need to validate
	// TODO: this is not implemented properly, fill-in with a stub/mock

	// We take the parent block number, we take its hash, we reply what's in the assertion
	// and then check the value with the one in the assertion.
	// (Tip from Simon: Specular uses VMRoot hash from go-ethereum)

	assertion, err := v.l1BridgeClient.GetAssertion(ctx, assertionID)
	if err != nil {
		return err
	}
	parentAssertion, err := v.l1BridgeClient.GetAssertion(ctx, parentAssertionID)
	if err != nil {
		return err
	}

	// Validate block number
	nextBlockNumb := big.NewInt(0)
	if nextBlockNumb.Add(parentAssertion.BlockNum, big.NewInt(1)).Cmp(assertion.BlockNum) == 0 {
		return fmt.Errorf("block number does not match with the parent assertion")
	}
	// Validate parent hash
	if parentAssertion.StateHash != assertion.StateHash {
		return fmt.Errorf("state hash does not match with the parent assertion")
	}
	
	// number, err := v.l2Client.BlockNumber()
	// if err != nil {
	// 	return err
	// }
	// v.l2Client.<Implement> from go-ethereum Client GetHead()->Header
	// if header.Hash == parentAssertion.Hash {
	// TODO: indicates we are out of sync
	// FIXME: that requires going back to a state of the last confirmed assertion, idea to handle:
	//        - wait for sync of the execution node (?)
	// Timeout, return sth... TBD
	// }
	// At this point confirmed we are in sync

	// TODO: Read appendTxBatch(bytes calldata txBatchData) from <Implement> or maybe L1Bridge client...
	// 	     anyway is to get the:
	//	     - calldata content
	//	     - parse
	//	     - get full transactions data
	//	     - apply to the parent state of the assertion (i.e. recompute the entire tree from the transaction hashes)
	// TODO: compare the state hash we end up with with the assertion hash
	// TODO: Significant optimization could be done here IMHO, sth around the partial proof (TBD)
	// 		 - the call data could include partial-proof

	// if assertionHash != parentHash {
	// TODO: Research if there is maybe some more work to be done before we challenge the assertion
	// TODO: Challenge assertion
	// err := IRollup challengeAssertion(
	//	players{
	//		us <- stakerAddress,
	//		who's them? <- probably the staker who submitted the assertion, right?
	//				This can come from:
	//				- IRollupStaker.AssertionID, so reversing this we can identify the Staker
	//				- IRollup.CreateAssertion releases an event which logs the Staker address
	//	}
	//	)
	//  if err != nil { // TODO: sth fundamental has gone wrong ... }
	//  return nil
	// }

	// We could confirm after all the checks
	cCtx, cancel := context.WithTimeout(ctx, transactTimeout)
	defer cancel()

	// FIXME: the last assertion may have changed in the mean-time,
	//        we should include a param to confirm specific assertion ID only (assertionID)
	_, err = v.l1TxMgr.ConfirmFirstUnresolvedAssertion(cCtx)
	if err != nil {
		return fmt.Errorf("failed to confirm assertion: %w", err)
	}
	log.Info("Confirmed assertion")
	return nil
}

// Rolls back local validator state, using the current L1 contract state as a checkpoint.
func (v *Validator) rollback(ctx context.Context) error {
	staker, err := v.l1BridgeClient.GetStaker(ctx, v.cfg.GetAccountAddr())
	if err != nil {
		return fmt.Errorf("failed to get staker: %w", err)
	}
	assertion, err := v.l1BridgeClient.GetAssertion(ctx, staker.AssertionID)
	if err != nil {
		return fmt.Errorf("failed to get assertion: %w", err)
	}
	v.lastCreatedAssertionAttrs = assertionAttributes{assertion.BlockNum.Uint64(), assertion.StateCommitment}
	return nil
}

// Gets the next assertion's attributes.
// We can relax this to get the next safe assertion's attributes but need to handle reorgs.
func (v *Validator) getNextAssertionAttrs(ctx context.Context) (assertionAttributes, error) {
	header, err := v.l2Client.HeaderByTag(ctx, eth.Safe)
	if err != nil {
		return assertionAttributes{}, fmt.Errorf("failed to get finalized assertion attrs: %w", err)
	}
	return assertionAttributes{header.Number.Uint64(), StateCommitment(&StateCommitmentV0{header.Hash()})}, nil
}

func (v *Validator) ensureStaked(ctx context.Context) error {
	staker, err := v.l1BridgeClient.GetStaker(ctx, v.cfg.GetAccountAddr())
	if err != nil {
		return fmt.Errorf("failed to get staker: %w", err)
	}
	if staker.IsStaked {
		log.Info("Already staked.")
		return nil
	}
	amount, err := v.l1BridgeClient.GetRequiredStakeAmount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get stake amount: %w", err)
	}
	_, err = v.l1TxMgr.Stake(ctx, amount)
	if err != nil {
		return fmt.Errorf("failed to stake: %w", err)
	}
	log.Info("Staked successfully.", "amount", amount)
	return nil
}

// TODO: refactor.
func (v *Validator) validateGenesis(ctx context.Context) error {
	assertion, err := v.l1BridgeClient.GetAssertion(ctx, common.Big0)
	if err != nil {
		return fmt.Errorf("failed to get genesis assertion: %w", err)
	}
	stateCommitment := assertion.StateCommitment
	// Check that the genesis assertion is correct.
	genesisBlock, err := v.l2Client.BlockByNumber(ctx, common.Big0)
	if err != nil {
		return fmt.Errorf("failed to get L2 genesis block: %w", err)
	}
	genesisStateCommitment := StateCommitment(&StateCommitmentV0{genesisBlock.Header().Hash()})
	if stateCommitment != genesisStateCommitment {
		return fmt.Errorf("mismatching genesis on L1=%s vs L2=%s", &stateCommitment, &genesisStateCommitment)
	}
	return nil
}
