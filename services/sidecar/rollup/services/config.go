package services

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth/txmgr"
	"github.com/urfave/cli/v2"
)

type SystemConfig struct {
	L1Config        `toml:"l1,omitempty"`
	L2Config        `toml:"l2,omitempty"`
	SequencerConfig `toml:"sequencer,omitempty"`
	ValidatorConfig `toml:"validator,omitempty"`
}

func (c *SystemConfig) L1() L1Config               { return c.L1Config }
func (c *SystemConfig) L2() L2Config               { return c.L2Config }
func (c *SystemConfig) Sequencer() SequencerConfig { return c.SequencerConfig }
func (c *SystemConfig) Validator() ValidatorConfig { return c.ValidatorConfig }

// Parses all CLI flags and returns a full system config.
func ParseSystemConfig(cliCtx *cli.Context) (*SystemConfig, error) {
	return parseFlags(cliCtx), nil
}

func parseFlags(cliCtx *cli.Context) *SystemConfig {
	utils.CheckExclusive(cliCtx, l1EndpointFlag, utils.MiningEnabledFlag)
	utils.CheckExclusive(cliCtx, l1EndpointFlag, utils.DeveloperFlag)
	var sequencerAddr common.Address
	if cliCtx.String(sequencerAddrFlag.Name) != "" {
		sequencerAddr = common.HexToAddress(cliCtx.String(sequencerAddrFlag.Name))
	}
	var validatorAddr common.Address
	if cliCtx.String(validatorAddrFlag.Name) != "" {
		validatorAddr = common.HexToAddress(cliCtx.String(validatorAddrFlag.Name))
	}
	var (
		l1ChainID         = big.NewInt(0).SetUint64(cliCtx.Uint64(l1ChainIDFlag.Name))
		sequencerTxMgrCfg = txmgr.NewConfigFromCLI(cliCtx, sequencerTxMgrNamespace, l1ChainID, sequencerAddr)
		validatorTxMgrCfg = txmgr.NewConfigFromCLI(cliCtx, validatorTxMgrNamespace, l1ChainID, validatorAddr)
	)
	return &SystemConfig{
		L1Config:        newL1ConfigFromCLI(cliCtx),
		L2Config:        newL2ConfigFromCLI(cliCtx),
		SequencerConfig: newSequencerConfigFromCLI(cliCtx, sequencerTxMgrCfg),
		ValidatorConfig: newValidatorConfigFromCLI(cliCtx, validatorTxMgrCfg),
	}
}

// L1 configuration
type L1Config struct {
	Endpoint           string         `toml:"endpoint,omitempty"`             // L1 API endpoint
	ChainID            uint64         `toml:"chainid,omitempty"`              // L1 chain ID
	RollupGenesisBlock uint64         `toml:"rollup_genesis,omitempty"`       // L1 Rollup genesis block
	SequencerInboxAddr common.Address `toml:"sequencer_inbox_addr,omitempty"` // L1 SequencerInbox contract address
	RollupAddr         common.Address `toml:"rollup_addr,omitempty"`          // L1 Rollup contract address
}

func newL1ConfigFromCLI(cliCtx *cli.Context) L1Config {
	return L1Config{
		Endpoint:           cliCtx.String(l1EndpointFlag.Name),
		ChainID:            cliCtx.Uint64(l1ChainIDFlag.Name),
		RollupGenesisBlock: cliCtx.Uint64(l1RollupGenesisBlockFlag.Name),
		SequencerInboxAddr: common.HexToAddress(cliCtx.String(l1SequencerInboxAddrFlag.Name)),
		RollupAddr:         common.HexToAddress(cliCtx.String(l1RollupAddrFlag.Name)),
	}
}

func (c L1Config) GetEndpoint() string                   { return c.Endpoint }
func (c L1Config) GetChainID() uint64                    { return c.ChainID }
func (c L1Config) GetRollupGenesisBlock() uint64         { return c.RollupGenesisBlock }
func (c L1Config) GetSequencerInboxAddr() common.Address { return c.SequencerInboxAddr }
func (c L1Config) GetRollupAddr() common.Address         { return c.RollupAddr }

// L2 configuration
type L2Config struct {
	Endpoint            string         `toml:"endpoint,omitempty"`                // L2 API endpoint
	ChainID             uint64         `toml:"chainid,omitempty"`                 // L2 chain ID
	L1FeeOverhead       int64          `toml:"l1_fee_overhead,omitempty"`         // Gas cost of sequencing a tx
	L1FeeMultiplier     float64        `toml:"l1_fee_multiplier,omitempty"`       // Scalar value to increase L1 fee
	L1OracleAddress     common.Address `toml:"l1_oracle_address,omitempty"`       // L2 Address of the L1Oracle
	L1OracleBaseFeeSlot common.Hash    `toml:"l1_oracle_base_fee_slot,omitempty"` // L1 basefee storage slot of the L1Oracle
}

func newL2ConfigFromCLI(cliCtx *cli.Context) L2Config {
	return L2Config{
		Endpoint:            cliCtx.String(l2EndpointFlag.Name),
		ChainID:             cliCtx.Uint64(l2ChainIDFlag.Name),
		L1FeeOverhead:       cliCtx.Int64(l2L1FeeOverheadFlag.Name),
		L1FeeMultiplier:     l2L1FeeMultiplierFlag.Value,
		L1OracleAddress:     common.HexToAddress(l2L1OracleAddressFlag.Value),
		L1OracleBaseFeeSlot: common.HexToHash(l2L1OracleBaseFeeSlotFlag.Value),
	}
}

