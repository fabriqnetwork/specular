package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

const methodNumBytes = 4

// Using singleton for convenience.
var serializationUtil *bridgeSerializationUtil

type bridgeSerializationUtil struct {
	inboxAbi     *abi.ABI
	challengeAbi *abi.ABI
}

func PackAppendTxBatchInput(contexts, txLengths []*big.Int, firstL2BlockNumber *big.Int, txs []byte) ([]byte, error) {
	if err := ensureUtilInit(); err != nil {
		return nil, fmt.Errorf("Failed to init serialization util: %w", err)
	}
	return serializationUtil.inboxAbi.Pack("appendTxBatch", contexts, txLengths, firstL2BlockNumber, txs)
}

func UnpackAppendTxBatchInput(tx *types.Transaction) ([]interface{}, error) {
	if err := ensureUtilInit(); err != nil {
		return nil, fmt.Errorf("Failed to init serialization util: %w", err)
	}
	return serializationUtil.inboxAbi.Methods["appendTxBatch"].Inputs.Unpack(tx.Data()[methodNumBytes:])
}

func UnpackBisectExecutionInput(tx *types.Transaction) ([]interface{}, error) {
	if err := ensureUtilInit(); err != nil {
		return nil, fmt.Errorf("Failed to init serialization util: %w", err)
	}
	return serializationUtil.challengeAbi.Methods["bisectExecution"].Inputs.Unpack(tx.Data()[methodNumBytes:])
}

func ensureUtilInit() error {
	if serializationUtil == nil {
		inboxAbi, err := bindings.ISequencerInboxMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("Failed to get ISequencerInbox ABI: %w", err)
		}
		challengeAbi, err := bindings.ISymChallengeMetaData.GetAbi()
		if err != nil {
			return fmt.Errorf("Failed to get ISymChallenge ABI: %w", err)
		}
		serializationUtil = &bridgeSerializationUtil{inboxAbi: inboxAbi, challengeAbi: challengeAbi}
	}
	return nil
}
