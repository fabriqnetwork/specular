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

// IRollupAssertion is an auto generated low-level Go binding around an user-defined struct.
type IRollupAssertion struct {
	StateHash      [32]byte
	InboxSize      *big.Int
	Parent         *big.Int
	Deadline       *big.Int
	ProposalTime   *big.Int
	NumStakers     *big.Int
	ChildInboxSize *big.Int
}

// IRollupStaker is an auto generated low-level Go binding around an user-defined struct.
type IRollupStaker struct {
	IsStaked         bool
	AmountStaked     *big.Int
	AssertionID      *big.Int
	CurrentChallenge common.Address
}

// RollupBaseMetaData contains all meta data concerning the RollupBase contract.
var RollupBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AssertionAlreadyResolved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AssertionOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChallengedStaker\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfirmationPeriodPending\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateAssertion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyAssertion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker1Challenge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"staker2Challenge\",\"type\":\"address\"}],\"name\":\"InDifferentChallenge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InboxReadLimitExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientStake\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfigChange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInboxSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidParent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MinimumAssertionPeriodNotPassed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoStaker\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoUnresolvedAssertion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAllStaked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"challenge\",\"type\":\"address\"}],\"name\":\"NotChallengeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInChallenge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotSiblings\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotStaked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParentAssertionUnstaked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PreviousStateHash\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakedOnUnconfirmedAssertion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakerStakedOnTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakersPresent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnproposedAssertion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"challengeAddr\",\"type\":\"address\"}],\"name\":\"AssertionChallenged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"}],\"name\":\"AssertionConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asserterAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"vmHash\",\"type\":\"bytes32\"}],\"name\":\"AssertionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"}],\"name\":\"AssertionRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakerAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"}],\"name\":\"StakerStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"}],\"name\":\"advanceStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"assertions\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"inboxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"parent\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proposalTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numStakers\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"childInboxSize\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[2]\",\"name\":\"players\",\"type\":\"address[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"assertionIDs\",\"type\":\"uint256[2]\"}],\"name\":\"challengeAssertion\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"loser\",\"type\":\"address\"}],\"name\":\"completeChallenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"confirmFirstUnresolvedAssertion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"confirmationPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"confirmedInboxSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"vmHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"inboxSize\",\"type\":\"uint256\"}],\"name\":\"createAssertion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRequiredStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daProvider\",\"outputs\":[{\"internalType\":\"contractIDAProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"}],\"name\":\"getAssertion\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"inboxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"parent\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proposalTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numStakers\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"childInboxSize\",\"type\":\"uint256\"}],\"internalType\":\"structIRollup.Assertion\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastConfirmedAssertionID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getStaker\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isStaked\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"amountStaked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"currentChallenge\",\"type\":\"address\"}],\"internalType\":\"structIRollup.Staker\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_daProvider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_verifier\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_confirmationPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_challengePeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minimumAssertionPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_baseStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_initialAssertionID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_initialInboxSize\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_initialVMhash\",\"type\":\"bytes32\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"stakerAddress\",\"type\":\"address\"}],\"name\":\"isStakedOnAssertion\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastConfirmedAssertionID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastCreatedAssertionID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastResolvedAssertionID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minimumAssertionPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numStakers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stakerAddress\",\"type\":\"address\"}],\"name\":\"rejectFirstUnresolvedAssertion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stakerAddress\",\"type\":\"address\"}],\"name\":\"removeStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requireFirstUnresolvedAssertionIsConfirmable\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stakerAddress\",\"type\":\"address\"}],\"name\":\"requireFirstUnresolvedAssertionIsRejectable\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"setBaseStakeAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPeriod\",\"type\":\"uint256\"}],\"name\":\"setChallengePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPeriod\",\"type\":\"uint256\"}],\"name\":\"setConfirmationPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newDAProvider\",\"type\":\"address\"}],\"name\":\"setDAProvider\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPeriod\",\"type\":\"uint256\"}],\"name\":\"setMinimumAssertionPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isStaked\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"amountStaked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"assertionID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"currentChallenge\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vault\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdrawableFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"zombies\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"stakerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"lastAssertionID\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RollupBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupBaseMetaData.ABI instead.
var RollupBaseABI = RollupBaseMetaData.ABI

// RollupBase is an auto generated Go binding around an Ethereum contract.
type RollupBase struct {
	RollupBaseCaller     // Read-only binding to the contract
	RollupBaseTransactor // Write-only binding to the contract
	RollupBaseFilterer   // Log filterer for contract events
}

// RollupBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupBaseSession struct {
	Contract     *RollupBase       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupBaseCallerSession struct {
	Contract *RollupBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// RollupBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupBaseTransactorSession struct {
	Contract     *RollupBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// RollupBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupBaseRaw struct {
	Contract *RollupBase // Generic contract binding to access the raw methods on
}

// RollupBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupBaseCallerRaw struct {
	Contract *RollupBaseCaller // Generic read-only contract binding to access the raw methods on
}

// RollupBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupBaseTransactorRaw struct {
	Contract *RollupBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollupBase creates a new instance of RollupBase, bound to a specific deployed contract.
func NewRollupBase(address common.Address, backend bind.ContractBackend) (*RollupBase, error) {
	contract, err := bindRollupBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollupBase{RollupBaseCaller: RollupBaseCaller{contract: contract}, RollupBaseTransactor: RollupBaseTransactor{contract: contract}, RollupBaseFilterer: RollupBaseFilterer{contract: contract}}, nil
}

// NewRollupBaseCaller creates a new read-only instance of RollupBase, bound to a specific deployed contract.
func NewRollupBaseCaller(address common.Address, caller bind.ContractCaller) (*RollupBaseCaller, error) {
	contract, err := bindRollupBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupBaseCaller{contract: contract}, nil
}

// NewRollupBaseTransactor creates a new write-only instance of RollupBase, bound to a specific deployed contract.
func NewRollupBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupBaseTransactor, error) {
	contract, err := bindRollupBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupBaseTransactor{contract: contract}, nil
}

// NewRollupBaseFilterer creates a new log filterer instance of RollupBase, bound to a specific deployed contract.
func NewRollupBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupBaseFilterer, error) {
	contract, err := bindRollupBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupBaseFilterer{contract: contract}, nil
}

// bindRollupBase binds a generic wrapper to an already deployed contract.
func bindRollupBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RollupBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupBase *RollupBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupBase.Contract.RollupBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupBase *RollupBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupBase.Contract.RollupBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupBase *RollupBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupBase.Contract.RollupBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupBase *RollupBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupBase *RollupBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupBase *RollupBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupBase.Contract.contract.Transact(opts, method, params...)
}

