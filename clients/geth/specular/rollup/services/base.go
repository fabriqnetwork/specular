package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

const timeInterval = 10 * time.Second

type BaseService struct {
	Config *Config

	Eth                Backend
	ProofBackend       proof.Backend
	Chain              *core.BlockChain
	L1                 *ethclient.Client
	TransactOpts       *bind.TransactOpts
	Inbox              *bindings.ISequencerInboxSession
	Rollup             *bindings.IRollupSession
	AssertionMap       *bindings.AssertionMapCallerSession
	BlockCh            chan types.Blocks
	ConfirmedIDCh      chan *big.Int
	PendingAssertionCh chan *rollupTypes.Assertion

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
		Config:             cfg,
		Eth:                eth,
		ProofBackend:       proofBackend,
		L1:                 l1,
		TransactOpts:       &transactOpts,
		Inbox:              inboxSession,
		Rollup:             rollupSession,
		AssertionMap:       assertionMapSession,
		BlockCh:            make(chan types.Blocks, 4096),
		ConfirmedIDCh:      make(chan *big.Int, 4096),
		PendingAssertionCh: make(chan *rollupTypes.Assertion, 4096),
		Ctx:                ctx,
		Cancel:             cancel,
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

func (b *BaseService) CreateDA(genesisRoot common.Hash) {
	defer b.Wg.Done()

	// Ticker
	var ticker = time.NewTicker(timeInterval)
	defer ticker.Stop()

	// Watch AssertionCreated event
	createdCh := make(chan *bindings.IRollupAssertionCreated, 4096)
	createdSub, err := b.Rollup.Contract.WatchAssertionCreated(&bind.WatchOpts{Context: b.Ctx}, createdCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer createdSub.Unsubscribe()

	// Current confirmed assertion, initalize it to genesis
	// TODO: sync from L1 Rollup
	confirmedAssertion := &rollupTypes.Assertion{
		ID:                    new(big.Int),
		VmHash:                genesisRoot,
		CumulativeGasUsed:     new(big.Int),
		InboxSize:             new(big.Int),
		Deadline:              new(big.Int),
		PrevCumulativeGasUsed: new(big.Int),
	}
	// Assertion created and pending for confirmation
	var pendingAssertion *rollupTypes.Assertion
	// Assertion to be created on L1 Rollup
	queuedAssertion := confirmedAssertion.Copy()
	queuedAssertion.StartBlock = 1

	// Create assertion on L1 Rollup
	commitAssertion := func() {
		pendingAssertion = queuedAssertion.Copy()
		queuedAssertion.StartBlock = queuedAssertion.EndBlock + 1
		queuedAssertion.PrevCumulativeGasUsed = new(big.Int).Set(queuedAssertion.CumulativeGasUsed)
		_, err = b.Rollup.CreateAssertion(
			pendingAssertion.VmHash,
			pendingAssertion.InboxSize,
			pendingAssertion.CumulativeGasUsed,
			confirmedAssertion.VmHash,
			confirmedAssertion.CumulativeGasUsed,
		)
		if err != nil {
			log.Error("Can not create DA", "error", err)
		}
	}

	// Blocks from the batchingLoop that will be sent to the inbox in the next tick
	var batchBlocks types.Blocks

	for {
		select {
		case <-ticker.C:
			if len(batchBlocks) == 0 {
				continue
			}
			batch := rollupTypes.NewTxBatch(batchBlocks, 0) // TODO: handle max batch size
			contexts, txLengths, txs, err := batch.SerializeToArgs()
			if err != nil {
				log.Error("Can not serialize batch", "error", err)
				continue
			}
			_, err = b.Inbox.AppendTxBatch(contexts, txLengths, txs)
			if err != nil {
				log.Error("Can not sequence batch", "error", err)
				continue
			}
			// Update queued assertion to latest batch
			queuedAssertion.VmHash = batch.LastBlockRoot()
			queuedAssertion.CumulativeGasUsed.Add(queuedAssertion.CumulativeGasUsed, batch.GasUsed)
			queuedAssertion.InboxSize.Add(queuedAssertion.InboxSize, batch.InboxSize())
			queuedAssertion.EndBlock = batch.LastBlockNumber()
			// If no assertion is pending, commit it
			if pendingAssertion == nil {
				commitAssertion()
			}
			batchBlocks = nil
		case blocks := <-b.BlockCh:
			// Add blocks
			batchBlocks = append(batchBlocks, blocks...)
		case ev := <-createdCh:
			// New assertion created on L1 Rollup
			if common.Address(ev.AsserterAddr) == b.Config.Coinbase {
				if ev.VmHash == pendingAssertion.VmHash {
					// If assertion is created by us, get ID and deadline
					pendingAssertion.ID = ev.AssertionID
					pendingAssertion.Deadline, err = b.AssertionMap.GetDeadline(ev.AssertionID)
					if err != nil {
						log.Error("Can not get DA deadline", "error", err)
						continue
					}
					// Send to confirmation goroutine to confirm it
					b.PendingAssertionCh <- pendingAssertion
				}
			}
		case id := <-b.ConfirmedIDCh:
			// New assertion confirmed
			if pendingAssertion.ID.Cmp(id) == 0 {
				confirmedAssertion = pendingAssertion
				if pendingAssertion.VmHash == queuedAssertion.VmHash {
					// We are done here, waiting for new batches
					pendingAssertion = nil
				} else {
					// Commit queued assertion
					commitAssertion()
				}
			} else {
				// TODO: decentralized sequencer
				// TODO: rewind blockchain, sync from L1, reset states
				log.Error("Confirmed ID is not current pending one", "get", id.String(), "expected", pendingAssertion.ID.String())
			}
		case <-b.Ctx.Done():
			return
		}
	}
}
