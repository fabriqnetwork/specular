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

// IAsymChallengeMetaData contains all meta data concerning the IAsymChallenge contract.
var IAsymChallengeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DeadlineNotPassed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotYourTurn\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"challengeState\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"challengedSegmentStart\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"challengedSegmentLength\",\"type\":\"uint256\"}],\"name\":\"Bisected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"loser\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIChallenge.CompletionReason\",\"name\":\"reason\",\"type\":\"uint8\"}],\"name\":\"Completed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"currentResponder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentResponderTimeLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IAsymChallengeABI is the input ABI used to generate the binding from.
// Deprecated: Use IAsymChallengeMetaData.ABI instead.
var IAsymChallengeABI = IAsymChallengeMetaData.ABI

// IAsymChallenge is an auto generated Go binding around an Ethereum contract.
type IAsymChallenge struct {
	IAsymChallengeCaller     // Read-only binding to the contract
	IAsymChallengeTransactor // Write-only binding to the contract
	IAsymChallengeFilterer   // Log filterer for contract events
}

// IAsymChallengeCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAsymChallengeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAsymChallengeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAsymChallengeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAsymChallengeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAsymChallengeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAsymChallengeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAsymChallengeSession struct {
	Contract     *IAsymChallenge   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAsymChallengeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAsymChallengeCallerSession struct {
	Contract *IAsymChallengeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IAsymChallengeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAsymChallengeTransactorSession struct {
	Contract     *IAsymChallengeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IAsymChallengeRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAsymChallengeRaw struct {
	Contract *IAsymChallenge // Generic contract binding to access the raw methods on
}

// IAsymChallengeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAsymChallengeCallerRaw struct {
	Contract *IAsymChallengeCaller // Generic read-only contract binding to access the raw methods on
}

// IAsymChallengeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAsymChallengeTransactorRaw struct {
	Contract *IAsymChallengeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAsymChallenge creates a new instance of IAsymChallenge, bound to a specific deployed contract.
func NewIAsymChallenge(address common.Address, backend bind.ContractBackend) (*IAsymChallenge, error) {
	contract, err := bindIAsymChallenge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAsymChallenge{IAsymChallengeCaller: IAsymChallengeCaller{contract: contract}, IAsymChallengeTransactor: IAsymChallengeTransactor{contract: contract}, IAsymChallengeFilterer: IAsymChallengeFilterer{contract: contract}}, nil
}

// NewIAsymChallengeCaller creates a new read-only instance of IAsymChallenge, bound to a specific deployed contract.
func NewIAsymChallengeCaller(address common.Address, caller bind.ContractCaller) (*IAsymChallengeCaller, error) {
	contract, err := bindIAsymChallenge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAsymChallengeCaller{contract: contract}, nil
}

// NewIAsymChallengeTransactor creates a new write-only instance of IAsymChallenge, bound to a specific deployed contract.
func NewIAsymChallengeTransactor(address common.Address, transactor bind.ContractTransactor) (*IAsymChallengeTransactor, error) {
	contract, err := bindIAsymChallenge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAsymChallengeTransactor{contract: contract}, nil
}

// NewIAsymChallengeFilterer creates a new log filterer instance of IAsymChallenge, bound to a specific deployed contract.
func NewIAsymChallengeFilterer(address common.Address, filterer bind.ContractFilterer) (*IAsymChallengeFilterer, error) {
	contract, err := bindIAsymChallenge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAsymChallengeFilterer{contract: contract}, nil
}

// bindIAsymChallenge binds a generic wrapper to an already deployed contract.
func bindIAsymChallenge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAsymChallengeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAsymChallenge *IAsymChallengeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAsymChallenge.Contract.IAsymChallengeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAsymChallenge *IAsymChallengeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAsymChallenge.Contract.IAsymChallengeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAsymChallenge *IAsymChallengeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAsymChallenge.Contract.IAsymChallengeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAsymChallenge *IAsymChallengeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAsymChallenge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAsymChallenge *IAsymChallengeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAsymChallenge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAsymChallenge *IAsymChallengeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAsymChallenge.Contract.contract.Transact(opts, method, params...)
}

