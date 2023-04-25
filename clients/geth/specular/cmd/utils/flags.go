// This file is modified for Specular under the terms of the GNU
// General Public License. Major modifications are marked with
// <specular modification><specular modification/>.

// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// Package utils contains internal helper functions for go-ethereum commands.
package utils

import (
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth"
	ethcatalyst "github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/les"
	lescatalyst "github.com/ethereum/go-ethereum/les/catalyst"
	"github.com/ethereum/go-ethereum/node"
	"github.com/specularl2/specular/clients/geth/specular/internal/ethapi"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	rollup "github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/urfave/cli/v2"
)

var (
	// <specular modification>
	// L1 config flags
	RollupL1EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l1-endpoint",
		Usage: "The api endpoint of L1 client",
		Value: "",
	}
	RollupL1ChainIDFlag = &cli.Uint64Flag{
		Name:  "rollup.l1-chainid",
		Usage: "The chain ID of L1 client",
		Value: 31337,
	}
	RollupL1RollupGenesisBlockFlag = &cli.Uint64Flag{
		Name:  "rollup.l1-rollup-genesis-block",
		Usage: "The block number of L1 rollup genesis block to sync from",
		Value: 0,
	}
	RollupSequencerInboxAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1-sequencer-inbox-addr",
		Usage: "The contract address of L1 sequencer inbox",
		Value: "",
	}
	RollupRollupAddrFlag = &cli.StringFlag{
		Name:  "rollup.l1-rollup-addr",
		Usage: "The contract address of L1 rollup",
		Value: "",
	}
	// Sequencer config flags
	RollupSequencerAddrFlag = &cli.StringFlag{
		Name:  "rollup.sequencer-addr",
		Usage: "The sequencer address to be unlocked (pass passphrash via --password)",
		Value: "",
	}
	RollupSequencerMinExecutionIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.sequencer-execution-interval",
		Usage: "Minimum time between block executions (seconds)",
		Value: 0,
	}
	RollupSequencerMaxExecutionIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.sequencer-execution-interval",
		Usage: "Maximum time between block executions (seconds)",
		Value: 1,
	}
	RollupSequencerSequencingIntervalFlag = &cli.Uint64Flag{
		Name:  "rollup.sequencer-sequencing-interval",
		Usage: "Time between batch sequencing attempts (seconds)",
		Value: 5,
	}
	// Validator config flags
	RollupValidatorAddrFlag = &cli.StringFlag{
		Name:  "rollup.validator-addr",
		Usage: "The validator address to be unlocked (pass passphrash via --password)",
		Value: "",
	}
	RollupValidatorIsActiveCreatorFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-active-creator",
		Usage: "Whether the validator should be an active assertion creator",
		Value: false,
	}
	RollupValidatorIsActiveChallengerFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-active-challenger",
		Usage: "Whether the validator should be an active challenger (i.e. issue challenges)",
		Value: false,
	}
	RollupValidatorIsResolverFlag = &cli.BoolFlag{
		Name:  "rollup.validator-is-resolver",
		Usage: "Whether the validator should resolve (confirm/reject) assertions",
		Value: false,
	}
	// TODO: read this from the contract
	RollupRollupStakeAmountFlag = &cli.Uint64Flag{
		Name:  "rollup.rollup-stake-amount",
		Usage: "Required staking amount",
		Value: 1000000000000000000,
	}

	// Indexer config flags
	RollupIndexerAddrFlag = &cli.StringFlag{
		Name:  "rollup.indexer-addr",
		Usage: "The indexer address to be unlocked (pass passphrash via --password)",
		Value: "",
	}
	// <specular modification/>
)

// RegisterEthService adds an Ethereum client to the stack.
// The second return value is the full node instance, which may be nil if the
// node is running as a light client.
func RegisterEthService(stack *node.Node, cfg *ethconfig.Config) (ethapi.Backend, *eth.Ethereum) {
	if cfg.SyncMode == downloader.LightSync {
		backend, err := les.New(stack, cfg)
		if err != nil {
			utils.Fatalf("Failed to register the Ethereum service: %v", err)
		}
		stack.RegisterAPIs(tracers.APIs(backend.ApiBackend))
		if err := lescatalyst.Register(stack, backend); err != nil {
			utils.Fatalf("Failed to register the Engine API service: %v", err)
		}
		return backend.ApiBackend, nil
	}
	backend, err := eth.New(stack, cfg)
	if err != nil {
		utils.Fatalf("Failed to register the Ethereum service: %v", err)
	}
	if cfg.LightServ > 0 {
		_, err := les.NewLesServer(stack, backend, cfg)
		if err != nil {
			utils.Fatalf("Failed to create the LES server: %v", err)
		}
	}
	if err := ethcatalyst.Register(stack, backend); err != nil {
		utils.Fatalf("Failed to register the Engine API service: %v", err)
	}
	stack.RegisterAPIs(tracers.APIs(backend.APIBackend))
	// <specular modification>
	stack.RegisterAPIs(proof.APIs(backend.APIBackend))
	// <specular modification/>
	return backend.APIBackend, backend
}

