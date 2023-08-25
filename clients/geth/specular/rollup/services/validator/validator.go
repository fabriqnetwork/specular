package validator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

type Validator struct {
	cfg      Config
	l1TxMgr  TxManager
	l2Client L2Client

	pendingConfirmation []assertion
}

type assertion struct {
	assertionID      *big.Int
	confirmationTime uint64 // L1 block number at which the assertion is eligible for confirmation.
}

func NewValidator(cfg Config, l1TxMgr TxManager, l2Client L2Client) *Validator {
	return &Validator{cfg: cfg, l1TxMgr: l1TxMgr, l2Client: l2Client}
}

func (v *Validator) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting batch disseminator...")
	if err := v.l2Client.EnsureDialed(ctx); err != nil {
		return fmt.Errorf("failed to create L2 client: %w", err)
	}
	eg.Go(func() error { return v.start(ctx) })
	log.Info("Validator started")
	return nil
}

func (v *Validator) start(ctx context.Context) error {
	var ticker = time.NewTicker(v.cfg.GetAssertInterval())
	defer ticker.Stop()
	v.step(ctx)
	for {
		select {
		case <-ticker.C:
			if err := v.step(ctx); err != nil {
				log.Errorf("Failed to step: %w", err)
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

// TODO: implement.
func (v *Validator) step(ctx context.Context) error {
	if err := v.createAssertion(); err != nil {
		return fmt.Errorf("failed to create assertion: %w", err)
	}
	// TODO: or reject, depending on circumstances.
	if err := v.confirmFirstUnresolvedAssertion(); err != nil {
		return fmt.Errorf("failed to confirm assertion: %w", err)
	}
	return nil
}

// If enough time has passed and txs have been sequenced to L1, create a new assertion.
// Add it to the queue for confirmation.
func (v *Validator) createAssertion() error { return nil }

// If the first unresolved assertion is eligible for confirmation, trigger its confirmation.
func (v *Validator) confirmFirstUnresolvedAssertion() error { return nil }
