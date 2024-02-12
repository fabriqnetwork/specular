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

// L1OracleMetaData contains all meta data concerning the L1Oracle contract.
var L1OracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"baseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1FeeOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1FeeScalar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"number\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_number\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_baseFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_l1FeeOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_l1FeeScalar\",\"type\":\"uint256\"}],\"name\":\"setL1OracleValues\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"stateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100e1565b600054610100900460ff161561008e5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116146100df576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6080516112c56101186000396000818161032701528181610370015281816104210152818161046101526104f401526112c56000f3fe6080604052600436106101145760003560e01c80638129fc1c116100a05780638da5cb5b116100645780638da5cb5b146102935780639588eca2146102bb5780639e8c4966146102d0578063b80777ea146102e7578063f2fde38b146102fd57600080fd5b80638129fc1c1461021c5780638381f58a146102315780638456cb59146102475780638b239f731461025c5780638b3a19f61461027357600080fd5b80634f1ef286116100e75780634f1ef286146101a657806352d1902d146101b95780635c975abb146101ce5780636ef25c3a146101f1578063715018a61461020757600080fd5b806309bd5a601461011957806313c3fb7b146101425780633659cfe61461016f5780633f4ba83a14610191575b600080fd5b34801561012557600080fd5b5061012f60ff5481565b6040519081526020015b60405180910390f35b34801561014e57600080fd5b5061012f61015d366004610f39565b60fb6020526000908152604090205481565b34801561017b57600080fd5b5061018f61018a366004610f78565b61031d565b005b34801561019d57600080fd5b5061018f610405565b61018f6101b4366004610fa9565b610417565b3480156101c557600080fd5b5061012f6104e7565b3480156101da57600080fd5b5060c95460ff166040519015158152602001610139565b3480156101fd57600080fd5b5061012f60fe5481565b34801561021357600080fd5b5061018f61059a565b34801561022857600080fd5b5061018f6105ac565b34801561023d57600080fd5b5061012f60fc5481565b34801561025357600080fd5b5061018f6106cc565b34801561026857600080fd5b5061012f6101005481565b34801561027f57600080fd5b5061018f61028e36600461106b565b6106dc565b34801561029f57600080fd5b506097546040516001600160a01b039091168152602001610139565b3480156102c757600080fd5b5061012f61080b565b3480156102dc57600080fd5b5061012f6101015481565b3480156102f357600080fd5b5061012f60fd5481565b34801561030957600080fd5b5061018f610318366004610f78565b61083b565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016300361036e5760405162461bcd60e51b8152600401610365906110b7565b60405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166103b7600080516020611249833981519152546001600160a01b031690565b6001600160a01b0316146103dd5760405162461bcd60e51b815260040161036590611103565b6103e6816108b1565b60408051600080825260208201909252610402918391906108c1565b50565b61040d610a31565b610415610a8b565b565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016300361045f5760405162461bcd60e51b8152600401610365906110b7565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166104a8600080516020611249833981519152546001600160a01b031690565b6001600160a01b0316146104ce5760405162461bcd60e51b815260040161036590611103565b6104d7826108b1565b6104e3828260016108c1565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146105875760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610365565b5060008051602061124983398151915290565b6105a2610a31565b6104156000610add565b600054610100900460ff16158080156105cc5750600054600160ff909116105b806105e65750303b1580156105e6575060005460ff166001145b6106495760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610365565b6000805460ff19166001179055801561066c576000805461ff0019166101001790555b610674610b2f565b61067c610b5e565b610684610b8d565b8015610402576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b6106d4610a31565b610415610bb4565b33411461073d5760405162461bcd60e51b815260206004820152602960248201527f4f6e6c792074686520636f696e626173652063616e2063616c6c207468697320604482015268333ab731ba34b7b71760b91b6064820152608401610365565b610745610bf1565b8660fc54106107bc5760405162461bcd60e51b815260206004820152603b60248201527f426c6f636b206e756d626572206d75737420626520677265617465722074686160448201527f6e207468652063757272656e7420626c6f636b206e756d6265722e00000000006064820152608401610365565b60fc87905560fd86905560fe85905560ff849055610100828155610101829055839060fb906000906107ee908b61114f565b60ff16815260208101919091526040016000205550505050505050565b600060fb600061010060fc54610821919061114f565b60ff1660ff16815260200190815260200160002054905090565b610843610a31565b6001600160a01b0381166108a85760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610365565b61040281610add565b6108b9610a31565b610402610c37565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156108f9576108f483610c80565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610953575060408051601f3d908101601f1916820190925261095091810190611171565b60015b6109b65760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610365565b6000805160206112498339815191528114610a255760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610365565b506108f4838383610d1c565b6097546001600160a01b031633146104155760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610365565b610a93610c37565b60c9805460ff191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516001600160a01b03909116815260200160405180910390a1565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff16610b565760405162461bcd60e51b81526004016103659061118a565b610415610d47565b600054610100900460ff16610b855760405162461bcd60e51b81526004016103659061118a565b610415610d77565b600054610100900460ff166104155760405162461bcd60e51b81526004016103659061118a565b610bbc610bf1565b60c9805460ff191660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258610ac03390565b60c95460ff16156104155760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b6044820152606401610365565b60c95460ff166104155760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b6044820152606401610365565b6001600160a01b0381163b610ced5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610365565b60008051602061124983398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b610d2583610daa565b600082511180610d325750805b156108f457610d418383610dea565b50505050565b600054610100900460ff16610d6e5760405162461bcd60e51b81526004016103659061118a565b61041533610add565b600054610100900460ff16610d9e5760405162461bcd60e51b81526004016103659061118a565b60c9805460ff19169055565b610db381610c80565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060610e0f838360405180606001604052806027815260200161126960279139610e16565b9392505050565b6060600080856001600160a01b031685604051610e3391906111f9565b600060405180830381855af49150503d8060008114610e6e576040519150601f19603f3d011682016040523d82523d6000602084013e610e73565b606091505b5091509150610e8486838387610e8e565b9695505050505050565b60608315610efd578251600003610ef6576001600160a01b0385163b610ef65760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610365565b5081610f07565b610f078383610f0f565b949350505050565b815115610f1f5781518083602001fd5b8060405162461bcd60e51b81526004016103659190611215565b600060208284031215610f4b57600080fd5b813560ff81168114610e0f57600080fd5b80356001600160a01b0381168114610f7357600080fd5b919050565b600060208284031215610f8a57600080fd5b610e0f82610f5c565b634e487b7160e01b600052604160045260246000fd5b60008060408385031215610fbc57600080fd5b610fc583610f5c565b9150602083013567ffffffffffffffff80821115610fe257600080fd5b818501915085601f830112610ff657600080fd5b81358181111561100857611008610f93565b604051601f8201601f19908116603f0116810190838211818310171561103057611030610f93565b8160405282815288602084870101111561104957600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600080600080600080600060e0888a03121561108657600080fd5b505085359760208701359750604087013596606081013596506080810135955060a0810135945060c0013592509050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b60008261116c57634e487b7160e01b600052601260045260246000fd5b500690565b60006020828403121561118357600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60005b838110156111f05781810151838201526020016111d8565b50506000910152565b6000825161120b8184602087016111d5565b9190910192915050565b60208152600082518060208401526112348160408501602087016111d5565b601f01601f1916919091016040019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212209792b00728e39caf5965cff6ab4b9c9e75dbdb7fc30fd4576354c5ff2bfc046764736f6c63430008110033",
}

