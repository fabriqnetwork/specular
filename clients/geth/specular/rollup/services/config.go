package services

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
	"github.com/urfave/cli/v2"
)

type systemConfig struct {
	L1Config        `toml:"l1,omitempty"`
	L2Config        `toml:"l2,omitempty"`
	SequencerConfig `toml:"sequencer,omitempty"`
	ValidatorConfig `toml:"validator,omitempty"`
	DriverConfig    `toml:"driver,omitempty"`
}

func (c *systemConfig) L1() L1Config               { return c.L1Config }
func (c *systemConfig) L2() L2Config               { return c.L2Config }
func (c *systemConfig) Sequencer() SequencerConfig { return c.SequencerConfig }
func (c *systemConfig) Validator() ValidatorConfig { return c.ValidatorConfig }
func (c *systemConfig) Driver() DriverConfig       { return c.DriverConfig }

// Parses all CLI flags and returns a full system config.
func ParseSystemConfig(cliCtx *cli.Context) (*systemConfig, error) {
	var (
		sysCfg  = parseFlags(cliCtx)
		cfgPath = cliCtx.String(cfgFileFlag.Name)
		err     error
	)
	// Override with config file if provided.
	if cfgPath != "" {
		log.Info("Parsing TOML config", "path", cfgPath)
		err = parseTOML(cfgPath, sysCfg)
	}
	return sysCfg, err
}

// Parses a TOML config file. Overwrites any fields that are already set in cfg.
// path: path to TOML config file.
// cfg: system config to decode into.
func parseTOML(path string, cfg *systemConfig) error {
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return fmt.Errorf("couldn't decode TOML config: %w", err)
	}
	// Print full config to stdout in TOML format for sanity.
	log.Info("Outputting full system config (including unconfigured defaults)...")
	if err := toml.NewEncoder(os.Stdout).Encode(cfg); err != nil {
		return fmt.Errorf("couldn't re-encode TOML: %w", err)
	}
	return nil
}

func parseFlags(cliCtx *cli.Context) *systemConfig {
	utils.CheckExclusive(cliCtx, l1EndpointFlag, utils.MiningEnabledFlag)
	utils.CheckExclusive(cliCtx, l1EndpointFlag, utils.DeveloperFlag)
	var (
		sequencerAddr, validatorAddr             common.Address
		sequencerPassphrase, validatorPassphrase string
		pwList                                   = utils.MakePasswordList(cliCtx)
		clefEndpoint                             = cliCtx.String(l2ClefEndpointFlag.Name)
	)
	if clefEndpoint == "" {
		if len(pwList) == 0 {
			utils.Fatalf("Failed to parse system config: no clef endpoint or password provided")
		}
		if cliCtx.String(sequencerAddrFlag.Name) != "" {
			sequencerAddr = common.HexToAddress(cliCtx.String(sequencerAddrFlag.Name))
			sequencerPassphrase = pwList[0]
			pwList = pwList[1:]
		}
		if cliCtx.String(validatorAddrFlag.Name) != "" {
			validatorAddr = common.HexToAddress(cliCtx.String(validatorAddrFlag.Name))
			validatorPassphrase = pwList[0]
		}
	}

	var (
		l1ChainID         = big.NewInt(0).SetUint64(cliCtx.Uint64(l1ChainIDFlag.Name))
		sequencerTxMgrCfg = txmgr.NewConfigFromCLI(cliCtx, sequencerTxMgrNamespace, l1ChainID, sequencerAddr)
		validatorTxMgrCfg = txmgr.NewConfigFromCLI(cliCtx, validatorTxMgrNamespace, l1ChainID, validatorAddr)
	)
	return &systemConfig{
		L1Config:        newL1ConfigFromCLI(cliCtx),
		L2Config:        newL2ConfigFromCLI(cliCtx),
		SequencerConfig: newSequencerConfigFromCLI(cliCtx, sequencerPassphrase, sequencerTxMgrCfg),
		ValidatorConfig: newValidatorConfigFromCLI(cliCtx, validatorPassphrase, validatorTxMgrCfg),
		DriverConfig:    newDriverConfigFromCLI(cliCtx),
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
		SequencerInboxAddr: common.HexToAddress(cliCtx.String(sequencerInboxAddrFlag.Name)),
		RollupAddr:         common.HexToAddress(cliCtx.String(rollupAddrFlag.Name)),
	}
}

