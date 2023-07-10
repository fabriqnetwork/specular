package engine

import "github.com/ethereum/go-ethereum/common"

type BuildPayloadAttributes struct {
	timestamp             uint64
	random                common.Hash
	suggestedFeeRecipient common.Address
	txs                   [][]byte
	noTxPool              bool
}

func NewBuildPayloadAttributes(
	timestamp uint64,
	random common.Hash,
	suggestedFeeRecipient common.Address,
	txs [][]byte,
	noTxPool bool,
) *BuildPayloadAttributes {
	return &BuildPayloadAttributes{
		timestamp:             timestamp,
		random:                random,
		suggestedFeeRecipient: suggestedFeeRecipient,
		txs:                   txs,
		noTxPool:              noTxPool,
	}
}

func (a *BuildPayloadAttributes) Timestamp() uint64   { return a.timestamp }
func (a *BuildPayloadAttributes) Random() common.Hash { return a.random }
func (a *BuildPayloadAttributes) SuggestedFeeRecipient() common.Address {
	return a.suggestedFeeRecipient
}
func (a *BuildPayloadAttributes) Txs() [][]byte  { return a.txs }
func (a *BuildPayloadAttributes) NoTxPool() bool { return a.noTxPool }
