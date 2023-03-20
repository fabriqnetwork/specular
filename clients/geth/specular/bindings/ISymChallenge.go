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
)

// ISymChallengeMetaData contains all meta data concerning the ISymChallenge contract.
var ISymChallengeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineNotPassed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotYourTurn\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"challengeState\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"challengedSegmentStart\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"challengedSegmentLength\",\"type\":\"uint256\"}],\"name\":\"Bisected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"loser\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIChallenge.CompletionReason\",\"name\":\"reason\",\"type\":\"uint8\"}],\"name\":\"Completed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"bisection\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"challengedSegmentIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"prevBisection\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"prevChallengedSegmentStart\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevChallengedSegmentLength\",\"type\":\"uint256\"}],\"name\":\"bisectExecution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentResponder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentResponderTimeLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_numSteps\",\"type\":\"uint256\"}],\"name\":\"initializeChallengeLength\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"oneStepProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"challengedStepIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"prevBisection\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"prevChallengedSegmentStart\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevChallengedSegmentLength\",\"type\":\"uint256\"}],\"name\":\"verifyOneStepProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ISymChallengeABI is the input ABI used to generate the binding from.
// Deprecated: Use ISymChallengeMetaData.ABI instead.
var ISymChallengeABI = ISymChallengeMetaData.ABI

// ISymChallenge is an auto generated Go binding around an Ethereum contract.
type ISymChallenge struct {
	ISymChallengeCaller     // Read-only binding to the contract
	ISymChallengeTransactor // Write-only binding to the contract
	ISymChallengeFilterer   // Log filterer for contract events
}

// ISymChallengeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISymChallengeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISymChallengeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISymChallengeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISymChallengeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISymChallengeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISymChallengeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISymChallengeSession struct {
	Contract     *ISymChallenge    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISymChallengeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISymChallengeCallerSession struct {
	Contract *ISymChallengeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ISymChallengeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISymChallengeTransactorSession struct {
	Contract     *ISymChallengeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ISymChallengeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISymChallengeRaw struct {
	Contract *ISymChallenge // Generic contract binding to access the raw methods on
}

// ISymChallengeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISymChallengeCallerRaw struct {
	Contract *ISymChallengeCaller // Generic read-only contract binding to access the raw methods on
}

// ISymChallengeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISymChallengeTransactorRaw struct {
	Contract *ISymChallengeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISymChallenge creates a new instance of ISymChallenge, bound to a specific deployed contract.
func NewISymChallenge(address common.Address, backend bind.ContractBackend) (*ISymChallenge, error) {
	contract, err := bindISymChallenge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISymChallenge{ISymChallengeCaller: ISymChallengeCaller{contract: contract}, ISymChallengeTransactor: ISymChallengeTransactor{contract: contract}, ISymChallengeFilterer: ISymChallengeFilterer{contract: contract}}, nil
}

// NewISymChallengeCaller creates a new read-only instance of ISymChallenge, bound to a specific deployed contract.
func NewISymChallengeCaller(address common.Address, caller bind.ContractCaller) (*ISymChallengeCaller, error) {
	contract, err := bindISymChallenge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISymChallengeCaller{contract: contract}, nil
}

// NewISymChallengeTransactor creates a new write-only instance of ISymChallenge, bound to a specific deployed contract.
func NewISymChallengeTransactor(address common.Address, transactor bind.ContractTransactor) (*ISymChallengeTransactor, error) {
	contract, err := bindISymChallenge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISymChallengeTransactor{contract: contract}, nil
}

// NewISymChallengeFilterer creates a new log filterer instance of ISymChallenge, bound to a specific deployed contract.
func NewISymChallengeFilterer(address common.Address, filterer bind.ContractFilterer) (*ISymChallengeFilterer, error) {
	contract, err := bindISymChallenge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISymChallengeFilterer{contract: contract}, nil
}

// bindISymChallenge binds a generic wrapper to an already deployed contract.
func bindISymChallenge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISymChallengeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISymChallenge *ISymChallengeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISymChallenge.Contract.ISymChallengeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISymChallenge *ISymChallengeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISymChallenge.Contract.ISymChallengeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISymChallenge *ISymChallengeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISymChallenge.Contract.ISymChallengeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISymChallenge *ISymChallengeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISymChallenge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISymChallenge *ISymChallengeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISymChallenge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISymChallenge *ISymChallengeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISymChallenge.Contract.contract.Transact(opts, method, params...)
}

