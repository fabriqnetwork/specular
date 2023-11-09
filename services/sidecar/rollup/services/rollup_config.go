package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	spTypes "github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
)

var (
	ErrBlockTimeZero                 = errors.New("block time cannot be 0")
	ErrMissingChannelTimeout         = errors.New("channel timeout must be set, this should cover at least a L1 block time")
	ErrInvalidSeqWindowSize          = errors.New("sequencing window size must at least be 2")
	ErrMissingGenesisL1Hash          = errors.New("genesis L1 hash cannot be empty")
	ErrMissingGenesisL2Hash          = errors.New("genesis L2 hash cannot be empty")
	ErrGenesisHashesSame             = errors.New("achievement get! rollup inception: L1 and L2 genesis cannot be the same")
	ErrMissingGenesisL2Time          = errors.New("missing L2 genesis time")
	ErrMissingBatcherAddr            = errors.New("missing genesis system config batcher address")
	ErrMissingOverhead               = errors.New("missing genesis system config overhead")
	ErrMissingScalar                 = errors.New("missing genesis system config scalar")
	ErrMissingGasLimit               = errors.New("missing genesis system config gas limit")
	ErrMissingBatchInboxAddress      = errors.New("missing batch inbox address")
	ErrMissingDepositContractAddress = errors.New("missing deposit contract address")
	ErrMissingL1ChainID              = errors.New("L1 chain ID must not be nil")
	ErrMissingL2ChainID              = errors.New("L2 chain ID must not be nil")
	ErrChainIDsSame                  = errors.New("L1 and L2 chain IDs must be different")
	ErrL1ChainIDNotPositive          = errors.New("L1 chain ID must be non-zero and positive")
	ErrL2ChainIDNotPositive          = errors.New("L2 chain ID must be non-zero and positive")
)

type Bytes32 [32]byte

func (b *Bytes32) UnmarshalJSON(text []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(b), text, b[:])
}

func (b *Bytes32) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("Bytes32", text, b[:])
}

func (b Bytes32) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

func (b Bytes32) String() string {
	return hexutil.Encode(b[:])
}

func (b Bytes32) TerminalString() string {
	return fmt.Sprintf("%x..%x", b[:3], b[29:])
}

type RollupConfig struct {
	// Genesis anchor point of the rollup
	Genesis Genesis `json:"genesis"`
	// Seconds per L2 block
	BlockTime uint64 `json:"block_time"`
	// Sequencer batches may not be more than MaxSequencerDrift seconds after
	// the L1 timestamp of the sequencing window end.
	//
	// Note: When L1 has many 1 second consecutive blocks, and L2 grows at fixed 2 seconds,
	// the L2 time may still grow beyond this difference.
	MaxSequencerDrift uint64 `json:"max_sequencer_drift"`
	// Number of epochs (L1 blocks) per sequencing window, including the epoch L1 origin block itself
	SeqWindowSize uint64 `json:"seq_window_size"`
	// Required to verify L1 signatures
	L1ChainID *big.Int `json:"l1_chain_id"`
	// Required to identify the L2 network and create p2p signatures unique for this chain.
	L2ChainID *big.Int `json:"l2_chain_id"`
	// Note: below addresses are part of the block-derivation process,
	// and required to be the same network-wide to stay in consensus.
	// L1 address that batches are sent to.
	BatchInboxAddress common.Address `json:"batch_inbox_address"`
}

type Genesis struct {
	// The L1 block that the rollup starts *after* (no derived transactions)
	L1 spTypes.BlockID `json:"l1"`
	// The L2 block the rollup starts from (no transactions, pre-configured state)
	L2 spTypes.BlockID `json:"l2"`
	// Timestamp of L2 block
	L2Time uint64 `json:"l2_time"`
	// Initial system configuration values.
	// The L2 genesis block may not include transactions, and thus cannot encode the config values,
	// unlike later L2 blocks.
	SystemConfig systemConfig `json:"system_config"`
}

// SystemConfig represents the rollup system configuration that carries over in every L2 block,
// and may be changed through L1 system config events.
// The initial SystemConfig at rollup genesis is embedded in the rollup configuration.
type systemConfig struct {
	// BatcherAddr identifies the batch-sender address used in batch-inbox data-transaction filtering.
	BatcherAddr common.Address `json:"batcherAddr"`
	// Overhead identifies the L1 fee overhead, and is passed through opaquely to op-geth.
	Overhead Bytes32 `json:"overhead"`
	// Scalar identifies the L1 fee scalar, and is passed through opaquely to op-geth.
	Scalar Bytes32 `json:"scalar"`
	// GasLimit identifies the L2 block gas limit
	GasLimit uint64 `json:"gasLimit"`
	// More fields can be added for future SystemConfig versions.
}

// NewDeployConfig reads a config file given a path on the filesystem.
func NewRollupConfig(path string) (*RollupConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("deploy config at %s not found: %w", path, err)
	}
	dec := json.NewDecoder(bytes.NewReader(file))
	// dec.DisallowUnknownFields()
	var config RollupConfig
	if err := dec.Decode(&config); err != nil {
		return nil, fmt.Errorf("cannot decode config: %w", err)
	}
	return &config, nil
}

// ValidateL1Config checks L1 config variables for errors.
func (cfg *RollupConfig) ValidateL1Config(ctx context.Context, client L1Client) error {
	// Validate the L1 Client Chain ID
	if err := cfg.CheckL1ChainID(ctx, client); err != nil {
		return err
	}
	// Validate the Rollup L1 Genesis Blockhash
	if err := cfg.CheckL1GenesisBlockHash(ctx, client); err != nil {
		return err
	}
	return nil
}

