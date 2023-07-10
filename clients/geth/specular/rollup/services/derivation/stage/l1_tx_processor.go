package stage

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/engine"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/da"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

const unset = -1

type L1TxProcessor struct {
	daSourceHandlers map[txFilterID]daSourceHandler
	rollupTxHandlers map[txFilterID]txHandler
	// state
	lastProcessedRelation types.BlockRelation
	numTxsRemaining       int
}

// TODO: support other DA sources.
type daSourceHandler func(ctx context.Context, l1BlockID types.BlockID, tx *ethTypes.Transaction) (types.BlockRelation, error)
type txHandler func(ctx context.Context, l1BlockID types.BlockID, tx *ethTypes.Transaction) error

type txFilterID struct {
	contractAddr common.Address
	methodID     string
}

func NewL1TxProcessor(
	daSourceHandlers map[txFilterID]daSourceHandler,
	rollupTxHandlers map[txFilterID]txHandler,
) *L1TxProcessor {
	return &L1TxProcessor{
		daSourceHandlers: daSourceHandlers,
		rollupTxHandlers: rollupTxHandlers,
		numTxsRemaining:  unset,
	}
}

func (s *L1TxProcessor) hasNext() bool {
	return s.lastProcessedRelation != types.EmptyRelation && s.numTxsRemaining > 0
}

func (s *L1TxProcessor) next() types.BlockRelation {
	next := s.lastProcessedRelation
	s.numTxsRemaining = unset
	s.lastProcessedRelation = types.EmptyRelation
	return next
}

// Processes transactions in the given filtered block, according to their method IDs.
func (s *L1TxProcessor) ingest(ctx context.Context, filteredL1Block filteredBlock) error {
	if len(filteredL1Block.txs) == 0 {
		return nil
	}
	log.Info("Ingesting L1 txs", "#txs", len(filteredL1Block.txs))
	// First time seeing this block.
	if s.numTxsRemaining == -1 {
		s.numTxsRemaining = len(filteredL1Block.txs)
	}
	numProcessed := len(filteredL1Block.txs) - s.numTxsRemaining
	// Skip txs already processed, if any.
	// This happens if there was an error on a previous call mid-block.
	for _, tx := range filteredL1Block.txs[numProcessed:] {
		log.Info("Processing L1 tx", "block#", filteredL1Block.blockID.Number)
		contractAddr := tx.To()
		if contractAddr == nil {
			return fmt.Errorf("tx.To is unexpectedly nil") // fatal
		}
		var (
			methodID      = bridge.TxMethodID(tx)
			daHandler     = s.daSourceHandlers[txFilterID{*contractAddr, methodID}]
			rollupHandler = s.rollupTxHandlers[txFilterID{*contractAddr, methodID}]
		)
		// Handle tx according to its method ID.
		if daHandler != nil {
			relation, err := daHandler(ctx, filteredL1Block.blockID, tx)
			if err != nil {
				return fmt.Errorf("DA handler failed for methodID: %w", err)
			}
			log.Info("Processed DA tx.")
			s.lastProcessedRelation = relation // todo: ft bug
		} else if rollupHandler != nil {
			err := rollupHandler(ctx, filteredL1Block.blockID, tx)
			if err != nil {
				return fmt.Errorf("rollup tx handler failed for methodID=%s: %w", methodID, err)
			}
			log.Info("Processed L1 rollup tx.")
		} else {
			return fmt.Errorf("no handler found for `methodID`=%s", methodID) // fatal
		}
		s.numTxsRemaining--
	}
	return nil
}

func (s *L1TxProcessor) recover(ctx context.Context, l1BlockID types.BlockID) error {
	s.numTxsRemaining = unset
	s.lastProcessedRelation = types.EmptyRelation
	// TODO: recover handlers
	return nil
}

// Builds payloads from AppendTxBatch transactions.
type payloadBuilder struct {
	l1Config    L1Config
	execBackend ExecutionBackend
	l2Client    L2Client
}

func NewPayloadBuilder(
	l1Config L1Config,
	execBackend ExecutionBackend,
	l2Client L2Client,
) *payloadBuilder {
	return &payloadBuilder{l1Config: l1Config, execBackend: execBackend, l2Client: l2Client}
}

func (b *payloadBuilder) BuildPayloads(
	ctx context.Context,
	l1BlockID types.BlockID,
	tx *ethTypes.Transaction,
) (types.BlockRelation, error) {
	if err := b.l2Client.EnsureDialed(ctx); err != nil {
		return types.EmptyRelation, NewRetryableError(fmt.Errorf("failed to create l2 client: %w", err))
	}
	// Parse tx calldata into l2 l2Blocks.
	l2Blocks, err := blocksFromTxCalldata(tx)
	if err != nil {
		return types.EmptyRelation, fmt.Errorf("Could not decode payloads from calldata: %w", err)
	}
	// Get current head to make sure we're not building old blocks.
	l2Head, err := b.l2Client.BlockNumber(ctx)
	if err != nil {
		return types.EmptyRelation, NewRetryableError(fmt.Errorf("failed to get latest l2 blockNumber: %w", err))
	}
	log.Info("Latest L2 block number", "number", l2Head)
	// Get sender address.
	signer := ethTypes.NewLondonSigner(big.NewInt(0).SetUint64(b.l1Config.GetChainID()))
	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		return types.EmptyRelation, fmt.Errorf("failed to get tx sender: %w", err)
	}
	// Build payloads.
	log.Info("Candidate payloads", "#blocks", len(l2Blocks), "from", from.Hex())
	var attrs *engine.BuildPayloadAttributes
	for _, l2Block := range l2Blocks {
		if l2Block.BlockNumber() <= l2Head {
			log.Warn("Skipping redundant l2 block", "number", l2Block.BlockNumber())
			continue
		}
		for i := uint64(0); i <= l2Head; i++ {
			a, _ := b.l2Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(i))
			log.Warn("BLONK", "#", a.Number, "hash", a.Hash())
		}
		log.Info("Building payload", "l2_block#", l2Block.BlockNumber(), "curr_head", l2Head)
		attrs = engine.NewBuildPayloadAttributes(l2Block.Timestamp(), common.Hash{}, from, l2Block.Txs(), true)
		_, err = b.execBackend.BuildPayload(ctx, attrs)
		if err != nil {
			// Current assumption: all batches are valid.
			// TODO: ignore/skip invalid batches.
			return types.EmptyRelation, fmt.Errorf("failed to build payload: %w", err)
		}
	}
	if attrs == nil {
		log.Info("Skipped all l2 blocks.")
		return types.EmptyRelation, nil
	}
	// Return last block relation.
	lastPayload := l2Blocks[len(l2Blocks)-1]
	header, err := b.l2Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(lastPayload.BlockNumber()))
	if err != nil {
		// TODO: retryable?
		return types.EmptyRelation, fmt.Errorf("failed to get header for last payload: %w", err)
	}
	relation := types.BlockRelation{L1BlockID: l1BlockID, L2BlockID: types.NewBlockIDFromHeader(header)}
	return relation, nil
}

// Stateless processor -- no-op.
// func (b *payloadBuilder) recover() error { return nil }

func blocksFromTxCalldata(tx *ethTypes.Transaction) ([]da.DerivationBlock, error) {
	// Decode input to appendTxBatch transaction.
	l1TxData, err := bridge.UnpackAppendTxBatchInput(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Construct blocks.
	blocks, err := da.BlocksFromData(l1TxData)
	// TODO: handle bad encoding (reject batch)
	if err != nil {
		return nil, fmt.Errorf("failed to split AppendTxBatch input into blocks, err: %w", err)
	}
	return blocks, nil
}
