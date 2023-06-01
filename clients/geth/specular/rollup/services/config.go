package services

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/txmgr"
	"github.com/urfave/cli/v2"
)

// TODO: use json tags and directly parse a config file
type systemConfig struct {
	l1Config
	l2Config
	sequencerConfig
	validatorConfig
	driverConfig
}

func (c *systemConfig) L1() l1Config               { return c.l1Config }
func (c *systemConfig) L2() l2Config               { return c.l2Config }
func (c *systemConfig) Sequencer() sequencerConfig { return c.sequencerConfig }
func (c *systemConfig) Validator() validatorConfig { return c.validatorConfig }
func (c *systemConfig) Driver() driverConfig       { return c.driverConfig }

// Parses all CLI flags and returns a full system config.
func ParseSystemConfig(cliCtx *cli.Context) *systemConfig {
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
			pwList = pwList[1:]
		}
	}

	var (
		l1ChainID         = big.NewInt(0).SetUint64(cliCtx.Uint64(l1ChainIDFlag.Name))
		sequencerTxMgrCfg = txmgr.NewConfigFromCLI(cliCtx, sequencerTxMgrNamespace, l1ChainID, sequencerAddr)
		validatorTxMgrCfg = txmgr.NewConfigFromCLI(cliCtx, validatorTxMgrNamespace, l1ChainID, validatorAddr)
	)
	return &systemConfig{
		l1Config:        newL1ConfigFromCLI(cliCtx),
		l2Config:        newL2ConfigFromCLI(cliCtx),
		sequencerConfig: newSequencerConfigFromCLI(cliCtx, sequencerPassphrase, sequencerTxMgrCfg),
		validatorConfig: newValidatorConfigFromCLI(cliCtx, validatorPassphrase, validatorTxMgrCfg),
		driverConfig:    newDriverConfigFromCLI(cliCtx),
	}
}

// L1 configuration
type l1Config struct {
	endpoint           string         // L1 API endpoint
	chainID            uint64         // L1 chain ID
	rollupGenesisBlock uint64         // L1 Rollup genesis block
	sequencerInboxAddr common.Address // L1 SequencerInbox contract address
	rollupAddr         common.Address // L1 Rollup contract address
}

func newL1ConfigFromCLI(cliCtx *cli.Context) l1Config {
	return l1Config{
		endpoint:           cliCtx.String(l1EndpointFlag.Name),
		chainID:            cliCtx.Uint64(l1ChainIDFlag.Name),
		rollupGenesisBlock: cliCtx.Uint64(l1RollupGenesisBlockFlag.Name),
		sequencerInboxAddr: common.HexToAddress(cliCtx.String(sequencerInboxAddrFlag.Name)),
		rollupAddr:         common.HexToAddress(cliCtx.String(rollupAddrFlag.Name)),
	}
}

func (c l1Config) Endpoint() string                   { return c.endpoint }
func (c l1Config) ChainID() uint64                    { return c.chainID }
func (c l1Config) RollupGenesisBlock() uint64         { return c.rollupGenesisBlock }
func (c l1Config) SequencerInboxAddr() common.Address { return c.sequencerInboxAddr }
func (c l1Config) RollupAddr() common.Address         { return c.rollupAddr }

// L2 configuration
type l2Config struct {
	endpoint     string // L2 API endpoint
	clefEndpoint string // The Clef Endpoint used for signing TXs
}

func newL2ConfigFromCLI(cliCtx *cli.Context) l2Config {
	return l2Config{
		endpoint:     cliCtx.String(l2EndpointFlag.Name),
		clefEndpoint: cliCtx.String(l2ClefEndpointFlag.Name),
	}
}

func (c l2Config) Endpoint() string     { return c.endpoint }
func (c l2Config) ClefEndpoint() string { return c.clefEndpoint }

// Sequencer node configuration
type sequencerConfig struct {
	accountAddr          common.Address
	passphrase           string        // The passphrase of the sequencer account
	minExecutionInterval time.Duration // Minimum time between block production. If 0, txs executed immediately -- FCFS.
	maxExecutionInterval time.Duration // Maximum time between block production. Must be >= `minExecutionInterval`.
	sequencingInterval   time.Duration // Time between batch sequencing attempts
	txMgrCfg             txmgr.Config
}

