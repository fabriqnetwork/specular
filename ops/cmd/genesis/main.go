package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularL2/specular/ops/genesis"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	app := cli.NewApp()
	app.Name = "spgenesis"
	app.Usage = "Generate specular genesis file"
	app.Action = GenerateSpecularGenesis
	app.Flags = Flags

	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "genesis-config",
		Usage:    "Path to the genesis config file",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "out",
		Usage:    "L2 genesis output file",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "l1-rpc-url",
		Usage:    "L1 RPC URL",
		Required: true,
	},
	&cli.StringFlag{
		Name:  "export-hash",
		Usage: "Genesis hash output file",
	},
	&cli.Uint64Flag{
		Name:  "l1-block",
		Usage: "L1 block number",
	},
}

type exportedHash struct {
	BlockHash common.Hash `json:"blockHash"`
	StateRoot common.Hash `json:"stateRoot"`

}

func GenerateSpecularGenesis(ctx *cli.Context) error {
	client, err := ethclient.Dial(ctx.String("l1-rpc-url"))
	if err != nil {
		return fmt.Errorf("cannot dial %s: %w", ctx.String("l1-rpc-url"), err)
	}

	var l1StartBlock *types.Block
	if ctx.IsSet("l1-block") {
		l1StartBlock, err = client.BlockByNumber(ctx.Context, big.NewInt(int64(ctx.Uint64("l1-block"))))
		if err != nil {
			return fmt.Errorf("cannot get block %d: %w", ctx.Uint64("l1-block"), err)
		}
	} else {
		l1StartBlock, err = client.BlockByNumber(ctx.Context, big.NewInt(rpc.SafeBlockNumber.Int64()))
		if err != nil {
			return fmt.Errorf("cannot get safe block: %w", err)
		}
	}

	genesisConfig := ctx.String("genesis-config")
	log.Info("Genesis config", "path", genesisConfig)
	config, err := genesis.NewGenesisConfig(genesisConfig)
	if err != nil {
		return err
	}

	l2Genesis, err := genesis.BuildL2Genesis(ctx.Context, config, l1StartBlock)
	if err != nil {
		return err
	}
	if err := writeGenesisFile(ctx.String("out"), l2Genesis); err != nil {
		return err
	}

	if ctx.IsSet("export-hash") {
		blockHash := l2Genesis.ToBlock().Hash()
		stateRoot := l2Genesis.ToBlock().Root()

		if err := writeGenesisFile(ctx.String("export-hash"), exportedHash{blockHash, stateRoot}); err != nil {
			return err
		}
	}
	return nil
}

func writeGenesisFile(outfile string, input any) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(input)
}
