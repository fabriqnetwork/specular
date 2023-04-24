package state

// import (
// 	"context"

// 	"github.com/ethereum/go-ethereum/core"
// 	"github.com/ethereum/go-ethereum/core/beacon"
// 	"github.com/ethereum/go-ethereum/core/types"
// )

// type FauxEngine struct {
// 	blockchain *core.BlockChain
// }

// type ForkChoiceState interface {
// 	Head() *types.Header
// 	Safe() *types.Header
// 	Finalized() *types.Header
// }

// // TODO: update Geth to use new Engine types
// type ForkchoiceState beacon.ForkchoiceStateV1
// type ForkChoiceResponse beacon.ForkChoiceResponse
// type PayloadAttributes beacon.PayloadAttributesV1
// type PayloadID beacon.PayloadID
// type ExecutionPayloadEnvelope beacon.ExecutableDataV1
// type ExecutableData beacon.ExecutableDataV1
// type PayloadStatus beacon.PayloadStatusV1

// func (e *FauxEngine) ForkChoiceUpdated(
// 	ctx context.Context,
// 	update ForkchoiceState,
// 	payloadAttributes *PayloadAttributes,
// ) (ForkChoiceResponse, error) {
// 	return ForkChoiceResponse{}, nil
// }

// func (e *FauxEngine) GetPayload(ctx context.Context, payloadID PayloadID) (ExecutionPayloadEnvelope, error) {
// 	return ExecutionPayloadEnvelope{}, nil
// }

// func (e *FauxEngine) NewPayload(ctx context.Context, params ExecutableData) (PayloadStatus, error) {
// 	return PayloadStatus{}, nil
// }
