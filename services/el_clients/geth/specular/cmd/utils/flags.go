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
	// RollupFlags
	RollupNodeFlag = &cli.StringFlag{
		Name:  "rollup.node",
		Usage: "Start node as rollup sequencer",
		Value: "",
	}
	RollupCoinBaseFlag = &cli.StringFlag{
		Name:  "rollup.coinbase",
		Usage: "The sequencer/validator address to be unlocked (pass passphrash via --password)",
		Value: "",
	}
	RollupClefEndpointFlag = &cli.StringFlag{
		Name:  "rollup.clefendpoint",
		Usage: "The Endpoint of the Clef instance that should be used as a signer)",
		Value: "",
	}
	RollupL1EndpointFlag = &cli.StringFlag{
		Name:  "rollup.l1endpoint",
		Usage: "The api endpoint of L1 client",
		Value: "",
	}
	RollupL1ChainIDFlag = &cli.Uint64Flag{
		Name:  "rollup.l1chainid",
		Usage: "The chain ID of L1 client",
		Value: 31337,
	}
	RollupSequencerAddrFlag = &cli.StringFlag{
		Name:  "rollup.sequencer-addr",
		Usage: "The account address of sequencer",
		Value: "",
	}
	RollupSequencerInboxAddrFlag = &cli.StringFlag{
		Name:  "rollup.sequencer-inbox-addr",
		Usage: "The contract address of L1 sequencer inbox",
		Value: "",
	}
	RollupRollupAddrFlag = &cli.StringFlag{
		Name:  "rollup.rollup-addr",
		Usage: "The contract address of L1 rollup",
		Value: "",
	}
	RollupL1RollupGenesisBlock = &cli.Uint64Flag{
		Name:  "rollup.l1-rollup-genesis-block",
		Usage: "The block number of L1 rollup genesis block to sync from",
		Value: 0,
	}
	RollupRollupStakeAmount = &cli.Uint64Flag{
		Name:  "rollup.rollup-stake-amount",
		Usage: "Required staking amount",
		Value: 1000000000000000000,
	}
	RollupL1FeeOverheadFlag = &cli.Int64Flag{
		Name:  "rollup.l1-fee-overhead",
		Usage: "Gas cost of sequencing a Tx",
		Value: 0,
	}
	RollupL1FeeMultiplierFlag = &cli.Float64Flag{
		Name:  "rollup.l1-fee-multiplier",
		Usage: "Scalar value to increase the L1 Fee",
		Value: 1.5,
	}
	RollupL1OracleAddressFlag = &cli.StringFlag{
		Name:  "rollup.l1-oracle-address",
		Usage: "The address of the L1Oracle contract",
		Value: "0xff00000000000000000000000000000000000002",
	}
	RollupL1OracleBaseFeeSlotFlag = &cli.StringFlag{
		Name:  "rollup.l1-oracle-base-fee-slot",
		Usage: "The address of the L1Oracle contract",
		Value: "0x18b94da8c18f49ac05520153402a0591c3c917271b9d13711fd6fdb213ded168", // keccak256("specular.basefee")
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

func MakeRollupConfig(ctx *cli.Context, cfg *ethconfig.Config) *rollup.Config {
	utils.CheckExclusive(ctx, RollupNodeFlag, utils.MiningEnabledFlag)
	utils.CheckExclusive(ctx, RollupNodeFlag, utils.DeveloperFlag)

	clefEndpoint := ctx.String(RollupClefEndpointFlag.Name)
	var passphrase string
	if list := utils.MakePasswordList(ctx); len(list) > 0 {
		passphrase = list[0]
	} else if clefEndpoint == "" {
		utils.Fatalf("Failed to register the Rollup service: coinbase account locked")
	}

	node := ctx.String(RollupNodeFlag.Name)
	coinbase := common.HexToAddress(ctx.String(RollupCoinBaseFlag.Name))
	sequencerAddr := common.HexToAddress(ctx.String(RollupSequencerAddrFlag.Name))
	if node == "sequencer" && sequencerAddr == (common.Address{}) {
		sequencerAddr = coinbase
	}
	if sequencerAddr == (common.Address{}) {
		utils.Fatalf("Failed to register the Rollup service: sequencer address not specified")
	}
	return &rollup.Config{
		Node:                 node,
		Coinbase:             coinbase,
		Passphrase:           passphrase,
		ClefEndpoint:         ctx.String(RollupClefEndpointFlag.Name),
		L1Endpoint:           ctx.String(RollupL1EndpointFlag.Name),
		L1ChainID:            ctx.Uint64(RollupL1ChainIDFlag.Name),
		L2ChainID:            cfg.NetworkId,
		SequencerAddr:        sequencerAddr,
		SequencerInboxAddr:   common.HexToAddress(ctx.String(RollupSequencerInboxAddrFlag.Name)),
		RollupAddr:           common.HexToAddress(ctx.String(RollupRollupAddrFlag.Name)),
		L1RollupGenesisBlock: ctx.Uint64(RollupL1RollupGenesisBlock.Name),
		RollupStakeAmount:    ctx.Uint64(RollupRollupStakeAmount.Name),
		L1FeeOverhead:        RollupL1FeeOverheadFlag.Value,
		L1FeeMultiplier:      RollupL1FeeMultiplierFlag.Value,
		L1OracleAddress:      common.HexToAddress(RollupL1OracleAddressFlag.Value),
		L1OracleBaseFeeSlot:  common.HexToHash(RollupL1OracleBaseFeeSlotFlag.Value),
	}
}