// L1OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use L1OracleMetaData.ABI instead.
var L1OracleABI = L1OracleMetaData.ABI

// L1OracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L1OracleMetaData.Bin instead.
var L1OracleBin = L1OracleMetaData.Bin

// DeployL1Oracle deploys a new Ethereum contract, binding an instance of L1Oracle to it.
func DeployL1Oracle(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L1Oracle, error) {
	parsed, err := L1OracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L1OracleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L1Oracle{L1OracleCaller: L1OracleCaller{contract: contract}, L1OracleTransactor: L1OracleTransactor{contract: contract}, L1OracleFilterer: L1OracleFilterer{contract: contract}}, nil
}

// L1Oracle is an auto generated Go binding around an Ethereum contract.
type L1Oracle struct {
	L1OracleCaller     // Read-only binding to the contract
	L1OracleTransactor // Write-only binding to the contract
	L1OracleFilterer   // Log filterer for contract events
}

// L1OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1OracleSession struct {
	Contract     *L1Oracle         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1OracleCallerSession struct {
	Contract *L1OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// L1OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1OracleTransactorSession struct {
	Contract     *L1OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// L1OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1OracleRaw struct {
	Contract *L1Oracle // Generic contract binding to access the raw methods on
}

// L1OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1OracleCallerRaw struct {
	Contract *L1OracleCaller // Generic read-only contract binding to access the raw methods on
}