func (c L1Config) GetEndpoint() string                   { return c.Endpoint }
func (c L1Config) GetChainID() uint64                    { return c.ChainID }
func (c L1Config) GetRollupGenesisBlock() uint64         { return c.RollupGenesisBlock }
func (c L1Config) GetSequencerInboxAddr() common.Address { return c.SequencerInboxAddr }
func (c L1Config) GetRollupAddr() common.Address         { return c.RollupAddr }

// L2 configuration
type L2Config struct {
	Endpoint              string                      `toml:"endpoint,omitempty"`      // L2 API endpoint
	ClefEndpoint          string                      `toml:"clef_endpoint,omitempty"` // The Clef Endpoint used for signing TXs
	BlockProductionPolicy BlockProductionPolicyConfig `toml:"block_production_policy,omitempty"`
}

// TODO: these should be configured in L1 contracts; we can read from there.
type BlockProductionPolicyConfig struct {
	TargetBlockTime time.Duration `toml:"target_block_time,omitempty"`
	IsEmptyAllowed  bool          `toml:"is_empty_allowed,omitempty"`
}

func (c BlockProductionPolicyConfig) GetTargetBlockTime() time.Duration { return c.TargetBlockTime }
func (c BlockProductionPolicyConfig) GetIsEmptyAllowed() bool           { return c.IsEmptyAllowed }

func newL2ConfigFromCLI(cliCtx *cli.Context) L2Config {
	return L2Config{
		Endpoint:     cliCtx.String(l2EndpointFlag.Name),
		ClefEndpoint: cliCtx.String(l2ClefEndpointFlag.Name),
		BlockProductionPolicy: BlockProductionPolicyConfig{
			TargetBlockTime: cliCtx.Duration(l2TargetBlockTimeFlag.Name),
		},
	}
}

func (c L2Config) GetEndpoint() string     { return c.Endpoint }
func (c L2Config) GetClefEndpoint() string { return c.ClefEndpoint }
func (c L2Config) GetBlockProductionPolicy() BlockProductionPolicyConfig {
	return c.BlockProductionPolicy
}

// Sequencer node configuration
type SequencerConfig struct {
	AccountAddr common.Address `toml:"account_addr,omitempty"`
	// The passphrase of the sequencer account
	Passphrase string `toml:"passphrase,omitempty"`
	// Max # of l2 blocks the safe head should lag behind the unsafe head. If <= 0, enforcement disabled.
	MaxSafeLag uint `toml:"max_safe_lag,omitempty"`
	// Time between batch dissemination (DA) attempts
	DisseminationInterval time.Duration `toml:"dissemination_interval,omitempty"`
	// Transaction manager configuration
	TxMgrCfg txmgr.Config `toml:"txmgr,omitempty"`
}

func (c SequencerConfig) GetAccountAddr() common.Address          { return c.AccountAddr }
func (c SequencerConfig) GetPassphrase() string                   { return c.Passphrase }
func (c SequencerConfig) GetMaxSafeLag() uint                     { return c.MaxSafeLag }
func (c SequencerConfig) GetDisseminationInterval() time.Duration { return c.DisseminationInterval }
func (c SequencerConfig) GetTxMgrCfg() txmgr.Config               { return c.TxMgrCfg }

func newSequencerConfigFromCLI(
	cliCtx *cli.Context,
	passphrase string,
	txMgrCfg txmgr.Config,
) SequencerConfig {
	return SequencerConfig{
		AccountAddr:           common.HexToAddress(cliCtx.String(sequencerAddrFlag.Name)),
		Passphrase:            passphrase,
		DisseminationInterval: time.Duration(cliCtx.Uint(sequencerSequencingIntervalFlag.Name)) * time.Second,
		TxMgrCfg:              txMgrCfg,
	}
}