// Assertions is a free data retrieval call binding the contract method 0x524232f6.
//
// Solidity: function assertions(uint256 ) view returns(bytes32 stateHash, uint256 inboxSize, uint256 parent, uint256 deadline, uint256 proposalTime, uint256 numStakers, uint256 childInboxSize)
func (_RollupBase *RollupBaseCaller) Assertions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	StateHash      [32]byte
	InboxSize      *big.Int
	Parent         *big.Int
	Deadline       *big.Int
	ProposalTime   *big.Int
	NumStakers     *big.Int
	ChildInboxSize *big.Int
}, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "assertions", arg0)

	outstruct := new(struct {
		StateHash      [32]byte
		InboxSize      *big.Int
		Parent         *big.Int
		Deadline       *big.Int
		ProposalTime   *big.Int
		NumStakers     *big.Int
		ChildInboxSize *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StateHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.InboxSize = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Parent = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Deadline = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.ProposalTime = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.NumStakers = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ChildInboxSize = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Assertions is a free data retrieval call binding the contract method 0x524232f6.
//
// Solidity: function assertions(uint256 ) view returns(bytes32 stateHash, uint256 inboxSize, uint256 parent, uint256 deadline, uint256 proposalTime, uint256 numStakers, uint256 childInboxSize)
func (_RollupBase *RollupBaseSession) Assertions(arg0 *big.Int) (struct {
	StateHash      [32]byte
	InboxSize      *big.Int
	Parent         *big.Int
	Deadline       *big.Int
	ProposalTime   *big.Int
	NumStakers     *big.Int
	ChildInboxSize *big.Int
}, error) {
	return _RollupBase.Contract.Assertions(&_RollupBase.CallOpts, arg0)
}

// Assertions is a free data retrieval call binding the contract method 0x524232f6.
//
// Solidity: function assertions(uint256 ) view returns(bytes32 stateHash, uint256 inboxSize, uint256 parent, uint256 deadline, uint256 proposalTime, uint256 numStakers, uint256 childInboxSize)
func (_RollupBase *RollupBaseCallerSession) Assertions(arg0 *big.Int) (struct {
	StateHash      [32]byte
	InboxSize      *big.Int
	Parent         *big.Int
	Deadline       *big.Int
	ProposalTime   *big.Int
	NumStakers     *big.Int
	ChildInboxSize *big.Int
}, error) {
	return _RollupBase.Contract.Assertions(&_RollupBase.CallOpts, arg0)
}

// BaseStakeAmount is a free data retrieval call binding the contract method 0x71129559.
//
// Solidity: function baseStakeAmount() view returns(uint256)
func (_RollupBase *RollupBaseCaller) BaseStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "baseStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseStakeAmount is a free data retrieval call binding the contract method 0x71129559.
//
// Solidity: function baseStakeAmount() view returns(uint256)
func (_RollupBase *RollupBaseSession) BaseStakeAmount() (*big.Int, error) {
	return _RollupBase.Contract.BaseStakeAmount(&_RollupBase.CallOpts)
}

// BaseStakeAmount is a free data retrieval call binding the contract method 0x71129559.
//
// Solidity: function baseStakeAmount() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) BaseStakeAmount() (*big.Int, error) {
	return _RollupBase.Contract.BaseStakeAmount(&_RollupBase.CallOpts)
}

// ChallengePeriod is a free data retrieval call binding the contract method 0xf3f480d9.
//
// Solidity: function challengePeriod() view returns(uint256)
func (_RollupBase *RollupBaseCaller) ChallengePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "challengePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChallengePeriod is a free data retrieval call binding the contract method 0xf3f480d9.
//
// Solidity: function challengePeriod() view returns(uint256)
func (_RollupBase *RollupBaseSession) ChallengePeriod() (*big.Int, error) {
	return _RollupBase.Contract.ChallengePeriod(&_RollupBase.CallOpts)
}

// ChallengePeriod is a free data retrieval call binding the contract method 0xf3f480d9.
//
// Solidity: function challengePeriod() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) ChallengePeriod() (*big.Int, error) {
	return _RollupBase.Contract.ChallengePeriod(&_RollupBase.CallOpts)
}

// ConfirmationPeriod is a free data retrieval call binding the contract method 0x0429b880.
//
// Solidity: function confirmationPeriod() view returns(uint256)
func (_RollupBase *RollupBaseCaller) ConfirmationPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "confirmationPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConfirmationPeriod is a free data retrieval call binding the contract method 0x0429b880.
//
// Solidity: function confirmationPeriod() view returns(uint256)
func (_RollupBase *RollupBaseSession) ConfirmationPeriod() (*big.Int, error) {
	return _RollupBase.Contract.ConfirmationPeriod(&_RollupBase.CallOpts)
}

// ConfirmationPeriod is a free data retrieval call binding the contract method 0x0429b880.
//
// Solidity: function confirmationPeriod() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) ConfirmationPeriod() (*big.Int, error) {
	return _RollupBase.Contract.ConfirmationPeriod(&_RollupBase.CallOpts)
}

// ConfirmedInboxSize is a free data retrieval call binding the contract method 0xc94b5847.
//
// Solidity: function confirmedInboxSize() view returns(uint256)
func (_RollupBase *RollupBaseCaller) ConfirmedInboxSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "confirmedInboxSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConfirmedInboxSize is a free data retrieval call binding the contract method 0xc94b5847.
//
// Solidity: function confirmedInboxSize() view returns(uint256)
func (_RollupBase *RollupBaseSession) ConfirmedInboxSize() (*big.Int, error) {
	return _RollupBase.Contract.ConfirmedInboxSize(&_RollupBase.CallOpts)
}

// ConfirmedInboxSize is a free data retrieval call binding the contract method 0xc94b5847.
//
// Solidity: function confirmedInboxSize() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) ConfirmedInboxSize() (*big.Int, error) {
	return _RollupBase.Contract.ConfirmedInboxSize(&_RollupBase.CallOpts)
}

// CurrentRequiredStake is a free data retrieval call binding the contract method 0x4d26732d.
//
// Solidity: function currentRequiredStake() view returns(uint256)
func (_RollupBase *RollupBaseCaller) CurrentRequiredStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "currentRequiredStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentRequiredStake is a free data retrieval call binding the contract method 0x4d26732d.
//
// Solidity: function currentRequiredStake() view returns(uint256)
func (_RollupBase *RollupBaseSession) CurrentRequiredStake() (*big.Int, error) {
	return _RollupBase.Contract.CurrentRequiredStake(&_RollupBase.CallOpts)
}

// CurrentRequiredStake is a free data retrieval call binding the contract method 0x4d26732d.
//
// Solidity: function currentRequiredStake() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) CurrentRequiredStake() (*big.Int, error) {
	return _RollupBase.Contract.CurrentRequiredStake(&_RollupBase.CallOpts)
}

// DaProvider is a free data retrieval call binding the contract method 0x8eb8198e.
//
// Solidity: function daProvider() view returns(address)
func (_RollupBase *RollupBaseCaller) DaProvider(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "daProvider")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DaProvider is a free data retrieval call binding the contract method 0x8eb8198e.
//
// Solidity: function daProvider() view returns(address)
func (_RollupBase *RollupBaseSession) DaProvider() (common.Address, error) {
	return _RollupBase.Contract.DaProvider(&_RollupBase.CallOpts)
}

// DaProvider is a free data retrieval call binding the contract method 0x8eb8198e.
//
// Solidity: function daProvider() view returns(address)
func (_RollupBase *RollupBaseCallerSession) DaProvider() (common.Address, error) {
	return _RollupBase.Contract.DaProvider(&_RollupBase.CallOpts)
}

// GetAssertion is a free data retrieval call binding the contract method 0x1d99e167.
//
// Solidity: function getAssertion(uint256 assertionID) view returns((bytes32,uint256,uint256,uint256,uint256,uint256,uint256))
func (_RollupBase *RollupBaseCaller) GetAssertion(opts *bind.CallOpts, assertionID *big.Int) (IRollupAssertion, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "getAssertion", assertionID)

	if err != nil {
		return *new(IRollupAssertion), err
	}

	out0 := *abi.ConvertType(out[0], new(IRollupAssertion)).(*IRollupAssertion)

	return out0, err

}

// GetAssertion is a free data retrieval call binding the contract method 0x1d99e167.
//
// Solidity: function getAssertion(uint256 assertionID) view returns((bytes32,uint256,uint256,uint256,uint256,uint256,uint256))
func (_RollupBase *RollupBaseSession) GetAssertion(assertionID *big.Int) (IRollupAssertion, error) {
	return _RollupBase.Contract.GetAssertion(&_RollupBase.CallOpts, assertionID)
}

// GetAssertion is a free data retrieval call binding the contract method 0x1d99e167.
//
// Solidity: function getAssertion(uint256 assertionID) view returns((bytes32,uint256,uint256,uint256,uint256,uint256,uint256))
func (_RollupBase *RollupBaseCallerSession) GetAssertion(assertionID *big.Int) (IRollupAssertion, error) {
	return _RollupBase.Contract.GetAssertion(&_RollupBase.CallOpts, assertionID)
}

// GetLastConfirmedAssertionID is a free data retrieval call binding the contract method 0x6afcc33c.
//
// Solidity: function getLastConfirmedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCaller) GetLastConfirmedAssertionID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "getLastConfirmedAssertionID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastConfirmedAssertionID is a free data retrieval call binding the contract method 0x6afcc33c.
//
// Solidity: function getLastConfirmedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseSession) GetLastConfirmedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.GetLastConfirmedAssertionID(&_RollupBase.CallOpts)
}

// GetLastConfirmedAssertionID is a free data retrieval call binding the contract method 0x6afcc33c.
//
// Solidity: function getLastConfirmedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) GetLastConfirmedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.GetLastConfirmedAssertionID(&_RollupBase.CallOpts)
}

