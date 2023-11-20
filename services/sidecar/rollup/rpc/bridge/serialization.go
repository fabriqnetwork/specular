package bridge

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/bindings"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
)

const (
	TxBatchAppendedEventName = "TxBatchAppended"
	// ISequencerInbox.sol functions
	AppendTxBatchFnName = "appendTxBatch"
	// IRollup.sol functions
	StakeFnName                           = "stake"
	AdvanceStakeFnName                    = "advanceStake"
	CreateAssertionFnName                 = "createAssertion"
	ConfirmFirstUnresolvedAssertionFnName = "confirmFirstUnresolvedAssertion"
	RejectFirstUnresolvedAssertionFnName  = "rejectFirstUnresolvedAssertion"
	// IChallenge.sol functions
	bisectExecutionFn = "bisectExecution"
	// IRollup.sol errors (TODO: figure out a work-around to hardcoding)
	NoUnresolvedAssertionErr     = "Error: VM Exception while processing transaction: reverted with custom error 'NoUnresolvedAssertion()'"
	ConfirmationPeriodPendingErr = "Error: VM Exception while processing transaction: reverted with custom error 'ConfirmationPeriodPending()'"
	// L1Oracle.sol functions
	SetL1OracleValues = "setL1OracleValues"

	MethodNumBytes = 4
)

// Using singleton for convenience.
var serializationUtil *bridgeSerializationUtil

type bridgeSerializationUtil struct {
	inboxAbi     *abi.ABI
	rollupAbi    *abi.ABI
	challengeAbi *abi.ABI
	l1OracleAbi  *abi.ABI
}

// ISequencerInbox.sol

func InboxEvent(name string) abi.Event { return serializationUtil.inboxAbi.Events[name] }

func UnpackAppendTxBatchInput(tx *types.Transaction) ([]any, error) {
	return serializationUtil.inboxAbi.Methods[AppendTxBatchFnName].Inputs.Unpack(tx.Data()[MethodNumBytes:])
}

func packAppendTxBatchInput(batch []byte) ([]byte, error) {
	return serializationUtil.inboxAbi.Pack(AppendTxBatchFnName, batch)
}

// IRollup.sol

func UnpackCreateAssertionInput(tx *types.Transaction) (common.Hash, *big.Int, error) {
	in, err := serializationUtil.rollupAbi.Methods[CreateAssertionFnName].Inputs.Unpack(tx.Data()[MethodNumBytes:])
	if err != nil {
		return common.Hash{}, nil, err
	}
	var (
		vmHash    = in[0].(common.Hash)
		inboxSize = in[1].(*big.Int)
	)
	return vmHash, inboxSize, err
}

func packStakeInput() ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(StakeFnName)
}

func packAdvanceStakeInput(assertionID *big.Int) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(AdvanceStakeFnName, assertionID)
}

func packCreateAssertionInput(vmHash common.Hash, blockNum *big.Int) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(CreateAssertionFnName, vmHash, blockNum)
}

func packConfirmFirstUnresolvedAssertionInput() ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(ConfirmFirstUnresolvedAssertionFnName)
}

func packRejectFirstUnresolvedAssertionInput(stakerAddress common.Address) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(RejectFirstUnresolvedAssertionFnName, stakerAddress)
}

// L1Oracle.sol

func UnpackL1OracleInput(tx *types.Transaction) (uint64, uint64, uint64, common.Hash, common.Hash, error) {
	in, err := serializationUtil.l1OracleAbi.Methods[SetL1OracleValues].Inputs.Unpack(tx.Data()[MethodNumBytes:])
	if err != nil {
		return 0, 0, 0, common.Hash{}, common.Hash{}, err
	}
	var (
		number       = in[0].(*big.Int).Uint64()
		timestamp    = in[1].(*big.Int).Uint64()
		baseFee      = in[2].(*big.Int).Uint64()
		hashRaw      = in[3].([32]byte)
		hash         = common.BytesToHash(hashRaw[:])
		stateRootRaw = in[4].([32]byte)
		stateRoot    = common.BytesToHash(stateRootRaw[:])
	)
	return number, timestamp, baseFee, hash, stateRoot, nil
}

// Ensures serializationUtil is initialized. Must be called prior to the methods above.
func ensureUtilInit() error {
	if serializationUtil == nil {
		inboxAbi, err := bindings.ISequencerInboxMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("failed to get ISequencerInbox ABI: %w", err)
		}
		rollupAbi, err := bindings.IRollupMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("failed to get IRollup ABI: %w", err)
		}
		challengeAbi, err := bindings.ISymChallengeMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("failed to get ISymChallenge ABI: %w", err)
		}
		l1OracleAbi, err := bindings.L1OracleMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("failed to get IL1Oracle ABI: %w", err)
		}
		serializationUtil = &bridgeSerializationUtil{
			inboxAbi:     inboxAbi,
			rollupAbi:    rollupAbi,
			challengeAbi: challengeAbi,
			l1OracleAbi:  l1OracleAbi,
		}
	}
	return nil
}
