package stage

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

// `Stage` defines a stage in a pipeline.
// TODO: split: Get + Advance to handle failures bette
type Stage[T any] interface {
	// `in` is the input from the downstream stage, if any.
	// Convention: use interface{} for T if no input and U if no output
	// Possible errors returned:
	// - RetryableError: Indicates caller should retry step.
	// - RecoverableError: Indicates caller should perform recovery.
	// - Unrecoverable fatal error (i.e. any other type): Unexpected. Indicates caller should not retry.
	Step(ctx context.Context) (T, error)
	// Recover from a re-org to the given L1 block number.
	Recover(ctx context.Context, l1BlockID l2types.BlockID) error
}

type ExecutionBackend interface {
	ForkchoiceUpdate(update services.ForkchoiceState) error
	BuildPayload(payload services.ExecutionPayload) error
}

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
}

type RetryableError struct{ err error }
type RecoverableError struct{ err error }

func (e RetryableError) Error() string   { return e.err.Error() }
func (e RecoverableError) Error() string { return e.err.Error() }

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
	l2Client EthClient,
	l1Client EthClient,
	l1State L1State,
) Stage[interface{}] {
	var (
		seqInboxABIMethods = bridge.InboxABIMethods()
		rollupABIMethods   = bridge.RollupABIMethods()
		payloadBuilder     = payloadBuilder{execBackend, l2Client}
	)
	// Define handlers for l1 tx processing.
	daHandlers := map[txFilterID]daSourceHandler{
		{cfg.SequencerInboxAddr(), string(seqInboxABIMethods[bridge.AppendTxBatchFnName].ID)}: payloadBuilder.Process,
	}
	rollupTxHandlers := map[txFilterID]txHandler{
		{cfg.RollupAddr(), string(rollupABIMethods[bridge.CreateAssertionFnName].ID)}:                 rollupState.OnAssertionCreated,
		{cfg.RollupAddr(), string(rollupABIMethods[bridge.ConfirmFirstUnresolvedAssertionFnName].ID)}: rollupState.OnAssertionConfirmed,
		{cfg.RollupAddr(), string(rollupABIMethods[bridge.RejectFirstUnresolvedAssertionFnName].ID)}:  rollupState.OnAssertionRejected,
	}
	// Define chained stages.
	l1HeaderRetrievalStage := L1HeaderRetrievalStage{l2types.BlockID{}, l1Client}
	l1TxRetrievalStage := L1TxRetrievalStage{prev: &l1HeaderRetrievalStage, filterFn: createTxFilterFn(daHandlers, rollupTxHandlers)}
	l1TxProcStage := L1TxProcessingStage{prev: &l1TxRetrievalStage, daSourceHandlers: daHandlers, rollupTxHandlers: rollupTxHandlers}
	l2ForkChoiceUpdateStage := NewL2ForkChoiceUpdateStage(&l1TxProcStage, execBackend, l1State)
	return l2ForkChoiceUpdateStage
}

func createTxFilterFn(
	daSourceHandlers map[txFilterID]daSourceHandler,
	rollupTxHandlers map[txFilterID]txHandler,
) func(*types.Transaction) bool {
	// Function returns true iff the tx is of a type handled by either a da source or rollup tx handler.
	filterFn := func(tx *types.Transaction) bool {
		to := tx.To()
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
	return filterFn
}
