package geth

import (
	"github.com/ethereum/go-ethereum/core/beacon"
)

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type PayloadStatus = beacon.PayloadStatusV1
type PayloadID = beacon.PayloadID

var InvalidForkChoiceState = beacon.InvalidForkChoiceState
var STATUS_INVALID = beacon.STATUS_INVALID

var (
	// VALID is returned by the engine API in the following calls:
	//   - forkchoiceUpdateV1: if the chain accepted the reorg (might ignore if it's stale)
	VALID = "VALID"

	// INVALID is returned by the engine API in the following calls:
	//   - forkchoiceUpdateV1: if the new head is unknown, pre-merge, or reorg to it fails
	INVALID = "INVALID"
)
