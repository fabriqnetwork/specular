package assertion

import (
	"context"
	"fmt"
	"math/big"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
)

const cacheSize = 64

type RollupState struct {
	LastCreatedAssertionID   *big.Int
	LastResolvedAssertionID  *big.Int
	LastConfirmedAssertionID *big.Int
}

type AssertionAux struct {
	l1Block      *big.Int
	l2BlockStart *big.Int
	l2BlockEnd   *big.Int
}

type RollupClient interface {
}

type AssertionManager struct {
	assertionCache *lru.Cache[*big.Int, *bindings.IRollupAssertion]
	assertionAux   *lru.Cache[*big.Int, *AssertionAux]

	client RollupClient
	// systemState *state.SystemState
}

func NewAssertionManager(client RollupClient) (*AssertionManager, error) {
	cache, err := lru.New[*big.Int, *bindings.IRollupAssertion](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("Failed to create assertion cache, err: %w", err)
	}
	auxCache, err := lru.New[*big.Int, *AssertionAux](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("Failed to create assertion aux cache, err: %w", err)
	}
	return &AssertionManager{
		cache,
		auxCache,
		client,
		// systemState,
	}, nil
}

func (m *AssertionManager) GetAssertion(
	ctx context.Context,
	assertionID *big.Int,
) (*bindings.IRollupAssertion, error) {
	if assertion, ok := m.assertionCache.Get(assertionID); ok {
		return assertion, nil
	}
	// if true {
	//
	// }
	return nil, nil
}

func (m *AssertionManager) GetAssertionAux(
	ctx context.Context,
	assertionID *big.Int,
) (*AssertionAux, error) {
	if aux, ok := m.assertionAux.Get(assertionID); ok {
		return aux, nil
	}
	return nil, nil
}
