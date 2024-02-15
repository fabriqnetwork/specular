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

// UUPSPlaceholderMetaData contains all meta data concerning the UUPSPlaceholder contract.
var UUPSPlaceholderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100dd565b600054610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116146100db576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b610aab806100ec6000396000f3fe6080604052600436106100705760003560e01c8063715018a61161004e578063715018a6146100ea5780638129fc1c146100ff5780638da5cb5b14610114578063f2fde38b1461013c57600080fd5b80633659cfe6146100755780634f1ef2861461009757806352d1902d146100aa575b600080fd5b34801561008157600080fd5b50610095610090366004610884565b61015c565b005b6100956100a53660046108b5565b610184565b3480156100b657600080fd5b506040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81526020015b60405180910390f35b3480156100f657600080fd5b5061009561019d565b34801561010b57600080fd5b506100956101b1565b34801561012057600080fd5b506097546040516001600160a01b0390911681526020016100e1565b34801561014857600080fd5b50610095610157366004610884565b6102ce565b61016581610344565b604080516000808252602082019092526101819183919061034c565b50565b61018d82610344565b6101998282600161034c565b5050565b6101a56104ce565b6101af6000610528565b565b600054610100900460ff16158080156101d15750600054600160ff909116105b806101eb5750303b1580156101eb575060005460ff166001145b6102535760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff191660011790558015610276576000805461ff0019166101001790555b61027e61057a565b6102866105a9565b8015610181576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b6102d66104ce565b6001600160a01b03811661033b5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161024a565b61018181610528565b6101816104ce565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156103845761037f836105d0565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156103de575060408051601f3d908101601f191682019092526103db91810190610977565b60015b6104415760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b606482015260840161024a565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81146104c25760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b606482015260840161024a565b5061037f83838361067e565b6097546001600160a01b031633146101af5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161024a565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166105a15760405162461bcd60e51b815260040161024a90610990565b6101af6106a9565b600054610100900460ff166101af5760405162461bcd60e51b815260040161024a90610990565b6001600160a01b0381163b61063d5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840161024a565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b610687836106d9565b6000825111806106945750805b1561037f576106a38383610719565b50505050565b600054610100900460ff166106d05760405162461bcd60e51b815260040161024a90610990565b6101af33610528565b6106e2816105d0565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061073e8383604051806060016040528060278152602001610a4f60279139610745565b9392505050565b6060600080856001600160a01b03168560405161076291906109ff565b600060405180830381855af49150503d806000811461079d576040519150601f19603f3d011682016040523d82523d6000602084013e6107a2565b606091505b50915091506107b3868383876107bd565b9695505050505050565b6060831561082c578251600003610825576001600160a01b0385163b6108255760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161024a565b5081610836565b610836838361083e565b949350505050565b81511561084e5781518083602001fd5b8060405162461bcd60e51b815260040161024a9190610a1b565b80356001600160a01b038116811461087f57600080fd5b919050565b60006020828403121561089657600080fd5b61073e82610868565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156108c857600080fd5b6108d183610868565b9150602083013567ffffffffffffffff808211156108ee57600080fd5b818501915085601f83011261090257600080fd5b8135818111156109145761091461089f565b604051601f8201601f19908116603f0116810190838211818310171561093c5761093c61089f565b8160405282815288602084870101111561095557600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60006020828403121561098957600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60005b838110156109f65781810151838201526020016109de565b50506000910152565b60008251610a118184602087016109db565b9190910192915050565b6020815260008251806020840152610a3a8160408501602087016109db565b601f01601f1916919091016040019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212204c505b40878a2a11d04a0ee978f9990c19a48b6b8c34a0e71ba0ba56be9e75ac64736f6c63430008110033",
}

// UUPSPlaceholderABI is the input ABI used to generate the binding from.
// Deprecated: Use UUPSPlaceholderMetaData.ABI instead.
var UUPSPlaceholderABI = UUPSPlaceholderMetaData.ABI

// UUPSPlaceholderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UUPSPlaceholderMetaData.Bin instead.
var UUPSPlaceholderBin = UUPSPlaceholderMetaData.Bin

