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
}

func (c *SystemConfig) L1() *L1Config               { return &c.L1Config }
func (c *SystemConfig) L2() *L2Config               { return &c.L2Config }
func (c *SystemConfig) Sequencer() *SequencerConfig { return &c.SequencerConfig }
func (c *SystemConfig) Validator() *ValidatorConfig { return &c.ValidatorConfig }

type L1Config struct {
	L1Endpoint           string         // L1 API endpoint
	L1ChainID            uint64         // L1 chain ID
	L1RollupGenesisBlock uint64         // L1 Rollup genesis block
	SequencerInboxAddr   common.Address // L1 SequencerInbox contract address
	RollupAddr           common.Address // L1 Rollup contract address
}

type L2Config struct {
	L2Endpoint     string // L2 API endpoint
	L2ClefEndpoint string // The Clef Endpoint used for signing TXs
}

type SequencerConfig struct {
	SequencerAccountAddr common.Address
	SequencerPassphrase  string        // The passphrase of the sequencer account
	MinExecutionInterval time.Duration // Minimum time between block production. If 0, txs executed immediately -- FCFS.
	MaxExecutionInterval time.Duration // Maximum time between block production. Must be >= `MinExecutionInterval`.
	SequencingInterval   time.Duration // Time between batch sequencing attempts
}

type ValidatorConfig struct {
	ValidatorAccountAddr common.Address
	ValidatorPassphrase  string // The passphrase of the validator account
	IsActiveStaker       bool   // Iff true, actively stakes on rollup contract
	IsActiveCreator      bool   // Iff true, actively tries to create new assertions (not just for a challenge).
	IsActiveChallenger   bool   // Iff true, actively issues challenges as challenger. *Defends* against challenges regardless.
	IsResolver           bool   // Iff true, attempts to resolve assertions (by confirming or rejecting)
	StakeAmount          uint64 // Size of stake to deposit to rollup contract
}
