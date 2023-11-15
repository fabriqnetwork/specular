package predeploys

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/ops/backends"
	"github.com/specularL2/specular/ops/bindings/bindings"
	"github.com/specularL2/specular/ops/deployer"
)

type StorageConfig struct {
	ImplValue  any
	ProxyValue any
}

type PredeployConfig struct {
	Proxied           bool
	ConstructorValues map[string]any
	Initializer       string
	InitializerValues map[string]any
	Storages          map[string]StorageConfig
}

type PredeployConfigs map[string]PredeployConfig

func (p PredeployConfigs) Check() error {
	checks := []struct {
		contract string
		initArg  []string
	}{
		{"L2Portal", []string{"_l1PortalAddress"}},
		{"L2StandardBridge", []string{"_otherBridge"}},
	}
	for _, check := range checks {
		if _, ok := p[check.contract]; !ok {
			return fmt.Errorf("missing predeploy for %s", check.contract)
		}
		for _, arg := range check.initArg {
			if _, ok := p[check.contract].InitializerValues[arg]; !ok {
				return fmt.Errorf("missing initializer arg %s for %s", arg, check.contract)
			}
		}
	}
	return nil
}

type DeploymentResult struct {
	Bytecode hexutil.Bytes
	Address  common.Address
}

// DeploymentResults represents the output of deploying each of the
// contracts so that the immutables can be set properly in the bytecode.
type DeploymentResults map[string]DeploymentResult

// BuildSpecular will deploy the L2 predeploys so that their immutables are set
// correctly.
func BuildSpecular(ctx context.Context, predeploy PredeployConfigs) (DeploymentResults, DeploymentResults, error) {
	if err := predeploy.Check(); err != nil {
		return nil, nil, err
	}

	backend := deployer.NewL2Backend()
	implDeploymentResults, err := BuildPredeployImpls(ctx, backend, predeploy)
	if err != nil {
		return nil, nil, err
	}
	proxyDeploymentResults, err := BuildPredeployProxies(ctx, backend, predeploy, implDeploymentResults)
	if err != nil {
		return nil, nil, err
	}

	return implDeploymentResults, proxyDeploymentResults, nil
}

func BuildPredeployImpls(ctx context.Context, backend *backends.SimulatedBackend, predeploys PredeployConfigs) (DeploymentResults, error) {
	implConstructors := []deployer.Constructor{
		{
			Name: "UUPSPlaceholder",
		},
		{
			Name: "L1Oracle",
		},
		{
			Name: "L2Portal",
		},
		{
			Name: "L2StandardBridge",
		},
		{
			Name: "L2BaseFeeVault",
		},
	}
	deployments, err := deployer.Deploy(backend, implConstructors, l2Deployer)
	if err != nil {
		return nil, err
	}
	results := make(DeploymentResults)
	for _, dep := range deployments {
		results[dep.Name] = DeploymentResult{
			Bytecode: dep.Bytecode,
			Address:  dep.Address,
		}
	}
	return results, nil
}

type initializeData struct {
	metaData *bind.MetaData
	argOrder []string
}

func (d initializeData) pack(predeploy PredeployConfig) ([]byte, error) {
	abi, err := d.metaData.GetAbi()
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, len(d.argOrder))
	for i, argName := range d.argOrder {
		args[i] = predeploy.InitializerValues[argName]
	}
	return abi.Pack(predeploy.Initializer, args...)
}

func BuildPredeployProxies(ctx context.Context, backend *backends.SimulatedBackend, predeploys PredeployConfigs, implDeploymentResults DeploymentResults) (DeploymentResults, error) {
	initData := map[string]initializeData{
		"UUPSPlaceholder": {
			metaData: bindings.UUPSPlaceholderMetaData,
		},
		"L1Oracle": {
			metaData: bindings.L1OracleMetaData,
		},
		"L2Portal": {
			metaData: bindings.L2PortalMetaData,
			argOrder: []string{"_l1PortalAddress"},
		},
		"L2StandardBridge": {
			metaData: bindings.L2StandardBridgeMetaData,
			argOrder: []string{"_otherBridge"},
		},
		"L2BaseFeeVault": {
			metaData: bindings.L2BaseFeeVaultMetaData,
		},
	}
	proxyConstructors := make([]deployer.Constructor, 0)
	for name, predeploy := range predeploys {
		if !predeploy.Proxied {
			continue
		}

		data, err := initData[name].pack(predeploy)
		if err != nil {
			return nil, err
		}
		proxyConstructors = append(proxyConstructors, deployer.Constructor{
			Name:     "ERC1967Proxy",
			ImplName: name,
			Args: []any{
				implDeploymentResults[name].Address,
				data,
			},
		})
	}

	deployments, err := deployer.Deploy(backend, proxyConstructors, l2Deployer)
	if err != nil {
		return nil, err
	}
	results := make(DeploymentResults)
	for _, dep := range deployments {
		results[dep.ImplName] = DeploymentResult{
			Bytecode: dep.Bytecode,
			Address:  dep.Address,
		}
	}
	return results, nil
}

func l2Deployer(backend *backends.SimulatedBackend, opts *bind.TransactOpts, deployment deployer.Constructor) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error
	switch deployment.Name {
	case "ERC1967Proxy":
		implAddr, ok := deployment.Args[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for implAddr")
		}
		data, ok := deployment.Args[1].([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid type for initializer data")
		}
		_, tx, _, err = bindings.DeployERC1967Proxy(opts, backend, implAddr, data)
	case "UUPSPlaceholder":
		_, tx, _, err = bindings.DeployUUPSPlaceholder(opts, backend)
	case "L1Oracle":
		_, tx, _, err = bindings.DeployL1Oracle(opts, backend)
	case "L2Portal":
		_, tx, _, err = bindings.DeployL2Portal(opts, backend)
	case "L2StandardBridge":
		_, tx, _, err = bindings.DeployL2StandardBridge(opts, backend)
	case "L2BaseFeeVault":
		_, tx, _, err = bindings.DeployL2BaseFeeVault(opts, backend)
	default:
		return tx, fmt.Errorf("unknown contract: %s", deployment.Name)
	}

	return tx, err
}
