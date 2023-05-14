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

type L1TxProcessingStage struct {
	prev             Stage[filteredBlock]
	daSourceHandlers map[txFilterID]daSourceHandler
	rollupTxHandlers map[txFilterID]txHandler
}

// TODO: support other DA sources.
type daSourceHandler func(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) (l2types.BlockRelation, error)
type txHandler func(ctx context.Context, l1BlockID l2types.BlockID, tx *types.Transaction) error

type txFilterID struct {
	contractAddr common.Address
	methodID     string
}

func (s *L1TxProcessingStage) Step(ctx context.Context) (l2types.BlockRelation, error) {
	filteredBlock, err := s.prev.Step(ctx)
	if err != nil {
		return l2types.BlockRelation{}, err
	}
	var relation l2types.BlockRelation
	for _, tx := range filteredBlock.txs {
		contractAddr := tx.To()
		if contractAddr == nil {
			return l2types.BlockRelation{}, fmt.Errorf("`tx.To` is unexpectedly nil")
		}
		methodID := string(tx.Data()[:bridge.MethodNumBytes])
		// Handle tx.
		daHandler := s.daSourceHandlers[txFilterID{*contractAddr, methodID}]
		if daHandler != nil {
			relation, err = daHandler(ctx, filteredBlock.blockID, tx)
			if err != nil {
				return l2types.BlockRelation{}, err
			}
		} else {
			handler := s.rollupTxHandlers[txFilterID{*contractAddr, methodID}]
			err = handler(ctx, filteredBlock.blockID, tx)
			if err != nil {
				return l2types.BlockRelation{}, err
			}
		}
	}
	return relation, err
}

func (s *L1TxProcessingStage) Recover(ctx context.Context, l1BlockID l2types.BlockID) error {
	// TODO: recover handlers
	return s.prev.Recover(ctx, l1BlockID)
}

// Builds payloads from AppendTxBatch transactions.
type payloadBuilder struct {
	execBackend ExecutionBackend
	l2Client    EthClient
}

// TODO: synchronize execBackend
func (b *payloadBuilder) Process(
	ctx context.Context,
	l1BlockID l2types.BlockID,
	tx *types.Transaction,
) (l2types.BlockRelation, error) {
	payloads, err := payloadsFromCalldata(tx)
	if err != nil {
		return l2types.BlockRelation{}, err
	}
	l2Head, err := b.l2Client.BlockNumber(ctx)
	if err != nil {
		return l2types.BlockRelation{}, fmt.Errorf("failed to get latest l2 blockNumber: %w", err)
	}
	for _, payload := range payloads {
		if payload.BlockNumber() >= l2Head {
			log.Warn("Skipping redundant payload", "number", payload.BlockNumber())
			continue
		}
		err = b.execBackend.BuildPayload(&payload)
		if err != nil {
			// TODO: ignore/skip bad batches
			return l2types.BlockRelation{}, fmt.Errorf("failed to build payload: %w", err)
		}
	}
	// Return last block relation.
	lastPayload := payloads[len(payloads)-1]
	header, err := b.l2Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(lastPayload.BlockNumber()))
	if err != nil {
		return l2types.BlockRelation{}, fmt.Errorf("failed to get header for last payload: %w", err)
	}
	relation := l2types.BlockRelation{L1BlockID: l1BlockID, L2BlockID: l2types.NewBlockIDFromHeader(header)}
	return relation, nil
}

func (b *payloadBuilder) recover() error { return nil }

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