// CurrentResponder is a free data retrieval call binding the contract method 0x8a8cd218.
//
// Solidity: function currentResponder() view returns(address)
func (_ISymChallenge *ISymChallengeCaller) CurrentResponder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ISymChallenge.contract.Call(opts, &out, "currentResponder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CurrentResponder is a free data retrieval call binding the contract method 0x8a8cd218.
//
// Solidity: function currentResponder() view returns(address)
func (_ISymChallenge *ISymChallengeSession) CurrentResponder() (common.Address, error) {
	return _ISymChallenge.Contract.CurrentResponder(&_ISymChallenge.CallOpts)
}

// CurrentResponder is a free data retrieval call binding the contract method 0x8a8cd218.
//
// Solidity: function currentResponder() view returns(address)
func (_ISymChallenge *ISymChallengeCallerSession) CurrentResponder() (common.Address, error) {
	return _ISymChallenge.Contract.CurrentResponder(&_ISymChallenge.CallOpts)
}

// CurrentResponderTimeLeft is a free data retrieval call binding the contract method 0xe87e3589.
//
// Solidity: function currentResponderTimeLeft() view returns(uint256)
func (_ISymChallenge *ISymChallengeCaller) CurrentResponderTimeLeft(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ISymChallenge.contract.Call(opts, &out, "currentResponderTimeLeft")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentResponderTimeLeft is a free data retrieval call binding the contract method 0xe87e3589.
//
// Solidity: function currentResponderTimeLeft() view returns(uint256)
func (_ISymChallenge *ISymChallengeSession) CurrentResponderTimeLeft() (*big.Int, error) {
	return _ISymChallenge.Contract.CurrentResponderTimeLeft(&_ISymChallenge.CallOpts)
}

// CurrentResponderTimeLeft is a free data retrieval call binding the contract method 0xe87e3589.
//
// Solidity: function currentResponderTimeLeft() view returns(uint256)
func (_ISymChallenge *ISymChallengeCallerSession) CurrentResponderTimeLeft() (*big.Int, error) {
	return _ISymChallenge.Contract.CurrentResponderTimeLeft(&_ISymChallenge.CallOpts)
}

// BisectExecution is a paid mutator transaction binding the contract method 0xcc8f6677.
//
// Solidity: function bisectExecution(bytes32[] bisection, uint256 challengedSegmentIndex, bytes32[] prevBisection, uint256 prevChallengedSegmentStart, uint256 prevChallengedSegmentLength) returns()
func (_ISymChallenge *ISymChallengeTransactor) BisectExecution(opts *bind.TransactOpts, bisection [][32]byte, challengedSegmentIndex *big.Int, prevBisection [][32]byte, prevChallengedSegmentStart *big.Int, prevChallengedSegmentLength *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.contract.Transact(opts, "bisectExecution", bisection, challengedSegmentIndex, prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength)
}

// BisectExecution is a paid mutator transaction binding the contract method 0xcc8f6677.
//
// Solidity: function bisectExecution(bytes32[] bisection, uint256 challengedSegmentIndex, bytes32[] prevBisection, uint256 prevChallengedSegmentStart, uint256 prevChallengedSegmentLength) returns()
func (_ISymChallenge *ISymChallengeSession) BisectExecution(bisection [][32]byte, challengedSegmentIndex *big.Int, prevBisection [][32]byte, prevChallengedSegmentStart *big.Int, prevChallengedSegmentLength *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.Contract.BisectExecution(&_ISymChallenge.TransactOpts, bisection, challengedSegmentIndex, prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength)
}

// BisectExecution is a paid mutator transaction binding the contract method 0xcc8f6677.
//
// Solidity: function bisectExecution(bytes32[] bisection, uint256 challengedSegmentIndex, bytes32[] prevBisection, uint256 prevChallengedSegmentStart, uint256 prevChallengedSegmentLength) returns()
func (_ISymChallenge *ISymChallengeTransactorSession) BisectExecution(bisection [][32]byte, challengedSegmentIndex *big.Int, prevBisection [][32]byte, prevChallengedSegmentStart *big.Int, prevChallengedSegmentLength *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.Contract.BisectExecution(&_ISymChallenge.TransactOpts, bisection, challengedSegmentIndex, prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength)
}

// InitializeChallengeLength is a paid mutator transaction binding the contract method 0x9909e0d9.
//
// Solidity: function initializeChallengeLength(uint256 _numSteps) returns()
func (_ISymChallenge *ISymChallengeTransactor) InitializeChallengeLength(opts *bind.TransactOpts, _numSteps *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.contract.Transact(opts, "initializeChallengeLength", _numSteps)
}

// InitializeChallengeLength is a paid mutator transaction binding the contract method 0x9909e0d9.
//
// Solidity: function initializeChallengeLength(uint256 _numSteps) returns()
func (_ISymChallenge *ISymChallengeSession) InitializeChallengeLength(_numSteps *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.Contract.InitializeChallengeLength(&_ISymChallenge.TransactOpts, _numSteps)
}

// InitializeChallengeLength is a paid mutator transaction binding the contract method 0x9909e0d9.
//
// Solidity: function initializeChallengeLength(uint256 _numSteps) returns()
func (_ISymChallenge *ISymChallengeTransactorSession) InitializeChallengeLength(_numSteps *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.Contract.InitializeChallengeLength(&_ISymChallenge.TransactOpts, _numSteps)
}

// Timeout is a paid mutator transaction binding the contract method 0x70dea79a.
//
// Solidity: function timeout() returns()
func (_ISymChallenge *ISymChallengeTransactor) Timeout(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISymChallenge.contract.Transact(opts, "timeout")
}

// Timeout is a paid mutator transaction binding the contract method 0x70dea79a.
//
// Solidity: function timeout() returns()
func (_ISymChallenge *ISymChallengeSession) Timeout() (*types.Transaction, error) {
	return _ISymChallenge.Contract.Timeout(&_ISymChallenge.TransactOpts)
}

// Timeout is a paid mutator transaction binding the contract method 0x70dea79a.
//
// Solidity: function timeout() returns()
func (_ISymChallenge *ISymChallengeTransactorSession) Timeout() (*types.Transaction, error) {
	return _ISymChallenge.Contract.Timeout(&_ISymChallenge.TransactOpts)
}

// VerifyOneStepProof is a paid mutator transaction binding the contract method 0x45b258a7.
//
// Solidity: function verifyOneStepProof(bytes oneStepProof, uint256 challengedStepIndex, bytes32[] prevBisection, uint256 prevChallengedSegmentStart, uint256 prevChallengedSegmentLength) returns()
func (_ISymChallenge *ISymChallengeTransactor) VerifyOneStepProof(opts *bind.TransactOpts, oneStepProof []byte, challengedStepIndex *big.Int, prevBisection [][32]byte, prevChallengedSegmentStart *big.Int, prevChallengedSegmentLength *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.contract.Transact(opts, "verifyOneStepProof", oneStepProof, challengedStepIndex, prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength)
}

// VerifyOneStepProof is a paid mutator transaction binding the contract method 0x45b258a7.
//
// Solidity: function verifyOneStepProof(bytes oneStepProof, uint256 challengedStepIndex, bytes32[] prevBisection, uint256 prevChallengedSegmentStart, uint256 prevChallengedSegmentLength) returns()
func (_ISymChallenge *ISymChallengeSession) VerifyOneStepProof(oneStepProof []byte, challengedStepIndex *big.Int, prevBisection [][32]byte, prevChallengedSegmentStart *big.Int, prevChallengedSegmentLength *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.Contract.VerifyOneStepProof(&_ISymChallenge.TransactOpts, oneStepProof, challengedStepIndex, prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength)
}

// VerifyOneStepProof is a paid mutator transaction binding the contract method 0x45b258a7.
//
// Solidity: function verifyOneStepProof(bytes oneStepProof, uint256 challengedStepIndex, bytes32[] prevBisection, uint256 prevChallengedSegmentStart, uint256 prevChallengedSegmentLength) returns()
func (_ISymChallenge *ISymChallengeTransactorSession) VerifyOneStepProof(oneStepProof []byte, challengedStepIndex *big.Int, prevBisection [][32]byte, prevChallengedSegmentStart *big.Int, prevChallengedSegmentLength *big.Int) (*types.Transaction, error) {
	return _ISymChallenge.Contract.VerifyOneStepProof(&_ISymChallenge.TransactOpts, oneStepProof, challengedStepIndex, prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength)
}

// ISymChallengeBisectedIterator is returned from FilterBisected and is used to iterate over the raw logs and unpacked data for Bisected events raised by the ISymChallenge contract.
type ISymChallengeBisectedIterator struct {
	Event *ISymChallengeBisected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ISymChallengeBisectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISymChallengeBisected)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ISymChallengeBisected)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ISymChallengeBisectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISymChallengeBisectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISymChallengeBisected represents a Bisected event raised by the ISymChallenge contract.
type ISymChallengeBisected struct {
	ChallengeState          [32]byte
	ChallengedSegmentStart  *big.Int
	ChallengedSegmentLength *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterBisected is a free log retrieval operation binding the contract event 0x8c3cfc522d91af51bb14f6db452f8c212ba664a426c79e5ef78872e7a1072074.
//
// Solidity: event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength)
func (_ISymChallenge *ISymChallengeFilterer) FilterBisected(opts *bind.FilterOpts) (*ISymChallengeBisectedIterator, error) {

	logs, sub, err := _ISymChallenge.contract.FilterLogs(opts, "Bisected")
	if err != nil {
		return nil, err
	}
	return &ISymChallengeBisectedIterator{contract: _ISymChallenge.contract, event: "Bisected", logs: logs, sub: sub}, nil
}

// WatchBisected is a free log subscription operation binding the contract event 0x8c3cfc522d91af51bb14f6db452f8c212ba664a426c79e5ef78872e7a1072074.
//
// Solidity: event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength)
func (_ISymChallenge *ISymChallengeFilterer) WatchBisected(opts *bind.WatchOpts, sink chan<- *ISymChallengeBisected) (event.Subscription, error) {

	logs, sub, err := _ISymChallenge.contract.WatchLogs(opts, "Bisected")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISymChallengeBisected)
				if err := _ISymChallenge.contract.UnpackLog(event, "Bisected", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBisected is a log parse operation binding the contract event 0x8c3cfc522d91af51bb14f6db452f8c212ba664a426c79e5ef78872e7a1072074.
//
// Solidity: event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength)
func (_ISymChallenge *ISymChallengeFilterer) ParseBisected(log types.Log) (*ISymChallengeBisected, error) {
	event := new(ISymChallengeBisected)
	if err := _ISymChallenge.contract.UnpackLog(event, "Bisected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ISymChallengeCompletedIterator is returned from FilterCompleted and is used to iterate over the raw logs and unpacked data for Completed events raised by the ISymChallenge contract.
type ISymChallengeCompletedIterator struct {
	Event *ISymChallengeCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ISymChallengeCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ISymChallengeCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ISymChallengeCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ISymChallengeCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ISymChallengeCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ISymChallengeCompleted represents a Completed event raised by the ISymChallenge contract.
type ISymChallengeCompleted struct {
	Winner common.Address
	Loser  common.Address
	Reason uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCompleted is a free log retrieval operation binding the contract event 0xa599fa89698188ea23144af5bd981dc904e4221ee98ed73883b509409808338d.
//
// Solidity: event Completed(address winner, address loser, uint8 reason)
func (_ISymChallenge *ISymChallengeFilterer) FilterCompleted(opts *bind.FilterOpts) (*ISymChallengeCompletedIterator, error) {

	logs, sub, err := _ISymChallenge.contract.FilterLogs(opts, "Completed")
	if err != nil {
		return nil, err
	}
	return &ISymChallengeCompletedIterator{contract: _ISymChallenge.contract, event: "Completed", logs: logs, sub: sub}, nil
}

// WatchCompleted is a free log subscription operation binding the contract event 0xa599fa89698188ea23144af5bd981dc904e4221ee98ed73883b509409808338d.
//
// Solidity: event Completed(address winner, address loser, uint8 reason)
func (_ISymChallenge *ISymChallengeFilterer) WatchCompleted(opts *bind.WatchOpts, sink chan<- *ISymChallengeCompleted) (event.Subscription, error) {

	logs, sub, err := _ISymChallenge.contract.WatchLogs(opts, "Completed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ISymChallengeCompleted)
				if err := _ISymChallenge.contract.UnpackLog(event, "Completed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCompleted is a log parse operation binding the contract event 0xa599fa89698188ea23144af5bd981dc904e4221ee98ed73883b509409808338d.
//
// Solidity: event Completed(address winner, address loser, uint8 reason)
func (_ISymChallenge *ISymChallengeFilterer) ParseCompleted(log types.Log) (*ISymChallengeCompleted, error) {
	event := new(ISymChallengeCompleted)
	if err := _ISymChallenge.contract.UnpackLog(event, "Completed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
