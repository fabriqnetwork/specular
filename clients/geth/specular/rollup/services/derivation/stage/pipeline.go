package stage

import (
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

// L1HeaderRetrieval -> L1TxRetrieval -> L1TxProcessing (RollupState + PayloadBuilder) -> L2ForkChoiceUpdate
func CreatePipeline(
	cfg DerivationConfig,
	execBackend ExecutionBackend,
	rollupState RollupState,
	l2Client L2Client,
	l1Client L1Client,
	l1State L1State,
) *TerminalStage[types.BlockRelation] {
	// Define and chain stages together.
	var (
		// Initialize processors
		daHandlers, rollupTxHandlers = createProcessors(cfg, execBackend, rollupState, l2Client)
		// Initialize stages
		l1HeaderRetrievalStage = L1HeaderRetrievalStage{cfg.GetGenesisL1BlockID(), l1Client}
		l1TxRetrievalStage     = NewStage[types.BlockID, filteredBlock](
			&l1HeaderRetrievalStage,
			NewL1TxRetriever(l1Client, createTxFilterFn(daHandlers, rollupTxHandlers)),
		)
		l1TxProcessingStage = NewStage[filteredBlock, types.BlockRelation](
			l1TxRetrievalStage,
			NewL1TxProcessor(daHandlers, rollupTxHandlers),
		)
		l2ForkChoiceUpdateStage = NewTerminalStage[types.BlockRelation](
			l1TxProcessingStage,
			NewL2ForkChoiceUpdater(cfg, execBackend, l2Client, l1State),
		)
	)
	return l2ForkChoiceUpdateStage
}

func createProcessors(
	cfg L1Config,
	execBackend ExecutionBackend,
	rollupState RollupState,
	l2Client L2Client,
) (map[txFilterID]daSourceHandler, map[txFilterID]txHandler) {
	var (
		seqInboxABIMethods = bridge.InboxABIMethods()
		rollupABIMethods   = bridge.RollupABIMethods()
		payloadBuilder     = NewPayloadBuilder(cfg, execBackend, l2Client)
	)
	// Define handlers for l1 tx processing.
	daHandlers := map[txFilterID]daSourceHandler{
		{cfg.GetSequencerInboxAddr(), string(seqInboxABIMethods[bridge.AppendTxBatchFnName].ID)}: payloadBuilder.BuildPayloads,
	}
	rollupTxHandlers := map[txFilterID]txHandler{
		{cfg.GetRollupAddr(), string(rollupABIMethods[bridge.CreateAssertionFnName].ID)}:                 rollupState.OnAssertionCreated,
		{cfg.GetRollupAddr(), string(rollupABIMethods[bridge.ConfirmFirstUnresolvedAssertionFnName].ID)}: rollupState.OnAssertionConfirmed,
		{cfg.GetRollupAddr(), string(rollupABIMethods[bridge.RejectFirstUnresolvedAssertionFnName].ID)}:  rollupState.OnAssertionRejected,
	}
	return daHandlers, rollupTxHandlers
}

func createTxFilterFn(
	daSourceHandlers map[txFilterID]daSourceHandler,
	rollupTxHandlers map[txFilterID]txHandler,
) func(*ethTypes.Transaction) bool {
	// Function returns true iff the tx is of a type handled by either a da source or rollup tx handler.
	return func(tx *ethTypes.Transaction) bool {
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
