package engine

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type PayloadAttrsBuilder struct {
	SuggestedFeeRecipient common.Address
}

func (b *PayloadAttrsBuilder) BuildForForkchoiceUpdate() *PayloadAttributes {
	return &PayloadAttributes{
		Timestamp:             uint64(time.Now().Unix()), // timestamp (now)
		Random:                common.Hash{},             // randomness (none)
		SuggestedFeeRecipient: b.SuggestedFeeRecipient,   // coinbase
	}
}

func (b *PayloadAttrsBuilder) BuildForBuildPayload() *BuildPayloadAttributes {
	return NewBuildPayloadAttributes(
		uint64(time.Now().Unix()), // timestamp (now)
		common.Hash{},             // randomness (none)
		b.SuggestedFeeRecipient,   // coinbase
		nil,                       // no txs to force-include
		false,                     // use tx pool
	)
}
