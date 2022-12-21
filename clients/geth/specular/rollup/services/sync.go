package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

const syncRange uint64 = 10000

// CommitBlocks executes and commits sequenced blocks to local blockchain
// TODO: this function shares a lot of codes with Batcher
// TODO: use StateProcessor::Process() instead
func (b *BaseService) CommitBlocks(blocks []*rollupTypes.SequenceBlock) error {
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
			GasLimit:   core.CalcGasLimit(parent.GasLimit(), ethconfig.Defaults.Miner.GasCeil), // TODO: this may cause problem
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

func batchEventToSequenceBlocks(l1 *ethclient.Client, abi *abi.ABI, ev *bindings.ISequencerInboxTxBatchAppended) ([]*rollupTypes.SequenceBlock, error) {
	// New appendTxBatch call
	tx, _, err := l1.TransactionByHash(context.Background(), ev.Raw.TxHash)
	if err != nil {
		return nil, err
	}
	// Decode input from abi
	decoded, err := abi.Methods["appendTxBatch"].Inputs.Unpack(tx.Data()[4:])
	if err != nil {
		return nil, err
	}
	// Construct batch
	batch, err := rollupTypes.TxBatchFromDecoded(decoded)
	if err != nil {
		return nil, err
	}
	// Split batch into blocks
	txNum := 0
	blocks := make([]*rollupTypes.SequenceBlock, 0, len(batch.Contexts))
	for _, ctx := range batch.Contexts {
		block := &rollupTypes.SequenceBlock{
			SequenceContext: ctx,
			Txs:             batch.Txs[txNum : txNum+int(ctx.NumTxs)],
		}
		blocks = append(blocks, block)
		txNum += int(ctx.NumTxs)
	}
	return blocks, nil
}

// SyncInbox syncs inbox events in L1 block range [start, end)
func (b *BaseService) SyncInbox(start, end uint64) error {
	log.Info("Syncing inbox", "start", start, "end", end)
	abi, err := bindings.ISequencerInboxMetaData.GetAbi()
	if err != nil {
		log.Crit("Failed to get ISequencerInbox ABI", "err", err)
	}
	currentBlock := start
	for currentBlock < end {
		currentEpochEnd := currentBlock + syncRange
		if currentEpochEnd >= end {
			currentEpochEnd = end - 1
		}
		log.Info("Syncing inbox", "currentBlock", currentBlock, "epoch end", currentEpochEnd)
		opts := &bind.FilterOpts{
			Start:   currentBlock,
			End:     &currentEpochEnd,
			Context: b.Ctx,
		}
		logIterator, err := b.Inbox.Contract.FilterTxBatchAppended(opts)
		if err != nil {
			log.Crit("Failed to get TxBatchAppended event", "err", err)
		}
		for logIterator.Next() {
			ev := logIterator.Event
			blocks, err := batchEventToSequenceBlocks(b.L1, abi, ev)
			if err != nil {
				log.Crit("Failed to convert batch event to sequence blocks", "err", err)
			}
			// Commit blocks to blockchain
			err = b.CommitBlocks(blocks)
			if err != nil {
				return err
			}
		}
		if err := logIterator.Error(); err != nil {
			log.Crit("Failed to get TxBatchAppended event", "err", err)
		}
		currentBlock = currentEpochEnd + 1
	}
	return nil
}

func (b *BaseService) SyncLoop(newBatchCh chan<- struct{}) {
	defer b.Wg.Done()

	// Get current block number
	currentBlock := b.Config.L1RollupGenesisBlock
	// Get L1 block head
	l1BlockHead, err := b.L1.BlockNumber(b.Ctx)
	if err != nil {
		log.Crit("Failed to get L1 block head", "err", err)
	}
	// Sync inbox untill the current l1 block head
	err = b.SyncInbox(currentBlock, l1BlockHead)
	if err != nil {
		log.Crit("Failed to sync inbox", "err", err)
	}
	currentBlock = l1BlockHead

	abi, err := bindings.ISequencerInboxMetaData.GetAbi()
	if err != nil {
		log.Crit("Failed to get ISequencerInbox ABI", "err", err)
	}

	// Listen to TxBatchAppendEvent
	batchEventCh := make(chan *bindings.ISequencerInboxTxBatchAppended, 4096)
	batchEventSub, err := b.Inbox.Contract.WatchTxBatchAppended(&bind.WatchOpts{Context: b.Ctx}, batchEventCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer batchEventSub.Unsubscribe()

	// Get L1 block head again and sync to it
	l1BlockHead, err = b.L1.BlockNumber(b.Ctx)
	if err != nil {
		log.Crit("Failed to get L1 block head", "err", err)
	}

	syncedCh := make(chan struct{})
	synced := false

	if l1BlockHead > currentBlock {
		go func() {
			err = b.SyncInbox(currentBlock, l1BlockHead)
			if err != nil {
				log.Crit("Failed to sync inbox", "err", err)
			}
			close(syncedCh)
		}()
	} else {
		close(syncedCh)
	}

	pendingBlocks := make([]*rollupTypes.SequenceBlock, 0)

	for {
		select {
		case <-syncedCh:
			syncedCh = nil
			synced = true
			// Commit pending blocks
			err = b.CommitBlocks(pendingBlocks)
			if err != nil {
				log.Crit("Failed to commit blocks", "err", err)
			}
		case ev := <-batchEventCh:
			blocks, err := batchEventToSequenceBlocks(b.L1, abi, ev)
			if err != nil {
				log.Crit("Failed to convert batch event to sequence blocks", "err", err)
			}
			if synced {
				// Commit blocks to blockchain
				err = b.CommitBlocks(blocks)
				if err != nil {
					log.Crit("Failed to commit blocks", "err", err)
				}
				newBatchCh <- struct{}{}
			} else {
				pendingBlocks = append(pendingBlocks, blocks...)
			}
		case <-b.Ctx.Done():
			return
		}
	}
}