// GetStaker is a free data retrieval call binding the contract method 0xa23c44b1.
//
// Solidity: function getStaker(address addr) view returns((bool,uint256,uint256,address))
func (_RollupBase *RollupBaseCaller) GetStaker(opts *bind.CallOpts, addr common.Address) (IRollupStaker, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "getStaker", addr)

	if err != nil {
		return *new(IRollupStaker), err
	}

	out0 := *abi.ConvertType(out[0], new(IRollupStaker)).(*IRollupStaker)

	return out0, err

}

// GetStaker is a free data retrieval call binding the contract method 0xa23c44b1.
//
// Solidity: function getStaker(address addr) view returns((bool,uint256,uint256,address))
func (_RollupBase *RollupBaseSession) GetStaker(addr common.Address) (IRollupStaker, error) {
	return _RollupBase.Contract.GetStaker(&_RollupBase.CallOpts, addr)
}

// GetStaker is a free data retrieval call binding the contract method 0xa23c44b1.
//
// Solidity: function getStaker(address addr) view returns((bool,uint256,uint256,address))
func (_RollupBase *RollupBaseCallerSession) GetStaker(addr common.Address) (IRollupStaker, error) {
	return _RollupBase.Contract.GetStaker(&_RollupBase.CallOpts, addr)
}

// IsStakedOnAssertion is a free data retrieval call binding the contract method 0xe58dda89.
//
// Solidity: function isStakedOnAssertion(uint256 assertionID, address stakerAddress) view returns(bool)
func (_RollupBase *RollupBaseCaller) IsStakedOnAssertion(opts *bind.CallOpts, assertionID *big.Int, stakerAddress common.Address) (bool, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "isStakedOnAssertion", assertionID, stakerAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStakedOnAssertion is a free data retrieval call binding the contract method 0xe58dda89.
//
// Solidity: function isStakedOnAssertion(uint256 assertionID, address stakerAddress) view returns(bool)
func (_RollupBase *RollupBaseSession) IsStakedOnAssertion(assertionID *big.Int, stakerAddress common.Address) (bool, error) {
	return _RollupBase.Contract.IsStakedOnAssertion(&_RollupBase.CallOpts, assertionID, stakerAddress)
}

// IsStakedOnAssertion is a free data retrieval call binding the contract method 0xe58dda89.
//
// Solidity: function isStakedOnAssertion(uint256 assertionID, address stakerAddress) view returns(bool)
func (_RollupBase *RollupBaseCallerSession) IsStakedOnAssertion(assertionID *big.Int, stakerAddress common.Address) (bool, error) {
	return _RollupBase.Contract.IsStakedOnAssertion(&_RollupBase.CallOpts, assertionID, stakerAddress)
}

// LastConfirmedAssertionID is a free data retrieval call binding the contract method 0xa56ba93b.
//
// Solidity: function lastConfirmedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCaller) LastConfirmedAssertionID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "lastConfirmedAssertionID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastConfirmedAssertionID is a free data retrieval call binding the contract method 0xa56ba93b.
//
// Solidity: function lastConfirmedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseSession) LastConfirmedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.LastConfirmedAssertionID(&_RollupBase.CallOpts)
}

// LastConfirmedAssertionID is a free data retrieval call binding the contract method 0xa56ba93b.
//
// Solidity: function lastConfirmedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) LastConfirmedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.LastConfirmedAssertionID(&_RollupBase.CallOpts)
}

// LastCreatedAssertionID is a free data retrieval call binding the contract method 0x107035a4.
//
// Solidity: function lastCreatedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCaller) LastCreatedAssertionID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "lastCreatedAssertionID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastCreatedAssertionID is a free data retrieval call binding the contract method 0x107035a4.
//
// Solidity: function lastCreatedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseSession) LastCreatedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.LastCreatedAssertionID(&_RollupBase.CallOpts)
}

// LastCreatedAssertionID is a free data retrieval call binding the contract method 0x107035a4.
//
// Solidity: function lastCreatedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) LastCreatedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.LastCreatedAssertionID(&_RollupBase.CallOpts)
}

// LastResolvedAssertionID is a free data retrieval call binding the contract method 0xb553ee84.
//
// Solidity: function lastResolvedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCaller) LastResolvedAssertionID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "lastResolvedAssertionID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastResolvedAssertionID is a free data retrieval call binding the contract method 0xb553ee84.
//
// Solidity: function lastResolvedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseSession) LastResolvedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.LastResolvedAssertionID(&_RollupBase.CallOpts)
}

// LastResolvedAssertionID is a free data retrieval call binding the contract method 0xb553ee84.
//
// Solidity: function lastResolvedAssertionID() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) LastResolvedAssertionID() (*big.Int, error) {
	return _RollupBase.Contract.LastResolvedAssertionID(&_RollupBase.CallOpts)
}

// MinimumAssertionPeriod is a free data retrieval call binding the contract method 0x45e38b64.
//
// Solidity: function minimumAssertionPeriod() view returns(uint256)
func (_RollupBase *RollupBaseCaller) MinimumAssertionPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "minimumAssertionPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumAssertionPeriod is a free data retrieval call binding the contract method 0x45e38b64.
//
// Solidity: function minimumAssertionPeriod() view returns(uint256)
func (_RollupBase *RollupBaseSession) MinimumAssertionPeriod() (*big.Int, error) {
	return _RollupBase.Contract.MinimumAssertionPeriod(&_RollupBase.CallOpts)
}

// MinimumAssertionPeriod is a free data retrieval call binding the contract method 0x45e38b64.
//
// Solidity: function minimumAssertionPeriod() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) MinimumAssertionPeriod() (*big.Int, error) {
	return _RollupBase.Contract.MinimumAssertionPeriod(&_RollupBase.CallOpts)
}

// NumStakers is a free data retrieval call binding the contract method 0x6c8b052a.
//
// Solidity: function numStakers() view returns(uint256)
func (_RollupBase *RollupBaseCaller) NumStakers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "numStakers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumStakers is a free data retrieval call binding the contract method 0x6c8b052a.
//
// Solidity: function numStakers() view returns(uint256)
func (_RollupBase *RollupBaseSession) NumStakers() (*big.Int, error) {
	return _RollupBase.Contract.NumStakers(&_RollupBase.CallOpts)
}