// ValidateL2Config checks L2 config variables for errors.
func (cfg *RollupConfig) ValidateL2Config(ctx context.Context, client L2Client) error {
	// Validate the L2 Client Chain ID
	if err := cfg.CheckL2ChainID(ctx, client); err != nil {
		return err
	}
	// Validate the Rollup L2 Genesis Blockhash
	if err := cfg.CheckL2GenesisBlockHash(ctx, client); err != nil {
		return err
	}
	return nil
}

func (cfg *RollupConfig) TargetBlockNumber(timestamp uint64) (num uint64, err error) {
	// subtract genesis time from timestamp to get the time elapsed since genesis, and then divide that
	// difference by the block time to get the expected L2 block number at the current time. If the
	// unsafe head does not have this block number, then there is a gap in the queue.
	genesisTimestamp := cfg.Genesis.L2Time
	if timestamp < genesisTimestamp {
		return 0, fmt.Errorf("did not reach genesis time (%d) yet", genesisTimestamp)
	}
	wallClockGenesisDiff := timestamp - genesisTimestamp
	// Note: round down, we should not request blocks into the future.
	blocksSinceGenesis := wallClockGenesisDiff / cfg.BlockTime
	return cfg.Genesis.L2.Number + blocksSinceGenesis, nil
}

type L1Client interface {
	ChainID(context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// CheckL1ChainID checks that the configured L1 chain ID matches the client's chain ID.
func (cfg *RollupConfig) CheckL1ChainID(ctx context.Context, client L1Client) error {
	id, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get L1 chain ID: %w", err)
	}
	if cfg.L1ChainID.Cmp(id) != 0 {
		return fmt.Errorf("incorrect L1 RPC chain id %d, expected %d", id, cfg.L1ChainID)
	}
	return nil
}

// CheckL1GenesisBlockHash checks that the configured L1 genesis block hash is valid for the given client.
func (cfg *RollupConfig) CheckL1GenesisBlockHash(ctx context.Context, client L1Client) error {
	l1GenesisBlockRef, err := client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(cfg.Genesis.L1.Number))
	if err != nil {
		return fmt.Errorf("failed to get L1 genesis blockhash: %w", err)
	}
	if l1GenesisBlockRef.Hash() != cfg.Genesis.L1.Hash {
		return fmt.Errorf("incorrect L1 genesis block hash %s, expected %s", l1GenesisBlockRef.Hash(), cfg.Genesis.L1.Hash)
	}
	return nil
}

type L2Client interface {
	ChainID(context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// CheckL2ChainID checks that the configured L2 chain ID matches the client's chain ID.
func (cfg *RollupConfig) CheckL2ChainID(ctx context.Context, client L2Client) error {
	id, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get L2 chain ID: %w", err)
	}
	if cfg.L2ChainID.Cmp(id) != 0 {
		return fmt.Errorf("incorrect L2 RPC chain id %d, expected %d", id, cfg.L2ChainID)
	}
	return nil
}

// CheckL2GenesisBlockHash checks that the configured L2 genesis block hash is valid for the given client.
func (cfg *RollupConfig) CheckL2GenesisBlockHash(ctx context.Context, client L2Client) error {
	l2GenesisBlockRef, err := client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(cfg.Genesis.L2.Number))
	if err != nil {
		return fmt.Errorf("failed to get L2 genesis blockhash: %w", err)
	}
	if l2GenesisBlockRef.Hash() != cfg.Genesis.L2.Hash {
		return fmt.Errorf("incorrect L2 genesis block hash %s, expected %s", l2GenesisBlockRef.Hash(), cfg.Genesis.L2.Hash)
	}
	return nil
}

// Check verifies that the given configuration makes sense
func (cfg *RollupConfig) Check() error {
	if cfg.BlockTime == 0 {
		return ErrBlockTimeZero
	}
	if cfg.SeqWindowSize < 2 {
		return ErrInvalidSeqWindowSize
	}
	if cfg.Genesis.L1.Hash == (common.Hash{}) {
		return ErrMissingGenesisL1Hash
	}
	if cfg.Genesis.L2.Hash == (common.Hash{}) {
		return ErrMissingGenesisL2Hash
	}
	if cfg.Genesis.L2.Hash == cfg.Genesis.L1.Hash {
		return ErrGenesisHashesSame
	}
	if cfg.Genesis.L2Time == 0 {
		return ErrMissingGenesisL2Time
	}
	if cfg.Genesis.SystemConfig.BatcherAddr == (common.Address{}) {
		return ErrMissingBatcherAddr
	}
	if cfg.Genesis.SystemConfig.Overhead == (Bytes32{}) {
		return ErrMissingOverhead
	}
	if cfg.Genesis.SystemConfig.Scalar == (Bytes32{}) {
		return ErrMissingScalar
	}
	if cfg.Genesis.SystemConfig.GasLimit == 0 {
		return ErrMissingGasLimit
	}
	if cfg.BatchInboxAddress == (common.Address{}) {
		return ErrMissingBatchInboxAddress
	}
	if cfg.L1ChainID == nil {
		return ErrMissingL1ChainID
	}
	if cfg.L2ChainID == nil {
		return ErrMissingL2ChainID
	}
	if cfg.L1ChainID.Cmp(cfg.L2ChainID) == 0 {
		return ErrChainIDsSame
	}
	if cfg.L1ChainID.Sign() < 1 {
		return ErrL1ChainIDNotPositive
	}
	if cfg.L2ChainID.Sign() < 1 {
		return ErrL2ChainIDNotPositive
	}
	return nil
}
