package sequencer

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	proof_mock "github.com/specularl2/specular/clients/geth/specular/proof/mock"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
	rollup_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
)

func mockSequencer(t *testing.T, ctrl *gomock.Controller, eth *rollup_mock.MockBackend) *MockSequencer {
	cfg := &services.Config{
		Node:               services.NODE_SEQUENCER,
		L1Endpoint:         "http://localhost:8545",
		L1ChainID:          31337,
		SequencerInboxAddr: common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"),
		RollupAddr:         common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F"),
	}
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123))
	//eth := rollup_mock.NewMockBackend(ctrl)
	proofBackend := proof_mock.NewMockBackend(ctrl)
	sequencer, err := NewMockSequencer(eth, proofBackend, cfg, auth)
	if err != nil {
		t.Errorf("Sequencer mock initialization failed. Error: %s", err)
	}
	return sequencer
}

func TestSequencerStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := mock.NewMockBackend(ctrl)
	mockBackend.EXPECT().BlockChain().Return(nil)
	mockSequencer(t, ctrl, mockBackend)
}

func TestSequencerStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// mockBackend := mock.NewMockBackend(ctrl)
	// sequencer := mockSequencer(t, ctrl, mockBackend)
	// time.Sleep(10 * time.Millisecond)
	// sequencer.Stop()
}
