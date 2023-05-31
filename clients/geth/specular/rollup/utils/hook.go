package utils

import (
	"bytes"
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

const TX_DATA_ZERO = 4
const TX_DATA_ONE = 16

// keccak256("specular.basefee")
const BASEFEE_SLOT = "0x18b94da8c18f49ac05520153402a0591c3c917271b9d13711fd6fdb213ded168"

type RollupConfig interface {
	GetCoinbase() common.Address
	GetL2ChainID() uint64
	GetL1FeeOverhead() uint64
	GetL1FeeMultiplier() float64
	GetL1OracleAddress() common.Address
}

func SpecularEVMHook(msg types.Message, evm *vm.EVM, cfg RollupConfig) error {
	tx := createTransaction(msg, cfg)
	fee, err := calculateL1Fee(tx, evm, cfg)
	if err != nil {
		return err
	}
	return chargeL1Fee(fee, msg, evm, cfg)
}

func createTransaction(msg types.Message, cfg RollupConfig) *types.Transaction {
	var txData types.TxData
	if msg.GasTipCap() != nil {
		txData = &types.DynamicFeeTx{
			ChainID:    big.NewInt(int64(cfg.GetL2ChainID())),
			Nonce:      msg.Nonce(),
			GasTipCap:  msg.GasTipCap(),
			GasFeeCap:  msg.GasFeeCap(),
			Gas:        msg.Gas(),
			To:         msg.To(),
			Value:      msg.Value(),
			Data:       msg.Data(),
			AccessList: msg.AccessList(),
		}
	} else if msg.AccessList() != nil {
		txData = &types.AccessListTx{
			ChainID:    big.NewInt(int64(cfg.GetL2ChainID())),
			Nonce:      msg.Nonce(),
			GasPrice:   msg.GasPrice(),
			Gas:        msg.Gas(),
			To:         msg.To(),
			Value:      msg.Value(),
			Data:       msg.Data(),
			AccessList: msg.AccessList(),
		}
	} else {
		txData = &types.LegacyTx{
			Nonce:    msg.Nonce(),
			GasPrice: msg.GasPrice(),
			Gas:      msg.Gas(),
			To:       msg.To(),
			Value:    msg.Value(),
			Data:     msg.Data(),
		}
	}
	return types.NewTx(txData)
}

func calculateL1Fee(tx *types.Transaction, evm *vm.EVM, cfg RollupConfig) (uint64, error) {
	// calculate L1 gas from RLP encoding
	buf := new(bytes.Buffer)
	err := tx.EncodeRLP(buf)
	if err != nil {
		return 0, err
	}
	bytes := buf.Bytes()
	// remove the last 3 bytes containing the signature
	// this mirrors the optimism implementation [1]
	// but contradicts the optimism spec [2]
	// [1] https://github.com/ethereum-optimism/optimism/blob/5d9a38dcd9dc79dce41a6d08f9b28ff850f77811/l2geth/rollup/fees/rollup_fee.go#L204
	// [2] https://github.com/ethereum-optimism/optimism/blob/develop/specs/exec-engine.md#l1-cost-fees-l1-fee-vault
	rlp := bytes[:len(bytes)-3]

	zeroes, ones := zeroesAndOnes(rlp)
	txDataGas := zeroes*TX_DATA_ZERO + ones*TX_DATA_ONE + cfg.GetL1FeeOverhead()
	basefee := readBasefee(evm, cfg)

	l1Fee := txDataGas * basefee
	scaledL1Fee := uint64(math.Ceil(float64(l1Fee) * cfg.GetL1FeeMultiplier()))

	return scaledL1Fee, nil
}

func chargeL1Fee(L1Fee uint64, msg types.Message, evm *vm.EVM, cfg RollupConfig) error {
	senderBalance := evm.StateDB.GetBalance(msg.From())
	if senderBalance.Uint64() < L1Fee {
		return errors.New("insufficient balance to cover L1 fee")
	}

	evm.StateDB.AddBalance(cfg.GetCoinbase(), big.NewInt(int64(L1Fee)))
	evm.StateDB.SubBalance(msg.From(), big.NewInt(int64(L1Fee)))

	log.Info("charged L1 Fee", "fee", L1Fee)
	return nil
}

func readBasefee(evm *vm.EVM, cfg RollupConfig) uint64 {
	basefee := evm.StateDB.GetState(cfg.GetL1OracleAddress(), common.HexToHash(BASEFEE_SLOT))
	log.Info("read basefee", "fee", basefee.Big().Uint64(), "oracle-addr", cfg.GetL1OracleAddress().String())
	return basefee.Big().Uint64()
}

// zeroesAndOnes counts the number of 0 bytes and non 0 bytes in a byte slice
func zeroesAndOnes(data []byte) (uint64, uint64) {
	var zeroes uint64
	var ones uint64
	for _, b := range data {
		if b == 0 {
			zeroes++
		} else {
			ones++
		}
	}
	return zeroes, ones
}
