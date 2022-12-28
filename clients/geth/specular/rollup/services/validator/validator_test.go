package validator

import (
	"math/big"
	"testing"
	"time"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	proof_mock "github.com/specularl2/specular/clients/geth/specular/proof/mock"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	rollup_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
)

func mockValidator(t *testing.T, ctrl *gomock.Controller) *Validator {
	cfg := &services.Config{
		Node:      services.NODE_VALIDATOR,
		L1ChainID: 123,
	}
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123))
	eth := rollup_mock.NewMockBackend(ctrl)
	proofBackend := proof_mock.NewMockBackend(ctrl)
	validator, err := NewValidator(eth, proofBackend, cfg, auth)
	if err != nil {
		t.Errorf("Validator mock initialization failed. Error: %s", err)
	}
	return validator
}

func TestValidatorStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator(t, ctrl)
}

func TestValidatorStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sequencer := mockValidator(t, ctrl)
	time.Sleep(10 * time.Millisecond)
	sequencer.Stop()
}
