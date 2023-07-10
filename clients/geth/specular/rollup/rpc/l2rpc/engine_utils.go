package l2rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/ethengine"
	"github.com/specularl2/specular/clients/geth/specular/utils"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

type engineClient interface {
	ForkchoiceUpdate(
		ctx context.Context,
		update *ForkChoiceState,
		attrs *PayloadAttributes,
	) (*ethengine.ForkChoiceResponse, error)
	NewPayload(ctx context.Context, payload *ExecutionPayload) (*PayloadStatus, error)
	GetPayload(ctx context.Context, payloadID PayloadID) (*ExecutionPayload, error)
	// BuildPayload(ctx context.Context, attrs BuildPayloadAttributes) (*BuildPayloadResponse, error)
}

type (
	PayloadID     = ethengine.PayloadID
	EngineError   = utils.CategorizedError[EngineErrType]
	EngineErrType uint
)

var EmptyPayloadID = ethengine.EmptyPayloadID

const (
	// BlockInsertTemporaryErr indicates that the insertion failed but may succeed at a later time without changes to the payload.
	BlockInsertTemporaryErr EngineErrType = iota
	// BlockInsertPrestateErr indicates that the pre-state to insert the payload could not be prepared, e.g. due to missing chain data.
	BlockInsertPrestateErr
	// BlockInsertPayloadErr indicates that the payload was invalid and cannot become canonical.
	BlockInsertPayloadErr
)

// StartPayload starts an execution payload building process in the provided Engine, with the given attributes.
// The severity of the error is distinguished to determine whether the same payload attributes may be re-attempted later.
func StartPayload(
	ctx context.Context,
	client engineClient,
	fc ForkChoiceState,
	attrs *PayloadAttributes,
) (PayloadID, *EngineError) {
	fcRes, err := client.ForkchoiceUpdate(ctx, &fc, attrs)
	if err != nil {
		if rpcErr, ok := err.(rpc.Error); ok {
			switch rpcErr.ErrorCode() {
			case ethengine.InvalidForkchoiceState:
				return EmptyPayloadID, &EngineError{
					Cat: BlockInsertPrestateErr,
					Err: fmt.Errorf("pre-block-creation forkchoice update was inconsistent with engine, need reset to resolve: %w", rpcErr),
				}
			case ethengine.InvalidPayloadAttributes:
				return EmptyPayloadID, &EngineError{
					Cat: BlockInsertPayloadErr,
					Err: fmt.Errorf("payload attributes are not valid, cannot build block: %w", rpcErr),
				}
			default:
				return EmptyPayloadID, &EngineError{
					Cat: BlockInsertPrestateErr,
					Err: fmt.Errorf("unexpected error code in forkchoice-updated response: %w", err),
				}
			}
		} else {
			return EmptyPayloadID, &EngineError{
				Cat: BlockInsertTemporaryErr,
				Err: fmt.Errorf("failed to create new block via forkchoice: %w", err),
			}
		}
	}
	switch ethengine.ExecutePayloadStatus(fcRes.PayloadStatus.Status) {
	case ethengine.ExecutionValid:
		id := fcRes.PayloadID
		if id == nil {
			return EmptyPayloadID, &EngineError{
				Cat: BlockInsertTemporaryErr,
				Err: fmt.Errorf("nil id in forkchoice result when expecting a valid ID"),
			}
		}
		return *id, nil
	// TODO(see optimism `StartPayload`): snap sync - specify explicit different error type if node is syncing
	case ethengine.ExecutionInvalid, ethengine.ExecutionInvalidBlockHash:
		return EmptyPayloadID, &EngineError{Cat: BlockInsertPayloadErr, Err: forkchoiceUpdateErr(fcRes.PayloadStatus)}
	default:
		return EmptyPayloadID, &EngineError{Cat: BlockInsertTemporaryErr, Err: forkchoiceUpdateErr(fcRes.PayloadStatus)}
	}
}