// NumStakers is a free data retrieval call binding the contract method 0x6c8b052a.
//
// Solidity: function numStakers() view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) NumStakers() (*big.Int, error) {
	return _RollupBase.Contract.NumStakers(&_RollupBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollupBase *RollupBaseCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollupBase *RollupBaseSession) Owner() (common.Address, error) {
	return _RollupBase.Contract.Owner(&_RollupBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollupBase *RollupBaseCallerSession) Owner() (common.Address, error) {
	return _RollupBase.Contract.Owner(&_RollupBase.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RollupBase *RollupBaseCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RollupBase *RollupBaseSession) ProxiableUUID() ([32]byte, error) {
	return _RollupBase.Contract.ProxiableUUID(&_RollupBase.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_RollupBase *RollupBaseCallerSession) ProxiableUUID() ([32]byte, error) {
	return _RollupBase.Contract.ProxiableUUID(&_RollupBase.CallOpts)
}

// RequireFirstUnresolvedAssertionIsConfirmable is a free data retrieval call binding the contract method 0x922a8807.
//
// Solidity: function requireFirstUnresolvedAssertionIsConfirmable() view returns()
func (_RollupBase *RollupBaseCaller) RequireFirstUnresolvedAssertionIsConfirmable(opts *bind.CallOpts) error {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "requireFirstUnresolvedAssertionIsConfirmable")

	if err != nil {
		return err
	}

	return err

}

// RequireFirstUnresolvedAssertionIsConfirmable is a free data retrieval call binding the contract method 0x922a8807.
//
// Solidity: function requireFirstUnresolvedAssertionIsConfirmable() view returns()
func (_RollupBase *RollupBaseSession) RequireFirstUnresolvedAssertionIsConfirmable() error {
	return _RollupBase.Contract.RequireFirstUnresolvedAssertionIsConfirmable(&_RollupBase.CallOpts)
}

// RequireFirstUnresolvedAssertionIsConfirmable is a free data retrieval call binding the contract method 0x922a8807.
//
// Solidity: function requireFirstUnresolvedAssertionIsConfirmable() view returns()
func (_RollupBase *RollupBaseCallerSession) RequireFirstUnresolvedAssertionIsConfirmable() error {
	return _RollupBase.Contract.RequireFirstUnresolvedAssertionIsConfirmable(&_RollupBase.CallOpts)
}

// RequireFirstUnresolvedAssertionIsRejectable is a free data retrieval call binding the contract method 0xc5122403.
//
// Solidity: function requireFirstUnresolvedAssertionIsRejectable(address stakerAddress) view returns()
func (_RollupBase *RollupBaseCaller) RequireFirstUnresolvedAssertionIsRejectable(opts *bind.CallOpts, stakerAddress common.Address) error {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "requireFirstUnresolvedAssertionIsRejectable", stakerAddress)

	if err != nil {
		return err
	}

	return err

}

// RequireFirstUnresolvedAssertionIsRejectable is a free data retrieval call binding the contract method 0xc5122403.
//
// Solidity: function requireFirstUnresolvedAssertionIsRejectable(address stakerAddress) view returns()
func (_RollupBase *RollupBaseSession) RequireFirstUnresolvedAssertionIsRejectable(stakerAddress common.Address) error {
	return _RollupBase.Contract.RequireFirstUnresolvedAssertionIsRejectable(&_RollupBase.CallOpts, stakerAddress)
}

// RequireFirstUnresolvedAssertionIsRejectable is a free data retrieval call binding the contract method 0xc5122403.
//
// Solidity: function requireFirstUnresolvedAssertionIsRejectable(address stakerAddress) view returns()
func (_RollupBase *RollupBaseCallerSession) RequireFirstUnresolvedAssertionIsRejectable(stakerAddress common.Address) error {
	return _RollupBase.Contract.RequireFirstUnresolvedAssertionIsRejectable(&_RollupBase.CallOpts, stakerAddress)
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers(address ) view returns(bool isStaked, uint256 amountStaked, uint256 assertionID, address currentChallenge)
func (_RollupBase *RollupBaseCaller) Stakers(opts *bind.CallOpts, arg0 common.Address) (struct {
	IsStaked         bool
	AmountStaked     *big.Int
	AssertionID      *big.Int
	CurrentChallenge common.Address
}, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "stakers", arg0)

	outstruct := new(struct {
		IsStaked         bool
		AmountStaked     *big.Int
		AssertionID      *big.Int
		CurrentChallenge common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsStaked = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.AmountStaked = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AssertionID = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CurrentChallenge = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers(address ) view returns(bool isStaked, uint256 amountStaked, uint256 assertionID, address currentChallenge)
func (_RollupBase *RollupBaseSession) Stakers(arg0 common.Address) (struct {
	IsStaked         bool
	AmountStaked     *big.Int
	AssertionID      *big.Int
	CurrentChallenge common.Address
}, error) {
	return _RollupBase.Contract.Stakers(&_RollupBase.CallOpts, arg0)
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers(address ) view returns(bool isStaked, uint256 amountStaked, uint256 assertionID, address currentChallenge)
func (_RollupBase *RollupBaseCallerSession) Stakers(arg0 common.Address) (struct {
	IsStaked         bool
	AmountStaked     *big.Int
	AssertionID      *big.Int
	CurrentChallenge common.Address
}, error) {
	return _RollupBase.Contract.Stakers(&_RollupBase.CallOpts, arg0)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_RollupBase *RollupBaseCaller) Vault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "vault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_RollupBase *RollupBaseSession) Vault() (common.Address, error) {
	return _RollupBase.Contract.Vault(&_RollupBase.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_RollupBase *RollupBaseCallerSession) Vault() (common.Address, error) {
	return _RollupBase.Contract.Vault(&_RollupBase.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_RollupBase *RollupBaseCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_RollupBase *RollupBaseSession) Verifier() (common.Address, error) {
	return _RollupBase.Contract.Verifier(&_RollupBase.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_RollupBase *RollupBaseCallerSession) Verifier() (common.Address, error) {
	return _RollupBase.Contract.Verifier(&_RollupBase.CallOpts)
}

// WithdrawableFunds is a free data retrieval call binding the contract method 0x2f30cabd.
//
// Solidity: function withdrawableFunds(address ) view returns(uint256)
func (_RollupBase *RollupBaseCaller) WithdrawableFunds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "withdrawableFunds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawableFunds is a free data retrieval call binding the contract method 0x2f30cabd.
//
// Solidity: function withdrawableFunds(address ) view returns(uint256)
func (_RollupBase *RollupBaseSession) WithdrawableFunds(arg0 common.Address) (*big.Int, error) {
	return _RollupBase.Contract.WithdrawableFunds(&_RollupBase.CallOpts, arg0)
}

// WithdrawableFunds is a free data retrieval call binding the contract method 0x2f30cabd.
//
// Solidity: function withdrawableFunds(address ) view returns(uint256)
func (_RollupBase *RollupBaseCallerSession) WithdrawableFunds(arg0 common.Address) (*big.Int, error) {
	return _RollupBase.Contract.WithdrawableFunds(&_RollupBase.CallOpts, arg0)
}

// Zombies is a free data retrieval call binding the contract method 0x2052465e.
//
// Solidity: function zombies(uint256 ) view returns(address stakerAddress, uint256 lastAssertionID)
func (_RollupBase *RollupBaseCaller) Zombies(opts *bind.CallOpts, arg0 *big.Int) (struct {
	StakerAddress   common.Address
	LastAssertionID *big.Int
}, error) {
	var out []interface{}
	err := _RollupBase.contract.Call(opts, &out, "zombies", arg0)

	outstruct := new(struct {
		StakerAddress   common.Address
		LastAssertionID *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakerAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.LastAssertionID = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Zombies is a free data retrieval call binding the contract method 0x2052465e.
//
// Solidity: function zombies(uint256 ) view returns(address stakerAddress, uint256 lastAssertionID)
func (_RollupBase *RollupBaseSession) Zombies(arg0 *big.Int) (struct {
	StakerAddress   common.Address
	LastAssertionID *big.Int
}, error) {
	return _RollupBase.Contract.Zombies(&_RollupBase.CallOpts, arg0)
}

// Zombies is a free data retrieval call binding the contract method 0x2052465e.
//
// Solidity: function zombies(uint256 ) view returns(address stakerAddress, uint256 lastAssertionID)
func (_RollupBase *RollupBaseCallerSession) Zombies(arg0 *big.Int) (struct {
	StakerAddress   common.Address
	LastAssertionID *big.Int
}, error) {
	return _RollupBase.Contract.Zombies(&_RollupBase.CallOpts, arg0)
}

// AdvanceStake is a paid mutator transaction binding the contract method 0x8821b2ae.
//
// Solidity: function advanceStake(uint256 assertionID) returns()
func (_RollupBase *RollupBaseTransactor) AdvanceStake(opts *bind.TransactOpts, assertionID *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "advanceStake", assertionID)
}

// AdvanceStake is a paid mutator transaction binding the contract method 0x8821b2ae.
//
// Solidity: function advanceStake(uint256 assertionID) returns()
func (_RollupBase *RollupBaseSession) AdvanceStake(assertionID *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.AdvanceStake(&_RollupBase.TransactOpts, assertionID)
}

// AdvanceStake is a paid mutator transaction binding the contract method 0x8821b2ae.
//
// Solidity: function advanceStake(uint256 assertionID) returns()
func (_RollupBase *RollupBaseTransactorSession) AdvanceStake(assertionID *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.AdvanceStake(&_RollupBase.TransactOpts, assertionID)
}

// ChallengeAssertion is a paid mutator transaction binding the contract method 0x2f06d1b0.
//
// Solidity: function challengeAssertion(address[2] players, uint256[2] assertionIDs) returns(address)
func (_RollupBase *RollupBaseTransactor) ChallengeAssertion(opts *bind.TransactOpts, players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "challengeAssertion", players, assertionIDs)
}

// ChallengeAssertion is a paid mutator transaction binding the contract method 0x2f06d1b0.
//
// Solidity: function challengeAssertion(address[2] players, uint256[2] assertionIDs) returns(address)
func (_RollupBase *RollupBaseSession) ChallengeAssertion(players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.ChallengeAssertion(&_RollupBase.TransactOpts, players, assertionIDs)
}

// ChallengeAssertion is a paid mutator transaction binding the contract method 0x2f06d1b0.
//
// Solidity: function challengeAssertion(address[2] players, uint256[2] assertionIDs) returns(address)
func (_RollupBase *RollupBaseTransactorSession) ChallengeAssertion(players [2]common.Address, assertionIDs [2]*big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.ChallengeAssertion(&_RollupBase.TransactOpts, players, assertionIDs)
}

// CompleteChallenge is a paid mutator transaction binding the contract method 0xfa7803e6.
//
// Solidity: function completeChallenge(address winner, address loser) returns()
func (_RollupBase *RollupBaseTransactor) CompleteChallenge(opts *bind.TransactOpts, winner common.Address, loser common.Address) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "completeChallenge", winner, loser)
}

// CompleteChallenge is a paid mutator transaction binding the contract method 0xfa7803e6.
//
// Solidity: function completeChallenge(address winner, address loser) returns()
func (_RollupBase *RollupBaseSession) CompleteChallenge(winner common.Address, loser common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.CompleteChallenge(&_RollupBase.TransactOpts, winner, loser)
}

// CompleteChallenge is a paid mutator transaction binding the contract method 0xfa7803e6.
//
// Solidity: function completeChallenge(address winner, address loser) returns()
func (_RollupBase *RollupBaseTransactorSession) CompleteChallenge(winner common.Address, loser common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.CompleteChallenge(&_RollupBase.TransactOpts, winner, loser)
}

// ConfirmFirstUnresolvedAssertion is a paid mutator transaction binding the contract method 0x2906040e.
//
// Solidity: function confirmFirstUnresolvedAssertion() returns()
func (_RollupBase *RollupBaseTransactor) ConfirmFirstUnresolvedAssertion(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "confirmFirstUnresolvedAssertion")
}

// ConfirmFirstUnresolvedAssertion is a paid mutator transaction binding the contract method 0x2906040e.
//
// Solidity: function confirmFirstUnresolvedAssertion() returns()
func (_RollupBase *RollupBaseSession) ConfirmFirstUnresolvedAssertion() (*types.Transaction, error) {
	return _RollupBase.Contract.ConfirmFirstUnresolvedAssertion(&_RollupBase.TransactOpts)
}

// ConfirmFirstUnresolvedAssertion is a paid mutator transaction binding the contract method 0x2906040e.
//
// Solidity: function confirmFirstUnresolvedAssertion() returns()
func (_RollupBase *RollupBaseTransactorSession) ConfirmFirstUnresolvedAssertion() (*types.Transaction, error) {
	return _RollupBase.Contract.ConfirmFirstUnresolvedAssertion(&_RollupBase.TransactOpts)
}

// CreateAssertion is a paid mutator transaction binding the contract method 0xb6da898f.
//
// Solidity: function createAssertion(bytes32 vmHash, uint256 inboxSize) returns()
func (_RollupBase *RollupBaseTransactor) CreateAssertion(opts *bind.TransactOpts, vmHash [32]byte, inboxSize *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "createAssertion", vmHash, inboxSize)
}

// CreateAssertion is a paid mutator transaction binding the contract method 0xb6da898f.
//
// Solidity: function createAssertion(bytes32 vmHash, uint256 inboxSize) returns()
func (_RollupBase *RollupBaseSession) CreateAssertion(vmHash [32]byte, inboxSize *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.CreateAssertion(&_RollupBase.TransactOpts, vmHash, inboxSize)
}

// CreateAssertion is a paid mutator transaction binding the contract method 0xb6da898f.
//
// Solidity: function createAssertion(bytes32 vmHash, uint256 inboxSize) returns()
func (_RollupBase *RollupBaseTransactorSession) CreateAssertion(vmHash [32]byte, inboxSize *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.CreateAssertion(&_RollupBase.TransactOpts, vmHash, inboxSize)
}

// Initialize is a paid mutator transaction binding the contract method 0xfb1b3337.
//
// Solidity: function initialize(address _vault, address _daProvider, address _verifier, uint256 _confirmationPeriod, uint256 _challengePeriod, uint256 _minimumAssertionPeriod, uint256 _baseStakeAmount, uint256 _initialAssertionID, uint256 _initialInboxSize, bytes32 _initialVMhash) returns()
func (_RollupBase *RollupBaseTransactor) Initialize(opts *bind.TransactOpts, _vault common.Address, _daProvider common.Address, _verifier common.Address, _confirmationPeriod *big.Int, _challengePeriod *big.Int, _minimumAssertionPeriod *big.Int, _baseStakeAmount *big.Int, _initialAssertionID *big.Int, _initialInboxSize *big.Int, _initialVMhash [32]byte) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "initialize", _vault, _daProvider, _verifier, _confirmationPeriod, _challengePeriod, _minimumAssertionPeriod, _baseStakeAmount, _initialAssertionID, _initialInboxSize, _initialVMhash)
}

// Initialize is a paid mutator transaction binding the contract method 0xfb1b3337.
//
// Solidity: function initialize(address _vault, address _daProvider, address _verifier, uint256 _confirmationPeriod, uint256 _challengePeriod, uint256 _minimumAssertionPeriod, uint256 _baseStakeAmount, uint256 _initialAssertionID, uint256 _initialInboxSize, bytes32 _initialVMhash) returns()
func (_RollupBase *RollupBaseSession) Initialize(_vault common.Address, _daProvider common.Address, _verifier common.Address, _confirmationPeriod *big.Int, _challengePeriod *big.Int, _minimumAssertionPeriod *big.Int, _baseStakeAmount *big.Int, _initialAssertionID *big.Int, _initialInboxSize *big.Int, _initialVMhash [32]byte) (*types.Transaction, error) {
	return _RollupBase.Contract.Initialize(&_RollupBase.TransactOpts, _vault, _daProvider, _verifier, _confirmationPeriod, _challengePeriod, _minimumAssertionPeriod, _baseStakeAmount, _initialAssertionID, _initialInboxSize, _initialVMhash)
}

// Initialize is a paid mutator transaction binding the contract method 0xfb1b3337.
//
// Solidity: function initialize(address _vault, address _daProvider, address _verifier, uint256 _confirmationPeriod, uint256 _challengePeriod, uint256 _minimumAssertionPeriod, uint256 _baseStakeAmount, uint256 _initialAssertionID, uint256 _initialInboxSize, bytes32 _initialVMhash) returns()
func (_RollupBase *RollupBaseTransactorSession) Initialize(_vault common.Address, _daProvider common.Address, _verifier common.Address, _confirmationPeriod *big.Int, _challengePeriod *big.Int, _minimumAssertionPeriod *big.Int, _baseStakeAmount *big.Int, _initialAssertionID *big.Int, _initialInboxSize *big.Int, _initialVMhash [32]byte) (*types.Transaction, error) {
	return _RollupBase.Contract.Initialize(&_RollupBase.TransactOpts, _vault, _daProvider, _verifier, _confirmationPeriod, _challengePeriod, _minimumAssertionPeriod, _baseStakeAmount, _initialAssertionID, _initialInboxSize, _initialVMhash)
}

// RejectFirstUnresolvedAssertion is a paid mutator transaction binding the contract method 0x042dca93.
//
// Solidity: function rejectFirstUnresolvedAssertion(address stakerAddress) returns()
func (_RollupBase *RollupBaseTransactor) RejectFirstUnresolvedAssertion(opts *bind.TransactOpts, stakerAddress common.Address) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "rejectFirstUnresolvedAssertion", stakerAddress)
}

