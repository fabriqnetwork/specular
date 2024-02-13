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

// IChallengeResultReceiverMetaData contains all meta data concerning the IChallengeResultReceiver contract.
var IChallengeResultReceiverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"loser\",\"type\":\"address\"}],\"name\":\"completeChallenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IChallengeResultReceiverABI is the input ABI used to generate the binding from.
// Deprecated: Use IChallengeResultReceiverMetaData.ABI instead.
var IChallengeResultReceiverABI = IChallengeResultReceiverMetaData.ABI

// IChallengeResultReceiver is an auto generated Go binding around an Ethereum contract.
type IChallengeResultReceiver struct {
	IChallengeResultReceiverCaller     // Read-only binding to the contract
	IChallengeResultReceiverTransactor // Write-only binding to the contract
	IChallengeResultReceiverFilterer   // Log filterer for contract events
}

// IChallengeResultReceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type IChallengeResultReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IChallengeResultReceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IChallengeResultReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IChallengeResultReceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IChallengeResultReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IChallengeResultReceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IChallengeResultReceiverSession struct {
	Contract     *IChallengeResultReceiver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IChallengeResultReceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IChallengeResultReceiverCallerSession struct {
	Contract *IChallengeResultReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// IChallengeResultReceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IChallengeResultReceiverTransactorSession struct {
	Contract     *IChallengeResultReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// IChallengeResultReceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type IChallengeResultReceiverRaw struct {
	Contract *IChallengeResultReceiver // Generic contract binding to access the raw methods on
}

// IChallengeResultReceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IChallengeResultReceiverCallerRaw struct {
	Contract *IChallengeResultReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// IChallengeResultReceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IChallengeResultReceiverTransactorRaw struct {
	Contract *IChallengeResultReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIChallengeResultReceiver creates a new instance of IChallengeResultReceiver, bound to a specific deployed contract.
func NewIChallengeResultReceiver(address common.Address, backend bind.ContractBackend) (*IChallengeResultReceiver, error) {
	contract, err := bindIChallengeResultReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IChallengeResultReceiver{IChallengeResultReceiverCaller: IChallengeResultReceiverCaller{contract: contract}, IChallengeResultReceiverTransactor: IChallengeResultReceiverTransactor{contract: contract}, IChallengeResultReceiverFilterer: IChallengeResultReceiverFilterer{contract: contract}}, nil
}

// NewIChallengeResultReceiverCaller creates a new read-only instance of IChallengeResultReceiver, bound to a specific deployed contract.
func NewIChallengeResultReceiverCaller(address common.Address, caller bind.ContractCaller) (*IChallengeResultReceiverCaller, error) {
	contract, err := bindIChallengeResultReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IChallengeResultReceiverCaller{contract: contract}, nil
}

// NewIChallengeResultReceiverTransactor creates a new write-only instance of IChallengeResultReceiver, bound to a specific deployed contract.
func NewIChallengeResultReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*IChallengeResultReceiverTransactor, error) {
	contract, err := bindIChallengeResultReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IChallengeResultReceiverTransactor{contract: contract}, nil
}

// NewIChallengeResultReceiverFilterer creates a new log filterer instance of IChallengeResultReceiver, bound to a specific deployed contract.
func NewIChallengeResultReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*IChallengeResultReceiverFilterer, error) {
	contract, err := bindIChallengeResultReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IChallengeResultReceiverFilterer{contract: contract}, nil
}

// bindIChallengeResultReceiver binds a generic wrapper to an already deployed contract.
func bindIChallengeResultReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IChallengeResultReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IChallengeResultReceiver *IChallengeResultReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IChallengeResultReceiver.Contract.IChallengeResultReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IChallengeResultReceiver *IChallengeResultReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IChallengeResultReceiver.Contract.IChallengeResultReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IChallengeResultReceiver *IChallengeResultReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IChallengeResultReceiver.Contract.IChallengeResultReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IChallengeResultReceiver *IChallengeResultReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IChallengeResultReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IChallengeResultReceiver *IChallengeResultReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IChallengeResultReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IChallengeResultReceiver *IChallengeResultReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IChallengeResultReceiver.Contract.contract.Transact(opts, method, params...)
}

// CompleteChallenge is a paid mutator transaction binding the contract method 0xfa7803e6.
//
// Solidity: function completeChallenge(address winner, address loser) returns()
func (_IChallengeResultReceiver *IChallengeResultReceiverTransactor) CompleteChallenge(opts *bind.TransactOpts, winner common.Address, loser common.Address) (*types.Transaction, error) {
	return _IChallengeResultReceiver.contract.Transact(opts, "completeChallenge", winner, loser)
}

// CompleteChallenge is a paid mutator transaction binding the contract method 0xfa7803e6.
//
// Solidity: function completeChallenge(address winner, address loser) returns()
func (_IChallengeResultReceiver *IChallengeResultReceiverSession) CompleteChallenge(winner common.Address, loser common.Address) (*types.Transaction, error) {
	return _IChallengeResultReceiver.Contract.CompleteChallenge(&_IChallengeResultReceiver.TransactOpts, winner, loser)
}

// CompleteChallenge is a paid mutator transaction binding the contract method 0xfa7803e6.
//
// Solidity: function completeChallenge(address winner, address loser) returns()
func (_IChallengeResultReceiver *IChallengeResultReceiverTransactorSession) CompleteChallenge(winner common.Address, loser common.Address) (*types.Transaction, error) {
	return _IChallengeResultReceiver.Contract.CompleteChallenge(&_IChallengeResultReceiver.TransactOpts, winner, loser)
}
