package stage

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type L1TxProcessor struct {
	daSourceHandlers map[txFilterID]daSourceHandler
	rollupTxHandlers map[txFilterID]txHandler
	// state
	lastProcessedRelation l2types.BlockRelation
	lastProcessedIdx      int
}

// TODO: support other DA sources.
type daSourceHandler func(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) (l2types.BlockRelation, error)
type txHandler func(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) error

type txFilterID struct {
	contractAddr common.Address
	methodID     string
}

func NewL1TxProcessor(
	daSourceHandlers map[txFilterID]daSourceHandler,
	rollupTxHandlers map[txFilterID]txHandler,
) *L1TxProcessor {
	return &L1TxProcessor{daSourceHandlers: daSourceHandlers, rollupTxHandlers: rollupTxHandlers}
}

func (s *L1TxProcessor) hasNext() bool {
	return s.lastProcessedRelation != l2types.BlockRelation{}
}

func (s *L1TxProcessor) next() l2types.BlockRelation {
	next := s.lastProcessedRelation
	s.lastProcessedIdx = 0
	s.lastProcessedRelation = l2types.BlockRelation{}
	return next
}

// Processes transactions in the given filtered block, according to their method IDs.
func (s *L1TxProcessor) ingest(ctx context.Context, filteredL1Block filteredBlock) error {
	for i, tx := range filteredL1Block.txs {
		// Skip txs already processed, if any.
		// This happens if there was an error on a previous call mid-block.
		if i <= s.lastProcessedIdx {
			continue
		}
		contractAddr := tx.To()
		if contractAddr == nil {
			return fmt.Errorf("`tx.To` is unexpectedly nil") // fatal
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
			s.lastProcessedRelation = relation
		} else if rollupHandler != nil {
			err := rollupHandler(ctx, filteredL1Block.blockID, tx)
			if err != nil {
				return fmt.Errorf("rollup tx handler failed for `methodID`=%s: %w", methodID, err)
			}
		} else {
			return fmt.Errorf("no handler found for `methodID`=%s", methodID) // fatal
		}
		s.lastProcessedIdx = i
	}
	return nil
}

func (s *L1TxProcessor) recover(ctx context.Context, l1BlockID l2types.BlockID) error {
	s.lastProcessedIdx = 0
	s.lastProcessedRelation = l2types.BlockRelation{}
	// TODO: recover handlers
	return nil
}

// Builds payloads from AppendTxBatch transactions.
type payloadBuilder struct {
	execBackend       ExecutionBackend
	l2ClientCreatorFn func(context.Context) (EthClient, error)
	l2Client          EthClient // lazily initialized
}

func NewPayloadBuilder(
	execBackend ExecutionBackend,
	l2ClientCreatorFn func(context.Context) (EthClient, error),
) *payloadBuilder {
	return &payloadBuilder{execBackend: execBackend, l2ClientCreatorFn: l2ClientCreatorFn}
}

// TODO: synchronize on execBackend
func (b *payloadBuilder) BuildPayloads(
	ctx context.Context,
	l1BlockID l2types.BlockID,
	tx *types.Transaction,
) (l2types.BlockRelation, error) {
	if err := b.ensureClientInit(ctx); err != nil {
		return l2types.BlockRelation{}, fmt.Errorf("failed to initialize payload builder: %w", err)
	}
	payloads, err := payloadsFromCalldata(tx)
	if err != nil {
		return l2types.BlockRelation{}, fmt.Errorf("Could not decode payloads from calldata: %w", err)
	}
	l2Head, err := b.l2Client.BlockNumber(ctx)
	if err != nil {
		return l2types.BlockRelation{}, RetryableError{fmt.Errorf("failed to get latest l2 blockNumber: %w", err)}
	}
	for _, payload := range payloads {
		if payload.BlockNumber() >= l2Head {
			log.Warn("Skipping redundant payload", "number", payload.BlockNumber())
			continue
		}
		err = b.execBackend.BuildPayload(&payload)
		if err != nil {
			// Current assumption: all batches are valid.
			// TODO: ignore/skip invalid batches.
			return l2types.BlockRelation{}, fmt.Errorf("failed to build payload: %w", err)
		}
	}
	// Return last block relation.
	lastPayload := payloads[len(payloads)-1]
	header, err := b.l2Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(lastPayload.BlockNumber()))
	if err != nil {
		// TODO: retryable?
		return l2types.BlockRelation{}, fmt.Errorf("failed to get header for last payload: %w", err)
	}
	relation := l2types.BlockRelation{L1BlockID: l1BlockID, L2BlockID: l2types.NewBlockIDFromHeader(header)}
	return relation, nil
}

func (b *payloadBuilder) ensureClientInit(ctx context.Context) error {
	if b.l2Client == nil {
		l2Client, err := b.l2ClientCreatorFn(ctx)
		if err != nil {
			return RetryableError{fmt.Errorf("failed to create l2 client: %w", err)} // TODO: retryable?
		}
		b.l2Client = l2Client
	}
	return nil
}

// Stateless processor -- no-op.
// func (b *payloadBuilder) recover() error { return nil }

func payloadsFromCalldata(tx *types.Transaction) ([]l2types.DerivationBlock, error) {
	// Decode input to appendTxBatch transaction.
	decoded, err := bridge.UnpackAppendTxBatchInput(tx)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Construct blocks.
	blocks, err := l2types.BlocksFromDecoded(decoded)
	// TODO: handle bad encoding (reject batch)
	if err != nil {
		return nil, fmt.Errorf("Failed to split AppendTxBatch input into blocks, err: %w", err)
	}
	return blocks, nil
}
