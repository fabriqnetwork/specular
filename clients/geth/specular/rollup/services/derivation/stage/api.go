package stage

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

// Represents a stage in a pipeline.
// Generic parameters:
// `T`: Stage output type.
type StageOps[T any] interface {
	// Possible errors returned:
	// - RetryableError: Indicates caller should retry step.
	// - RecoverableError: Indicates caller should perform recovery.
	// - Unrecoverable fatal error (i.e. any other type): Unexpected. Indicates caller should not retry.
	Pull(ctx context.Context) (T, error)
	// Recovers from a re-org to the given L1 block number.
	Recover(ctx context.Context, l1BlockID l2types.BlockID) error
}

type ExecutionBackend interface {
	ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(payload services.ExecutionPayload) error
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
}

type L1State interface {
	Head() l2types.BlockID
	Safe() l2types.BlockID
	Finalized() l2types.BlockID
}

type L1Config interface {
	SequencerInboxAddr() common.Address
	RollupAddr() common.Address
	RollupGenesisBlock() uint64
}

type RetryableError struct{ Err error }
type RecoverableError struct{ Err error }

func (e RetryableError) Error() string   { return e.Err.Error() }
func (e RecoverableError) Error() string { return e.Err.Error() }

type RollupState interface {
	OnAssertionCreated(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) error
	OnAssertionConfirmed(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) error
	OnAssertionRejected(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) error
}

// L1HeaderRetrieval -> L1TxRetrieval -> L1TxProcessing (RollupState + PayloadBuilder) -> L2ForkChoiceUpdate
func CreatePipeline(
	cfg L1Config,
	execBackend ExecutionBackend,
	rollupState RollupState,
	l2ClientCreatorFn func(ctx context.Context) (EthClient, error),
	l1Client EthClient,
	l1State L1State,
) *TerminalStage[l2types.BlockRelation] {
	// Define and chain stages together.
	var (
		// Processors
		daHandlers, rollupTxHandlers = createProcessors(cfg, execBackend, rollupState, l2ClientCreatorFn)
		// Stages
		genesisBlockID         = l2types.NewBlockID(cfg.RollupGenesisBlock(), common.Hash{})
		l1HeaderRetrievalStage = L1HeaderRetrievalStage{genesisBlockID, l1Client}
		l1TxRetrievalStage     = NewStage[l2types.BlockID, filteredBlock](
			&l1HeaderRetrievalStage,
			NewL1TxRetriever(l1Client, createTxFilterFn(daHandlers, rollupTxHandlers)),
		)
		l1TxProcessingStage = NewStage[filteredBlock, l2types.BlockRelation](
			l1TxRetrievalStage,
			NewL1TxProcessor(daHandlers, rollupTxHandlers),
		)
		l2ForkChoiceUpdateStage = NewTerminalStage[l2types.BlockRelation](
			l1TxProcessingStage,
			NewL2ForkChoiceUpdater(execBackend, l1State),
		)
	)
	return l2ForkChoiceUpdateStage
}

func createProcessors(
	cfg L1Config,
	execBackend ExecutionBackend,
	rollupState RollupState,
	l2ClientCreatorFn func(ctx context.Context) (EthClient, error),
) (map[txFilterID]daSourceHandler, map[txFilterID]txHandler) {
	var (
		seqInboxABIMethods = bridge.InboxABIMethods()
		rollupABIMethods   = bridge.RollupABIMethods()
		payloadBuilder     = payloadBuilder{execBackend: execBackend, l2ClientCreatorFn: l2ClientCreatorFn}
	)
	// Define handlers for l1 tx processing.
	daHandlers := map[txFilterID]daSourceHandler{
		{cfg.SequencerInboxAddr(), string(seqInboxABIMethods[bridge.AppendTxBatchFnName].ID)}: payloadBuilder.BuildPayloads,
	}
	rollupTxHandlers := map[txFilterID]txHandler{
		{cfg.RollupAddr(), string(rollupABIMethods[bridge.CreateAssertionFnName].ID)}:                 rollupState.OnAssertionCreated,
		{cfg.RollupAddr(), string(rollupABIMethods[bridge.ConfirmFirstUnresolvedAssertionFnName].ID)}: rollupState.OnAssertionConfirmed,
		{cfg.RollupAddr(), string(rollupABIMethods[bridge.RejectFirstUnresolvedAssertionFnName].ID)}:  rollupState.OnAssertionRejected,
	}
	return daHandlers, rollupTxHandlers
}

func createTxFilterFn(
	daSourceHandlers map[txFilterID]daSourceHandler,
	rollupTxHandlers map[txFilterID]txHandler,
) func(*types.Transaction) bool {
	// Function returns true iff the tx is of a type handled by either a da source or rollup tx handler.
	return func(tx *types.Transaction) bool {
		var to = tx.To()
		if to == nil {
			return false
		}
		var (
			addr     = *to
			methodID = tx.Data()[:bridge.MethodNumBytes]
			filterID = txFilterID{addr, string(methodID)}
		)
		if _, ok := daSourceHandlers[filterID]; ok {
			return true
		}
		if _, ok := rollupTxHandlers[filterID]; ok {
			return true
		}
		return false
	}
}