// L1OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1OracleTransactorRaw struct {
	Contract *L1OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1Oracle creates a new instance of L1Oracle, bound to a specific deployed contract.
func NewL1Oracle(address common.Address, backend bind.ContractBackend) (*L1Oracle, error) {
	contract, err := bindL1Oracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1Oracle{L1OracleCaller: L1OracleCaller{contract: contract}, L1OracleTransactor: L1OracleTransactor{contract: contract}, L1OracleFilterer: L1OracleFilterer{contract: contract}}, nil
}

// NewL1OracleCaller creates a new read-only instance of L1Oracle, bound to a specific deployed contract.
func NewL1OracleCaller(address common.Address, caller bind.ContractCaller) (*L1OracleCaller, error) {
	contract, err := bindL1Oracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1OracleCaller{contract: contract}, nil
}

// NewL1OracleTransactor creates a new write-only instance of L1Oracle, bound to a specific deployed contract.
func NewL1OracleTransactor(address common.Address, transactor bind.ContractTransactor) (*L1OracleTransactor, error) {
	contract, err := bindL1Oracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1OracleTransactor{contract: contract}, nil
}

// NewL1OracleFilterer creates a new log filterer instance of L1Oracle, bound to a specific deployed contract.
func NewL1OracleFilterer(address common.Address, filterer bind.ContractFilterer) (*L1OracleFilterer, error) {
	contract, err := bindL1Oracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1OracleFilterer{contract: contract}, nil
}

// bindL1Oracle binds a generic wrapper to an already deployed contract.
func bindL1Oracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1OracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Oracle *L1OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Oracle.Contract.L1OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Oracle *L1OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Oracle.Contract.L1OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Oracle *L1OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Oracle.Contract.L1OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Oracle *L1OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Oracle *L1OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Oracle *L1OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Oracle.Contract.contract.Transact(opts, method, params...)
}

// BaseFee is a free data retrieval call binding the contract method 0x6ef25c3a.
//
// Solidity: function baseFee() view returns(uint256)
func (_L1Oracle *L1OracleCaller) BaseFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "baseFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseFee is a free data retrieval call binding the contract method 0x6ef25c3a.
//
// Solidity: function baseFee() view returns(uint256)
func (_L1Oracle *L1OracleSession) BaseFee() (*big.Int, error) {
	return _L1Oracle.Contract.BaseFee(&_L1Oracle.CallOpts)
}

// BaseFee is a free data retrieval call binding the contract method 0x6ef25c3a.
//
// Solidity: function baseFee() view returns(uint256)
func (_L1Oracle *L1OracleCallerSession) BaseFee() (*big.Int, error) {
	return _L1Oracle.Contract.BaseFee(&_L1Oracle.CallOpts)
}

// Hash is a free data retrieval call binding the contract method 0x09bd5a60.
//
// Solidity: function hash() view returns(bytes32)
func (_L1Oracle *L1OracleCaller) Hash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "hash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Hash is a free data retrieval call binding the contract method 0x09bd5a60.
//
// Solidity: function hash() view returns(bytes32)
func (_L1Oracle *L1OracleSession) Hash() ([32]byte, error) {
	return _L1Oracle.Contract.Hash(&_L1Oracle.CallOpts)
}

// Hash is a free data retrieval call binding the contract method 0x09bd5a60.
//
// Solidity: function hash() view returns(bytes32)
func (_L1Oracle *L1OracleCallerSession) Hash() ([32]byte, error) {
	return _L1Oracle.Contract.Hash(&_L1Oracle.CallOpts)
}

