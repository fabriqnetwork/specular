package sequencer

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	proof_mock "github.com/specularl2/specular/clients/geth/specular/proof/mock"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	rollup_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
)

func mockSequencer(t *testing.T, ctrl *gomock.Controller) *Sequencer {
	cfg := &services.Config{
		Node:      services.NODE_SEQUENCER,
		L1ChainID: 123,
	}
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123))
	eth := rollup_mock.NewMockBackend(ctrl)
	proofBackend := proof_mock.NewMockBackend(ctrl)
	sequencer, err := NewSequencer(eth, proofBackend, cfg, auth)
	if err != nil {
		t.Errorf("Sequencer mock initialization failed. Error: %s", err)
	}
	return sequencer
}

func TestSequencerStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSequencer(t, ctrl)
}

func TestSequencerStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sequencer := mockSequencer(t, ctrl)
	time.Sleep(10 * time.Millisecond)
	sequencer.Stop()
}
