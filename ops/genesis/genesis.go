package genesis

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/specularL2/specular/ops/bindings/bindings"
	"github.com/specularL2/specular/ops/predeploys"
	"github.com/specularL2/specular/ops/state"
)

// defaultGasLimit represents the default gas limit for a genesis block.
const defaultGasLimit = 30_000_000

func NewL2EmptyGenesis(config *GenesisConfig, block *types.Block) (*core.Genesis, error) {
	if config.L2ChainID == 0 {
		return nil, errors.New("must define L2 ChainID")
	}

	specularChainConfig := params.ChainConfig{
		ChainID:                       new(big.Int).SetUint64(config.L2ChainID),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            big.NewInt(0),
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
		EnableL2EngineApi:             true,
		L2BaseFeeRecipient:            predeploys.L2BaseFeeVaultAddr,
	}

	gasLimit := config.L2GenesisBlockGasLimit
	if gasLimit == 0 {
		gasLimit = defaultGasLimit
	}
	baseFee := config.L2GenesisBlockBaseFeePerGas
	if baseFee == nil {
		baseFee = newHexBig(params.InitialBaseFee)
	}
	difficulty := config.L2GenesisBlockDifficulty
	if difficulty == nil {
		difficulty = newHexBig(0)
	}

	extraData := config.L2GenesisBlockExtraData
	// Ensure that the extradata is valid
	if size := len(extraData); size > 32 {
		return nil, fmt.Errorf("transition block extradata too long: %d", size)
	}

	return &core.Genesis{
		Config:     &specularChainConfig,
		Nonce:      uint64(config.L2GenesisBlockNonce),
		Timestamp:  block.Time(),
		ExtraData:  extraData,
		GasLimit:   uint64(gasLimit),
		Difficulty: difficulty.ToInt(),
		Mixhash:    config.L2GenesisBlockMixHash,
		Coinbase:   config.L2GenesisBlockCoinbase,
		Number:     uint64(config.L2GenesisBlockNumber),
		GasUsed:    uint64(config.L2GenesisBlockGasUsed),
		ParentHash: config.L2GenesisBlockParentHash,
		BaseFee:    baseFee.ToInt(),
		Alloc:      map[common.Address]core.GenesisAccount{},
	}, nil
}

func BuildL2Genesis(ctx context.Context, config *GenesisConfig, l1StartBlock *types.Block) (*core.Genesis, error) {
	genesis, err := NewL2EmptyGenesis(config, l1StartBlock)
	if err != nil {
		return nil, err
	}
	predeployConfigs := GeneratePredeployConfig(config, l1StartBlock)

	db := state.NewMemoryStateDB(genesis)

	implDeployments, proxyDeployments, err := predeploys.BuildSpecular(ctx, predeployConfigs)
	if err != nil {
		return nil, err
	}

	hasPredeploy := make(map[common.Address]struct{})
	for name, config := range predeployConfigs {
		hasPredeploy[*predeploys.Predeploys[name]] = struct{}{}
		if err := setupPredeploy(ctx, db, name, config, implDeployments[name]); err != nil {
			return nil, err
		}
		if err := setupProxy(ctx, db, name, config, proxyDeployments[name]); err != nil {
			return nil, err
		}
	}
	for i := uint64(0); i <= predeploys.PredeployProxyCount; i++ {
		bigAddr := new(big.Int).Or(predeploys.BigL2PredeployNamespace, new(big.Int).SetUint64(i))
		addr := common.BigToAddress(bigAddr)
		if _, ok := hasPredeploy[addr]; ok {
			continue
		}
		if err := setupEmptyProxy(ctx, db, addr, predeployConfigs["UUPSPlaceholder"]); err != nil {
			return nil, err
		}
	}

	genesisWithPredeploy := db.Genesis()
	setupAllocs(genesisWithPredeploy, config.Alloc)

	return genesisWithPredeploy, nil
}

func setupPredeploy(ctx context.Context, db vm.StateDB, name string, config predeploys.PredeployConfig, implDep predeploys.DeploymentResult) error {
	implAddr := *predeploys.Predeploys[name]
	if config.Proxied {
		var err error
		implAddr, err = predeploys.AddressToCodeNamespace(implAddr)
		if err != nil {
			return err
		}
	}
	db.CreateAccount(implAddr)
	db.SetCode(implAddr, implDep.Bytecode)
	implStorageValues := make(state.StorageValues)
	for label, value := range config.Storages {
		if value.ImplValue != nil {
			implStorageValues[label] = value.ImplValue
		}
	}
	log.Debug("Setting impl storage", "name", name, "address", implAddr)
	if err := state.SetStorage(name, implAddr, implStorageValues, db); err != nil {
		return err
	}
	return nil
}

func setupProxy(ctx context.Context, db vm.StateDB, name string, config predeploys.PredeployConfig, proxyDep predeploys.DeploymentResult) error {
	if !config.Proxied {
		return nil
	}
	proxyAddr := *predeploys.Predeploys[name]
	implAddr, err := predeploys.AddressToCodeNamespace(proxyAddr)
	if err != nil {
		return err
	}
	db.CreateAccount(proxyAddr)
	db.SetCode(proxyAddr, proxyDep.Bytecode)
	db.SetState(proxyAddr, predeploys.ImplementationSlot, state.AddressAsLeftPaddedHash(implAddr))
	proxyStorageValues := make(state.StorageValues)
	for label, value := range config.Storages {
		if value.ProxyValue != nil {
			proxyStorageValues[label] = value.ProxyValue
		}
	}
	log.Debug("Setting proxy storage", "name", name, "address", proxyAddr)
	if err := state.SetStorage(name, proxyAddr, proxyStorageValues, db); err != nil {
		return err
	}
	return nil
}

func setupEmptyProxy(ctx context.Context, db vm.StateDB, addr common.Address, placeholderConfig predeploys.PredeployConfig) error {
	proxyCode, err := bindings.GetDeployedBytecode("ERC1967Proxy")
	if err != nil {
		return err
	}
	db.CreateAccount(addr)
	db.SetCode(addr, proxyCode)
	db.SetState(addr, predeploys.ImplementationSlot, state.AddressAsLeftPaddedHash(predeploys.UUPSPlaceholderAddr))
	proxyStorageValues := make(state.StorageValues)
	for label, value := range placeholderConfig.Storages {
		if value.ProxyValue != nil {
			proxyStorageValues[label] = value.ProxyValue
		}
	}
	log.Debug("Setting empty proxy storage", "address", addr)
	if err := state.SetStorage("UUPSPlaceholder", addr, proxyStorageValues, db); err != nil {
		return err
	}
	return nil
}

func setupAllocs(genesis *core.Genesis, allocs map[common.Address]core.GenesisAccount) {
	for addr, account := range allocs {
		if existAccount, ok := genesis.Alloc[addr]; ok {
			log.Warn("Overwriting existing genesis account", "address", addr)
			if account.Balance != nil || account.Balance.Cmp(common.Big0) > 0 {
				existAccount.Balance = account.Balance
			}
			if len(account.Code) > 0 {
				existAccount.Code = account.Code
			}
			if len(account.Storage) > 0 {
				existAccount.Storage = account.Storage
			}
			if account.Nonce != 0 {
				existAccount.Nonce = account.Nonce
			}
			genesis.Alloc[addr] = existAccount
		} else {
			genesis.Alloc[addr] = account
		}
	}
}
