package bridge

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

const (
	TxBatchAppendedEventName = "TxBatchAppended"
	// ISequencerInbox.sol
	AppendTxBatchFnName = "appendTxBatch"
	// IRollup.sol
	StakeFnName                           = "stake"
	AdvanceStakeFnName                    = "advanceStake"
	CreateAssertionFnName                 = "createAssertion"
	ConfirmFirstUnresolvedAssertionFnName = "confirmFirstUnresolvedAssertion"
	RejectFirstUnresolvedAssertionFnName  = "rejectFirstUnresolvedAssertion"
	// IChallenge.sol
	bisectExecutionFn = "bisectExecution"

	MethodNumBytes = 4
)

// Using singleton for convenience.
var serializationUtil *bridgeSerializationUtil

type bridgeSerializationUtil struct {
	inboxAbi     *abi.ABI
	rollupAbi    *abi.ABI
	challengeAbi *abi.ABI
}

func InboxABIMethods() map[string]abi.Method     { return serializationUtil.inboxAbi.Methods }
func RollupABIMethods() map[string]abi.Method    { return serializationUtil.rollupAbi.Methods }
func ChallengeABIMethods() map[string]abi.Method { return serializationUtil.challengeAbi.Methods }
func TxMethodID(tx *types.Transaction) string    { return string(tx.Data()[:MethodNumBytes]) }

// ISequencerInbox.sol

func InboxEvent(name string) abi.Event { return serializationUtil.inboxAbi.Events[name] }

func UnpackAppendTxBatchInput(tx *types.Transaction) ([]any, error) {
	return serializationUtil.inboxAbi.Methods[AppendTxBatchFnName].Inputs.Unpack(tx.Data()[MethodNumBytes:])
}

func packAppendTxBatchInput(contexts, txLengths []*big.Int, firstL2BlockNumber *big.Int, txs []byte) ([]byte, error) {
	return serializationUtil.inboxAbi.Pack(AppendTxBatchFnName, contexts, txLengths, firstL2BlockNumber, txs)
}

// IRollup.sol

func UnpackCreateAssertionInput(tx *types.Transaction) (common.Hash, *big.Int, error) {
	in, err := serializationUtil.rollupAbi.Unpack(CreateAssertionFnName, tx.Data()[MethodNumBytes:])
	if err != nil {
		return common.Hash{}, nil, err
	}
	vmHash := in[0].(common.Hash)
	inboxSize := in[1].(*big.Int)
	return vmHash, inboxSize, err
}

func UnpackBisectExecutionInput(tx *types.Transaction) ([]any, error) {
	return serializationUtil.challengeAbi.Methods[bisectExecutionFn].Inputs.Unpack(tx.Data()[MethodNumBytes:])
}

func packStakeInput(stakeAmount *big.Int) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(StakeFnName, stakeAmount)
}

func packAdvanceStakeInput(assertionID *big.Int) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(AdvanceStakeFnName, assertionID)
}

func packCreateAssertionInput(vmHash common.Hash, inboxSize *big.Int) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(CreateAssertionFnName, vmHash, inboxSize)
}

func packConfirmFirstUnresolvedAssertionInput() ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(ConfirmFirstUnresolvedAssertionFnName)
}

func packRejectFirstUnresolvedAssertionInput(stakerAddress common.Address) ([]byte, error) {
	return serializationUtil.rollupAbi.Pack(RejectFirstUnresolvedAssertionFnName, stakerAddress)
}

// Ensures serializationUtil is initialized. Must be called prior to the methods above.
func ensureUtilInit() error {
	if serializationUtil == nil {
		inboxAbi, err := bindings.ISequencerInboxMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("Failed to get ISequencerInbox ABI: %w", err)
		}
		rollupAbi, err := bindings.IRollupMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("Failed to get IRollup ABI: %w", err)
		}
		challengeAbi, err := bindings.ISymChallengeMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("Failed to get ISymChallenge ABI: %w", err)
		}
		serializationUtil = &bridgeSerializationUtil{
			inboxAbi:     inboxAbi,
			rollupAbi:    rollupAbi,
			challengeAbi: challengeAbi,
		}
	}
	return nil
}
