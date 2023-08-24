package rollup

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// TODO: maybe move these constants to the rollup config or a separate file
const (
	txDataZero          = 4
	txDataOne           = 16
	txSignatureOverhead = 68
)

type RollupConfig interface {
	GetCoinbase() common.Address         // recipient of the L1 Fee
	GetL2ChainID() uint64                // chain ID of the specular rollup
	GetL1FeeOverhead() int64             // fixed cost of submitting a tx to L1
	GetL1FeeMultiplier() float64         // value to scale the L1 Fee
	GetL1OracleAddress() common.Address  // contract providing the L1 basefee
	GetL1OracleBaseFeeSlot() common.Hash // L1 basefee storage slot
}

// MakeSpecularEVMPreTransferHook creates specular's vm.EVMHook function
// which is injected into the EVM and runs before every transfer
// currently this is only used to calculate & charge the L1 Fee
func MakeSpecularEVMPreTransferHook(cfg RollupConfig) vm.EVMHook {
	log.Info("Injected Specular EVM hook")
	log.Info("L1Oracle config", "address", cfg.GetL1OracleAddress(), "baseFeeSlot", cfg.GetL1OracleBaseFeeSlot())
	return func(msg types.Message, evm *vm.EVM) error {
		tx := transactionFromMessage(msg, cfg)
		fee, err := calculateL1Fee(tx, evm, cfg)
		if err != nil {
			return err
		}
		return chargeL1Fee(fee, msg, evm, cfg)
	}
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
	if err := tx.EncodeRLP(buf); err != nil {
		return common.Big0, err
	}
	bytes := buf.Bytes()
	// remove the last 3 bytes containing the signature
	// this mirrors the optimism implementation [1]
	// but contradicts the optimism spec [2]
	// [1] https://github.com/ethereum-optimism/optimism/blob/5d9a38dcd9dc79dce41a6d08f9b28ff850f77811/l2geth/rollup/fees/rollup_fee.go#L204
	// [2] https://github.com/ethereum-optimism/optimism/blob/develop/specs/exec-engine.md#l1-cost-fees-l1-fee-vault
	rlp := bytes[:len(bytes)-3]

	var (
		zeroes, ones = zeroesAndOnes(rlp)

		txDataGas    = big.NewInt(zeroes*txDataZero + (ones+txSignatureOverhead)*txDataOne + cfg.GetL1FeeOverhead())
		basefee      = readStorageSlot(evm, cfg.GetL1OracleAddress(), cfg.GetL1OracleBaseFeeSlot())
		feeMutiplier = cfg.GetL1FeeMultiplier()

		l1Fee       = new(big.Int).Mul(txDataGas, basefee)
		scaledL1Fee = ScaleBigInt(l1Fee, feeMutiplier)
	)

	log.Trace(
		"calculated l1 fee",
		"txDataGas", txDataGas,
		"basefee", basefee,
		"l1Fee", l1Fee,
		"feeMutiplier", feeMutiplier,
		"scaledL1Fee", scaledL1Fee,
	)

	return scaledL1Fee, nil
}

// multiply a big.Int with a float
// only the first 3 decimal places of the scalar are used to guarantee precision
func ScaleBigInt(num *big.Int, scalar float64) *big.Int {
	var (
		f    = new(big.Float).SetInt(num)
		s, _ = new(big.Float).SetString(fmt.Sprintf("%.3f", scalar))

		scaledNum     = new(big.Float).Mul(f, s)
		roundedNum, _ = scaledNum.Int(nil)
	)

	if !scaledNum.IsInt() {
		roundedNum = roundedNum.Add(roundedNum, common.Big1)
	}

	return roundedNum
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
	var zeroes, ones int64
	for _, b := range data {
		if b == 0 {
			zeroes++
		} else {
			ones++
		}
	}
	return zeroes, ones
}
