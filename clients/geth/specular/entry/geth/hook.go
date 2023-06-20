package entry

import (
	"bytes"
	"errors"
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
	GetCoinbase() common.Address        // recipient of the L1 Fee
	GetL2ChainID() uint64               // chain ID of the specular rollup
	GetL1FeeOverhead() int64            // fixed cost of submitting a tx to L1
	GetL1FeeMultiplier() float64        // value to scale the L1 Fee
	GetL1OracleAddress() common.Address // contract providing the L1 basefee
}

// SpecularEVMPreTransferHook is injected into the EVM and runs before every transfer
// currently this is only used to calculate & charge the L1 Fee
func SpecularEVMPreTransferHook(msg types.Message, evm *vm.EVM, cfg RollupConfig) error {
	tx := transactionFromMessage(msg, cfg)
	fee, err := calculateL1Fee(tx, evm, cfg)
	if err != nil {
		return err
	}
	return chargeL1Fee(fee, msg, evm, cfg)
}

// creates a Transaction from a transaction
// the Tx Type is reconstructed and the signature is left empty
func transactionFromMessage(msg types.Message, cfg RollupConfig) *types.Transaction {
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

// calculates the L1 fee for a transaction using the following formula:
// L1Fee = L1FeeMultiplier * L1BaseFee * (TxDataGas + L1OverheadGas)
// L1BaseFee is dynamically set by the L1Oracle
// L1FeeMultiplier & L1OverheadGas are set in the rollup configuration
func calculateL1Fee(tx *types.Transaction, evm *vm.EVM, cfg RollupConfig) (*big.Int, error) {
	// calculate L1 gas from RLP encoding
	buf := new(bytes.Buffer)
	err := tx.EncodeRLP(buf)
	if err != nil {
		return common.Big0, err
	}
	bytes := buf.Bytes()
	// remove the last 3 bytes containing the signature
	// this mirrors the optimism implementation [1]
	// but contradicts the optimism spec [2]
	// [1] https://github.com/ethereum-optimism/optimism/blob/5d9a38dcd9dc79dce41a6d08f9b28ff850f77811/l2geth/rollup/fees/rollup_fee.go#L204
	// [2] https://github.com/ethereum-optimism/optimism/blob/develop/specs/exec-engine.md#l1-cost-fees-l1-fee-vault
	rlp := bytes[:len(bytes)-3]

	zeroes, ones := zeroesAndOnes(rlp)

	txDataGas := big.NewInt(zeroes*TX_DATA_ZERO + ones*TX_DATA_ONE + cfg.GetL1FeeOverhead())
	basefee := readStorageSlot(evm, cfg.GetL1OracleAddress(), common.HexToHash(BASEFEE_SLOT))

	l1Fee := new(big.Float).SetInt(new(big.Int).Mul(txDataGas, basefee))
	feeMutiplier := new(big.Float).SetFloat64(cfg.GetL1FeeMultiplier())

	scaledL1Fee := new(big.Float).Mul(l1Fee, feeMutiplier)
	roundedL1Fee, _ := new(big.Float).Add(scaledL1Fee, big.NewFloat(0.5)).Int(nil)

	log.Trace(
		"calculated l1 fee",
		"txDataGas", txDataGas,
		"basefee", basefee,
		"l1Fee", l1Fee,
		"feeMutiplier", feeMutiplier,
		"scaledL1Fee", scaledL1Fee,
		"roundedL1Fee", roundedL1Fee,
	)

	return roundedL1Fee, nil
}

// subtract the L1 Fee from the sender of the Tx
// add the Fee to the balance of the coinbase address
func chargeL1Fee(l1Fee *big.Int, msg types.Message, evm *vm.EVM, cfg RollupConfig) error {
	senderBalance := evm.StateDB.GetBalance(msg.From())
	if senderBalance.Cmp(l1Fee) < 0 {
		return errors.New("insufficient balance to cover L1 fee")
	}

	evm.StateDB.AddBalance(cfg.GetCoinbase(), l1Fee)
	evm.StateDB.SubBalance(msg.From(), l1Fee)

	log.Info("charged L1 Fee", "fee", l1Fee.Uint64())
	return nil
}

// read the value from a given address / storage slot
func readStorageSlot(evm *vm.EVM, address common.Address, slot common.Hash) *big.Int {
	return evm.StateDB.GetState(address, slot).Big()
}

// zeroesAndOnes counts the number of 0 bytes and non 0 bytes in a byte slice
func zeroesAndOnes(data []byte) (int64, int64) {
	var zeroes int64
	var ones int64
	for _, b := range data {
		if b == 0 {
			zeroes++
		} else {
			ones++
		}
	}
	return zeroes, ones
}
