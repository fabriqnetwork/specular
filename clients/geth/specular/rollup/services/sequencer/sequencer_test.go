package sequencer

import (
	"context"
	//"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/ethclient"
	//ethclient_test "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	//"github.com/ethereum/go-ethereum/log"
	"github.com/golang/mock/gomock"
	proof_mock "github.com/specularl2/specular/clients/geth/specular/proof/mock"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	rollup_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
	testing_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/testutils"
)

// type mockL1Client struct {
// 	expectedL1Endpoint string
// }

// func (m *mockL1Client) DialContext(ctx context.Context, endpoint string) (*ethclient.Client, error) {
// 	if m.expectedL1Endpoint != endpoint {
// 		return nil, fmt.Errorf("incorrect endpoint, expected %s but got %s", m.expectedL1Endpoint, endpoint)
// 	}
// 	// return a mock instance of ethclient.Client here
// 	return &ethclient.Client{}, nil
// }

// func (m *mockL1Client) HeaderByNumber(ctx context.Context, blockNr *big.Int) (*types.Header, error) {
// 	// implement the ethclient.Client interface
// 	return nil, nil
// }

// *MockSequencer
func mockSequencer(t *testing.T, ctrl *gomock.Controller) *Sequencer {
	// Use l1Client in below cfg
	cfg := &services.Config{
		Node:               services.NODE_SEQUENCER,
		L1Endpoint:         "http://localhost:8545",
		L1ChainID:          31337,
		SequencerInboxAddr: common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"),
		RollupAddr:         common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F"),
		UseMockedL1Client:  true,
	}
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123))

	//eth := rollup_mock.NewMockBackend(ctrl)

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

	// Create mocked L1 client
	// var mockedL1Client ethclient.Client = &mockL1Client{expectedL1Endpoint: "http://localhost:8545"}
	// mockedL1Client := &ethclient.Client{}
	//mockedL1Client := &testing_mock.MockL1Client{}

	l1config := &testing_mock.L1ClientConfig{
		HeaderByNumber: func(ctx context.Context, blockNr *big.Int) (*types.Header, error) {
			return nil, nil
		},
		BalanceAt: func(ctx context.Context, address common.Address, blockNumber *big.Int) (*big.Int, error) {
			return nil, nil
		},
		EstimateGas: func(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
			return 21000, nil
		},
		NonceAt: func(ctx context.Context, address common.Address, blockNumber *big.Int) (uint64, error) {
			return 0, nil
		},
		SendTransaction: func(ctx context.Context, signedTx *types.Transaction) error {
			return nil
		},
		TransactionReceipt: func(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
			return nil, nil
		},
	}

	mockedL1Client := rollup_mock.NewMockEthClient(ctrl)
	// mockedL1Client := testing_mock.NewL1Client(*l1config)

	//mockedL1Client := ethclient_test.TestEthClient(t)
	//mockedL1Client := testing_mock.TestEthClient(t)
	sequencer, err := NewSequencer(backend, proofBackend, cfg, auth, *mockedL1Client)

	if err != nil {
		t.Errorf("Sequencer mock initialization failed. Error: %s", err)
	}
	return sequencer
}

func TestSequencerMock(t *testing.T) {
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

// func TestSequencerStart(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	sequencer := mockSequencer(t, ctrl)
// 	time.Sleep(1000 * time.Millisecond)
// 	sequencer.Start()
// 	time.Sleep(5000 * time.Millisecond)
// 	sequencer.Stop()
// }
