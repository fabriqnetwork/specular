package sequencer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/ethengine"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/engine"
	"github.com/specularl2/specular/clients/geth/specular/utils"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

const sealingDurationEstimate = 50 * time.Millisecond

// Responsible for ordering and executing new transactions.
// TODO: Support:
// - PBS-style ordering
// - remote ordering + "weak DA" in single call (systems conflate these roles)
type Sequencer struct {
	engineMgr    EngineManager
	attrsBuilder PayloadAttributesBuilder
}

type PayloadAttributesBuilder interface {
	BuildPayloadAttributes() (*ethengine.PayloadAttributes, error)
}

type L2Client interface {
	TxPoolStatus(ctx context.Context) (map[string]hexutil.Uint, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}

type (
	sequencerStepStatus uint
	SequencerError      = utils.CategorizedError[sequencerStepStatus]
)

const (
	starting   sequencerStepStatus = iota // sequencer is trying to start building an unsafe payload
	busy                                  // sequencer is blocked by an ongoing operation
	confirming                            // sequencer is trying to confirm an unsafe payload
)

func NewSequencer(engineMgr EngineManager, attrsBuilder PayloadAttributesBuilder) *Sequencer {
	return &Sequencer{engineMgr: engineMgr, attrsBuilder: attrsBuilder}
}

func (s *Sequencer) EngineManager() EngineManager { return s.engineMgr }

// Responsible for building payloads.
func (s *Sequencer) Step(ctx context.Context) *SequencerError {
	if curr := s.engineMgr.CurrentBuildJob(); curr != nil {
		// Skip confirmation if engine is building a safe payload.
		if curr.Type() == engine.Safe {
			log.Warn(
				"avoiding sequencing to not interrupt safe-head changes",
				"onto", curr.Onto(),
				"onto_time", curr.Onto().GetTime(),
			)
			return &SequencerError{Cat: busy, Err: fmt.Errorf("engine is building a safe payload")}
		}
		// Confirm unsafe payload.
		if _, err := s.engineMgr.ConfirmPayload(ctx); err != nil {
			return &SequencerError{Cat: confirming, Err: err}
		}
		return nil
	}
	if err := s.startBuildingPayload(ctx); err != nil {
		return &SequencerError{Cat: starting, Err: err}
	}
	return nil
}

func (s *Sequencer) startBuildingPayload(ctx context.Context) error {
	l2Head := s.engineMgr.L2State().Head()
	attrs, err := s.attrsBuilder.BuildPayloadAttributes()
	if err != nil {
		return err
	}
	// Start a payload building process.
	if err := s.engineMgr.StartPayload(ctx, l2Head, attrs, false); err != nil {
		return fmt.Errorf("failed to start building on top of L2 chain %s, err: %w", l2Head, err)
	}
	return nil
}
