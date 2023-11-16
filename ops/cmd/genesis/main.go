package main

import (
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularL2/specular/ops/genesis"
	"github.com/urfave/cli/v2"
)

func main() {
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
}

func GenerateSpecularGenesis(ctx *cli.Context) error {
	fakeHeader := &types.Header{
		Number:  big.NewInt(0),
		BaseFee: big.NewInt(0),
	}
	log.Info("header", "header", fakeHeader)
	fakeBlock := types.NewBlockWithHeader(fakeHeader)

	genesisConfig := ctx.String("genesis-config")
	log.Info("Genesis config", "path", genesisConfig)
	config, err := genesis.NewGenesisConfig(genesisConfig)
	if err != nil {
		return err
	}

	l2Genesis, err := genesis.BuildL2Genesis(ctx.Context, config, fakeBlock)
	if err != nil {
		return err
	}
	if err := writeGenesisFile(ctx.String("out"), l2Genesis); err != nil {
		return err
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
