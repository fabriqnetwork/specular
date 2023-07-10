package api

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
)

// Defines interface between Specular services and the underlying execution backend.
// TODO: generalize to better support clients other than Geth.
type ExecutionBackend interface {
	ForkchoiceUpdate(ctx context.Context, update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(ctx context.Context, attrs BuildPayloadAttributes) (*BuildPayloadResponse, error)
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type BuildPayloadResponse = beacon.ForkChoiceResponse

type BuildPayloadAttributes interface {
	// Existing payload attributes
	Timestamp() uint64
	Random() common.Hash
	SuggestedFeeRecipient() common.Address
	// TODO: uncomment after upgrading geth
	// Withdrawals() []*types.Withdrawal
	// Attributes necessary for rollup functionality.
	Txs() [][]byte
	NoTxPool() bool
}
