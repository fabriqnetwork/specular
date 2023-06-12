package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type genesisConfig struct {
	genesisBlockID types.BlockID
}

func NewGenesisConfig(genesisBlockID types.BlockID) genesisConfig {
	return genesisConfig{genesisBlockID: genesisBlockID}
}

func (c genesisConfig) GetGenesisL1BlockID() types.BlockID { return c.genesisBlockID }

func getGenesisL1BlockID(ctx context.Context, cfg L1Config, l1Client *eth.EthClient) (types.BlockID, error) {
	genesisHeader, err := l1Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(cfg.GetRollupGenesisBlock()))
	if err != nil {
		return types.BlockID{}, fmt.Errorf("failed to get genesis header: %v", err)
	}
	return types.NewBlockID(cfg.GetRollupGenesisBlock(), genesisHeader.Hash()), nil
}
