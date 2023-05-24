package services

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type SystemConfig struct {
	L1Config
	L2Config
	SequencerConfig
	ValidatorConfig
	DriverConfig
}

func NewSystemConfig(
	// l1
	l1Endpoint string,
	l1ChainID uint64,
	l1RollupGenesisBlock uint64,
	l1SequencerInboxAddr common.Address,
	l1RollupAddr common.Address,
	// l2
	l2Endpoint string,
	l2ClefEndpoint string,
	// sequencer
	sequencerAddr common.Address,
	sequencerPassphrase string,
	minExecutionInterval time.Duration,
	maxExecutionInterval time.Duration,
	sequencingInterval time.Duration,
	// validator
	validatorAddr common.Address,
	validatorPassphrase string,
	validatorIsActiveStaker bool,
	validatorIsActiveCreator bool,
	validatorIsActiveChallenger bool,
	validatorIsActiveResolver bool,
	validatorStakeAmount uint64,
	// driver
	driverStepInterval time.Duration,
	driverRetryDelay time.Duration,
	driverNumAttempts uint,
) *SystemConfig {
	return &SystemConfig{
		L1Config:        L1Config{l1Endpoint, l1ChainID, l1RollupGenesisBlock, l1SequencerInboxAddr, l1RollupAddr},
		L2Config:        L2Config{l2Endpoint, l2ClefEndpoint},
		SequencerConfig: SequencerConfig{sequencerAddr, sequencerPassphrase, minExecutionInterval, maxExecutionInterval, sequencingInterval},
		ValidatorConfig: ValidatorConfig{validatorAddr, validatorPassphrase, validatorIsActiveStaker, validatorIsActiveCreator, validatorIsActiveChallenger, validatorIsActiveResolver, validatorStakeAmount},
		DriverConfig:    DriverConfig{driverStepInterval, driverRetryDelay, driverNumAttempts},
	}
}

func (c *SystemConfig) L1() *L1Config               { return &c.L1Config }
func (c *SystemConfig) L2() *L2Config               { return &c.L2Config }
func (c *SystemConfig) Sequencer() *SequencerConfig { return &c.SequencerConfig }
func (c *SystemConfig) Validator() *ValidatorConfig { return &c.ValidatorConfig }
func (c *SystemConfig) Driver() *DriverConfig       { return &c.DriverConfig }

// L1 blockchain configuration
type L1Config struct {
	endpoint           string         // L1 API endpoint
	chainID            uint64         // L1 chain ID
	rollupGenesisBlock uint64         // L1 Rollup genesis block
	sequencerInboxAddr common.Address // L1 SequencerInbox contract address
	rollupAddr         common.Address // L1 Rollup contract address
}

func (c *L1Config) Endpoint() string                   { return c.endpoint }
func (c *L1Config) ChainID() uint64                    { return c.chainID }
func (c *L1Config) RollupGenesisBlock() uint64         { return c.rollupGenesisBlock }
func (c *L1Config) SequencerInboxAddr() common.Address { return c.sequencerInboxAddr }
func (c *L1Config) RollupAddr() common.Address         { return c.rollupAddr }

// L2 blockchain configuration
type L2Config struct {
	endpoint     string // L2 API endpoint
	clefEndpoint string // The Clef Endpoint used for signing TXs
}

func (c *L2Config) Endpoint() string     { return c.endpoint }
func (c *L2Config) ClefEndpoint() string { return c.clefEndpoint }

type SequencerConfig struct {
	accountAddr          common.Address
	passphrase           string        // The passphrase of the sequencer account
	minExecutionInterval time.Duration // Minimum time between block production. If 0, txs executed immediately -- FCFS.
	maxExecutionInterval time.Duration // Maximum time between block production. Must be >= `MinExecutionInterval`.
	sequencingInterval   time.Duration // Time between batch sequencing attempts
}

func (c *SequencerConfig) AccountAddr() common.Address         { return c.accountAddr }
func (c *SequencerConfig) Passphrase() string                  { return c.passphrase }
func (c *SequencerConfig) MinExecutionInterval() time.Duration { return c.minExecutionInterval }
func (c *SequencerConfig) MaxExecutionInterval() time.Duration { return c.maxExecutionInterval }
func (c *SequencerConfig) SequencingInterval() time.Duration   { return c.sequencingInterval }

type ValidatorConfig struct {
	accountAddr        common.Address
	passphrase         string // The passphrase of the validator account
	isActiveStaker     bool   // Iff true, actively stakes on rollup contract
	isActiveCreator    bool   // Iff true, actively tries to create new assertions (not just for a challenge).
	isActiveChallenger bool   // Iff true, actively issues challenges as challenger. *Defends* against challenges regardless.
	isResolver         bool   // Iff true, attempts to resolve assertions (by confirming or rejecting)
	stakeAmount        uint64 // Size of stake to deposit to rollup contract
}

func (c *ValidatorConfig) AccountAddr() common.Address { return c.accountAddr }
func (c *ValidatorConfig) Passphrase() string          { return c.passphrase }
func (c *ValidatorConfig) IsActiveStaker() bool        { return c.isActiveStaker }
func (c *ValidatorConfig) IsActiveCreator() bool       { return c.isActiveCreator }
func (c *ValidatorConfig) IsActiveChallenger() bool    { return c.isActiveChallenger }
func (c *ValidatorConfig) IsResolver() bool            { return c.isResolver }
func (c *ValidatorConfig) StakeAmount() uint64         { return c.stakeAmount }

type DriverConfig struct {
	stepInterval time.Duration // Time between driver steps (in steady state; failures may trigger a longer backoff)
	retryDelay   time.Duration // Time to wait before retrying a step attempt
	numAttempts  uint          // Number of attempts to attempt driver step before catastrophically failing. Must be > 0.
}

func (c *DriverConfig) StepInterval() time.Duration { return c.stepInterval }
func (c *DriverConfig) RetryDelay() time.Duration   { return c.retryDelay }
func (c *DriverConfig) NumAttempts() uint           { return c.numAttempts }
