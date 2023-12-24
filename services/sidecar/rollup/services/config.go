package services

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/specularL2/specular/services/sidecar/utils/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"

	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth/txmgr"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
)

// TODO: rename due to naming conflict
type SystemConfig struct {
	ProtocolConfig     `toml:"protocol,omitempty"`
	L1Config           `toml:"l1,omitempty"`
	L2Config           `toml:"l2,omitempty"`
	DisseminatorConfig `toml:"disseminator,omitempty"`
	ValidatorConfig    `toml:"validator,omitempty"`
	Verbosity          log.Lvl `toml:"verbosity,omitempty"`
}

func (c *SystemConfig) Protocol() ProtocolConfig         { return c.ProtocolConfig }
func (c *SystemConfig) L1() L1Config                     { return c.L1Config }
func (c *SystemConfig) L2() L2Config                     { return c.L2Config }
func (c *SystemConfig) Disseminator() DisseminatorConfig { return c.DisseminatorConfig }
func (c *SystemConfig) Validator() ValidatorConfig       { return c.ValidatorConfig }

func (c *SystemConfig) validate() error {
	if !(c.DisseminatorConfig.IsEnabled || c.ValidatorConfig.IsEnabled) {
		return fmt.Errorf("at least one of disseminator and validator must be enabled")
	}
	if err := c.DisseminatorConfig.validate(); err != nil {
		return fmt.Errorf("disseminator config invalid: %w", err)
	}
	if err := c.ValidatorConfig.validate(); err != nil {
		return fmt.Errorf("validator config invalid: %w", err)
	}
	return nil
}

// Parses all CLI flags and returns a full system config.
func ParseSystemConfig(cliCtx *cli.Context) (*SystemConfig, error) {
	protocolCfg, err := newProtocolConfigFromCLI(cliCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse protocol config: %w", err)
	}
	var (
		l1ChainID = protocolCfg.GetRollup().L1ChainID
		cfg       = &SystemConfig{
			ProtocolConfig:     protocolCfg,
			L1Config:           newL1ConfigFromCLI(cliCtx),
			L2Config:           newL2ConfigFromCLI(cliCtx),
			DisseminatorConfig: newDisseminatorConfigFromCLI(cliCtx, l1ChainID),
			ValidatorConfig:    newValidatorConfigFromCLI(cliCtx, l1ChainID),
			Verbosity:          log.Lvl(cliCtx.Int(VerbosityFlag.Name)),
		}
	)
	// Validate.
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}
	return cfg, nil
}

// Protocol configuration
type ProtocolConfig struct {
	Rollup RollupConfig `toml:"rollup,omitempty"`
}

func newProtocolConfigFromCLI(cliCtx *cli.Context) (ProtocolConfig, error) {
	rollupCfg, err := NewRollupConfig(cliCtx.String(protocolRollupCfgPathFlag.Name))
	if err != nil {
		return ProtocolConfig{}, err
	}
	return ProtocolConfig{Rollup: *rollupCfg}, nil
}

// TODO: cleanup (consider: exposing parameters via getters in `c.Rollup` directly).
func (c ProtocolConfig) GetRollup() RollupConfig               { return c.Rollup }
func (c ProtocolConfig) GetRollupAddr() common.Address         { return c.Rollup.RollupAddress }
func (c ProtocolConfig) GetSeqWindowSize() uint64              { return c.Rollup.SeqWindowSize }
func (c ProtocolConfig) GetSequencerInboxAddr() common.Address { return c.Rollup.BatchInboxAddress }
func (c ProtocolConfig) GetL1ChainID() uint64                  { return c.Rollup.L1ChainID.Uint64() }
func (c ProtocolConfig) GetL2ChainID() uint64                  { return c.Rollup.L2ChainID.Uint64() }
func (c ProtocolConfig) GetL1OracleAddr() common.Address {
	// TODO: import from package or config
	return common.HexToAddress("0x2A00000000000000000000000000000000000010")
}

// L1 configuration
type L1Config struct {
	Endpoint string `toml:"endpoint,omitempty"` // L1 API endpoint
}

func newL1ConfigFromCLI(cliCtx *cli.Context) L1Config {
	return L1Config{Endpoint: cliCtx.String(l1EndpointFlag.Name)}
}

func (c L1Config) GetEndpoint() string { return c.Endpoint }

// L2 configuration
type L2Config struct {
	Endpoint string `toml:"endpoint,omitempty"` // L2 API endpoint
	ChainID  uint64 `toml:"chainid,omitempty"`  // L2 chain ID
}

func newL2ConfigFromCLI(cliCtx *cli.Context) L2Config {
	return L2Config{Endpoint: cliCtx.String(l2EndpointFlag.Name)}
}

func (c L2Config) GetEndpoint() string { return c.Endpoint }

// Sequencer node configuration
type DisseminatorConfig struct {
	// Whether this node is a sequencer
	IsEnabled bool `toml:"enabled,omitempty"`
	// The address of this sequencer
	AccountAddr common.Address `toml:"account_addr,omitempty"`
	// The private key for AccountAddr
	PrivateKey *ecdsa.PrivateKey
	// The Clef Endpoint used for signing txs
	ClefEndpoint string `toml:"clef_endpoint,omitempty"`
	// Time between batch dissemination (DA) steps
	DisseminationInterval time.Duration `toml:"dissemination_interval,omitempty"`
	// The safety margin for batch tx submission (in # of L1 blocks)
	SubSafetyMargin uint64 `toml:"sub_safety_margin,omitempty"`
	// The target size of a batch tx submitted to L1 (bytes).
	TargetBatchSize uint64 `toml:"max_l1_tx_size,omitempty"`
	// Transaction manager configuration
	TxMgrCfg txmgr.Config `toml:"txmgr,omitempty"`
}