// <specular modification>
func MakeRollupConfig(ctx *cli.Context) *rollup.SystemConfig {
	utils.CheckExclusive(ctx, RollupL1EndpointFlag, utils.MiningEnabledFlag)
	utils.CheckExclusive(ctx, RollupL1EndpointFlag, utils.DeveloperFlag)
	// Indexer must run standalone (for now).
	utils.CheckExclusive(ctx, RollupSequencerAddrFlag, RollupIndexerAddrFlag)
	utils.CheckExclusive(ctx, RollupValidatorAddrFlag, RollupIndexerAddrFlag)

	pwList := utils.MakePasswordList(ctx)
	if len(pwList) == 0 {
		utils.Fatalf("Failed to register rollup services: no password provided")
	}

	var sequencerAddr common.Address
	var sequencerPassphrase string
	if ctx.String(RollupSequencerAddrFlag.Name) != "" {
		sequencerAddr = common.HexToAddress(ctx.String(RollupSequencerAddrFlag.Name))
		sequencerPassphrase = pwList[0]
		pwList = pwList[1:]
	}
	var validatorAddr common.Address
	var validatorPassphrase string
	if ctx.String(RollupValidatorAddrFlag.Name) != "" {
		validatorAddr = common.HexToAddress(ctx.String(RollupValidatorAddrFlag.Name))
		validatorPassphrase = pwList[0]
		pwList = pwList[1:]
	}
	var indexerAddr common.Address
	var indexerPassphrase string
	if ctx.String(RollupIndexerAddrFlag.Name) != "" {
		indexerAddr = common.HexToAddress(ctx.String(RollupIndexerAddrFlag.Name))
		indexerPassphrase = pwList[0]
		pwList = pwList[1:]
	}

	return &rollup.SystemConfig{
		L1Config: rollup.L1Config{
			L1Endpoint:           ctx.String(RollupL1EndpointFlag.Name),
			L1ChainID:            ctx.Uint64(RollupL1ChainIDFlag.Name),
			L1RollupGenesisBlock: ctx.Uint64(RollupL1RollupGenesisBlockFlag.Name),
			SequencerInboxAddr:   common.HexToAddress(ctx.String(RollupSequencerInboxAddrFlag.Name)),
			RollupAddr:           common.HexToAddress(ctx.String(RollupRollupAddrFlag.Name)),
		},
		SequencerConfig: rollup.SequencerConfig{
			SequencerAccountAddr: sequencerAddr,
			SequencerPassphrase:  sequencerPassphrase,
			MinExecutionInterval: time.Duration(ctx.Uint64(RollupSequencerMinExecutionIntervalFlag.Name)) * time.Second,
			MaxExecutionInterval: time.Duration(ctx.Uint64(RollupSequencerMaxExecutionIntervalFlag.Name)) * time.Second,
			SequencingInterval:   time.Duration(ctx.Uint64(RollupSequencerSequencingIntervalFlag.Name)) * time.Second,
		},
		ValidatorConfig: rollup.ValidatorConfig{
			ValidatorAccountAddr: validatorAddr,
			ValidatorPassphrase:  validatorPassphrase,
			RollupStakeAmount:    ctx.Uint64(RollupRollupStakeAmountFlag.Name),
			IsActiveCreator:      ctx.Bool(RollupValidatorIsActiveCreatorFlag.Name),
			IsActiveChallenger:   ctx.Bool(RollupValidatorIsActiveChallengerFlag.Name),
			IsResolver:           ctx.Bool(RollupValidatorIsResolverFlag.Name),
		},
		IndexerConfig: rollup.IndexerConfig{
			IndexerAccountAddr: indexerAddr,
			IndexerPassphrase:  indexerPassphrase,
		},
	}
}

// <specular modification/>
