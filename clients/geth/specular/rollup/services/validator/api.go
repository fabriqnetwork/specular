package validator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

type ValidatorServiceConfig interface {
	Validator() *services.ValidatorConfig
	L1() *services.L1Config
}

type AssertionManager interface {
	GetAssertion(ctx context.Context, assertionID *big.Int) (*bindings.IRollupAssertion, error)
}

type TxManager interface {
	Send(ctx context.Context, candidate txmgr.TxCandidate) (*types.Receipt, error)
}

type L2Client interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	Close()
}
