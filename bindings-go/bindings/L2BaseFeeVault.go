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

// L2BaseFeeVaultMetaData contains all meta data concerning the L2BaseFeeVault contract.
var L2BaseFeeVaultMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minWithdrawalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalProcessed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051610fbb61004c600039600081816101f00152818161023901528181610486015281816104c601526105550152610fbb6000f3fe6080604052600436106100a05760003560e01c80638129fc1c116100645780638129fc1c146101335780638312f1491461014857806384411d651461015e5780638da5cb5b14610174578063f2bcd022146101a6578063f2fde38b146101c657600080fd5b80633659cfe6146100ac5780633ccfd60b146100ce5780634f1ef286146100e357806352d1902d146100f6578063715018a61461011e57600080fd5b366100a757005b600080fd5b3480156100b857600080fd5b506100cc6100c7366004610cbb565b6101e6565b005b3480156100da57600080fd5b506100cc6102ce565b6100cc6100f1366004610cec565b61047c565b34801561010257600080fd5b5061010b610548565b6040519081526020015b60405180910390f35b34801561012a57600080fd5b506100cc6105fb565b34801561013f57600080fd5b506100cc61060f565b34801561015457600080fd5b5061010b60c95481565b34801561016a57600080fd5b5061010b60cb5481565b34801561018057600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610115565b3480156101b257600080fd5b5060ca5461018e906001600160a01b031681565b3480156101d257600080fd5b506100cc6101e1366004610cbb565b610727565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036102375760405162461bcd60e51b815260040161022e90610dae565b60405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610280600080516020610f3f833981519152546001600160a01b031690565b6001600160a01b0316146102a65760405162461bcd60e51b815260040161022e90610dfa565b6102af8161079d565b604080516000808252602082019092526102cb918391906107a5565b50565b60c9544710156103595760405162461bcd60e51b815260206004820152604a60248201527f4665655661756c743a207769746864726177616c20616d6f756e74206d75737460448201527f2062652067726561746572207468616e206d696e696d756d20776974686472616064820152691dd85b08185b5bdd5b9d60b21b608482015260a40161022e565b60004790508060cb60008282546103709190610e46565b909155505060ca54604080518381526001600160a01b0390921660208301523382820152517fc8a211cc64b6ed1b50595a9fcb1932b6d1e5a6e8ef15b60e5b1f988ea9086bba9181900360600190a160ca546040516000916001600160a01b03169083908381818185875af1925050503d806000811461040c576040519150601f19603f3d011682016040523d82523d6000602084013e610411565b606091505b50509050806104785760405162461bcd60e51b815260206004820152602d60248201527f4665655661756c743a206661696c656420746f2073656e642045544820746f2060448201526c199959481c9958da5c1a595b9d609a1b606482015260840161022e565b5050565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036104c45760405162461bcd60e51b815260040161022e90610dae565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661050d600080516020610f3f833981519152546001600160a01b031690565b6001600160a01b0316146105335760405162461bcd60e51b815260040161022e90610dfa565b61053c8261079d565b610478828260016107a5565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146105e85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c0000000000000000606482015260840161022e565b50600080516020610f3f83398151915290565b610603610915565b61060d600061096f565b565b600054610100900460ff161580801561062f5750600054600160ff909116105b806106495750303b158015610649575060005460ff166001145b6106ac5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b606482015260840161022e565b6000805460ff1916600117905580156106cf576000805461ff0019166101001790555b6106d76109c1565b6106df6109f0565b80156102cb576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b61072f610915565b6001600160a01b0381166107945760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161022e565b6102cb8161096f565b6102cb610915565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156107dd576107d883610a17565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610837575060408051601f3d908101601f1916820190925261083491810190610e67565b60015b61089a5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b606482015260840161022e565b600080516020610f3f83398151915281146109095760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b606482015260840161022e565b506107d8838383610ab3565b6097546001600160a01b0316331461060d5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161022e565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166109e85760405162461bcd60e51b815260040161022e90610e80565b61060d610ade565b600054610100900460ff1661060d5760405162461bcd60e51b815260040161022e90610e80565b6001600160a01b0381163b610a845760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840161022e565b600080516020610f3f83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b610abc83610b0e565b600082511180610ac95750805b156107d857610ad88383610b4e565b50505050565b600054610100900460ff16610b055760405162461bcd60e51b815260040161022e90610e80565b61060d3361096f565b610b1781610a17565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060610b738383604051806060016040528060278152602001610f5f60279139610b7c565b90505b92915050565b6060600080856001600160a01b031685604051610b999190610eef565b600060405180830381855af49150503d8060008114610bd4576040519150601f19603f3d011682016040523d82523d6000602084013e610bd9565b606091505b5091509150610bea86838387610bf4565b9695505050505050565b60608315610c63578251600003610c5c576001600160a01b0385163b610c5c5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161022e565b5081610c6d565b610c6d8383610c75565b949350505050565b815115610c855781518083602001fd5b8060405162461bcd60e51b815260040161022e9190610f0b565b80356001600160a01b0381168114610cb657600080fd5b919050565b600060208284031215610ccd57600080fd5b610b7382610c9f565b634e487b7160e01b600052604160045260246000fd5b60008060408385031215610cff57600080fd5b610d0883610c9f565b9150602083013567ffffffffffffffff80821115610d2557600080fd5b818501915085601f830112610d3957600080fd5b813581811115610d4b57610d4b610cd6565b604051601f8201601f19908116603f01168101908382118183101715610d7357610d73610cd6565b81604052828152886020848701011115610d8c57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b80820180821115610b7657634e487b7160e01b600052601160045260246000fd5b600060208284031215610e7957600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60005b83811015610ee6578181015183820152602001610ece565b50506000910152565b60008251610f01818460208701610ecb565b9190910192915050565b6020815260008251806020840152610f2a816040850160208701610ecb565b601f01601f1916919091016040019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212205e4634541a8ee1b5feb834a6ef1df7ac681f6dcdde98f6f6890ac511241aa67264736f6c63430008110033",
}

