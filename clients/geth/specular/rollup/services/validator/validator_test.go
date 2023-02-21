package validator

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	//"github.com/ethereum/go-ethereum/log"
	"github.com/golang/mock/gomock"
	proof_mock "github.com/specularl2/specular/clients/geth/specular/proof/mock"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	testing_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/testutils"
)

func mockValidator(t *testing.T, ctrl *gomock.Controller) *Validator {
	cfg := &services.Config{
		Node:               services.NODE_VALIDATOR,
		L1Endpoint:         "http://localhost:8545",
		L1ChainID:          31337,
		SequencerInboxAddr: common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"),
		RollupAddr:         common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F"),
	}
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123))
	//eth := rollup_mock.NewMockBackend(ctrl)
	//proofBackend := proof_mock.NewMockBackend(ctrl)

	// Create chainConfig
	chainDB := rawdb.NewMemoryDatabase()
	genesis := core.DeveloperGenesisBlock(15, 11_500_000, common.HexToAddress("12345"))
	chainConfig, _, err := core.SetupGenesisBlock(chainDB, genesis)
	if err != nil {
		t.Fatalf("can't create new chain config: %v", err)
	}

	// Create consensus engine
	engine := clique.New(chainConfig.Clique, chainDB)

	// Create Ethereum backend
	bc, err := core.NewBlockChain(chainDB, nil, chainConfig, engine, vm.Config{}, nil, nil)

	if err != nil {
		t.Fatalf("can't create new chain %v", err)
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(chainDB), nil)
	blockchain := &testing_mock.TestBlockChain{statedb, 10000000, new(event.Feed)}

	pool := core.NewTxPool(core.DefaultTxPoolConfig, chainConfig, blockchain)
	backend := testing_mock.NewMockBackend(bc, pool)

	proofBackend := proof_mock.NewMockBackend(ctrl)

	validator, err := NewValidator(backend, proofBackend, cfg, auth)
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
	validator := mockValidator(t, ctrl)
	time.Sleep(10 * time.Millisecond)
	validator.Stop()
}