// ConfirmPayload ends an execution payload building process in the provided Engine, and persists the payload as the canonical head.
// If updateSafe is true, then the payload will also be recognized as safe-head at the same time.
// The severity of the error is distinguished to determine whether the payload was valid and can become canonical.
func ConfirmPayload(
	ctx context.Context,
	client engineClient,
	fc ForkChoiceState,
	id PayloadID,
	updateSafe bool,
) (*ExecutionPayload, *EngineError) {
	payload, err := client.GetPayload(ctx, id)
	if err != nil {
		// even if it is an input-error (unknown payload ID), it is temporary,
		// since we will re-attempt the full payload building, not just the retrieval of the payload.
		return nil, &EngineError{Cat: BlockInsertTemporaryErr, Err: fmt.Errorf("failed to get execution payload: %w", err)}
	}
	// TODO: sanityCheckPayload(payload)

	status, err := client.NewPayload(ctx, payload)
	if err != nil {
		return nil, &EngineError{Cat: BlockInsertTemporaryErr, Err: fmt.Errorf("failed to insert execution payload: %w", err)}
	}
	switch status.Status {
	case ethengine.ExecutionValid:
	case ethengine.ExecutionInvalid, ethengine.ExecutionInvalidBlockHash:
		return nil, &EngineError{Cat: BlockInsertPayloadErr, Err: newPayloadErr(payload, status)}
	default:
		return nil, &EngineError{Cat: BlockInsertTemporaryErr, Err: newPayloadErr(payload, status)}
	}

	fc.HeadBlockHash = payload.BlockHash
	if updateSafe {
		fc.SafeBlockHash = payload.BlockHash
	}
	fcRes, err := client.ForkchoiceUpdate(ctx, &fc, nil)
	if err != nil {
		if rpcErr, ok := err.(rpc.Error); ok {
			switch rpcErr.ErrorCode() {
			case ethengine.InvalidForkchoiceState:
				// if we succeed to update the forkchoice pre-payload, but fail post-payload, then it is a payload error
				return nil, &EngineError{
					Cat: BlockInsertPayloadErr,
					Err: fmt.Errorf("post-block-creation forkchoice update was inconsistent with engine, need reset to resolve: %w", rpcErr),
				}
			default:
				return nil, &EngineError{
					Cat: BlockInsertPrestateErr,
					Err: fmt.Errorf("unexpected error code in forkchoice-updated response: %w", err),
				}
			}
		} else {
			return nil, &EngineError{Cat: BlockInsertTemporaryErr, Err: fmt.Errorf("failed to make the new l2 block canonical via forkchoice: %w", err)}
		}
	}
	if fcRes.PayloadStatus.Status != ethengine.ExecutionValid {
		return nil, &EngineError{Cat: BlockInsertPayloadErr, Err: forkchoiceUpdateErr(fcRes.PayloadStatus)}
	}
	log.Info(
		"inserted block",
		"hash", payload.BlockHash,
		"number", uint64(payload.Number),
		"state_root", payload.StateRoot,
		"timestamp", uint64(payload.Timestamp),
		"parent", payload.ParentHash,
		// "prev_randao", payload.PrevRandao,
		"fee_recipient", payload.FeeRecipient,
		"txs", len(payload.Transactions),
		"update_safe", updateSafe,
	)
	return payload, nil
}

func forkchoiceUpdateErr(payloadStatus PayloadStatus) error {
	switch ethengine.ExecutePayloadStatus(payloadStatus.Status) {
	case ethengine.ExecutionValid:
		panic("function should not be called with valid status")
	case ethengine.ExecutionSyncing:
		return fmt.Errorf("updated forkchoice, but node is syncing")
	case ethengine.ExecutionAccepted, ethengine.ExecutionInvalidTerminalBlock, ethengine.ExecutionInvalidBlockHash:
		// ACCEPTED, INVALID_TERMINAL_BLOCK, INVALID_BLOCK_HASH are only for execution
		return fmt.Errorf("unexpected %s status, could not update forkchoice", payloadStatus.Status)
	case ethengine.ExecutionInvalid:
		return fmt.Errorf("cannot update forkchoice, block is invalid")
	default:
		return fmt.Errorf("unknown forkchoice status: %q", string(payloadStatus.Status))
	}
}

func newPayloadErr(payload *ExecutionPayload, payloadStatus *PayloadStatus) error {
	switch payloadStatus.Status {
	case ethengine.ExecutionValid:
		panic("function should not be called with valid status")
	case ethengine.ExecutionSyncing:
		return fmt.Errorf("failed to execute payload %s, node is syncing", payload.BlockID())
	case ethengine.ExecutionInvalid:
		return fmt.Errorf(
			"execution payload %s was INVALID! Latest valid hash is %s, ignoring bad block: %v",
			payload.BlockID(), payloadStatus.LatestValidHash, payloadStatus.ValidationError,
		)
	case ethengine.ExecutionInvalidBlockHash:
		return fmt.Errorf("execution payload %s has INVALID BLOCKHASH! %v", payload.BlockHash, payloadStatus.ValidationError)
	case ethengine.ExecutionInvalidTerminalBlock:
		return fmt.Errorf(
			"engine is misconfigured. Received invalid-terminal-block error while engine API should be active at genesis. err: %v",
			payloadStatus.ValidationError,
		)
	case ethengine.ExecutionAccepted:
		return fmt.Errorf("execution payload cannot be validated yet, latest valid hash is %s", payloadStatus.LatestValidHash)
	default:
		return fmt.Errorf("unknown execution status on %s: %q, ", payload.BlockID(), string(payloadStatus.Status))
	}
}
