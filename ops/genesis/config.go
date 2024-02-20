package genesis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/ops/predeploys"
)

// initialzedValue represents the `Initializable` contract value.
const InitializedValue uint8 = 1
const MaxInitializedValue uint8 = 255

var defaultL1Sender = common.HexToAddress("0x000000000000000000000000000000000000dEaD")

type GenesisConfig struct {
	// L2ChainID is the chain ID of the L2 chain.
	L2ChainID uint64 `json:"l2ChainID"`

	L2GenesisBlockNonce         hexutil.Uint64 `json:"l2GenesisBlockNonce"`
	L2GenesisBlockGasLimit      hexutil.Uint64 `json:"l2GenesisBlockGasLimit"`
	L2GenesisBlockDifficulty    *hexutil.Big   `json:"l2GenesisBlockDifficulty"`
	L2GenesisBlockMixHash       common.Hash    `json:"l2GenesisBlockMixHash"`
	L2GenesisBlockCoinbase      common.Address `json:"l2GenesisBlockCoinbase"`
	L2GenesisBlockNumber        hexutil.Uint64 `json:"l2GenesisBlockNumber"`
	L2GenesisBlockGasUsed       hexutil.Uint64 `json:"l2GenesisBlockGasUsed"`
	L2GenesisBlockParentHash    common.Hash    `json:"l2GenesisBlockParentHash"`
	L2GenesisBlockBaseFeePerGas *hexutil.Big   `json:"l2GenesisBlockBaseFeePerGas"`
	L2GenesisBlockExtraData     hexutil.Bytes  `json:"l2GenesisBlockExtraData"`

	L2PredeployOwner         common.Address `json:"l2PredeployOwner"`
	L1PortalAddress          common.Address `json:"l1PortalAddress,omitempty"`
	L1StandardBridgeAddress  common.Address `json:"l1StandardBridgeAddress,omitempty"`
	L2FeesWithdrawalAddress  common.Address `json:"l2FeesWithdrawalAddress"`
	L2FeesMinWithdrwalAmount *hexutil.Big   `json:"l2FeesMinWithdrwalAmount"`

	L1FeeOverhead *hexutil.Big `json:"l1FeeOverhead"`
	L1FeeScalar   *hexutil.Big `json:"l1FeeScalar"`

	Alloc core.GenesisAlloc `json:"alloc"`
}

func GeneratePredeployConfig(config *GenesisConfig, block *types.Block) predeploys.PredeployConfigs {
	predeployConfigs := predeploys.PredeployConfigs{
		"UUPSPlaceholder": {
			Proxied:     false,
			Initializer: "initialize",
			Storages: map[string]predeploys.StorageConfig{
				"_initialized":  {ProxyValue: InitializedValue, ImplValue: MaxInitializedValue},
				"_initializing": {ProxyValue: false, ImplValue: false},
				"_owner":        {ProxyValue: config.L2PredeployOwner},
			},
		},
		"L1Oracle": {
			Proxied:     true,
			Initializer: "initialize",
			Storages: map[string]predeploys.StorageConfig{
				"_initialized":  {ProxyValue: InitializedValue, ImplValue: MaxInitializedValue},
				"_initializing": {ProxyValue: false, ImplValue: false},
				"_owner":        {ProxyValue: config.L2PredeployOwner},
				"number":        {ProxyValue: block.Number()},
				"timestamp":     {ProxyValue: block.Time()},
				"baseFee":       {ProxyValue: block.BaseFee()},
				"hash":          {ProxyValue: block.Hash()},
				"l1FeeOverhead": {ProxyValue: (*big.Int)(config.L1FeeOverhead)},
				"l1FeeScalar":   {ProxyValue: (*big.Int)(config.L1FeeScalar)},
			},
		},
		"L2Portal": {
			Proxied:     true,
			Initializer: "initialize",
			InitializerValues: map[string]any{
				"_l1PortalAddress": config.L1PortalAddress,
			},
			Storages: map[string]predeploys.StorageConfig{
				"_initialized":    {ProxyValue: InitializedValue, ImplValue: MaxInitializedValue},
				"_initializing":   {ProxyValue: false, ImplValue: false},
				"_owner":          {ProxyValue: config.L2PredeployOwner},
				"l1PortalAddress": {ProxyValue: config.L1PortalAddress},
				"l1Sender":        {ProxyValue: defaultL1Sender},
			},
		},
		"L2StandardBridge": {
			Proxied:     true,
			Initializer: "initialize",
			InitializerValues: map[string]any{
				"_otherBridge": config.L1StandardBridgeAddress,
			},
			Storages: map[string]predeploys.StorageConfig{
				"_initialized":   {ProxyValue: InitializedValue, ImplValue: MaxInitializedValue},
				"_initializing":  {ProxyValue: false, ImplValue: false},
				"_owner":         {ProxyValue: config.L2PredeployOwner},
				"OTHER_BRIDGE":   {ProxyValue: config.L1StandardBridgeAddress},
				"PORTAL_ADDRESS": {ProxyValue: *predeploys.Predeploys["L2Portal"]},
			},
		},
		"L1FeeVault": {
			Proxied:     true,
			Initializer: "initialize",
			Storages: map[string]predeploys.StorageConfig{
				"_initialized":        {ProxyValue: InitializedValue, ImplValue: MaxInitializedValue},
				"_initializing":       {ProxyValue: false, ImplValue: false},
				"_owner":              {ProxyValue: config.L2PredeployOwner},
				"withdrawalAddress":   {ProxyValue: config.L2FeesWithdrawalAddress},
				"minWithdrawalAmount": {ProxyValue: (*big.Int)(config.L2FeesMinWithdrwalAmount)},
			},
		},
		"L2BaseFeeVault": {
			Proxied:     true,
			Initializer: "initialize",
			Storages: map[string]predeploys.StorageConfig{
				"_initialized":        {ProxyValue: InitializedValue, ImplValue: MaxInitializedValue},
				"_initializing":       {ProxyValue: false, ImplValue: false},
				"_owner":              {ProxyValue: config.L2PredeployOwner},
				"withdrawalAddress":   {ProxyValue: config.L2FeesWithdrawalAddress},
				"minWithdrawalAmount": {ProxyValue: (*big.Int)(config.L2FeesMinWithdrwalAmount)},
			},
		},
		"MintableERC20Factory": {
			Proxied: false,
			ConstructorValues: map[string]any{
				"_bridge": predeploys.MintableERC20FactoryAddr,
			},
		},
	}
	return predeployConfigs
}

func NewGenesisConfig(path string) (*GenesisConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("deploy config at %s not found: %w", path, err)
	}

	dec := json.NewDecoder(bytes.NewReader(file))
	dec.DisallowUnknownFields()

	var config GenesisConfig
	if err := dec.Decode(&config); err != nil {
		return nil, fmt.Errorf("cannot unmarshal deploy config: %w", err)
	}

	return &config, nil
}