func (c sequencerConfig) AccountAddr() common.Address         { return c.accountAddr }
func (c sequencerConfig) Passphrase() string                  { return c.passphrase }
func (c sequencerConfig) MinExecutionInterval() time.Duration { return c.minExecutionInterval }
func (c sequencerConfig) MaxExecutionInterval() time.Duration { return c.maxExecutionInterval }
func (c sequencerConfig) SequencingInterval() time.Duration   { return c.sequencingInterval }
func (c sequencerConfig) TxMgrCfg() txmgr.Config              { return c.txMgrCfg }

func newSequencerConfigFromCLI(
	cliCtx *cli.Context,
	passphrase string,
	txMgrCfg txmgr.Config,
) sequencerConfig {
	return sequencerConfig{
		accountAddr:          common.HexToAddress(cliCtx.String(sequencerAddrFlag.Name)),
		passphrase:           passphrase,
		minExecutionInterval: time.Duration(cliCtx.Uint(sequencerMinExecIntervalFlag.Name)) * time.Second,
		maxExecutionInterval: time.Duration(cliCtx.Uint(sequencerMaxExecIntervalFlag.Name)) * time.Second,
		sequencingInterval:   time.Duration(cliCtx.Uint(sequencerSequencingIntervalFlag.Name)) * time.Second,
		txMgrCfg:             txMgrCfg,
	}
}

// Validator node configuration
type validatorConfig struct {
	accountAddr        common.Address
	passphrase         string // The passphrase of the validator account
	isActiveStaker     bool   // Iff true, actively stakes on rollup contract
	isActiveCreator    bool   // Iff true, actively tries to create new assertions (not just for a challenge).
	isActiveChallenger bool   // Iff true, actively issues challenges as challenger. *Defends* against challenges regardless.
	isResolver         bool   // Iff true, attempts to resolve assertions (by confirming or rejecting)
	stakeAmount        uint64 // Size of stake to deposit to rollup contract
	txMgrCfg           txmgr.Config
}

func (c validatorConfig) AccountAddr() common.Address { return c.accountAddr }
func (c validatorConfig) Passphrase() string          { return c.passphrase }
func (c validatorConfig) IsActiveStaker() bool        { return c.isActiveStaker }
func (c validatorConfig) IsActiveCreator() bool       { return c.isActiveCreator }
func (c validatorConfig) IsActiveChallenger() bool    { return c.isActiveChallenger }
func (c validatorConfig) IsResolver() bool            { return c.isResolver }
func (c validatorConfig) StakeAmount() uint64         { return c.stakeAmount }
func (c validatorConfig) TxMgrCfg() txmgr.Config      { return c.txMgrCfg }

func newValidatorConfigFromCLI(
	ctx *cli.Context,
	passphrase string,
	txMgrCfg txmgr.Config,
) validatorConfig {
	return validatorConfig{
		accountAddr:        common.HexToAddress(ctx.String(validatorAddrFlag.Name)),
		passphrase:         passphrase,
		isActiveStaker:     ctx.Bool(validatorIsActiveStakerFlag.Name),
		isActiveCreator:    ctx.Bool(validatorIsActiveCreatorFlag.Name),
		isActiveChallenger: ctx.Bool(validatorIsActiveChallengerFlag.Name),
		isResolver:         ctx.Bool(validatorIsResolverFlag.Name),
		stakeAmount:        ctx.Uint64(validatorStakeAmountFlag.Name),
		txMgrCfg:           txMgrCfg,
	}
}

// Driver configuration
type driverConfig struct {
	stepInterval time.Duration // Time between driver steps (in steady state; failures may trigger a longer backoff)
	retryDelay   time.Duration // Time to wait before retrying a step attempt
	numAttempts  uint          // Number of attempts to attempt driver step before catastrophically failing. Must be > 0.
}

func (c driverConfig) StepInterval() time.Duration { return c.stepInterval }
func (c driverConfig) RetryDelay() time.Duration   { return c.retryDelay }
func (c driverConfig) NumAttempts() uint           { return c.numAttempts }

func newDriverConfigFromCLI(cliCtx *cli.Context) driverConfig {
	return driverConfig{
		stepInterval: time.Duration(cliCtx.Uint(driverStepIntervalFlag.Name)) * time.Second,
		retryDelay:   time.Duration(cliCtx.Uint(driverRetryDelayFlag.Name)) * time.Second,
		numAttempts:  cliCtx.Uint(driverNumAttemptsFlag.Name),
	}
}