// CurrentResponder is a free data retrieval call binding the contract method 0x8a8cd218.
//
// Solidity: function currentResponder() view returns(address)
func (_IAsymChallenge *IAsymChallengeCaller) CurrentResponder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAsymChallenge.contract.Call(opts, &out, "currentResponder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CurrentResponder is a free data retrieval call binding the contract method 0x8a8cd218.
//
// Solidity: function currentResponder() view returns(address)
func (_IAsymChallenge *IAsymChallengeSession) CurrentResponder() (common.Address, error) {
	return _IAsymChallenge.Contract.CurrentResponder(&_IAsymChallenge.CallOpts)
}

// CurrentResponder is a free data retrieval call binding the contract method 0x8a8cd218.
//
// Solidity: function currentResponder() view returns(address)
func (_IAsymChallenge *IAsymChallengeCallerSession) CurrentResponder() (common.Address, error) {
	return _IAsymChallenge.Contract.CurrentResponder(&_IAsymChallenge.CallOpts)
}

// CurrentResponderTimeLeft is a free data retrieval call binding the contract method 0xe87e3589.
//
// Solidity: function currentResponderTimeLeft() view returns(uint256)
func (_IAsymChallenge *IAsymChallengeCaller) CurrentResponderTimeLeft(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAsymChallenge.contract.Call(opts, &out, "currentResponderTimeLeft")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentResponderTimeLeft is a free data retrieval call binding the contract method 0xe87e3589.
//
// Solidity: function currentResponderTimeLeft() view returns(uint256)
func (_IAsymChallenge *IAsymChallengeSession) CurrentResponderTimeLeft() (*big.Int, error) {
	return _IAsymChallenge.Contract.CurrentResponderTimeLeft(&_IAsymChallenge.CallOpts)
}

// CurrentResponderTimeLeft is a free data retrieval call binding the contract method 0xe87e3589.
//
// Solidity: function currentResponderTimeLeft() view returns(uint256)
func (_IAsymChallenge *IAsymChallengeCallerSession) CurrentResponderTimeLeft() (*big.Int, error) {
	return _IAsymChallenge.Contract.CurrentResponderTimeLeft(&_IAsymChallenge.CallOpts)
}

// Timeout is a paid mutator transaction binding the contract method 0x70dea79a.
//
// Solidity: function timeout() returns()
func (_IAsymChallenge *IAsymChallengeTransactor) Timeout(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAsymChallenge.contract.Transact(opts, "timeout")
}

// Timeout is a paid mutator transaction binding the contract method 0x70dea79a.
//
// Solidity: function timeout() returns()
func (_IAsymChallenge *IAsymChallengeSession) Timeout() (*types.Transaction, error) {
	return _IAsymChallenge.Contract.Timeout(&_IAsymChallenge.TransactOpts)
}

// Timeout is a paid mutator transaction binding the contract method 0x70dea79a.
//
// Solidity: function timeout() returns()
func (_IAsymChallenge *IAsymChallengeTransactorSession) Timeout() (*types.Transaction, error) {
	return _IAsymChallenge.Contract.Timeout(&_IAsymChallenge.TransactOpts)
}

// IAsymChallengeBisectedIterator is returned from FilterBisected and is used to iterate over the raw logs and unpacked data for Bisected events raised by the IAsymChallenge contract.
type IAsymChallengeBisectedIterator struct {
	Event *IAsymChallengeBisected // Event containing the contract specifics and raw log

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
func (it *IAsymChallengeBisectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAsymChallengeBisected)
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
		it.Event = new(IAsymChallengeBisected)
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
func (it *IAsymChallengeBisectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAsymChallengeBisectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAsymChallengeBisected represents a Bisected event raised by the IAsymChallenge contract.
type IAsymChallengeBisected struct {
	ChallengeState          [32]byte
	ChallengedSegmentStart  *big.Int
	ChallengedSegmentLength *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterBisected is a free log retrieval operation binding the contract event 0x8c3cfc522d91af51bb14f6db452f8c212ba664a426c79e5ef78872e7a1072074.
//
// Solidity: event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength)
func (_IAsymChallenge *IAsymChallengeFilterer) FilterBisected(opts *bind.FilterOpts) (*IAsymChallengeBisectedIterator, error) {

	logs, sub, err := _IAsymChallenge.contract.FilterLogs(opts, "Bisected")
	if err != nil {
		return nil, err
	}
	return &IAsymChallengeBisectedIterator{contract: _IAsymChallenge.contract, event: "Bisected", logs: logs, sub: sub}, nil
}

// WatchBisected is a free log subscription operation binding the contract event 0x8c3cfc522d91af51bb14f6db452f8c212ba664a426c79e5ef78872e7a1072074.
//
// Solidity: event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength)
func (_IAsymChallenge *IAsymChallengeFilterer) WatchBisected(opts *bind.WatchOpts, sink chan<- *IAsymChallengeBisected) (event.Subscription, error) {

	logs, sub, err := _IAsymChallenge.contract.WatchLogs(opts, "Bisected")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAsymChallengeBisected)
				if err := _IAsymChallenge.contract.UnpackLog(event, "Bisected", log); err != nil {
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
func (_IAsymChallenge *IAsymChallengeFilterer) ParseBisected(log types.Log) (*IAsymChallengeBisected, error) {
	event := new(IAsymChallengeBisected)
	if err := _IAsymChallenge.contract.UnpackLog(event, "Bisected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAsymChallengeCompletedIterator is returned from FilterCompleted and is used to iterate over the raw logs and unpacked data for Completed events raised by the IAsymChallenge contract.
type IAsymChallengeCompletedIterator struct {
	Event *IAsymChallengeCompleted // Event containing the contract specifics and raw log

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
func (it *IAsymChallengeCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAsymChallengeCompleted)
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
		it.Event = new(IAsymChallengeCompleted)
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
func (it *IAsymChallengeCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAsymChallengeCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAsymChallengeCompleted represents a Completed event raised by the IAsymChallenge contract.
type IAsymChallengeCompleted struct {
	Winner common.Address
	Loser  common.Address
	Reason uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCompleted is a free log retrieval operation binding the contract event 0xa599fa89698188ea23144af5bd981dc904e4221ee98ed73883b509409808338d.
//
// Solidity: event Completed(address winner, address loser, uint8 reason)
func (_IAsymChallenge *IAsymChallengeFilterer) FilterCompleted(opts *bind.FilterOpts) (*IAsymChallengeCompletedIterator, error) {

	logs, sub, err := _IAsymChallenge.contract.FilterLogs(opts, "Completed")
	if err != nil {
		return nil, err
	}
	return &IAsymChallengeCompletedIterator{contract: _IAsymChallenge.contract, event: "Completed", logs: logs, sub: sub}, nil
}

// WatchCompleted is a free log subscription operation binding the contract event 0xa599fa89698188ea23144af5bd981dc904e4221ee98ed73883b509409808338d.
//
// Solidity: event Completed(address winner, address loser, uint8 reason)
func (_IAsymChallenge *IAsymChallengeFilterer) WatchCompleted(opts *bind.WatchOpts, sink chan<- *IAsymChallengeCompleted) (event.Subscription, error) {

	logs, sub, err := _IAsymChallenge.contract.WatchLogs(opts, "Completed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAsymChallengeCompleted)
				if err := _IAsymChallenge.contract.UnpackLog(event, "Completed", log); err != nil {
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
func (_IAsymChallenge *IAsymChallengeFilterer) ParseCompleted(log types.Log) (*IAsymChallengeCompleted, error) {
	event := new(IAsymChallengeCompleted)
	if err := _IAsymChallenge.contract.UnpackLog(event, "Completed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
