package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type BaseService struct {
	Config *Config

	Eth          Backend
	ProofBackend proof.Backend
	Chain        *core.BlockChain
	L1Client     client.L1BridgeClient

	Cancel context.CancelFunc
	Wg     sync.WaitGroup
}

func NewBaseService(eth Backend, proofBackend proof.Backend, l1Client client.L1BridgeClient, cfg *Config) (*BaseService, error) {
	if eth == nil {
		return nil, fmt.Errorf("can not use light client with rollup")
	}
	return &BaseService{
		Config:       cfg,
		Eth:          eth,
		ProofBackend: proofBackend,
		L1Client:     l1Client,
		Chain:        eth.BlockChain(),
	}, nil
}

// Starts the rollup service.
func (b *BaseService) Start() (context.Context, error) {
	ctx, cancel := context.WithCancel(context.Background())
	b.Cancel = cancel
	// Check if we are at genesis
	if b.Eth.BlockChain().CurrentBlock().NumberU64() != 0 {
		return nil, fmt.Errorf("Rollup service can only start from clean history")
	}
	return ctx, nil
}

func (b *BaseService) Stop() error {
	log.Info("Stopping service...")
	b.Cancel()
	b.Wg.Wait()
	log.Info("Service stopped.")
	return nil
}

func (b *BaseService) Stake(ctx context.Context) error {
	isStaked, err := b.L1Client.IsStaked(ctx)
	if err != nil {
		return fmt.Errorf("Failed to stake, err: %w", err)
	}
	if !isStaked {
		err = b.L1Client.Stake(big.NewInt(int64(b.Config.RollupStakeAmount)))
	}
	if err != nil {
		return fmt.Errorf("Failed to stake, err: %w", err)
	}
	return nil
}

// Sync to current L1 block head and commit blocks.
func (b *BaseService) SyncL2ChainToL1Head(ctx context.Context, start uint64) (uint64, error) {
	l1BlockHead, err := b.L1Client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	// start, l1BlockHead
	opts := bind.FilterOpts{Start: start, End: &l1BlockHead, Context: ctx}
	eventsIter, err := b.L1Client.FilterTxBatchAppendedEvents(&opts)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	blocks, err := b.processEvents(ctx, eventsIter)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	if err := b.commitBlocks(blocks); err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	return l1BlockHead, nil
}

func (b *BaseService) SyncLoop(ctx context.Context, newBatchCh chan<- struct{}) {
	defer b.Wg.Done()
	batchEventCh := make(chan *bindings.ISequencerInboxTxBatchAppended)
	opts := bind.WatchOpts{Start: &b.Config.L1RollupGenesisBlock, Context: ctx}
	batchEventSub, err := b.L1Client.WatchTxBatchAppended(&opts, batchEventCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer batchEventSub.Unsubscribe()
	// Process TxBatchAppended events.
	for {
		select {
		case ev := <-batchEventCh:
			blocks, err := b.processEvent(ctx, ev)
			if err != nil {
				log.Crit("Failed to process event", "err", err)
			}
			// Commit blocks to blockchain
			if err = b.commitBlocks(blocks); err != nil {
				log.Crit("Failed to commit blocks", "err", err)
			}
			if newBatchCh != nil {
				newBatchCh <- struct{}{}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *BaseService) processEvents(
	ctx context.Context,
	eventsIter *bindings.ISequencerInboxTxBatchAppendedIterator,
) ([]*rollupTypes.SequenceBlock, error) {
	blocks := make([]*rollupTypes.SequenceBlock, 0)
	for eventsIter.Next() {
		eventBlocks, err := b.processEvent(ctx, eventsIter.Event)
		if err != nil {
			return nil, fmt.Errorf("Failed to process event, err: %w", err)
		}
		blocks = append(blocks, eventBlocks...)
	}
	if err := eventsIter.Error(); err != nil {
		return nil, fmt.Errorf("Failed to iterate through events, err: %w", err)
	}
	return blocks, nil
}

// Reads tx data associated with batch event and returns corresponding list of L2 blocks.
func (b *BaseService) processEvent(
	ctx context.Context,
	ev *bindings.ISequencerInboxTxBatchAppended,
) ([]*rollupTypes.SequenceBlock, error) {
	tx, _, err := b.L1Client.TransactionByHash(ctx, ev.Raw.TxHash)
	if err != nil {
		return nil, fmt.Errorf("Failed to get transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Decode input to appendTxBatch transaction.
	decoded, err := b.L1Client.DecodeAppendTxBatchInput(tx)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Construct batch. TODO: decode into blocks directly.
	batch, err := rollupTypes.TxBatchFromDecoded(decoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to split AppendTxBatch input into batches, err: %w", err)
	}
	return batch.SplitToBlocks(), nil
}

// commitBlocks executes and commits sequenced blocks to local blockchain
// TODO: this function shares a lot of codes with Batcher
// TODO: use StateProcessor::Process() instead
func (b *BaseService) commitBlocks(blocks []*rollupTypes.SequenceBlock) error {
	if len(blocks) == 0 {
		return nil
	}
	chainConfig := b.Chain.Config()
	parent := b.Chain.CurrentBlock()
	if parent == nil {
		return fmt.Errorf("missing parent")
	}
	num := parent.Number()
	if num.Uint64() != blocks[0].BlockNumber-1 {
		return fmt.Errorf("rollup services unsynced")
	}
	state, err := b.Chain.StateAt(parent.Root())
	if err != nil {
		return err
	}
	state.StartPrefetcher("rollup")
	defer state.StopPrefetcher()

	for _, sblock := range blocks {
		header := &types.Header{
			ParentHash: parent.Hash(),
			Number:     new(big.Int).SetUint64(sblock.BlockNumber),
			GasLimit:   core.CalcGasLimit(parent.GasLimit(), ethconfig.Defaults.Miner.GasCeil), // TODO: this may cause problem if the gas limit generated on sequencer side mismatch with this one
			Time:       sblock.Timestamp,
			Coinbase:   b.Config.SequencerAddr,
			Difficulty: common.Big1, // Fake difficulty. Avoid use 0 here because it means the merge happened
		}
		gasPool := new(core.GasPool).AddGas(header.GasLimit)
		var receipts []*types.Receipt
		for idx, tx := range sblock.Txs {
			state.Prepare(tx.Hash(), idx)
			receipt, err := core.ApplyTransaction(chainConfig, b.Chain, &b.Config.SequencerAddr, gasPool, state, header, tx, &header.GasUsed, *b.Chain.GetVMConfig())
			if err != nil {
				return err
			}
			receipts = append(receipts, receipt)
		}
		// Finalize header
		header.Root = state.IntermediateRoot(b.Chain.Config().IsEIP158(header.Number))
		header.UncleHash = types.CalcUncleHash(nil)
		// Assemble block
		block := types.NewBlock(header, sblock.Txs, nil, receipts, trie.NewStackTrie(nil))
		hash := block.Hash()
		// Finalize receipts and logs
		var logs []*types.Log
		for i, receipt := range receipts {
			// Add block location fields
			receipt.BlockHash = hash
			receipt.BlockNumber = block.Number()
			receipt.TransactionIndex = uint(i)

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for _, log := range receipt.Logs {
				log.BlockHash = hash
			}
			logs = append(logs, receipt.Logs...)
		}
		_, err := b.Chain.WriteBlockAndSetHead(block, receipts, logs, state, true)
		if err != nil {
			return err
		}
	}
	return nil
}