// RejectFirstUnresolvedAssertion is a paid mutator transaction binding the contract method 0x042dca93.
//
// Solidity: function rejectFirstUnresolvedAssertion(address stakerAddress) returns()
func (_RollupBase *RollupBaseSession) RejectFirstUnresolvedAssertion(stakerAddress common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.RejectFirstUnresolvedAssertion(&_RollupBase.TransactOpts, stakerAddress)
}

// RejectFirstUnresolvedAssertion is a paid mutator transaction binding the contract method 0x042dca93.
//
// Solidity: function rejectFirstUnresolvedAssertion(address stakerAddress) returns()
func (_RollupBase *RollupBaseTransactorSession) RejectFirstUnresolvedAssertion(stakerAddress common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.RejectFirstUnresolvedAssertion(&_RollupBase.TransactOpts, stakerAddress)
}

// RemoveStake is a paid mutator transaction binding the contract method 0xfe2ba848.
//
// Solidity: function removeStake(address stakerAddress) returns()
func (_RollupBase *RollupBaseTransactor) RemoveStake(opts *bind.TransactOpts, stakerAddress common.Address) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "removeStake", stakerAddress)
}

// RemoveStake is a paid mutator transaction binding the contract method 0xfe2ba848.
//
// Solidity: function removeStake(address stakerAddress) returns()
func (_RollupBase *RollupBaseSession) RemoveStake(stakerAddress common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.RemoveStake(&_RollupBase.TransactOpts, stakerAddress)
}