// L1FeeOverhead is a free data retrieval call binding the contract method 0x8b239f73.
//
// Solidity: function l1FeeOverhead() view returns(uint256)
func (_L1Oracle *L1OracleCaller) L1FeeOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "l1FeeOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L1FeeOverhead is a free data retrieval call binding the contract method 0x8b239f73.
//
// Solidity: function l1FeeOverhead() view returns(uint256)
func (_L1Oracle *L1OracleSession) L1FeeOverhead() (*big.Int, error) {
	return _L1Oracle.Contract.L1FeeOverhead(&_L1Oracle.CallOpts)
}

// L1FeeOverhead is a free data retrieval call binding the contract method 0x8b239f73.
//
// Solidity: function l1FeeOverhead() view returns(uint256)
func (_L1Oracle *L1OracleCallerSession) L1FeeOverhead() (*big.Int, error) {
	return _L1Oracle.Contract.L1FeeOverhead(&_L1Oracle.CallOpts)
}

// L1FeeScalar is a free data retrieval call binding the contract method 0x9e8c4966.
//
// Solidity: function l1FeeScalar() view returns(uint256)
func (_L1Oracle *L1OracleCaller) L1FeeScalar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "l1FeeScalar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L1FeeScalar is a free data retrieval call binding the contract method 0x9e8c4966.
//
// Solidity: function l1FeeScalar() view returns(uint256)
func (_L1Oracle *L1OracleSession) L1FeeScalar() (*big.Int, error) {
	return _L1Oracle.Contract.L1FeeScalar(&_L1Oracle.CallOpts)
}

