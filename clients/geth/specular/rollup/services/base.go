package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
)

type BaseService struct {
	Config *Config

	Eth          Backend
	ProofBackend proof.Backend
	Chain        *core.BlockChain
	L1           *ethclient.Client
	TransactOpts *bind.TransactOpts
	Inbox        *bindings.ISequencerInboxSession
	Rollup       *bindings.IRollupSession
	AssertionMap *bindings.AssertionMapCallerSession

	Ctx    context.Context
	Cancel context.CancelFunc
	Wg     sync.WaitGroup
}

func NewBaseService(eth Backend, proofBackend proof.Backend, cfg *Config, auth *bind.TransactOpts) (*BaseService, error) {
	if eth == nil {
		return nil, fmt.Errorf("can not use light client with rollup")
	}
	ctx, cancel := context.WithCancel(context.Background())
	l1, err := ethclient.DialContext(ctx, cfg.L1Endpoint)
	if err != nil {
		cancel()
		return nil, err
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

// Starts the rollup service.
// If cleanL1 is true, the service will only start from a clean L1 history.
// If stake is true, the service will ensure the service is already staked or will stake.
func (b *BaseService) Start(cleanL1, stake bool) error {
	// Check if we are at genesis
	// TODO: if not, sync from L1
	currentBlock := b.Eth.BlockChain().CurrentBlock()
	if currentBlock.NumberU64() != 0 {
		log.Crit("Rollup service can only start from clean history")
	}
	log.Info("Genesis root", "root", currentBlock.Root())
	// TODO: handle existing L1.
	if cleanL1 {
		inboxSize, err := b.Inbox.GetInboxSize()
		if err != nil {
			log.Crit("Failed to get initial inbox size", "err", err)
		}
		if inboxSize.Cmp(common.Big0) != 0 {
			log.Crit("Rollup service can only start from genesis")
		}
	}
	// Ensure node is staked.
	if stake {
		isStaked, err := b.Rollup.Contract.IsStaked(&bind.CallOpts{Pending: true, Context: b.Ctx}, b.TransactOpts.From)
		if err != nil {
			return err
		}
		if !isStaked {
			stakeOpts := b.Rollup.TransactOpts
			stakeOpts.Value = big.NewInt(int64(b.Config.RollupStakeAmount))
			_, err := b.Rollup.Contract.Stake(&stakeOpts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
