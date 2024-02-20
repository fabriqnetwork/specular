// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// L2PortalDeterministicStorageMetaData contains all meta data concerning the L2PortalDeterministicStorage contract.
var L2PortalDeterministicStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"initiatedWithdrawals\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// L2PortalDeterministicStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use L2PortalDeterministicStorageMetaData.ABI instead.
var L2PortalDeterministicStorageABI = L2PortalDeterministicStorageMetaData.ABI

// L2PortalDeterministicStorage is an auto generated Go binding around an Ethereum contract.
type L2PortalDeterministicStorage struct {
	L2PortalDeterministicStorageCaller     // Read-only binding to the contract
	L2PortalDeterministicStorageTransactor // Write-only binding to the contract
	L2PortalDeterministicStorageFilterer   // Log filterer for contract events
}

// L2PortalDeterministicStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2PortalDeterministicStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2PortalDeterministicStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2PortalDeterministicStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2PortalDeterministicStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2PortalDeterministicStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2PortalDeterministicStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2PortalDeterministicStorageSession struct {
	Contract     *L2PortalDeterministicStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// L2PortalDeterministicStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2PortalDeterministicStorageCallerSession struct {
	Contract *L2PortalDeterministicStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// L2PortalDeterministicStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2PortalDeterministicStorageTransactorSession struct {
	Contract     *L2PortalDeterministicStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// L2PortalDeterministicStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2PortalDeterministicStorageRaw struct {
	Contract *L2PortalDeterministicStorage // Generic contract binding to access the raw methods on
}

// L2PortalDeterministicStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2PortalDeterministicStorageCallerRaw struct {
	Contract *L2PortalDeterministicStorageCaller // Generic read-only contract binding to access the raw methods on
}

// L2PortalDeterministicStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2PortalDeterministicStorageTransactorRaw struct {
	Contract *L2PortalDeterministicStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2PortalDeterministicStorage creates a new instance of L2PortalDeterministicStorage, bound to a specific deployed contract.
func NewL2PortalDeterministicStorage(address common.Address, backend bind.ContractBackend) (*L2PortalDeterministicStorage, error) {
	contract, err := bindL2PortalDeterministicStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2PortalDeterministicStorage{L2PortalDeterministicStorageCaller: L2PortalDeterministicStorageCaller{contract: contract}, L2PortalDeterministicStorageTransactor: L2PortalDeterministicStorageTransactor{contract: contract}, L2PortalDeterministicStorageFilterer: L2PortalDeterministicStorageFilterer{contract: contract}}, nil
}

// NewL2PortalDeterministicStorageCaller creates a new read-only instance of L2PortalDeterministicStorage, bound to a specific deployed contract.
func NewL2PortalDeterministicStorageCaller(address common.Address, caller bind.ContractCaller) (*L2PortalDeterministicStorageCaller, error) {
	contract, err := bindL2PortalDeterministicStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2PortalDeterministicStorageCaller{contract: contract}, nil
}

// NewL2PortalDeterministicStorageTransactor creates a new write-only instance of L2PortalDeterministicStorage, bound to a specific deployed contract.
func NewL2PortalDeterministicStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*L2PortalDeterministicStorageTransactor, error) {
	contract, err := bindL2PortalDeterministicStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2PortalDeterministicStorageTransactor{contract: contract}, nil
}

// NewL2PortalDeterministicStorageFilterer creates a new log filterer instance of L2PortalDeterministicStorage, bound to a specific deployed contract.
func NewL2PortalDeterministicStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*L2PortalDeterministicStorageFilterer, error) {
	contract, err := bindL2PortalDeterministicStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2PortalDeterministicStorageFilterer{contract: contract}, nil
}

// bindL2PortalDeterministicStorage binds a generic wrapper to an already deployed contract.
func bindL2PortalDeterministicStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2PortalDeterministicStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2PortalDeterministicStorage.Contract.L2PortalDeterministicStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2PortalDeterministicStorage.Contract.L2PortalDeterministicStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2PortalDeterministicStorage.Contract.L2PortalDeterministicStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2PortalDeterministicStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2PortalDeterministicStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2PortalDeterministicStorage.Contract.contract.Transact(opts, method, params...)
}

// InitiatedWithdrawals is a free data retrieval call binding the contract method 0xbf286e39.
//
// Solidity: function initiatedWithdrawals(bytes32 ) view returns(bool)
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageCaller) InitiatedWithdrawals(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _L2PortalDeterministicStorage.contract.Call(opts, &out, "initiatedWithdrawals", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InitiatedWithdrawals is a free data retrieval call binding the contract method 0xbf286e39.
//
// Solidity: function initiatedWithdrawals(bytes32 ) view returns(bool)
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageSession) InitiatedWithdrawals(arg0 [32]byte) (bool, error) {
	return _L2PortalDeterministicStorage.Contract.InitiatedWithdrawals(&_L2PortalDeterministicStorage.CallOpts, arg0)
}

// InitiatedWithdrawals is a free data retrieval call binding the contract method 0xbf286e39.
//
// Solidity: function initiatedWithdrawals(bytes32 ) view returns(bool)
func (_L2PortalDeterministicStorage *L2PortalDeterministicStorageCallerSession) InitiatedWithdrawals(arg0 [32]byte) (bool, error) {
	return _L2PortalDeterministicStorage.Contract.InitiatedWithdrawals(&_L2PortalDeterministicStorage.CallOpts, arg0)
}
