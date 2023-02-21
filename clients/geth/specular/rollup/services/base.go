package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	rollup_mock "github.com/specularl2/specular/clients/geth/specular/rollup/services/mock"
	testutils "github.com/specularl2/specular/clients/geth/specular/rollup/services/testutils"
)

type BaseService struct {
	Config *Config

	Eth          Backend
	ProofBackend proof.Backend
	Chain        *core.BlockChain
	L1           *testutils.EthClient
	TransactOpts *bind.TransactOpts
	Inbox        *bindings.ISequencerInboxSession
	Rollup       *bindings.IRollupSession
	AssertionMap *bindings.AssertionMapCallerSession

	Ctx    context.Context
	Cancel context.CancelFunc
	Wg     sync.WaitGroup
}

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

type L1Client interface {
	// required methods from *ethclient.Client
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	// additional methods for testing purposes
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

type ClientWithMock interface {
	*ethclient.Client
	rollup_mock.MockEthClient
	//bind.ContractBackend
}

// TOOD: double check gomock, does it mock for functions not defined in types?
func NewBaseService(eth Backend, proofBackend proof.Backend, cfg *Config, auth *bind.TransactOpts, L1Client testutils.EthClient) (*BaseService, error) {
	if eth == nil {
		return nil, fmt.Errorf("can not use light client with rollup")
	}
	ctx, cancel := context.WithCancel(context.Background())

	// var l1 rollup_mock.MockEthClient
	var l1 testutils.Client
	//var l1 ClientWithMock
	var err error
	if cfg.UseMockedL1Client {
		// l1, err = mockedL1Client.DialContext(ctx, "")
		// l1 = mockedL1Client
		l1 = mockedL1Client
		//l1. = mockedL1Client
		//log.info("this is err: ", "err", err)
	} else {
		// TODO: make this passed in
		// l1, err = ethclient.DialContext(ctx, cfg.L1Endpoint)
		ethClient, err := ethclient.DialContext(ctx, cfg.L1Endpoint)
		if err != nil {
			cancel()
			return nil, err
		}
	}

	callOpts := bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	transactOpts := bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: big.NewInt(800000000),
		Context:  ctx,
	}
	inbox, err := bindings.NewISequencerInbox(common.Address(cfg.SequencerInboxAddr), l1)
	if err != nil {
		cancel()
		return nil, err
	}
	inboxSession := &bindings.ISequencerInboxSession{
		Contract:     inbox,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	rollup, err := bindings.NewIRollup(common.Address(cfg.RollupAddr), l1)
	if err != nil {
		cancel()
		return nil, err
	}
	rollupSession := &bindings.IRollupSession{
		Contract:     rollup,
		CallOpts:     callOpts,
		TransactOpts: transactOpts,
	}
	assertionMapAddr, err := rollupSession.Assertions()
	if err != nil {
		cancel()
		return nil, err
	}
	assertionMap, err := bindings.NewAssertionMapCaller(assertionMapAddr, l1)
	if err != nil {
		cancel()
		return nil, err
	}
	assertionMapSession := &bindings.AssertionMapCallerSession{
		Contract: assertionMap,
		CallOpts: callOpts,
	}
	b := &BaseService{
		Config:       cfg,
		Eth:          eth,
		ProofBackend: proofBackend,
		L1:           l1,
		TransactOpts: &transactOpts,
		Inbox:        inboxSession,
		Rollup:       rollupSession,
		AssertionMap: assertionMapSession,
		Ctx:          ctx,
		Cancel:       cancel,
	}
	b.Chain = eth.BlockChain()
	return b, nil
}

// Start starts the rollup service
// If cleanL1 is true, the service will only start from a clean L1 history
// If stake is true, the service will try to stake on start
// Returns the genesis block
func (b *BaseService) Start(cleanL1, stake bool) *types.Block {
	// Check if we are at genesis
	// TODO: if not, sync from L1
	genesis := b.Eth.BlockChain().CurrentBlock()
	if genesis.NumberU64() != 0 {
		log.Crit("Rollup service can only start from clean history")
	}
	log.Info("Genesis root", "root", genesis.Root())

	if cleanL1 {
		inboxSize, err := b.Inbox.GetInboxSize()
		if err != nil {
			log.Crit("Failed to get initial inbox size", "err", err)
		}
		if inboxSize.Cmp(common.Big0) != 0 {
			log.Crit("Rollup service can only start from genesis")
		}
	}

	if stake {
		// Initial staking
		// TODO: sync L1 staking status
		stakeOpts := b.Rollup.TransactOpts
		stakeOpts.Value = big.NewInt(int64(b.Config.RollupStakeAmount))
		_, err := b.Rollup.Contract.Stake(&stakeOpts)
		if err != nil {
			log.Crit("Failed to stake", "err", err)
		}
	}
	return genesis
}