// Validator node configuration
type ValidatorConfig struct {
	AccountAddr common.Address `toml:"account_addr,omitempty"`
	// The passphrase of the validator account
	Passphrase string `toml:"passphrase,omitempty"`
	// True iff actively stakes on rollup contract
	IsActiveStaker bool `toml:"is_active_staker,omitempty"`
	// True iff actively tries to create new assertions (not just for a challenge).
	IsActiveCreator bool `toml:"is_active_creator,omitempty"`
	// True iff actively issues challenges as challenger. *Defends* against challenges regardless.
	IsActiveChallenger bool `toml:"is_active_challenger,omitempty"`
	// True iff attempts to resolve assertions (by confirming or rejecting)
	IsResolver bool `toml:"is_resolver,omitempty"`
	// Size of stake to deposit to rollup contract
	StakeAmount uint64 `toml:"stake_amount,omitempty"`
	// Transaction manager configuration
	TxMgrCfg txmgr.Config `toml:"txmgr,omitempty"`
}

func (c ValidatorConfig) GetAccountAddr() common.Address { return c.AccountAddr }
func (c ValidatorConfig) GetPassphrase() string          { return c.Passphrase }
func (c ValidatorConfig) GetIsActiveStaker() bool        { return c.IsActiveStaker }
func (c ValidatorConfig) GetIsActiveCreator() bool       { return c.IsActiveCreator }
func (c ValidatorConfig) GetIsActiveChallenger() bool    { return c.IsActiveChallenger }
func (c ValidatorConfig) GetIsResolver() bool            { return c.IsResolver }
func (c ValidatorConfig) GetStakeAmount() uint64         { return c.StakeAmount }
func (c ValidatorConfig) GetTxMgrCfg() txmgr.Config      { return c.TxMgrCfg }

func newValidatorConfigFromCLI(
	ctx *cli.Context,
	passphrase string,
	txMgrCfg txmgr.Config,
) ValidatorConfig {
	return ValidatorConfig{
		AccountAddr:        common.HexToAddress(ctx.String(validatorAddrFlag.Name)),
		Passphrase:         passphrase,
		IsActiveStaker:     ctx.Bool(validatorIsActiveStakerFlag.Name),
		IsActiveCreator:    ctx.Bool(validatorIsActiveCreatorFlag.Name),
		IsActiveChallenger: ctx.Bool(validatorIsActiveChallengerFlag.Name),
		IsResolver:         ctx.Bool(validatorIsResolverFlag.Name),
		StakeAmount:        ctx.Uint64(validatorStakeAmountFlag.Name),
		TxMgrCfg:           txMgrCfg,
	}
}

// Driver configuration
type DriverConfig struct {
	// Time between driver steps (in steady state; failures may trigger a longer backoff)
	StepInterval time.Duration `toml:"step_interval,omitempty"`
	// Time to wait before retrying a step attempt
	RetryDelay time.Duration `toml:"retry_delay,omitempty"`
	// # attempts to attempt driver step before catastrophically failing. Must be > 0.
	NumAttempts uint `toml:"num_attempts,omitempty"`
}

func (c DriverConfig) GetStepInterval() time.Duration { return c.StepInterval }
func (c DriverConfig) GetRetryDelay() time.Duration   { return c.RetryDelay }
func (c DriverConfig) GetNumAttempts() uint           { return c.NumAttempts }

func newDriverConfigFromCLI(cliCtx *cli.Context) DriverConfig {
	return DriverConfig{
		StepInterval: time.Duration(cliCtx.Uint(driverStepIntervalFlag.Name)) * time.Second,
		RetryDelay:   time.Duration(cliCtx.Uint(driverRetryDelayFlag.Name)) * time.Second,
		NumAttempts:  cliCtx.Uint(driverNumAttemptsFlag.Name),
	}
}
