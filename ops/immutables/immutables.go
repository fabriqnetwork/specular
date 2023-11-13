// Adapted from Optimism's `op-chain-ops`

package immutables

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/specularL2/specular/ops/backends"
	"github.com/specularL2/specular/ops/bindings/bindings"
	"github.com/specularL2/specular/ops/deployer"
)

// ImmutableValues represents the values to be set in immutable code.
// The key is the name of the variable and the value is the value to set in
// immutable code.
type ImmutableValues map[string]any

// ImmutableConfig represents the immutable configuration for the L2 predeploy
// contracts.
type ImmutableConfig map[string]ImmutableValues

// Check does a sanity check that the specific values that
// Specular uses are set inside of the ImmutableConfig.
func (i ImmutableConfig) Check() error {
	if _, ok := i["general"]["deployer"]; !ok {
		return errors.New("general: deployer not set")
	}
	if _, ok := i["L2Portal"]["otherMessenger"]; !ok {
		return errors.New("L2CrossDomainMessenger otherMessenger not set")
	}
	if _, ok := i["L2StandardBridge"]["otherBridge"]; !ok {
		return errors.New("L2StandardBridge otherBridge not set")
	}

	return nil
}

// DeploymentResults represents the output of deploying each of the
// contracts so that the immutables can be set properly in the bytecode.
type DeploymentResults map[string]hexutil.Bytes

// BuildSpecular will deploy the L2 predeploys so that their immutables are set
// correctly.
func BuildSpecular(immutable ImmutableConfig) (DeploymentResults, error) {
	if err := immutable.Check(); err != nil {
		return DeploymentResults{}, err
	}

	deployments := []deployer.Constructor{
		{
			Name: "UUPSPlaceholder",
		},
		{
			Name: "L1Block",
		},
		{
			Name: "L2Portal",
		},
		{
			Name: "L2StandardBridge",
		},
	}
	return BuildL2(deployments)
}

// BuildL2 will deploy contracts to a simulated backend so that their immutables
// can be properly set. The bytecode returned in the results is suitable to be
// inserted into the state via state surgery.
func BuildL2(constructors []deployer.Constructor) (DeploymentResults, error) {
	log.Info("Creating L2 state")
	deployments, err := deployer.Deploy(deployer.NewL2Backend(), constructors, l2Deployer)
	if err != nil {
		return nil, err
	}
	results := make(DeploymentResults)
	for _, dep := range deployments {
		results[dep.Name] = dep.Bytecode
	}
	return results, nil
}

func l2Deployer(backend *backends.SimulatedBackend, opts *bind.TransactOpts, deployment deployer.Constructor) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error
	switch deployment.Name {
	case "UUPSPlaceholder":
		_, tx, _, err = bindings.DeployUUPSPlaceholder(opts, backend)
	case "L1Oracle":
		_, tx, _, err = bindings.DeployL1Oracle(opts, backend)
	case "L2Portal":
		_, tx, _, err = bindings.DeployL2Portal(opts, backend)
	case "L2StandardBridge":
		_, tx, _, err = bindings.DeployL2StandardBridge(opts, backend)
	default:
		return tx, fmt.Errorf("unknown contract: %s", deployment.Name)
	}

	return tx, err
}