// RemoveStake is a paid mutator transaction binding the contract method 0xfe2ba848.
//
// Solidity: function removeStake(address stakerAddress) returns()
func (_RollupBase *RollupBaseTransactorSession) RemoveStake(stakerAddress common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.RemoveStake(&_RollupBase.TransactOpts, stakerAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollupBase *RollupBaseTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollupBase *RollupBaseSession) RenounceOwnership() (*types.Transaction, error) {
	return _RollupBase.Contract.RenounceOwnership(&_RollupBase.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollupBase *RollupBaseTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RollupBase.Contract.RenounceOwnership(&_RollupBase.TransactOpts)
}

// SetBaseStakeAmount is a paid mutator transaction binding the contract method 0x3986e6fc.
//
// Solidity: function setBaseStakeAmount(uint256 newAmount) returns()
func (_RollupBase *RollupBaseTransactor) SetBaseStakeAmount(opts *bind.TransactOpts, newAmount *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "setBaseStakeAmount", newAmount)
}

// SetBaseStakeAmount is a paid mutator transaction binding the contract method 0x3986e6fc.
//
// Solidity: function setBaseStakeAmount(uint256 newAmount) returns()
func (_RollupBase *RollupBaseSession) SetBaseStakeAmount(newAmount *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetBaseStakeAmount(&_RollupBase.TransactOpts, newAmount)
}

// SetBaseStakeAmount is a paid mutator transaction binding the contract method 0x3986e6fc.
//
// Solidity: function setBaseStakeAmount(uint256 newAmount) returns()
func (_RollupBase *RollupBaseTransactorSession) SetBaseStakeAmount(newAmount *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetBaseStakeAmount(&_RollupBase.TransactOpts, newAmount)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseTransactor) SetChallengePeriod(opts *bind.TransactOpts, newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "setChallengePeriod", newPeriod)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseSession) SetChallengePeriod(newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetChallengePeriod(&_RollupBase.TransactOpts, newPeriod)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseTransactorSession) SetChallengePeriod(newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetChallengePeriod(&_RollupBase.TransactOpts, newPeriod)
}

// SetConfirmationPeriod is a paid mutator transaction binding the contract method 0xbea50ae3.
//
// Solidity: function setConfirmationPeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseTransactor) SetConfirmationPeriod(opts *bind.TransactOpts, newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "setConfirmationPeriod", newPeriod)
}

// SetConfirmationPeriod is a paid mutator transaction binding the contract method 0xbea50ae3.
//
// Solidity: function setConfirmationPeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseSession) SetConfirmationPeriod(newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetConfirmationPeriod(&_RollupBase.TransactOpts, newPeriod)
}

// SetConfirmationPeriod is a paid mutator transaction binding the contract method 0xbea50ae3.
//
// Solidity: function setConfirmationPeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseTransactorSession) SetConfirmationPeriod(newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetConfirmationPeriod(&_RollupBase.TransactOpts, newPeriod)
}

// SetDAProvider is a paid mutator transaction binding the contract method 0xf397e38e.
//
// Solidity: function setDAProvider(address newDAProvider) returns()
func (_RollupBase *RollupBaseTransactor) SetDAProvider(opts *bind.TransactOpts, newDAProvider common.Address) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "setDAProvider", newDAProvider)
}

// SetDAProvider is a paid mutator transaction binding the contract method 0xf397e38e.
//
// Solidity: function setDAProvider(address newDAProvider) returns()
func (_RollupBase *RollupBaseSession) SetDAProvider(newDAProvider common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.SetDAProvider(&_RollupBase.TransactOpts, newDAProvider)
}

// SetDAProvider is a paid mutator transaction binding the contract method 0xf397e38e.
//
// Solidity: function setDAProvider(address newDAProvider) returns()
func (_RollupBase *RollupBaseTransactorSession) SetDAProvider(newDAProvider common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.SetDAProvider(&_RollupBase.TransactOpts, newDAProvider)
}

// SetMinimumAssertionPeriod is a paid mutator transaction binding the contract method 0x948d6588.
//
// Solidity: function setMinimumAssertionPeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseTransactor) SetMinimumAssertionPeriod(opts *bind.TransactOpts, newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "setMinimumAssertionPeriod", newPeriod)
}

// SetMinimumAssertionPeriod is a paid mutator transaction binding the contract method 0x948d6588.
//
// Solidity: function setMinimumAssertionPeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseSession) SetMinimumAssertionPeriod(newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetMinimumAssertionPeriod(&_RollupBase.TransactOpts, newPeriod)
}