func (c DisseminatorConfig) GetIsEnabled() bool                      { return c.IsEnabled }
func (c DisseminatorConfig) GetAccountAddr() common.Address          { return c.AccountAddr }
func (c DisseminatorConfig) GetPrivateKey() *ecdsa.PrivateKey        { return c.PrivateKey }
func (c DisseminatorConfig) GetClefEndpoint() string                 { return c.ClefEndpoint }
func (c DisseminatorConfig) GetDisseminationInterval() time.Duration { return c.DisseminationInterval }
func (c DisseminatorConfig) GetSubSafetyMargin() uint64              { return c.SubSafetyMargin }
func (c DisseminatorConfig) GetTargetBatchSize() uint64              { return c.TargetBatchSize }
func (c DisseminatorConfig) GetTxMgrCfg() txmgr.Config               { return c.TxMgrCfg }

// Validates the configuration.
func (c DisseminatorConfig) validate() error {
	if !c.IsEnabled {
		return nil
	}
	if c.PrivateKey == nil && c.ClefEndpoint == "" {
		return fmt.Errorf("missing both private key and clef endpoint (require at least one)")
	}
	if c.PrivateKey != nil && c.AccountAddr != crypto.PubkeyToAddress(c.PrivateKey.PublicKey) {
		return fmt.Errorf("private key does not correspond to account address")
	}
	// Enforce sensible values.
	if c.TargetBatchSize < 128 {
		return fmt.Errorf("target batch size must be at least 128B")
	}
	return c.TxMgrCfg.Validate()
}

func newDisseminatorConfigFromCLI(cliCtx *cli.Context, l1ChainID *big.Int) DisseminatorConfig {
	var (
		privateKey = toPrivateKey(cliCtx.String(disseminatorPrivateKeyFlag.Name))
		address    = crypto.PubkeyToAddress(privateKey.PublicKey)
		txMgrCfg   = txmgr.NewConfigFromCLI(cliCtx, disseminatorTxMgrNamespace, l1ChainID, address)
	)
	return DisseminatorConfig{
		IsEnabled:             cliCtx.Bool(disseminatorEnableFlag.Name),
		AccountAddr:           address,
		PrivateKey:            privateKey,
		ClefEndpoint:          cliCtx.String(disseminatorClefEndpointFlag.Name),
		DisseminationInterval: time.Duration(cliCtx.Uint(disseminatorIntervalFlag.Name)) * time.Second,
		SubSafetyMargin:       cliCtx.Uint64(disseminatorSubSafetyMarginFlag.Name),
		TargetBatchSize:       cliCtx.Uint64(disseminatorTargetBatchSizeFlag.Name),
		TxMgrCfg:              txMgrCfg,
	}
}

type ValidatorConfig struct {
	// Whether this node is a validator
	IsEnabled bool `toml:"enabled,omitempty"`
	// The address of this validator
	AccountAddr common.Address `toml:"account_addr,omitempty"`
	// The private key for AccountAddr
	PrivateKey *ecdsa.PrivateKey
	// The Clef Endpoint used for signing txs
	ClefEndpoint string `toml:"clef_endpoint,omitempty"`
	// Time between validation steps
	ValidationInterval time.Duration `toml:"validation_interval,omitempty"`
	// Transaction manager configuration
	TxMgrCfg txmgr.Config `toml:"txmgr,omitempty"`
}

func (c ValidatorConfig) GetIsEnabled() bool                   { return c.IsEnabled }
func (c ValidatorConfig) GetAccountAddr() common.Address       { return c.AccountAddr }
func (c ValidatorConfig) GetPrivateKey() *ecdsa.PrivateKey     { return c.PrivateKey }
func (c ValidatorConfig) GetClefEndpoint() string              { return c.ClefEndpoint }
func (c ValidatorConfig) GetValidationInterval() time.Duration { return c.ValidationInterval }
func (c ValidatorConfig) GetTxMgrCfg() txmgr.Config            { return c.TxMgrCfg }

// Validates the configuration.
func (c ValidatorConfig) validate() error {
	if !c.IsEnabled {
		return nil
	}
	if c.PrivateKey == nil && c.ClefEndpoint == "" {
		return fmt.Errorf("missing both private key and clef endpoint (require at least one)")
	}
	if c.PrivateKey != nil && c.AccountAddr != crypto.PubkeyToAddress(c.PrivateKey.PublicKey) {
		return fmt.Errorf("private key does not correspond to account address")
	}
	return c.TxMgrCfg.Validate()
}

func newValidatorConfigFromCLI(cliCtx *cli.Context, l1ChainID *big.Int) ValidatorConfig {
	var (
		privateKey = toPrivateKey(cliCtx.String(validatorPrivateKeyFlag.Name))
		address    = crypto.PubkeyToAddress(privateKey.PublicKey)
		txMgrCfg   = txmgr.NewConfigFromCLI(cliCtx, validatorTxMgrNamespace, l1ChainID, address)
	)
	return ValidatorConfig{
		IsEnabled:          cliCtx.Bool(validatorEnableFlag.Name),
		AccountAddr:        address,
		PrivateKey:         privateKey,
		ClefEndpoint:       cliCtx.String(validatorClefEndpointFlag.Name),
		ValidationInterval: time.Duration(cliCtx.Uint(validatorValidationIntervalFlag.Name)) * time.Second,
		TxMgrCfg:           txMgrCfg,
	}
}

func toPrivateKey(keyStr string) *ecdsa.PrivateKey {
	if keyStr == "" {
		return nil
	}
	secretKey, err := crypto.HexToECDSA(keyStr[2:])
	if err != nil {
		panic("failed to parse secret key: " + err.Error())
	}
	return secretKey
}