// L2BaseFeeVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use L2BaseFeeVaultMetaData.ABI instead.
var L2BaseFeeVaultABI = L2BaseFeeVaultMetaData.ABI

// L2BaseFeeVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L2BaseFeeVaultMetaData.Bin instead.
var L2BaseFeeVaultBin = L2BaseFeeVaultMetaData.Bin

// DeployL2BaseFeeVault deploys a new Ethereum contract, binding an instance of L2BaseFeeVault to it.
func DeployL2BaseFeeVault(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L2BaseFeeVault, error) {
	parsed, err := L2BaseFeeVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L2BaseFeeVaultBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L2BaseFeeVault{L2BaseFeeVaultCaller: L2BaseFeeVaultCaller{contract: contract}, L2BaseFeeVaultTransactor: L2BaseFeeVaultTransactor{contract: contract}, L2BaseFeeVaultFilterer: L2BaseFeeVaultFilterer{contract: contract}}, nil
}

// L2BaseFeeVault is an auto generated Go binding around an Ethereum contract.
type L2BaseFeeVault struct {
	L2BaseFeeVaultCaller     // Read-only binding to the contract
	L2BaseFeeVaultTransactor // Write-only binding to the contract
	L2BaseFeeVaultFilterer   // Log filterer for contract events
}

// L2BaseFeeVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2BaseFeeVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2BaseFeeVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2BaseFeeVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2BaseFeeVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2BaseFeeVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2BaseFeeVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2BaseFeeVaultSession struct {
	Contract     *L2BaseFeeVault   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L2BaseFeeVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2BaseFeeVaultCallerSession struct {
	Contract *L2BaseFeeVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// L2BaseFeeVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2BaseFeeVaultTransactorSession struct {
	Contract     *L2BaseFeeVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// L2BaseFeeVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2BaseFeeVaultRaw struct {
	Contract *L2BaseFeeVault // Generic contract binding to access the raw methods on
}

// L2BaseFeeVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2BaseFeeVaultCallerRaw struct {
	Contract *L2BaseFeeVaultCaller // Generic read-only contract binding to access the raw methods on
}

// L2BaseFeeVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2BaseFeeVaultTransactorRaw struct {
	Contract *L2BaseFeeVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2BaseFeeVault creates a new instance of L2BaseFeeVault, bound to a specific deployed contract.
func NewL2BaseFeeVault(address common.Address, backend bind.ContractBackend) (*L2BaseFeeVault, error) {
	contract, err := bindL2BaseFeeVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVault{L2BaseFeeVaultCaller: L2BaseFeeVaultCaller{contract: contract}, L2BaseFeeVaultTransactor: L2BaseFeeVaultTransactor{contract: contract}, L2BaseFeeVaultFilterer: L2BaseFeeVaultFilterer{contract: contract}}, nil
}

// NewL2BaseFeeVaultCaller creates a new read-only instance of L2BaseFeeVault, bound to a specific deployed contract.
func NewL2BaseFeeVaultCaller(address common.Address, caller bind.ContractCaller) (*L2BaseFeeVaultCaller, error) {
	contract, err := bindL2BaseFeeVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultCaller{contract: contract}, nil
}