// SetMinimumAssertionPeriod is a paid mutator transaction binding the contract method 0x948d6588.
//
// Solidity: function setMinimumAssertionPeriod(uint256 newPeriod) returns()
func (_RollupBase *RollupBaseTransactorSession) SetMinimumAssertionPeriod(newPeriod *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.SetMinimumAssertionPeriod(&_RollupBase.TransactOpts, newPeriod)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_RollupBase *RollupBaseTransactor) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "stake")
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_RollupBase *RollupBaseSession) Stake() (*types.Transaction, error) {
	return _RollupBase.Contract.Stake(&_RollupBase.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_RollupBase *RollupBaseTransactorSession) Stake() (*types.Transaction, error) {
	return _RollupBase.Contract.Stake(&_RollupBase.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollupBase *RollupBaseTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollupBase *RollupBaseSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.TransferOwnership(&_RollupBase.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollupBase *RollupBaseTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.TransferOwnership(&_RollupBase.TransactOpts, newOwner)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 stakeAmount) returns()
func (_RollupBase *RollupBaseTransactor) Unstake(opts *bind.TransactOpts, stakeAmount *big.Int) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "unstake", stakeAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 stakeAmount) returns()
func (_RollupBase *RollupBaseSession) Unstake(stakeAmount *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.Unstake(&_RollupBase.TransactOpts, stakeAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 stakeAmount) returns()
func (_RollupBase *RollupBaseTransactorSession) Unstake(stakeAmount *big.Int) (*types.Transaction, error) {
	return _RollupBase.Contract.Unstake(&_RollupBase.TransactOpts, stakeAmount)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_RollupBase *RollupBaseTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_RollupBase *RollupBaseSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.UpgradeTo(&_RollupBase.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_RollupBase *RollupBaseTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _RollupBase.Contract.UpgradeTo(&_RollupBase.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RollupBase *RollupBaseTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RollupBase *RollupBaseSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RollupBase.Contract.UpgradeToAndCall(&_RollupBase.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_RollupBase *RollupBaseTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _RollupBase.Contract.UpgradeToAndCall(&_RollupBase.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_RollupBase *RollupBaseTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupBase.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_RollupBase *RollupBaseSession) Withdraw() (*types.Transaction, error) {
	return _RollupBase.Contract.Withdraw(&_RollupBase.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_RollupBase *RollupBaseTransactorSession) Withdraw() (*types.Transaction, error) {
	return _RollupBase.Contract.Withdraw(&_RollupBase.TransactOpts)
}

// RollupBaseAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the RollupBase contract.
type RollupBaseAdminChangedIterator struct {
	Event *RollupBaseAdminChanged // Event containing the contract specifics and raw log

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
func (it *RollupBaseAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseAdminChanged)
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
		it.Event = new(RollupBaseAdminChanged)
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
func (it *RollupBaseAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseAdminChanged represents a AdminChanged event raised by the RollupBase contract.
type RollupBaseAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_RollupBase *RollupBaseFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*RollupBaseAdminChangedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &RollupBaseAdminChangedIterator{contract: _RollupBase.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_RollupBase *RollupBaseFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *RollupBaseAdminChanged) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseAdminChanged)
				if err := _RollupBase.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_RollupBase *RollupBaseFilterer) ParseAdminChanged(log types.Log) (*RollupBaseAdminChanged, error) {
	event := new(RollupBaseAdminChanged)
	if err := _RollupBase.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseAssertionChallengedIterator is returned from FilterAssertionChallenged and is used to iterate over the raw logs and unpacked data for AssertionChallenged events raised by the RollupBase contract.
type RollupBaseAssertionChallengedIterator struct {
	Event *RollupBaseAssertionChallenged // Event containing the contract specifics and raw log

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
func (it *RollupBaseAssertionChallengedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseAssertionChallenged)
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
		it.Event = new(RollupBaseAssertionChallenged)
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
func (it *RollupBaseAssertionChallengedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseAssertionChallengedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseAssertionChallenged represents a AssertionChallenged event raised by the RollupBase contract.
type RollupBaseAssertionChallenged struct {
	AssertionID   *big.Int
	ChallengeAddr common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAssertionChallenged is a free log retrieval operation binding the contract event 0xd0ebe74b4f7d89a9b0fdc9d95f887a7b925c6c7300b5c4b2c3304d97925840fa.
//
// Solidity: event AssertionChallenged(uint256 assertionID, address challengeAddr)
func (_RollupBase *RollupBaseFilterer) FilterAssertionChallenged(opts *bind.FilterOpts) (*RollupBaseAssertionChallengedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "AssertionChallenged")
	if err != nil {
		return nil, err
	}
	return &RollupBaseAssertionChallengedIterator{contract: _RollupBase.contract, event: "AssertionChallenged", logs: logs, sub: sub}, nil
}

// WatchAssertionChallenged is a free log subscription operation binding the contract event 0xd0ebe74b4f7d89a9b0fdc9d95f887a7b925c6c7300b5c4b2c3304d97925840fa.
//
// Solidity: event AssertionChallenged(uint256 assertionID, address challengeAddr)
func (_RollupBase *RollupBaseFilterer) WatchAssertionChallenged(opts *bind.WatchOpts, sink chan<- *RollupBaseAssertionChallenged) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "AssertionChallenged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseAssertionChallenged)
				if err := _RollupBase.contract.UnpackLog(event, "AssertionChallenged", log); err != nil {
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

// ParseAssertionChallenged is a log parse operation binding the contract event 0xd0ebe74b4f7d89a9b0fdc9d95f887a7b925c6c7300b5c4b2c3304d97925840fa.
//
// Solidity: event AssertionChallenged(uint256 assertionID, address challengeAddr)
func (_RollupBase *RollupBaseFilterer) ParseAssertionChallenged(log types.Log) (*RollupBaseAssertionChallenged, error) {
	event := new(RollupBaseAssertionChallenged)
	if err := _RollupBase.contract.UnpackLog(event, "AssertionChallenged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseAssertionConfirmedIterator is returned from FilterAssertionConfirmed and is used to iterate over the raw logs and unpacked data for AssertionConfirmed events raised by the RollupBase contract.
type RollupBaseAssertionConfirmedIterator struct {
	Event *RollupBaseAssertionConfirmed // Event containing the contract specifics and raw log

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
func (it *RollupBaseAssertionConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseAssertionConfirmed)
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
		it.Event = new(RollupBaseAssertionConfirmed)
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
func (it *RollupBaseAssertionConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseAssertionConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseAssertionConfirmed represents a AssertionConfirmed event raised by the RollupBase contract.
type RollupBaseAssertionConfirmed struct {
	AssertionID *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssertionConfirmed is a free log retrieval operation binding the contract event 0x453430d123684340024ae0a229704bdab39c93dc48bb5a0b4bc83142d95d48ef.
//
// Solidity: event AssertionConfirmed(uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) FilterAssertionConfirmed(opts *bind.FilterOpts) (*RollupBaseAssertionConfirmedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "AssertionConfirmed")
	if err != nil {
		return nil, err
	}
	return &RollupBaseAssertionConfirmedIterator{contract: _RollupBase.contract, event: "AssertionConfirmed", logs: logs, sub: sub}, nil
}

// WatchAssertionConfirmed is a free log subscription operation binding the contract event 0x453430d123684340024ae0a229704bdab39c93dc48bb5a0b4bc83142d95d48ef.
//
// Solidity: event AssertionConfirmed(uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) WatchAssertionConfirmed(opts *bind.WatchOpts, sink chan<- *RollupBaseAssertionConfirmed) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "AssertionConfirmed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseAssertionConfirmed)
				if err := _RollupBase.contract.UnpackLog(event, "AssertionConfirmed", log); err != nil {
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

// ParseAssertionConfirmed is a log parse operation binding the contract event 0x453430d123684340024ae0a229704bdab39c93dc48bb5a0b4bc83142d95d48ef.
//
// Solidity: event AssertionConfirmed(uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) ParseAssertionConfirmed(log types.Log) (*RollupBaseAssertionConfirmed, error) {
	event := new(RollupBaseAssertionConfirmed)
	if err := _RollupBase.contract.UnpackLog(event, "AssertionConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseAssertionCreatedIterator is returned from FilterAssertionCreated and is used to iterate over the raw logs and unpacked data for AssertionCreated events raised by the RollupBase contract.
type RollupBaseAssertionCreatedIterator struct {
	Event *RollupBaseAssertionCreated // Event containing the contract specifics and raw log

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
func (it *RollupBaseAssertionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseAssertionCreated)
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
		it.Event = new(RollupBaseAssertionCreated)
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
func (it *RollupBaseAssertionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseAssertionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseAssertionCreated represents a AssertionCreated event raised by the RollupBase contract.
type RollupBaseAssertionCreated struct {
	AssertionID  *big.Int
	AsserterAddr common.Address
	VmHash       [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAssertionCreated is a free log retrieval operation binding the contract event 0xf41917cc5ddc34dc57b3ea71e866801af6a254bddeadaffd1177ad8e46cb0d6b.
//
// Solidity: event AssertionCreated(uint256 assertionID, address asserterAddr, bytes32 vmHash)
func (_RollupBase *RollupBaseFilterer) FilterAssertionCreated(opts *bind.FilterOpts) (*RollupBaseAssertionCreatedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "AssertionCreated")
	if err != nil {
		return nil, err
	}
	return &RollupBaseAssertionCreatedIterator{contract: _RollupBase.contract, event: "AssertionCreated", logs: logs, sub: sub}, nil
}

// WatchAssertionCreated is a free log subscription operation binding the contract event 0xf41917cc5ddc34dc57b3ea71e866801af6a254bddeadaffd1177ad8e46cb0d6b.
//
// Solidity: event AssertionCreated(uint256 assertionID, address asserterAddr, bytes32 vmHash)
func (_RollupBase *RollupBaseFilterer) WatchAssertionCreated(opts *bind.WatchOpts, sink chan<- *RollupBaseAssertionCreated) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "AssertionCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseAssertionCreated)
				if err := _RollupBase.contract.UnpackLog(event, "AssertionCreated", log); err != nil {
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

// ParseAssertionCreated is a log parse operation binding the contract event 0xf41917cc5ddc34dc57b3ea71e866801af6a254bddeadaffd1177ad8e46cb0d6b.
//
// Solidity: event AssertionCreated(uint256 assertionID, address asserterAddr, bytes32 vmHash)
func (_RollupBase *RollupBaseFilterer) ParseAssertionCreated(log types.Log) (*RollupBaseAssertionCreated, error) {
	event := new(RollupBaseAssertionCreated)
	if err := _RollupBase.contract.UnpackLog(event, "AssertionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseAssertionRejectedIterator is returned from FilterAssertionRejected and is used to iterate over the raw logs and unpacked data for AssertionRejected events raised by the RollupBase contract.
type RollupBaseAssertionRejectedIterator struct {
	Event *RollupBaseAssertionRejected // Event containing the contract specifics and raw log

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
func (it *RollupBaseAssertionRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseAssertionRejected)
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
		it.Event = new(RollupBaseAssertionRejected)
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
func (it *RollupBaseAssertionRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseAssertionRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseAssertionRejected represents a AssertionRejected event raised by the RollupBase contract.
type RollupBaseAssertionRejected struct {
	AssertionID *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssertionRejected is a free log retrieval operation binding the contract event 0x5b24ab8ceb442373727ac5c559a027521cb52db451c74710ebed9faa5fe15a7c.
//
// Solidity: event AssertionRejected(uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) FilterAssertionRejected(opts *bind.FilterOpts) (*RollupBaseAssertionRejectedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "AssertionRejected")
	if err != nil {
		return nil, err
	}
	return &RollupBaseAssertionRejectedIterator{contract: _RollupBase.contract, event: "AssertionRejected", logs: logs, sub: sub}, nil
}

// WatchAssertionRejected is a free log subscription operation binding the contract event 0x5b24ab8ceb442373727ac5c559a027521cb52db451c74710ebed9faa5fe15a7c.
//
// Solidity: event AssertionRejected(uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) WatchAssertionRejected(opts *bind.WatchOpts, sink chan<- *RollupBaseAssertionRejected) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "AssertionRejected")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseAssertionRejected)
				if err := _RollupBase.contract.UnpackLog(event, "AssertionRejected", log); err != nil {
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

// ParseAssertionRejected is a log parse operation binding the contract event 0x5b24ab8ceb442373727ac5c559a027521cb52db451c74710ebed9faa5fe15a7c.
//
// Solidity: event AssertionRejected(uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) ParseAssertionRejected(log types.Log) (*RollupBaseAssertionRejected, error) {
	event := new(RollupBaseAssertionRejected)
	if err := _RollupBase.contract.UnpackLog(event, "AssertionRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the RollupBase contract.
type RollupBaseBeaconUpgradedIterator struct {
	Event *RollupBaseBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *RollupBaseBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseBeaconUpgraded)
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
		it.Event = new(RollupBaseBeaconUpgraded)
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
func (it *RollupBaseBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseBeaconUpgraded represents a BeaconUpgraded event raised by the RollupBase contract.
type RollupBaseBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_RollupBase *RollupBaseFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*RollupBaseBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &RollupBaseBeaconUpgradedIterator{contract: _RollupBase.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_RollupBase *RollupBaseFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *RollupBaseBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseBeaconUpgraded)
				if err := _RollupBase.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_RollupBase *RollupBaseFilterer) ParseBeaconUpgraded(log types.Log) (*RollupBaseBeaconUpgraded, error) {
	event := new(RollupBaseBeaconUpgraded)
	if err := _RollupBase.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseConfigChangedIterator is returned from FilterConfigChanged and is used to iterate over the raw logs and unpacked data for ConfigChanged events raised by the RollupBase contract.
type RollupBaseConfigChangedIterator struct {
	Event *RollupBaseConfigChanged // Event containing the contract specifics and raw log

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
func (it *RollupBaseConfigChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseConfigChanged)
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
		it.Event = new(RollupBaseConfigChanged)
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
func (it *RollupBaseConfigChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseConfigChanged represents a ConfigChanged event raised by the RollupBase contract.
type RollupBaseConfigChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterConfigChanged is a free log retrieval operation binding the contract event 0xb9b6902016bd1219d5fa6161243b61e7e9f7f959526dd94ef8fa3e403bf881c3.
//
// Solidity: event ConfigChanged()
func (_RollupBase *RollupBaseFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*RollupBaseConfigChangedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &RollupBaseConfigChangedIterator{contract: _RollupBase.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

// WatchConfigChanged is a free log subscription operation binding the contract event 0xb9b6902016bd1219d5fa6161243b61e7e9f7f959526dd94ef8fa3e403bf881c3.
//
// Solidity: event ConfigChanged()
func (_RollupBase *RollupBaseFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *RollupBaseConfigChanged) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseConfigChanged)
				if err := _RollupBase.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

// ParseConfigChanged is a log parse operation binding the contract event 0xb9b6902016bd1219d5fa6161243b61e7e9f7f959526dd94ef8fa3e403bf881c3.
//
// Solidity: event ConfigChanged()
func (_RollupBase *RollupBaseFilterer) ParseConfigChanged(log types.Log) (*RollupBaseConfigChanged, error) {
	event := new(RollupBaseConfigChanged)
	if err := _RollupBase.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RollupBase contract.
type RollupBaseInitializedIterator struct {
	Event *RollupBaseInitialized // Event containing the contract specifics and raw log

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
func (it *RollupBaseInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseInitialized)
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
		it.Event = new(RollupBaseInitialized)
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
func (it *RollupBaseInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseInitialized represents a Initialized event raised by the RollupBase contract.
type RollupBaseInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RollupBase *RollupBaseFilterer) FilterInitialized(opts *bind.FilterOpts) (*RollupBaseInitializedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RollupBaseInitializedIterator{contract: _RollupBase.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_RollupBase *RollupBaseFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RollupBaseInitialized) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseInitialized)
				if err := _RollupBase.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_RollupBase *RollupBaseFilterer) ParseInitialized(log types.Log) (*RollupBaseInitialized, error) {
	event := new(RollupBaseInitialized)
	if err := _RollupBase.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RollupBase contract.
type RollupBaseOwnershipTransferredIterator struct {
	Event *RollupBaseOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RollupBaseOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseOwnershipTransferred)
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
		it.Event = new(RollupBaseOwnershipTransferred)
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
func (it *RollupBaseOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseOwnershipTransferred represents a OwnershipTransferred event raised by the RollupBase contract.
type RollupBaseOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollupBase *RollupBaseFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RollupBaseOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RollupBaseOwnershipTransferredIterator{contract: _RollupBase.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollupBase *RollupBaseFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RollupBaseOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseOwnershipTransferred)
				if err := _RollupBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_RollupBase *RollupBaseFilterer) ParseOwnershipTransferred(log types.Log) (*RollupBaseOwnershipTransferred, error) {
	event := new(RollupBaseOwnershipTransferred)
	if err := _RollupBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseStakerStakedIterator is returned from FilterStakerStaked and is used to iterate over the raw logs and unpacked data for StakerStaked events raised by the RollupBase contract.
type RollupBaseStakerStakedIterator struct {
	Event *RollupBaseStakerStaked // Event containing the contract specifics and raw log

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
func (it *RollupBaseStakerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseStakerStaked)
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
		it.Event = new(RollupBaseStakerStaked)
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
func (it *RollupBaseStakerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseStakerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseStakerStaked represents a StakerStaked event raised by the RollupBase contract.
type RollupBaseStakerStaked struct {
	StakerAddr  common.Address
	AssertionID *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakerStaked is a free log retrieval operation binding the contract event 0x617d31491414a4ab2bd831e566a31837fa7fb6582921c91dffbbe83fbca789f3.
//
// Solidity: event StakerStaked(address stakerAddr, uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) FilterStakerStaked(opts *bind.FilterOpts) (*RollupBaseStakerStakedIterator, error) {

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "StakerStaked")
	if err != nil {
		return nil, err
	}
	return &RollupBaseStakerStakedIterator{contract: _RollupBase.contract, event: "StakerStaked", logs: logs, sub: sub}, nil
}

// WatchStakerStaked is a free log subscription operation binding the contract event 0x617d31491414a4ab2bd831e566a31837fa7fb6582921c91dffbbe83fbca789f3.
//
// Solidity: event StakerStaked(address stakerAddr, uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) WatchStakerStaked(opts *bind.WatchOpts, sink chan<- *RollupBaseStakerStaked) (event.Subscription, error) {

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "StakerStaked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseStakerStaked)
				if err := _RollupBase.contract.UnpackLog(event, "StakerStaked", log); err != nil {
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

// ParseStakerStaked is a log parse operation binding the contract event 0x617d31491414a4ab2bd831e566a31837fa7fb6582921c91dffbbe83fbca789f3.
//
// Solidity: event StakerStaked(address stakerAddr, uint256 assertionID)
func (_RollupBase *RollupBaseFilterer) ParseStakerStaked(log types.Log) (*RollupBaseStakerStaked, error) {
	event := new(RollupBaseStakerStaked)
	if err := _RollupBase.contract.UnpackLog(event, "StakerStaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupBaseUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the RollupBase contract.
type RollupBaseUpgradedIterator struct {
	Event *RollupBaseUpgraded // Event containing the contract specifics and raw log

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
func (it *RollupBaseUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupBaseUpgraded)
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
		it.Event = new(RollupBaseUpgraded)
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
func (it *RollupBaseUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupBaseUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupBaseUpgraded represents a Upgraded event raised by the RollupBase contract.
type RollupBaseUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RollupBase *RollupBaseFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RollupBaseUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RollupBase.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RollupBaseUpgradedIterator{contract: _RollupBase.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_RollupBase *RollupBaseFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RollupBaseUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _RollupBase.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupBaseUpgraded)
				if err := _RollupBase.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_RollupBase *RollupBaseFilterer) ParseUpgraded(log types.Log) (*RollupBaseUpgraded, error) {
	event := new(RollupBaseUpgraded)
	if err := _RollupBase.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