func (c L2Config) GetEndpoint() string                 { return c.Endpoint }
func (c L2Config) GetChainID() uint64                  { return c.ChainID }
func (c L2Config) GetL1FeeOverhead() int64             { return c.L1FeeOverhead }
func (c L2Config) GetL1FeeMultiplier() float64         { return c.L1FeeMultiplier }
func (c L2Config) GetL1OracleAddress() common.Address  { return c.L1OracleAddress }
func (c L2Config) GetL1OracleBaseFeeSlot() common.Hash { return c.L1OracleBaseFeeSlot }

// Sequencer node configuration
type SequencerConfig struct {
	// Whether this node is a sequencer
	IsEnabled bool `toml:"enabled,omitempty"`
	// The address of this sequencer
	AccountAddr common.Address `toml:"account_addr,omitempty"`
	// The secret key for AccountAddr
	SecretKey *ecdsa.PrivateKey
	// The Clef Endpoint used for signing txs
	ClefEndpoint string `toml:"clef_endpoint,omitempty"`
	// Time between batch dissemination (DA) steps
	DisseminationInterval time.Duration `toml:"dissemination_interval,omitempty"`
	// Transaction manager configuration
	TxMgrCfg txmgr.Config `toml:"txmgr,omitempty"`
}

func (c SequencerConfig) GetIsEnabled() bool                      { return c.IsEnabled }
func (c SequencerConfig) GetAccountAddr() common.Address          { return c.AccountAddr }
func (c SequencerConfig) GetSecretKey() *ecdsa.PrivateKey         { return c.SecretKey }
func (c SequencerConfig) GetClefEndpoint() string                 { return c.ClefEndpoint }
func (c SequencerConfig) GetDisseminationInterval() time.Duration { return c.DisseminationInterval }
func (c SequencerConfig) GetTxMgrCfg() txmgr.Config               { return c.TxMgrCfg }

func newSequencerConfigFromCLI(
	cliCtx *cli.Context,
	txMgrCfg txmgr.Config,
) SequencerConfig {
	return SequencerConfig{
		IsEnabled:             cliCtx.Bool(sequencerEnableSequencerFlag.Name),
		AccountAddr:           common.HexToAddress(cliCtx.String(sequencerAddrFlag.Name)),
		SecretKey:             toSecretKey(cliCtx.String(sequencerSecretKeyFlag.Name)),
		ClefEndpoint:          cliCtx.String(sequencerClefEndpointFlag.Name),
		DisseminationInterval: time.Duration(cliCtx.Uint(sequencerSequencingIntervalFlag.Name)) * time.Second,
		TxMgrCfg:              txMgrCfg,
	}
}

type ValidatorConfig struct {
	// Whether this node is a validator
	IsEnabled bool `toml:"enabled,omitempty"`
	// The address of this validator
	AccountAddr common.Address `toml:"account_addr,omitempty"`
	// The secret key for AccountAddr
	SecretKey *ecdsa.PrivateKey
	// The Clef Endpoint used for signing txs
	ClefEndpoint string `toml:"clef_endpoint,omitempty"`
	// Time between validation steps
	ValidationInterval time.Duration `toml:"validation_interval,omitempty"`
	// Transaction manager configuration
	TxMgrCfg txmgr.Config `toml:"txmgr,omitempty"`
}

func (c ValidatorConfig) GetIsEnabled() bool                   { return c.IsEnabled }
func (c ValidatorConfig) GetAccountAddr() common.Address       { return c.AccountAddr }
func (c ValidatorConfig) GetSecretKey() *ecdsa.PrivateKey      { return c.SecretKey }
func (c ValidatorConfig) GetClefEndpoint() string              { return c.ClefEndpoint }
func (c ValidatorConfig) GetValidationInterval() time.Duration { return c.ValidationInterval }
func (c ValidatorConfig) GetTxMgrCfg() txmgr.Config            { return c.TxMgrCfg }

func newValidatorConfigFromCLI(
	cliCtx *cli.Context,
	txMgrCfg txmgr.Config,
) ValidatorConfig {
	return ValidatorConfig{
		IsEnabled:          cliCtx.Bool(validatorEnableValidatorFlag.Name),
		AccountAddr:        common.HexToAddress(cliCtx.String(validatorAddrFlag.Name)),
		SecretKey:          toSecretKey(cliCtx.String(validatorSecretKeyFlag.Name)),
		ClefEndpoint:       cliCtx.String(validatorClefEndpointFlag.Name),
		ValidationInterval: time.Duration(cliCtx.Uint(validatorValidationIntervalFlag.Name)) * time.Second,
		TxMgrCfg:           txMgrCfg,
	}
}

func toSecretKey(keyStr string) *ecdsa.PrivateKey {
	if keyStr == "" {
		return nil
	}
	secretKey, err := crypto.HexToECDSA(keyStr[2:])
	if err != nil {
		panic("failed to parse secret key: " + err.Error())
	}
	return secretKey
}
