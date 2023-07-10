package l2rpc

import (
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/ethengine"
)

// TODO: Use EngineClient in place of ExecutionBackend
// TODO: upgrade Geth to use new Engine types

type EngineClient struct{ ethengine.Client }

type L2Config interface{ Endpoint() string }

type (
	PayloadStatus = ethengine.PayloadStatus
	// ForkChoice
	ForkChoiceResponse = ethengine.ForkChoiceResponse
	ForkChoiceState    = ethengine.ForkChoiceState
	PayloadAttributes  = ethengine.PayloadAttributes
	// NewPayload and GetPayload
	ExecutionPayload = ethengine.ExecutionPayload
)

func NewEngineClient(c ethengine.Client) *EngineClient { return &EngineClient{c} }

// func DialWithRetry(ctx context.Context, cfg L2Config) (*EngineClient, error) {
// 	l2Client, err := eth.DialWithRetry(ctx, cfg.Endpoint(), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to dial L2 client: %v", err)
// 	}
// 	return &EngineClient{l2Client}, nil
// }

// func (c *EngineClient) ForkchoiceUpdate(
// 	ctx context.Context,
// 	forkChoice *ForkchoiceState,
// 	attributes *PayloadAttributes,
// ) (*ForkChoiceResponse, error) {
// 	return c.EngineClient.ForkchoiceUpdate(ctx, forkChoice, attributes)
// }

// func (c *EngineClient) NewPayload(ctx context.Context, payload *ExecutionPayload) (*PayloadStatus, error) {
// }
