package engine

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/ethengine"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/l2rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// Engine manager.
type Manager struct {
	client  EngineClient
	l2State EthState
	curr    *EngineBuildJob
}

// type EthState interface {
// 	Head() types.L2BlockRef
// 	SafeHead() types.L2BlockRef
// 	Finalized() types.L2BlockRef
// }

type (
	ForkChoiceState   = ethengine.ForkChoiceState
	PayloadAttributes = ethengine.PayloadAttributes
	PayloadID         = ethengine.PayloadID
	PayloadStatus     = ethengine.PayloadStatus
	ExecutionPayload  = ethengine.ExecutionPayload
)

type EngineClient interface {
	ForkchoiceUpdate(
		ctx context.Context,
		update *ForkChoiceState,
		attrs *PayloadAttributes,
	) (*ethengine.ForkChoiceResponse, error)
	NewPayload(ctx context.Context, payload *ExecutionPayload) (*PayloadStatus, error)
	GetPayload(ctx context.Context, payloadID PayloadID) (*ExecutionPayload, error)
	// BuildPayload(ctx context.Context, attrs BuildPayloadAttributes) (*BuildPayloadResponse, error)
}

type EthState struct {
	head      types.L2BlockRef
	safeHead  types.L2BlockRef
	finalized types.L2BlockRef
}

func (s *EthState) Head() types.L2BlockRef      { return s.head }
func (s *EthState) Safe() types.L2BlockRef      { return s.safeHead }
func (s *EthState) Finalized() types.L2BlockRef { return s.finalized }

type EngineBuildJob struct {
	id      ethengine.PayloadID
	onto    types.L2BlockRef
	jobType JobType
}

func (j *EngineBuildJob) ID() ethengine.PayloadID { return j.id }
func (j *EngineBuildJob) Onto() types.L2BlockRef  { return j.onto }
func (j *EngineBuildJob) Type() JobType           { return j.jobType }

type JobType bool

const (
	Unsafe JobType = false
	Safe   JobType = true
)

func NewManager(client EngineClient) *Manager { return &Manager{client: client} }

func (m *Manager) CurrentBuildJob() *EngineBuildJob { return m.curr }
func (m *Manager) L2State() EthState                { return m.l2State }
func (m *Manager) IsBuilding() bool                 { return m.curr != nil }
func (m *Manager) IsBuildingConsistently() bool {
	return m.curr != nil && m.curr.onto.BlockID == m.l2State.Head().BlockID
}

// StartPayload requests the engine to start building a block with the given attributes.
// If updateSafe, the resulting block will be marked as a safe block.
func (m *Manager) StartPayload(
	ctx context.Context,
	parent types.L2BlockRef,
	attrs *PayloadAttributes,
	jobType JobType,
) error {
	if m.curr != nil {
		log.Warn(
			"did not finish previous block building, starting new building now",
			"prev_onto", m.curr.onto,
			"prev_payload_id", m.curr.id,
			"new_onto", parent,
		)
		// TODO (see optimism): maybe worth it to force-cancel the old payload ID here.
	}
	fc := ForkChoiceState{
		HeadBlockHash:      parent.Hash,                // block we're building on
		SafeBlockHash:      m.l2State.Safe().Hash,      // no change
		FinalizedBlockHash: m.l2State.Finalized().Hash, // no change
	}
	id, err := l2rpc.StartPayload(ctx, m.client, fc, attrs)
	if err != nil {
		return err
	}
	m.curr = &EngineBuildJob{id: id, jobType: jobType, onto: parent}
	return nil
}

// ConfirmPayload requests the engine to complete the current block.
// If no block is being built, or if it fails, an error is returned.
func (m *Manager) ConfirmPayload(ctx context.Context) (*ExecutionPayload, error) {
	if m.curr == nil {
		return nil, fmt.Errorf("cannot complete payload building: not currently building a payload") // BlockInsertPrestateErr,
	}
	if m.curr.onto.Hash != m.l2State.Head().Hash { // E.g. when safe-attributes consolidation fails, it will drop the existing work.
		log.Warn(
			"engine is building block that reorgs previous unsafe head",
			"onto", m.curr.onto,
			"unsafe", m.l2State.Head(),
		)
	}
	fc := ForkChoiceState{
		HeadBlockHash:      common.Hash{},         // gets overridden
		SafeBlockHash:      m.l2State.Safe().Hash, // gets overridden if jobType == Safe
		FinalizedBlockHash: m.l2State.Finalized().Hash,
	}
	payload, confirmErr := l2rpc.ConfirmPayload(ctx, m.client, fc, m.curr.id, m.curr.jobType == Safe)
	if confirmErr != nil {
		return nil, fmt.Errorf(
			"failed to complete building on top of L2 chain %s, id: %s, error: %w",
			m.curr.onto, m.curr.id, confirmErr,
		)
	}
	ref := payloadToBlockRef(payload) // default to &m.cfg.Genesis
	// if err != nil {
	// 	return nil NewResetError(fmt.Errorf("failed to decode L2 block ref from payload: %w", err)) // BlockInsertPayloadErr
	// }

	// Update l2 heads.
	m.l2State.head = ref
	if m.curr.jobType == Safe {
		m.l2State.safeHead = ref
		// m.postProcessSafeL2()
	}
	m.curr = nil
	return payload, nil
}

func (m *Manager) Reset(_ context.Context) {
	m.curr = nil
}

func payloadToBlockRef(payload *ExecutionPayload) types.L2BlockRef {
	l1Origin := types.BlockID{} // TODO: placeholder
	return types.NewL2BlockRef(
		payload.Number,
		payload.BlockHash,
		payload.ParentHash,
		l1Origin,
		payload.Timestamp,
	)
}

// func (s *sequencer) buildPayloadFromTxPool(ctx context.Context) error {
// 	var (
// 		attrs       = createBuildPayloadAttributes(s.cfg.GetAccountAddr())
// 		header, err = s.l2Client.HeaderByTag(ctx, eth.Latest)
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch latest header: %w", err)
// 	}
// 	log.Info("Building payload...", "latest", types.NewBlockIDFromHeader(header))
// 	// Build payload (currently synchronous).
// 	_, err = s.engineMgr.BuildPayload(ctx, attrs)
// 	if err != nil {
// 		return fmt.Errorf("failed to build payload: %w", err)
// 	}
// 	header, err = s.l2Client.HeaderByTag(ctx, eth.Latest)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch latest header: %w", err)
// 	}
// 	log.Info("Built payload", "latest", types.NewBlockIDFromHeader(header))
// 	return nil
// }

// Enforces a payload is built only if there are pending txs. Useful for testing.
// func (s *sequencer) buildNonEmptyPayloadFromTxPool(ctx context.Context) error {
// 	status, err := s.l2Client.TxPoolStatus(ctx)
// 	if err != nil {
// 		return fmt.Errorf("Failed to fetch tx pool status: %w", err)
// 	}
// 	// Check if there are pending txs to build a payload from.
// 	numQueued, numPending := uint64(status["queued"]), uint64(status["pending"])
// 	log.Trace("Tx pool status.", "#queued", numQueued, "#pending", numPending)
// 	if numPending <= 0 {
// 		log.Info("Nothing to publish.")
// 		return nil
// 	}
// 	if err := s.buildPayloadFromTxPool(ctx); err != nil {
// 		return err
// 	}
// 	return nil
// }
