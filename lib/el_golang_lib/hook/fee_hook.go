package hook

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularL2/specular/ops/bindings/bindings"
)

// TODO: maybe move these constants to the rollup config or a separate file
const (
	txDataZero          = 4
	txDataOne           = 16
	txSignatureOverhead = 68
)

var L1OracleAddress = common.HexToAddress("0x2A00000000000000000000000000000000000010")

type RollupConfig interface {
	GetL1FeeRecipient() common.Address // recipient of the L1 Fee
	GetL2ChainID() uint64              // chain ID of the specular rollup
}

// MakeSpecularEVMPreTransferHook creates specular's vm.EVMHook function
// which is injected into the EVM and runs before every transfer
// currently this is only used to calculate & charge the L1 Fee
func MakeSpecularEVMPreTransferHook(l2ChainId uint64, l1FeeRecipient common.Address) vm.EVMHook {
	log.Info("Injected Specular EVM hook")

	feeStorageSlots := getStorageSlots()
	log.Info("L1Oracle config", "address", L1OracleAddress, "overhead", feeStorageSlots.overheadSlot, "baseFeeSlot", feeStorageSlots.baseFeeSlot, "scalarSlot", feeStorageSlots.scalarSlot)

	return func(msg vm.MessageInterface, db vm.StateDB) error {
		tx := transactionFromMessage(msg, l2ChainId)
		fee, err := calculateL1Fee(tx, db)
		if err != nil {
			return err
		}

		return chargeL1Fee(fee, msg, db, l1FeeRecipient)
	}
}

// MakeSpecularL1FeeReader creates specular's vm.EVMReader function
// which is injected into the EVM and can be used to return the L1Fee of a transaction.
// This is a read only method and does not change the state.
func MakeSpecularL1FeeReader(l2ChainId uint64) vm.EVMReader {
	log.Info("Injected Specular EVM reader")
	log.Info("L1Oracle config", "address", L1OracleAddress)
	return func(tx *types.Transaction, db vm.StateDB) (*big.Int, error) {
		return calculateL1Fee(tx, db)
	}
}

// creates a Transaction from a transaction
// the Tx Type is reconstructed and the signature is left empty
func transactionFromMessage(msg vm.MessageInterface, l2ChainId uint64) *types.Transaction {
	var txData types.TxData
	if msg.GetGasTipCap() != nil {
		txData = &types.DynamicFeeTx{
			ChainID:    new(big.Int).SetUint64(l2ChainId),
			Nonce:      msg.GetNonce(),
			GasTipCap:  msg.GetGasTipCap(),
			GasFeeCap:  msg.GetGasFeeCap(),
			Gas:        msg.GetGasLimit(),
			To:         msg.GetTo(),
			Value:      msg.GetValue(),
			Data:       msg.GetData(),
			AccessList: msg.GetAccessList(),
		}
	} else if msg.GetAccessList() != nil {
		txData = &types.AccessListTx{
			ChainID:    new(big.Int).SetUint64(l2ChainId),
			Nonce:      msg.GetNonce(),
			GasPrice:   msg.GetGasPrice(),
			Gas:        msg.GetGasLimit(),
			To:         msg.GetTo(),
			Value:      msg.GetValue(),
			Data:       msg.GetData(),
			AccessList: msg.GetAccessList(),
		}
	} else {
		txData = &types.LegacyTx{
			Nonce:    msg.GetNonce(),
			GasPrice: msg.GetGasPrice(),
			Gas:      msg.GetGasLimit(),
			To:       msg.GetTo(),
			Value:    msg.GetValue(),
			Data:     msg.GetData(),
		}
	}
	return types.NewTx(txData)
}

// calculates the L1 fee for a transaction using the following formula:
// L1Fee = L1FeeMultiplier * L1BaseFee * (TxDataGas + L1OverheadGas)
// L1BaseFee is dynamically set by the L1Oracle
// L1FeeMultiplier & L1OverheadGas are set in the rollup configuration
func calculateL1Fee(tx *types.Transaction, db vm.StateDB) (*big.Int, error) {
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
		zeroes, ones  = zeroesAndOnes(rlp)
		rollupDataGas = zeroes*txDataZero + (ones+txSignatureOverhead)*txDataOne

		feeStorageSlots = getStorageSlots()

		overhead = readStorageSlot(db, L1OracleAddress, feeStorageSlots.overheadSlot)
		basefee  = readStorageSlot(db, L1OracleAddress, feeStorageSlots.baseFeeSlot)
		scalar   = readStorageSlot(db, L1OracleAddress, feeStorageSlots.scalarSlot)
	)

	log.Info(
		"calculated l1 fee",
		"rollupDataGas", rollupDataGas,
		"overhead", overhead,
		"basefee", basefee,
		"scalar", scalar,
	)

	l1GasUsed := new(big.Int).SetUint64(rollupDataGas)
	l1GasUsed = l1GasUsed.Add(l1GasUsed, overhead)
	l1Cost := l1GasUsed.Mul(l1GasUsed, basefee)
	l1Cost = l1Cost.Mul(l1Cost, scalar)
	return l1Cost.Div(l1Cost, big.NewInt(1_000_000)), nil
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
func chargeL1Fee(l1Fee *big.Int, msg vm.MessageInterface, db vm.StateDB, l1FeeRecipient common.Address) error {
	senderBalance := db.GetBalance(msg.GetFrom())

	if senderBalance.Cmp(l1Fee) < 0 {
		return errors.New("insufficient balance to cover L1 fee")
	}

	db.AddBalance(l1FeeRecipient, l1Fee)
	db.SubBalance(msg.GetFrom(), l1Fee)

	log.Info("charged L1 Fee", "fee", l1Fee.Uint64())
	return nil
}

// read the value from a given address / storage slot
func readStorageSlot(db vm.StateDB, address common.Address, slot common.Hash) *big.Int {
	return db.GetState(address, slot).Big()
}

// zeroesAndOnes counts the number of 0 bytes and non 0 bytes in a byte slice
func zeroesAndOnes(data []byte) (uint64, uint64) {
	var zeroes, ones uint64
	for _, b := range data {
		if b == 0 {
			zeroes++
		} else {
			ones++
		}
	}
	return zeroes, ones
}

type feeStorageSlots struct {
	baseFeeSlot  common.Hash
	overheadSlot common.Hash
	scalarSlot   common.Hash
}

func getStorageSlots() feeStorageSlots {
	layout, err := bindings.GetStorageLayout("L1Oracle")
	if err != nil {
		panic("could not get storage layout for L1Oracle")
	}

	baseFeeEntry, err := layout.GetStorageLayoutEntry("baseFee")
	if err != nil {
		panic("could not get basefee storage slot")
	}
	overheadEntry, err := layout.GetStorageLayoutEntry("l1FeeOverhead")
	if err != nil {
		panic("could not get overhead storage slot")
	}
	scalarEntry, err := layout.GetStorageLayoutEntry("l1FeeScalar")
	if err != nil {
		panic("could not get scalar storage slot")
	}

	return feeStorageSlots{
		baseFeeSlot:  common.BigToHash(new(big.Int).SetUint64(uint64(baseFeeEntry.Slot))),
		overheadSlot: common.BigToHash(new(big.Int).SetUint64(uint64(overheadEntry.Slot))),
		scalarSlot:   common.BigToHash(new(big.Int).SetUint64(uint64(scalarEntry.Slot))),
	}
}