// L1FeeScalar is a free data retrieval call binding the contract method 0x9e8c4966.
//
// Solidity: function l1FeeScalar() view returns(uint256)
func (_L1Oracle *L1OracleCallerSession) L1FeeScalar() (*big.Int, error) {
	return _L1Oracle.Contract.L1FeeScalar(&_L1Oracle.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_L1Oracle *L1OracleCaller) Number(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "number")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_L1Oracle *L1OracleSession) Number() (*big.Int, error) {
	return _L1Oracle.Contract.Number(&_L1Oracle.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint256)
func (_L1Oracle *L1OracleCallerSession) Number() (*big.Int, error) {
	return _L1Oracle.Contract.Number(&_L1Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Oracle *L1OracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Oracle *L1OracleSession) Owner() (common.Address, error) {
	return _L1Oracle.Contract.Owner(&_L1Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Oracle *L1OracleCallerSession) Owner() (common.Address, error) {
	return _L1Oracle.Contract.Owner(&_L1Oracle.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1Oracle *L1OracleCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1Oracle *L1OracleSession) Paused() (bool, error) {
	return _L1Oracle.Contract.Paused(&_L1Oracle.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1Oracle *L1OracleCallerSession) Paused() (bool, error) {
	return _L1Oracle.Contract.Paused(&_L1Oracle.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L1Oracle *L1OracleCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L1Oracle *L1OracleSession) ProxiableUUID() ([32]byte, error) {
	return _L1Oracle.Contract.ProxiableUUID(&_L1Oracle.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L1Oracle *L1OracleCallerSession) ProxiableUUID() ([32]byte, error) {
	return _L1Oracle.Contract.ProxiableUUID(&_L1Oracle.CallOpts)
}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (_L1Oracle *L1OracleCaller) StateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "stateRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (_L1Oracle *L1OracleSession) StateRoot() ([32]byte, error) {
	return _L1Oracle.Contract.StateRoot(&_L1Oracle.CallOpts)
}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(bytes32)
func (_L1Oracle *L1OracleCallerSession) StateRoot() ([32]byte, error) {
	return _L1Oracle.Contract.StateRoot(&_L1Oracle.CallOpts)
}

// StateRoots is a free data retrieval call binding the contract method 0x13c3fb7b.
//
// Solidity: function stateRoots(uint8 ) view returns(bytes32)
func (_L1Oracle *L1OracleCaller) StateRoots(opts *bind.CallOpts, arg0 uint8) ([32]byte, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "stateRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateRoots is a free data retrieval call binding the contract method 0x13c3fb7b.
//
// Solidity: function stateRoots(uint8 ) view returns(bytes32)
func (_L1Oracle *L1OracleSession) StateRoots(arg0 uint8) ([32]byte, error) {
	return _L1Oracle.Contract.StateRoots(&_L1Oracle.CallOpts, arg0)
}

// StateRoots is a free data retrieval call binding the contract method 0x13c3fb7b.
//
// Solidity: function stateRoots(uint8 ) view returns(bytes32)
func (_L1Oracle *L1OracleCallerSession) StateRoots(arg0 uint8) ([32]byte, error) {
	return _L1Oracle.Contract.StateRoots(&_L1Oracle.CallOpts, arg0)
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() view returns(uint256)
func (_L1Oracle *L1OracleCaller) Timestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Oracle.contract.Call(opts, &out, "timestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() view returns(uint256)
func (_L1Oracle *L1OracleSession) Timestamp() (*big.Int, error) {
	return _L1Oracle.Contract.Timestamp(&_L1Oracle.CallOpts)
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() view returns(uint256)
func (_L1Oracle *L1OracleCallerSession) Timestamp() (*big.Int, error) {
	return _L1Oracle.Contract.Timestamp(&_L1Oracle.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_L1Oracle *L1OracleTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_L1Oracle *L1OracleSession) Initialize() (*types.Transaction, error) {
	return _L1Oracle.Contract.Initialize(&_L1Oracle.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_L1Oracle *L1OracleTransactorSession) Initialize() (*types.Transaction, error) {
	return _L1Oracle.Contract.Initialize(&_L1Oracle.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1Oracle *L1OracleTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1Oracle *L1OracleSession) Pause() (*types.Transaction, error) {
	return _L1Oracle.Contract.Pause(&_L1Oracle.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1Oracle *L1OracleTransactorSession) Pause() (*types.Transaction, error) {
	return _L1Oracle.Contract.Pause(&_L1Oracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1Oracle *L1OracleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1Oracle *L1OracleSession) RenounceOwnership() (*types.Transaction, error) {
	return _L1Oracle.Contract.RenounceOwnership(&_L1Oracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1Oracle *L1OracleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L1Oracle.Contract.RenounceOwnership(&_L1Oracle.TransactOpts)
}

// SetL1OracleValues is a paid mutator transaction binding the contract method 0x8b3a19f6.
//
// Solidity: function setL1OracleValues(uint256 _number, uint256 _timestamp, uint256 _baseFee, bytes32 _hash, bytes32 _stateRoot, uint256 _l1FeeOverhead, uint256 _l1FeeScalar) returns()
func (_L1Oracle *L1OracleTransactor) SetL1OracleValues(opts *bind.TransactOpts, _number *big.Int, _timestamp *big.Int, _baseFee *big.Int, _hash [32]byte, _stateRoot [32]byte, _l1FeeOverhead *big.Int, _l1FeeScalar *big.Int) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "setL1OracleValues", _number, _timestamp, _baseFee, _hash, _stateRoot, _l1FeeOverhead, _l1FeeScalar)
}

// SetL1OracleValues is a paid mutator transaction binding the contract method 0x8b3a19f6.
//
// Solidity: function setL1OracleValues(uint256 _number, uint256 _timestamp, uint256 _baseFee, bytes32 _hash, bytes32 _stateRoot, uint256 _l1FeeOverhead, uint256 _l1FeeScalar) returns()
func (_L1Oracle *L1OracleSession) SetL1OracleValues(_number *big.Int, _timestamp *big.Int, _baseFee *big.Int, _hash [32]byte, _stateRoot [32]byte, _l1FeeOverhead *big.Int, _l1FeeScalar *big.Int) (*types.Transaction, error) {
	return _L1Oracle.Contract.SetL1OracleValues(&_L1Oracle.TransactOpts, _number, _timestamp, _baseFee, _hash, _stateRoot, _l1FeeOverhead, _l1FeeScalar)
}

// SetL1OracleValues is a paid mutator transaction binding the contract method 0x8b3a19f6.
//
// Solidity: function setL1OracleValues(uint256 _number, uint256 _timestamp, uint256 _baseFee, bytes32 _hash, bytes32 _stateRoot, uint256 _l1FeeOverhead, uint256 _l1FeeScalar) returns()
func (_L1Oracle *L1OracleTransactorSession) SetL1OracleValues(_number *big.Int, _timestamp *big.Int, _baseFee *big.Int, _hash [32]byte, _stateRoot [32]byte, _l1FeeOverhead *big.Int, _l1FeeScalar *big.Int) (*types.Transaction, error) {
	return _L1Oracle.Contract.SetL1OracleValues(&_L1Oracle.TransactOpts, _number, _timestamp, _baseFee, _hash, _stateRoot, _l1FeeOverhead, _l1FeeScalar)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Oracle *L1OracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Oracle *L1OracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Oracle.Contract.TransferOwnership(&_L1Oracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Oracle *L1OracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Oracle.Contract.TransferOwnership(&_L1Oracle.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1Oracle *L1OracleTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1Oracle *L1OracleSession) Unpause() (*types.Transaction, error) {
	return _L1Oracle.Contract.Unpause(&_L1Oracle.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1Oracle *L1OracleTransactorSession) Unpause() (*types.Transaction, error) {
	return _L1Oracle.Contract.Unpause(&_L1Oracle.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L1Oracle *L1OracleTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L1Oracle *L1OracleSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L1Oracle.Contract.UpgradeTo(&_L1Oracle.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_L1Oracle *L1OracleTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _L1Oracle.Contract.UpgradeTo(&_L1Oracle.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L1Oracle *L1OracleTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1Oracle.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L1Oracle *L1OracleSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1Oracle.Contract.UpgradeToAndCall(&_L1Oracle.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L1Oracle *L1OracleTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1Oracle.Contract.UpgradeToAndCall(&_L1Oracle.TransactOpts, newImplementation, data)
}

// L1OracleAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the L1Oracle contract.
type L1OracleAdminChangedIterator struct {
	Event *L1OracleAdminChanged // Event containing the contract specifics and raw log

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
func (it *L1OracleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OracleAdminChanged)
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
		it.Event = new(L1OracleAdminChanged)
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
func (it *L1OracleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OracleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OracleAdminChanged represents a AdminChanged event raised by the L1Oracle contract.
type L1OracleAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L1Oracle *L1OracleFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*L1OracleAdminChangedIterator, error) {

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &L1OracleAdminChangedIterator{contract: _L1Oracle.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_L1Oracle *L1OracleFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *L1OracleAdminChanged) (event.Subscription, error) {

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OracleAdminChanged)
				if err := _L1Oracle.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_L1Oracle *L1OracleFilterer) ParseAdminChanged(log types.Log) (*L1OracleAdminChanged, error) {
	event := new(L1OracleAdminChanged)
	if err := _L1Oracle.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1OracleBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the L1Oracle contract.
type L1OracleBeaconUpgradedIterator struct {
	Event *L1OracleBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *L1OracleBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OracleBeaconUpgraded)
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
		it.Event = new(L1OracleBeaconUpgraded)
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
func (it *L1OracleBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OracleBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OracleBeaconUpgraded represents a BeaconUpgraded event raised by the L1Oracle contract.
type L1OracleBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_L1Oracle *L1OracleFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*L1OracleBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &L1OracleBeaconUpgradedIterator{contract: _L1Oracle.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_L1Oracle *L1OracleFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *L1OracleBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OracleBeaconUpgraded)
				if err := _L1Oracle.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_L1Oracle *L1OracleFilterer) ParseBeaconUpgraded(log types.Log) (*L1OracleBeaconUpgraded, error) {
	event := new(L1OracleBeaconUpgraded)
	if err := _L1Oracle.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1OracleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the L1Oracle contract.
type L1OracleInitializedIterator struct {
	Event *L1OracleInitialized // Event containing the contract specifics and raw log

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
func (it *L1OracleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OracleInitialized)
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
		it.Event = new(L1OracleInitialized)
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
func (it *L1OracleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OracleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OracleInitialized represents a Initialized event raised by the L1Oracle contract.
type L1OracleInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_L1Oracle *L1OracleFilterer) FilterInitialized(opts *bind.FilterOpts) (*L1OracleInitializedIterator, error) {

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &L1OracleInitializedIterator{contract: _L1Oracle.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_L1Oracle *L1OracleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *L1OracleInitialized) (event.Subscription, error) {

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OracleInitialized)
				if err := _L1Oracle.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_L1Oracle *L1OracleFilterer) ParseInitialized(log types.Log) (*L1OracleInitialized, error) {
	event := new(L1OracleInitialized)
	if err := _L1Oracle.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1OracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L1Oracle contract.
type L1OracleOwnershipTransferredIterator struct {
	Event *L1OracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L1OracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OracleOwnershipTransferred)
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
		it.Event = new(L1OracleOwnershipTransferred)
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
func (it *L1OracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OracleOwnershipTransferred represents a OwnershipTransferred event raised by the L1Oracle contract.
type L1OracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Oracle *L1OracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1OracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1OracleOwnershipTransferredIterator{contract: _L1Oracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Oracle *L1OracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L1OracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OracleOwnershipTransferred)
				if err := _L1Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L1Oracle *L1OracleFilterer) ParseOwnershipTransferred(log types.Log) (*L1OracleOwnershipTransferred, error) {
	event := new(L1OracleOwnershipTransferred)
	if err := _L1Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1OraclePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the L1Oracle contract.
type L1OraclePausedIterator struct {
	Event *L1OraclePaused // Event containing the contract specifics and raw log

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
func (it *L1OraclePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OraclePaused)
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
		it.Event = new(L1OraclePaused)
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
func (it *L1OraclePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OraclePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OraclePaused represents a Paused event raised by the L1Oracle contract.
type L1OraclePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_L1Oracle *L1OracleFilterer) FilterPaused(opts *bind.FilterOpts) (*L1OraclePausedIterator, error) {

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &L1OraclePausedIterator{contract: _L1Oracle.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_L1Oracle *L1OracleFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *L1OraclePaused) (event.Subscription, error) {

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OraclePaused)
				if err := _L1Oracle.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_L1Oracle *L1OracleFilterer) ParsePaused(log types.Log) (*L1OraclePaused, error) {
	event := new(L1OraclePaused)
	if err := _L1Oracle.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1OracleUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the L1Oracle contract.
type L1OracleUnpausedIterator struct {
	Event *L1OracleUnpaused // Event containing the contract specifics and raw log

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
func (it *L1OracleUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OracleUnpaused)
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
		it.Event = new(L1OracleUnpaused)
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
func (it *L1OracleUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OracleUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OracleUnpaused represents a Unpaused event raised by the L1Oracle contract.
type L1OracleUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_L1Oracle *L1OracleFilterer) FilterUnpaused(opts *bind.FilterOpts) (*L1OracleUnpausedIterator, error) {

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &L1OracleUnpausedIterator{contract: _L1Oracle.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_L1Oracle *L1OracleFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *L1OracleUnpaused) (event.Subscription, error) {

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OracleUnpaused)
				if err := _L1Oracle.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_L1Oracle *L1OracleFilterer) ParseUnpaused(log types.Log) (*L1OracleUnpaused, error) {
	event := new(L1OracleUnpaused)
	if err := _L1Oracle.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1OracleUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the L1Oracle contract.
type L1OracleUpgradedIterator struct {
	Event *L1OracleUpgraded // Event containing the contract specifics and raw log

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
func (it *L1OracleUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1OracleUpgraded)
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
		it.Event = new(L1OracleUpgraded)
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
func (it *L1OracleUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1OracleUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1OracleUpgraded represents a Upgraded event raised by the L1Oracle contract.
type L1OracleUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L1Oracle *L1OracleFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*L1OracleUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L1Oracle.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &L1OracleUpgradedIterator{contract: _L1Oracle.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L1Oracle *L1OracleFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *L1OracleUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L1Oracle.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1OracleUpgraded)
				if err := _L1Oracle.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_L1Oracle *L1OracleFilterer) ParseUpgraded(log types.Log) (*L1OracleUpgraded, error) {
	event := new(L1OracleUpgraded)
	if err := _L1Oracle.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