// NewL2BaseFeeVaultTransactor creates a new write-only instance of L2BaseFeeVault, bound to a specific deployed contract.
func NewL2BaseFeeVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*L2BaseFeeVaultTransactor, error) {
	contract, err := bindL2BaseFeeVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultTransactor{contract: contract}, nil
}

// NewL2BaseFeeVaultFilterer creates a new log filterer instance of L2BaseFeeVault, bound to a specific deployed contract.
func NewL2BaseFeeVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*L2BaseFeeVaultFilterer, error) {
	contract, err := bindL2BaseFeeVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultFilterer{contract: contract}, nil
}

// bindL2BaseFeeVault binds a generic wrapper to an already deployed contract.
func bindL2BaseFeeVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2BaseFeeVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2BaseFeeVault *L2BaseFeeVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2BaseFeeVault.Contract.L2BaseFeeVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2BaseFeeVault *L2BaseFeeVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.L2BaseFeeVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2BaseFeeVault *L2BaseFeeVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.L2BaseFeeVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2BaseFeeVault *L2BaseFeeVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2BaseFeeVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.contract.Transact(opts, method, params...)
}

// MinWithdrawalAmount is a free data retrieval call binding the contract method 0x8312f149.
//
// Solidity: function minWithdrawalAmount() view returns(uint256)
func (_L2BaseFeeVault *L2BaseFeeVaultCaller) MinWithdrawalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2BaseFeeVault.contract.Call(opts, &out, "minWithdrawalAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinWithdrawalAmount is a free data retrieval call binding the contract method 0x8312f149.
//
// Solidity: function minWithdrawalAmount() view returns(uint256)
func (_L2BaseFeeVault *L2BaseFeeVaultSession) MinWithdrawalAmount() (*big.Int, error) {
	return _L2BaseFeeVault.Contract.MinWithdrawalAmount(&_L2BaseFeeVault.CallOpts)
}

// MinWithdrawalAmount is a free data retrieval call binding the contract method 0x8312f149.
//
// Solidity: function minWithdrawalAmount() view returns(uint256)
func (_L2BaseFeeVault *L2BaseFeeVaultCallerSession) MinWithdrawalAmount() (*big.Int, error) {
	return _L2BaseFeeVault.Contract.MinWithdrawalAmount(&_L2BaseFeeVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2BaseFeeVault *L2BaseFeeVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2BaseFeeVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2BaseFeeVault *L2BaseFeeVaultSession) Owner() (common.Address, error) {
	return _L2BaseFeeVault.Contract.Owner(&_L2BaseFeeVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2BaseFeeVault *L2BaseFeeVaultCallerSession) Owner() (common.Address, error) {
	return _L2BaseFeeVault.Contract.Owner(&_L2BaseFeeVault.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L2BaseFeeVault *L2BaseFeeVaultCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L2BaseFeeVault.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L2BaseFeeVault *L2BaseFeeVaultSession) ProxiableUUID() ([32]byte, error) {
	return _L2BaseFeeVault.Contract.ProxiableUUID(&_L2BaseFeeVault.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L2BaseFeeVault *L2BaseFeeVaultCallerSession) ProxiableUUID() ([32]byte, error) {
	return _L2BaseFeeVault.Contract.ProxiableUUID(&_L2BaseFeeVault.CallOpts)
}

// TotalProcessed is a free data retrieval call binding the contract method 0x84411d65.
//
// Solidity: function totalProcessed() view returns(uint256)
func (_L2BaseFeeVault *L2BaseFeeVaultCaller) TotalProcessed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2BaseFeeVault.contract.Call(opts, &out, "totalProcessed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalProcessed is a free data retrieval call binding the contract method 0x84411d65.
//
// Solidity: function totalProcessed() view returns(uint256)
func (_L2BaseFeeVault *L2BaseFeeVaultSession) TotalProcessed() (*big.Int, error) {
	return _L2BaseFeeVault.Contract.TotalProcessed(&_L2BaseFeeVault.CallOpts)
}

// TotalProcessed is a free data retrieval call binding the contract method 0x84411d65.
//
// Solidity: function totalProcessed() view returns(uint256)
func (_L2BaseFeeVault *L2BaseFeeVaultCallerSession) TotalProcessed() (*big.Int, error) {
	return _L2BaseFeeVault.Contract.TotalProcessed(&_L2BaseFeeVault.CallOpts)
}

// WithdrawalAddress is a free data retrieval call binding the contract method 0xf2bcd022.
//
// Solidity: function withdrawalAddress() view returns(address)
func (_L2BaseFeeVault *L2BaseFeeVaultCaller) WithdrawalAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2BaseFeeVault.contract.Call(opts, &out, "withdrawalAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WithdrawalAddress is a free data retrieval call binding the contract method 0xf2bcd022.
//
// Solidity: function withdrawalAddress() view returns(address)
func (_L2BaseFeeVault *L2BaseFeeVaultSession) WithdrawalAddress() (common.Address, error) {
	return _L2BaseFeeVault.Contract.WithdrawalAddress(&_L2BaseFeeVault.CallOpts)
}

// WithdrawalAddress is a free data retrieval call binding the contract method 0xf2bcd022.
//
// Solidity: function withdrawalAddress() view returns(address)
func (_L2BaseFeeVault *L2BaseFeeVaultCallerSession) WithdrawalAddress() (common.Address, error) {
	return _L2BaseFeeVault.Contract.WithdrawalAddress(&_L2BaseFeeVault.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) Initialize() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.Initialize(&_L2BaseFeeVault.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) Initialize() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.Initialize(&_L2BaseFeeVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.RenounceOwnership(&_L2BaseFeeVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.RenounceOwnership(&_L2BaseFeeVault.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.TransferOwnership(&_L2BaseFeeVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.TransferOwnership(&_L2BaseFeeVault.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.UpgradeTo(&_L2BaseFeeVault.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.UpgradeTo(&_L2BaseFeeVault.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.UpgradeToAndCall(&_L2BaseFeeVault.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.UpgradeToAndCall(&_L2BaseFeeVault.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) Withdraw() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.Withdraw(&_L2BaseFeeVault.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) Withdraw() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.Withdraw(&_L2BaseFeeVault.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2BaseFeeVault.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2BaseFeeVault *L2BaseFeeVaultSession) Receive() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.Receive(&_L2BaseFeeVault.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2BaseFeeVault *L2BaseFeeVaultTransactorSession) Receive() (*types.Transaction, error) {
	return _L2BaseFeeVault.Contract.Receive(&_L2BaseFeeVault.TransactOpts)
}

// L2BaseFeeVaultAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultAdminChangedIterator struct {
	Event *L2BaseFeeVaultAdminChanged // Event containing the contract specifics and raw log

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
func (it *L2BaseFeeVaultAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2BaseFeeVaultAdminChanged)
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
		it.Event = new(L2BaseFeeVaultAdminChanged)
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
func (it *L2BaseFeeVaultAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2BaseFeeVaultAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2BaseFeeVaultAdminChanged represents a AdminChanged event raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*L2BaseFeeVaultAdminChangedIterator, error) {

	logs, sub, err := _L2BaseFeeVault.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultAdminChangedIterator{contract: _L2BaseFeeVault.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *L2BaseFeeVaultAdminChanged) (event.Subscription, error) {

	logs, sub, err := _L2BaseFeeVault.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2BaseFeeVaultAdminChanged)
				if err := _L2BaseFeeVault.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) ParseAdminChanged(log types.Log) (*L2BaseFeeVaultAdminChanged, error) {
	event := new(L2BaseFeeVaultAdminChanged)
	if err := _L2BaseFeeVault.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2BaseFeeVaultBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultBeaconUpgradedIterator struct {
	Event *L2BaseFeeVaultBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *L2BaseFeeVaultBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2BaseFeeVaultBeaconUpgraded)
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
		it.Event = new(L2BaseFeeVaultBeaconUpgraded)
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
func (it *L2BaseFeeVaultBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2BaseFeeVaultBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2BaseFeeVaultBeaconUpgraded represents a BeaconUpgraded event raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*L2BaseFeeVaultBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _L2BaseFeeVault.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultBeaconUpgradedIterator{contract: _L2BaseFeeVault.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *L2BaseFeeVaultBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _L2BaseFeeVault.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2BaseFeeVaultBeaconUpgraded)
				if err := _L2BaseFeeVault.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) ParseBeaconUpgraded(log types.Log) (*L2BaseFeeVaultBeaconUpgraded, error) {
	event := new(L2BaseFeeVaultBeaconUpgraded)
	if err := _L2BaseFeeVault.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2BaseFeeVaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultInitializedIterator struct {
	Event *L2BaseFeeVaultInitialized // Event containing the contract specifics and raw log

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
func (it *L2BaseFeeVaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2BaseFeeVaultInitialized)
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
		it.Event = new(L2BaseFeeVaultInitialized)
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
func (it *L2BaseFeeVaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2BaseFeeVaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2BaseFeeVaultInitialized represents a Initialized event raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*L2BaseFeeVaultInitializedIterator, error) {

	logs, sub, err := _L2BaseFeeVault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultInitializedIterator{contract: _L2BaseFeeVault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *L2BaseFeeVaultInitialized) (event.Subscription, error) {

	logs, sub, err := _L2BaseFeeVault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2BaseFeeVaultInitialized)
				if err := _L2BaseFeeVault.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) ParseInitialized(log types.Log) (*L2BaseFeeVaultInitialized, error) {
	event := new(L2BaseFeeVaultInitialized)
	if err := _L2BaseFeeVault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2BaseFeeVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultOwnershipTransferredIterator struct {
	Event *L2BaseFeeVaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L2BaseFeeVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2BaseFeeVaultOwnershipTransferred)
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
		it.Event = new(L2BaseFeeVaultOwnershipTransferred)
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
func (it *L2BaseFeeVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2BaseFeeVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2BaseFeeVaultOwnershipTransferred represents a OwnershipTransferred event raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L2BaseFeeVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2BaseFeeVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultOwnershipTransferredIterator{contract: _L2BaseFeeVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L2BaseFeeVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2BaseFeeVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2BaseFeeVaultOwnershipTransferred)
				if err := _L2BaseFeeVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) ParseOwnershipTransferred(log types.Log) (*L2BaseFeeVaultOwnershipTransferred, error) {
	event := new(L2BaseFeeVaultOwnershipTransferred)
	if err := _L2BaseFeeVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2BaseFeeVaultUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultUpgradedIterator struct {
	Event *L2BaseFeeVaultUpgraded // Event containing the contract specifics and raw log

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
func (it *L2BaseFeeVaultUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2BaseFeeVaultUpgraded)
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
		it.Event = new(L2BaseFeeVaultUpgraded)
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
func (it *L2BaseFeeVaultUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2BaseFeeVaultUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2BaseFeeVaultUpgraded represents a Upgraded event raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*L2BaseFeeVaultUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L2BaseFeeVault.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultUpgradedIterator{contract: _L2BaseFeeVault.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *L2BaseFeeVaultUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L2BaseFeeVault.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2BaseFeeVaultUpgraded)
				if err := _L2BaseFeeVault.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) ParseUpgraded(log types.Log) (*L2BaseFeeVaultUpgraded, error) {
	event := new(L2BaseFeeVaultUpgraded)
	if err := _L2BaseFeeVault.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2BaseFeeVaultWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultWithdrawalIterator struct {
	Event *L2BaseFeeVaultWithdrawal // Event containing the contract specifics and raw log

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
func (it *L2BaseFeeVaultWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2BaseFeeVaultWithdrawal)
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
		it.Event = new(L2BaseFeeVaultWithdrawal)
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
func (it *L2BaseFeeVaultWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2BaseFeeVaultWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2BaseFeeVaultWithdrawal represents a Withdrawal event raised by the L2BaseFeeVault contract.
type L2BaseFeeVaultWithdrawal struct {
	Value *big.Int
	To    common.Address
	From  common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0xc8a211cc64b6ed1b50595a9fcb1932b6d1e5a6e8ef15b60e5b1f988ea9086bba.
//
// Solidity: event Withdrawal(uint256 value, address to, address from)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*L2BaseFeeVaultWithdrawalIterator, error) {

	logs, sub, err := _L2BaseFeeVault.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &L2BaseFeeVaultWithdrawalIterator{contract: _L2BaseFeeVault.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0xc8a211cc64b6ed1b50595a9fcb1932b6d1e5a6e8ef15b60e5b1f988ea9086bba.
//
// Solidity: event Withdrawal(uint256 value, address to, address from)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *L2BaseFeeVaultWithdrawal) (event.Subscription, error) {

	logs, sub, err := _L2BaseFeeVault.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2BaseFeeVaultWithdrawal)
				if err := _L2BaseFeeVault.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0xc8a211cc64b6ed1b50595a9fcb1932b6d1e5a6e8ef15b60e5b1f988ea9086bba.
//
// Solidity: event Withdrawal(uint256 value, address to, address from)
func (_L2BaseFeeVault *L2BaseFeeVaultFilterer) ParseWithdrawal(log types.Log) (*L2BaseFeeVaultWithdrawal, error) {
	event := new(L2BaseFeeVaultWithdrawal)
	if err := _L2BaseFeeVault.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
