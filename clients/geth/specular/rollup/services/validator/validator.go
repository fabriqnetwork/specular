package validator

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/assertion"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

const interval = 10 * time.Second

type Validator struct {
	cfg               Config
	l2ClientCreatorFn l2ClientCreatorFn
	l2Client          L2Client
	l1TxMgr           TxManager
	rollupState       RollupState
	proofBackend      proof.Backend
}

type l2ClientCreatorFn func(ctx context.Context) (L2Client, error)

func NewValidator(
	cfg Config,
	l2ClientCreatorFn l2ClientCreatorFn,
	l1TxMgr TxManager,
	proofBackend proof.Backend,
	rollupState RollupState,
) *Validator {
	return &Validator{
		cfg:          cfg,
		l2Client:     nil, // Initialized in `Start()`
		l1TxMgr:      l1TxMgr,
		proofBackend: proofBackend,
		rollupState:  rollupState,
	}
}

func (v *Validator) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting validator...")
	// Connect to L2 client.
	l2Client, err := v.l2ClientCreatorFn(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create L2 client: %w", err)
	}
	v.l2Client = l2Client
	if v.cfg.GetIsActiveStaker() {
		if err := v.Stake(ctx); err != nil {
			return fmt.Errorf("Failed to stake: %w", err)
		}
	}
	// end, err := v.SyncL2ChainToL1Head(ctx, v.Config.L1RollupGenesisBlock)
	// go v.SyncLoop(ctx, end+1, nil)
	if v.cfg.GetIsActiveStaker() {
		eg.Go(func() error { return v.eventLoop(ctx) })
	}
	if v.cfg.GetIsActiveChallenger() {
		eg.Go(func() error { return v.challengeLoop(ctx) })
	}
	return nil
}

func (v *Validator) eventLoop(ctx context.Context) error {
	var createdCh = make(chan *bindings.IRollupAssertionCreated)

	var ticker = time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Note: assertions are created/resolved relatively infrequently,
			// compared to batch sequencing. However, the queue may grow larger
			// particularly during disputes.
			if v.cfg.GetIsResolver() {
				err := v.tryResolve(ctx)
				if err != nil {
					return fmt.Errorf("Failed to resolve assertions: %w", err)
				}
			}
			if v.cfg.GetIsActiveCreator() {
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
			if common.Address(ev.AsserterAddr) == v.cfg.GetAccountAddr() {
				// No need to validate, advance stake or resolve (already done above).
				return nil
			}
			err := v.validate()
			if err != nil {
				log.Error("Failed to validate assertion", "error", err)
				if v.cfg.GetIsActiveChallenger() {
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
	_, err := v.l1TxMgr.CreateAssertion(ctx, vmHash, inboxSize)
	if err != nil {
		return nil, err
	}
	staker, err := v.rollupState.GetStaker(ctx, v.cfg.GetAccountAddr())
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
	// Update queued assertion to latest batch
	// vmHash := batch.LastBlockRoot()
	// inboxSize.Add(queuedAssertion.InboxSize, batch.Size())
	// queuedAssertion.EndBlock = batch.LastBlockNumber()
	return common.Hash{}, nil
}

func (v *Validator) Stake(ctx context.Context) error {
	staker, err := v.rollupState.GetStaker(ctx, v.cfg.GetAccountAddr())
	if err != nil {
		return fmt.Errorf("Failed to get staker, to stake, err: %w", err)
	}
	if !staker.IsStaked {
		_, err = v.l1TxMgr.Stake(ctx, big.NewInt(int64(v.cfg.GetStakeAmount())))
		if err != nil {
			return fmt.Errorf("Failed to stake, err: %w", err)
		}
	}
	return nil
}