// DeployUUPSPlaceholder deploys a new Ethereum contract, binding an instance of UUPSPlaceholder to it.
func DeployUUPSPlaceholder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UUPSPlaceholder, error) {
	parsed, err := UUPSPlaceholderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UUPSPlaceholderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UUPSPlaceholder{UUPSPlaceholderCaller: UUPSPlaceholderCaller{contract: contract}, UUPSPlaceholderTransactor: UUPSPlaceholderTransactor{contract: contract}, UUPSPlaceholderFilterer: UUPSPlaceholderFilterer{contract: contract}}, nil
}

// UUPSPlaceholder is an auto generated Go binding around an Ethereum contract.
type UUPSPlaceholder struct {
	UUPSPlaceholderCaller     // Read-only binding to the contract
	UUPSPlaceholderTransactor // Write-only binding to the contract
	UUPSPlaceholderFilterer   // Log filterer for contract events
}

// UUPSPlaceholderCaller is an auto generated read-only Go binding around an Ethereum contract.
type UUPSPlaceholderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSPlaceholderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UUPSPlaceholderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSPlaceholderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UUPSPlaceholderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSPlaceholderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UUPSPlaceholderSession struct {
	Contract     *UUPSPlaceholder  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UUPSPlaceholderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UUPSPlaceholderCallerSession struct {
	Contract *UUPSPlaceholderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// UUPSPlaceholderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UUPSPlaceholderTransactorSession struct {
	Contract     *UUPSPlaceholderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UUPSPlaceholderRaw is an auto generated low-level Go binding around an Ethereum contract.
type UUPSPlaceholderRaw struct {
	Contract *UUPSPlaceholder // Generic contract binding to access the raw methods on
}

// UUPSPlaceholderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UUPSPlaceholderCallerRaw struct {
	Contract *UUPSPlaceholderCaller // Generic read-only contract binding to access the raw methods on
}

// UUPSPlaceholderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UUPSPlaceholderTransactorRaw struct {
	Contract *UUPSPlaceholderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUUPSPlaceholder creates a new instance of UUPSPlaceholder, bound to a specific deployed contract.
func NewUUPSPlaceholder(address common.Address, backend bind.ContractBackend) (*UUPSPlaceholder, error) {
	contract, err := bindUUPSPlaceholder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholder{UUPSPlaceholderCaller: UUPSPlaceholderCaller{contract: contract}, UUPSPlaceholderTransactor: UUPSPlaceholderTransactor{contract: contract}, UUPSPlaceholderFilterer: UUPSPlaceholderFilterer{contract: contract}}, nil
}

// NewUUPSPlaceholderCaller creates a new read-only instance of UUPSPlaceholder, bound to a specific deployed contract.
func NewUUPSPlaceholderCaller(address common.Address, caller bind.ContractCaller) (*UUPSPlaceholderCaller, error) {
	contract, err := bindUUPSPlaceholder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderCaller{contract: contract}, nil
}

// NewUUPSPlaceholderTransactor creates a new write-only instance of UUPSPlaceholder, bound to a specific deployed contract.
func NewUUPSPlaceholderTransactor(address common.Address, transactor bind.ContractTransactor) (*UUPSPlaceholderTransactor, error) {
	contract, err := bindUUPSPlaceholder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderTransactor{contract: contract}, nil
}

// NewUUPSPlaceholderFilterer creates a new log filterer instance of UUPSPlaceholder, bound to a specific deployed contract.
func NewUUPSPlaceholderFilterer(address common.Address, filterer bind.ContractFilterer) (*UUPSPlaceholderFilterer, error) {
	contract, err := bindUUPSPlaceholder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderFilterer{contract: contract}, nil
}

// bindUUPSPlaceholder binds a generic wrapper to an already deployed contract.
func bindUUPSPlaceholder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UUPSPlaceholderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UUPSPlaceholder *UUPSPlaceholderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UUPSPlaceholder.Contract.UUPSPlaceholderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UUPSPlaceholder *UUPSPlaceholderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.UUPSPlaceholderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UUPSPlaceholder *UUPSPlaceholderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.UUPSPlaceholderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UUPSPlaceholder *UUPSPlaceholderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UUPSPlaceholder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UUPSPlaceholder *UUPSPlaceholderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UUPSPlaceholder *UUPSPlaceholderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UUPSPlaceholder *UUPSPlaceholderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UUPSPlaceholder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UUPSPlaceholder *UUPSPlaceholderSession) Owner() (common.Address, error) {
	return _UUPSPlaceholder.Contract.Owner(&_UUPSPlaceholder.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UUPSPlaceholder *UUPSPlaceholderCallerSession) Owner() (common.Address, error) {
	return _UUPSPlaceholder.Contract.Owner(&_UUPSPlaceholder.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSPlaceholder *UUPSPlaceholderCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UUPSPlaceholder.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSPlaceholder *UUPSPlaceholderSession) ProxiableUUID() ([32]byte, error) {
	return _UUPSPlaceholder.Contract.ProxiableUUID(&_UUPSPlaceholder.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSPlaceholder *UUPSPlaceholderCallerSession) ProxiableUUID() ([32]byte, error) {
	return _UUPSPlaceholder.Contract.ProxiableUUID(&_UUPSPlaceholder.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSPlaceholder.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_UUPSPlaceholder *UUPSPlaceholderSession) Initialize() (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.Initialize(&_UUPSPlaceholder.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactorSession) Initialize() (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.Initialize(&_UUPSPlaceholder.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSPlaceholder.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UUPSPlaceholder *UUPSPlaceholderSession) RenounceOwnership() (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.RenounceOwnership(&_UUPSPlaceholder.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.RenounceOwnership(&_UUPSPlaceholder.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UUPSPlaceholder.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UUPSPlaceholder *UUPSPlaceholderSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.TransferOwnership(&_UUPSPlaceholder.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.TransferOwnership(&_UUPSPlaceholder.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSPlaceholder.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSPlaceholder *UUPSPlaceholderSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.UpgradeTo(&_UUPSPlaceholder.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.UpgradeTo(&_UUPSPlaceholder.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSPlaceholder.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSPlaceholder *UUPSPlaceholderSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.UpgradeToAndCall(&_UUPSPlaceholder.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSPlaceholder *UUPSPlaceholderTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSPlaceholder.Contract.UpgradeToAndCall(&_UUPSPlaceholder.TransactOpts, newImplementation, data)
}

// UUPSPlaceholderAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the UUPSPlaceholder contract.
type UUPSPlaceholderAdminChangedIterator struct {
	Event *UUPSPlaceholderAdminChanged // Event containing the contract specifics and raw log

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
func (it *UUPSPlaceholderAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSPlaceholderAdminChanged)
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
		it.Event = new(UUPSPlaceholderAdminChanged)
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
func (it *UUPSPlaceholderAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSPlaceholderAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSPlaceholderAdminChanged represents a AdminChanged event raised by the UUPSPlaceholder contract.
type UUPSPlaceholderAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*UUPSPlaceholderAdminChangedIterator, error) {

	logs, sub, err := _UUPSPlaceholder.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderAdminChangedIterator{contract: _UUPSPlaceholder.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *UUPSPlaceholderAdminChanged) (event.Subscription, error) {

	logs, sub, err := _UUPSPlaceholder.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSPlaceholderAdminChanged)
				if err := _UUPSPlaceholder.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) ParseAdminChanged(log types.Log) (*UUPSPlaceholderAdminChanged, error) {
	event := new(UUPSPlaceholderAdminChanged)
	if err := _UUPSPlaceholder.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UUPSPlaceholderBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the UUPSPlaceholder contract.
type UUPSPlaceholderBeaconUpgradedIterator struct {
	Event *UUPSPlaceholderBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *UUPSPlaceholderBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSPlaceholderBeaconUpgraded)
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
		it.Event = new(UUPSPlaceholderBeaconUpgraded)
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
func (it *UUPSPlaceholderBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSPlaceholderBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSPlaceholderBeaconUpgraded represents a BeaconUpgraded event raised by the UUPSPlaceholder contract.
type UUPSPlaceholderBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*UUPSPlaceholderBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _UUPSPlaceholder.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderBeaconUpgradedIterator{contract: _UUPSPlaceholder.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *UUPSPlaceholderBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _UUPSPlaceholder.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSPlaceholderBeaconUpgraded)
				if err := _UUPSPlaceholder.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) ParseBeaconUpgraded(log types.Log) (*UUPSPlaceholderBeaconUpgraded, error) {
	event := new(UUPSPlaceholderBeaconUpgraded)
	if err := _UUPSPlaceholder.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UUPSPlaceholderInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the UUPSPlaceholder contract.
type UUPSPlaceholderInitializedIterator struct {
	Event *UUPSPlaceholderInitialized // Event containing the contract specifics and raw log

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
func (it *UUPSPlaceholderInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSPlaceholderInitialized)
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
		it.Event = new(UUPSPlaceholderInitialized)
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
func (it *UUPSPlaceholderInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSPlaceholderInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSPlaceholderInitialized represents a Initialized event raised by the UUPSPlaceholder contract.
type UUPSPlaceholderInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) FilterInitialized(opts *bind.FilterOpts) (*UUPSPlaceholderInitializedIterator, error) {

	logs, sub, err := _UUPSPlaceholder.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderInitializedIterator{contract: _UUPSPlaceholder.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *UUPSPlaceholderInitialized) (event.Subscription, error) {

	logs, sub, err := _UUPSPlaceholder.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSPlaceholderInitialized)
				if err := _UUPSPlaceholder.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) ParseInitialized(log types.Log) (*UUPSPlaceholderInitialized, error) {
	event := new(UUPSPlaceholderInitialized)
	if err := _UUPSPlaceholder.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UUPSPlaceholderOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UUPSPlaceholder contract.
type UUPSPlaceholderOwnershipTransferredIterator struct {
	Event *UUPSPlaceholderOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *UUPSPlaceholderOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSPlaceholderOwnershipTransferred)
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
		it.Event = new(UUPSPlaceholderOwnershipTransferred)
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
func (it *UUPSPlaceholderOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSPlaceholderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSPlaceholderOwnershipTransferred represents a OwnershipTransferred event raised by the UUPSPlaceholder contract.
type UUPSPlaceholderOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UUPSPlaceholderOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UUPSPlaceholder.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderOwnershipTransferredIterator{contract: _UUPSPlaceholder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UUPSPlaceholderOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UUPSPlaceholder.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSPlaceholderOwnershipTransferred)
				if err := _UUPSPlaceholder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) ParseOwnershipTransferred(log types.Log) (*UUPSPlaceholderOwnershipTransferred, error) {
	event := new(UUPSPlaceholderOwnershipTransferred)
	if err := _UUPSPlaceholder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UUPSPlaceholderUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the UUPSPlaceholder contract.
type UUPSPlaceholderUpgradedIterator struct {
	Event *UUPSPlaceholderUpgraded // Event containing the contract specifics and raw log

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
func (it *UUPSPlaceholderUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSPlaceholderUpgraded)
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
		it.Event = new(UUPSPlaceholderUpgraded)
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
func (it *UUPSPlaceholderUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSPlaceholderUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSPlaceholderUpgraded represents a Upgraded event raised by the UUPSPlaceholder contract.
type UUPSPlaceholderUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*UUPSPlaceholderUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UUPSPlaceholder.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &UUPSPlaceholderUpgradedIterator{contract: _UUPSPlaceholder.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UUPSPlaceholderUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UUPSPlaceholder.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSPlaceholderUpgraded)
				if err := _UUPSPlaceholder.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSPlaceholder *UUPSPlaceholderFilterer) ParseUpgraded(log types.Log) (*UUPSPlaceholderUpgraded, error) {
	event := new(UUPSPlaceholderUpgraded)
	if err := _UUPSPlaceholder.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
