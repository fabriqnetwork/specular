package ethengine

import (
	"context"

	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// Engine client for vanilla Ethereum
type Client struct{ rpcClient }

type rpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

func NewEngineClient(c rpcClient) *Client { return &Client{c} }

func (c *Client) ForkchoiceUpdate(
	ctx context.Context,
	forkchoice *ForkChoiceState,
	attributes *PayloadAttributes,
) (*ForkChoiceResponse, error) {
	var result ForkChoiceResponse
	log.Trace("Sharing forkchoice-updated signal")
	if err := c.CallContext(ctx, &result, "engine_forkchoiceUpdatedV1", forkchoice, attributes); err != nil {
		return nil, fmt.Errorf("failed to update fork-choice: %w", err)
	}
	return &result, nil
}

// NewPayload executes a full block on the execution engine.
func (c *Client) NewPayload(ctx context.Context, payload *ExecutionPayload) (*PayloadStatus, error) {
	log.Trace("Sending payload for execution")
	var (
		result PayloadStatus
		err    = c.CallContext(ctx, &result, "engine_newPayloadV1", payload)
	)
	log.Trace(
		"Received payload execution result",
		"status", result.Status,
		"latestValidHash", result.LatestValidHash,
		"message", result.ValidationError,
	)
	return &result, err
}

// GetPayload gets the execution payload associated with the PayloadId.
func (c *Client) GetPayload(ctx context.Context, payloadId PayloadID) (*ExecutionPayload, error) {
	log.Trace("Getting payload")
	var (
		result ExecutionPayload
		err    = c.CallContext(ctx, &result, "engine_getPayloadV1", payloadId)
	)
	if err != nil {
		log.Warn("Failed to get payload", "payload_id", payloadId, "err", err)
	} else {
		log.Trace("Received payload")
	}
	return &result, err
}
