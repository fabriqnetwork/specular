package validator

import (
	"context"
	"fmt"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

type Validator struct {
	cfg      Config
	l1TxMgr  TxManager
	l2Client L2Client
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
	return nil
}
