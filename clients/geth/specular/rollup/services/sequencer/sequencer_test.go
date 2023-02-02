package sequencer

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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	//"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/golang/mock/gomock"
	proof_mock "github.com/specularl2/specular/clients/geth/specular/proof/mock"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	//rollup_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
)

type mockBackend struct {
	bc     *core.BlockChain
	txPool *core.TxPool
}

func NewMockBackend(bc *core.BlockChain, txPool *core.TxPool) *mockBackend {
	return &mockBackend{
		bc:     bc,
		txPool: txPool,
	}
}

func (m *mockBackend) BlockChain() *core.BlockChain {
	return m.bc
}

func (m *mockBackend) TxPool() *core.TxPool {
	return m.txPool
}

func (m *mockBackend) StateAtBlock(block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error) {
	return nil, nil
}

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.gasLimit,
	}, nil, nil, nil, trie.NewStackTrie(nil))
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

// *MockSequencer
func mockSequencer(t *testing.T, ctrl *gomock.Controller) *Sequencer {
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
	blockchain := &testBlockChain{statedb, 10000000, new(event.Feed)}

	pool := core.NewTxPool(core.DefaultTxPoolConfig, chainConfig, blockchain)
	backend := NewMockBackend(bc, pool)

	proofBackend := proof_mock.NewMockBackend(ctrl)

	sequencer, err := NewSequencer(backend, proofBackend, cfg, auth)

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
