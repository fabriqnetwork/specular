package services

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	NODE_SEQUENCER = "sequencer"
	NODE_VALIDATOR = "validator"
	NODE_INDEXER   = "indexer"
)

// Config is the configuration of rollup services
type Config struct {
	Node                 string         // Rollup node type, either sequencer or validator
	Coinbase             common.Address // The account used for L1 and L2 activity
	Passphrase           string         // The passphrase of the coinbase account
	ClefEndpoint         string         // The Clef Endpoint used for signing TXs
	L1Endpoint           string         // L1 API endpoint
	L1ChainID            uint64         // L1 chain ID
	L2ChainID            uint64         // L2 chain ID
	SequencerAddr        common.Address // Validator only
	SequencerInboxAddr   common.Address // L1 SequencerInbox contract address
	RollupAddr           common.Address // L1 Rollup contract address
	L1RollupGenesisBlock uint64         // L1 Rollup genesis block
	RollupStakeAmount    uint64         // Amount of stake
	L1FeeOverhead        int64          // Gas cost of sequencing a Tx
	L1FeeMultiplier      float64        // Scalar value to increase L1 fee
	L1OracleAddress      common.Address // L2 Address of the L1Oracle
	L1OracleBaseFeeSlot  common.Hash    // L1 basefee storage slot of the L1Oracle
}

func (c *Config) GetCoinbase() common.Address {
	return c.Coinbase
}

func (c *Config) GetL2ChainID() uint64 {
	return c.L2ChainID
}

func (c *Config) GetL1FeeOverhead() int64 {
	return c.L1FeeOverhead
}

func (c *Config) GetL1FeeMultiplier() float64 {
	return c.L1FeeMultiplier
}

func (c *Config) GetL1OracleAddress() common.Address {
	return c.L1OracleAddress
}

func (c *Config) GetL1OracleBaseFeeSlot() common.Hash {
	return c.L1OracleBaseFeeSlot
}
