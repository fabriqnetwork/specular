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
	"github.com/ethereum/go-ethereum/eth"
	ethcatalyst "github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/les"
	lescatalyst "github.com/ethereum/go-ethereum/les/catalyst"
	"github.com/ethereum/go-ethereum/node"
	"github.com/specularL2/specular/services/sidecar/internal/ethapi"
	"github.com/specularL2/specular/services/sidecar/proof"
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
