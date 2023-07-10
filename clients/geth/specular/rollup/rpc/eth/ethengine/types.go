package ethengine

import (
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type (
	PayloadStatus = beacon.PayloadStatusV1
	PayloadID     = beacon.PayloadID
	// ForkChoice
	ForkChoiceResponse = beacon.ForkChoiceResponse
	ForkChoiceState    = beacon.ForkchoiceStateV1
	PayloadAttributes  = beacon.PayloadAttributesV1
	// GetPayload / NewPayload
	ExecutionPayload beacon.ExecutableDataV1
)

func (p *ExecutionPayload) BlockID() types.BlockID { return types.NewBlockID(p.Number, p.BlockHash) }

var EmptyPayloadID PayloadID = PayloadID{}

type ExecutePayloadStatus = string
type ErrorCode = int

var (
	// given payload is valid
	ExecutionValid ExecutePayloadStatus = ExecutePayloadStatus(beacon.VALID)
	// given payload is invalid
	ExecutionInvalid ExecutePayloadStatus = ExecutePayloadStatus(beacon.INVALID)
	// sync process is in progress
	ExecutionSyncing ExecutePayloadStatus = ExecutePayloadStatus(beacon.SYNCING)
	// returned if the payload is not fully validated, and does not extend the canonical chain,
	// but will be remembered for later (on reorgs or sync updates and such)
	ExecutionAccepted ExecutePayloadStatus = ExecutePayloadStatus(beacon.ACCEPTED)
	// if the block-hash in the payload is not correct
	ExecutionInvalidBlockHash ExecutePayloadStatus = ExecutePayloadStatus(beacon.INVALIDBLOCKHASH)
	// proof-of-stake transition only, not used in rollup
	ExecutionInvalidTerminalBlock ExecutePayloadStatus = "INVALID_TERMINAL_BLOCK"

	UnknownPayload           ErrorCode = beacon.UnknownPayload.ErrorCode()           // Payload does not exist / is not available.
	InvalidForkchoiceState   ErrorCode = beacon.InvalidForkChoiceState.ErrorCode()   // Forkchoice state is invalid / inconsistent.
	InvalidPayloadAttributes ErrorCode = beacon.InvalidPayloadAttributes.ErrorCode() // Payload attributes are invalid / inconsistent.
)
