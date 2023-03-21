package services

import (
	"context"
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
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

type BaseService struct {
	Config *Config

	Eth          Backend
	ProofBackend proof.Backend
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
	}, nil
}

// Starts the rollup service.
func (b *BaseService) Start() (context.Context, error) {
	ctx, cancel := context.WithCancel(context.Background())
	b.Cancel = cancel
	return ctx, nil
}

func (b *BaseService) Stop() error {
	log.Info("Stopping service...")
	b.Cancel()
	b.Wg.Wait()
	log.Info("Service stopped.")
	return nil
}

func (b *BaseService) Chain() *core.BlockChain {
	return b.Eth.BlockChain()
}

// Gets the last validated assertion.
func (b *BaseService) GetLastValidatedAssertion(ctx context.Context) (*rollupTypes.Assertion, error) {
	opts := bind.FilterOpts{Start: b.Config.L1RollupGenesisBlock, Context: ctx}
	assertionID, err := b.L1Client.GetLastValidatedAssertionID(&opts)

	var assertionCreatedEvent *bindings.IRollupAssertionCreated
	var lastValidatedAssertion bindings.IRollupAssertion
	if err != nil {
		// If no assertion was validated (or other errors encountered), try to use the genesis assertion.
		log.Warn("No validated assertions found, using genesis assertion", "err", err)
		assertionCreatedEvent, err = b.L1Client.GetGenesisAssertionCreated(&opts)
		if err != nil {
			return nil, fmt.Errorf("Failed to get `AssertionCreated` event for last validated assertion, err: %w", err)
		}
		// Check that the genesis assertion is correct.
		vmHash := common.BytesToHash(assertionCreatedEvent.VmHash[:])
		genesisRoot := b.Eth.BlockChain().GetBlockByNumber(0).Root()
		if vmHash != genesisRoot {
			return nil, fmt.Errorf("Mismatching genesis %s vs %s", vmHash, genesisRoot.String())
		}
		log.Info("Genesis assertion found", "assertionID", assertionCreatedEvent.AssertionID)
		// Get assertion.
		lastValidatedAssertion, err = b.L1Client.GetAssertion(assertionCreatedEvent.AssertionID)
	} else {
		// If an assertion was validated, use it.
		log.Info("Last validated assertion ID found", "assertionID", assertionID)
		lastValidatedAssertion, err = b.L1Client.GetAssertion(assertionID)
		if err != nil {
			return nil, fmt.Errorf("Failed to get last validated assertion, err: %w", err)
		}
		opts = bind.FilterOpts{Start: lastValidatedAssertion.ProposalTime.Uint64(), Context: ctx}
		assertionCreatedEvent, err = b.L1Client.FilterAssertionCreated(&opts, assertionID)
		if err != nil {
			return nil, fmt.Errorf("Failed to get `AssertionCreated` event for last validated assertion, err: %w", err)
		}
	}
	// Initialize assertion.
	assertion := NewAssertionFrom(&lastValidatedAssertion, assertionCreatedEvent)
	// Set its boundaries using parent. TODO: move this out. Use local caching.
	opts = bind.FilterOpts{Start: b.Config.L1RollupGenesisBlock, Context: ctx}
	parentAssertionCreatedEvent, err := b.L1Client.FilterAssertionCreated(&opts, lastValidatedAssertion.Parent)
	if err != nil {
		return nil, fmt.Errorf("Failed to get `AssertionCreated` event for parent assertion, err: %w", err)
	}
	err = b.setL2BlockBoundaries(assertion, parentAssertionCreatedEvent)
	if err != nil {
		return nil, fmt.Errorf("Failed to set L2 block boundaries for last validated assertion, err: %w", err)
	}
	return assertion, nil
}

func (b *BaseService) Stake(ctx context.Context) error {
	staker, err := b.L1Client.GetStaker()
	if err != nil {
		return fmt.Errorf("Failed to get staker, to stake, err: %w", err)
	}
	if !staker.IsStaked {
		err = b.L1Client.Stake(big.NewInt(int64(b.Config.RollupStakeAmount)))
	}
	if err != nil {
		return fmt.Errorf("Failed to stake, err: %w", err)
	}
	return nil
}

// Sync to current L1 block head and commit blocks.
// `start` is the block number to start syncing from.
// Returns the last synced block number (inclusive).
func (b *BaseService) SyncL2ChainToL1Head(ctx context.Context, start uint64) (uint64, error) {
	l1BlockHead, err := b.L1Client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	opts := bind.FilterOpts{Start: start, End: &l1BlockHead, Context: ctx}
	eventsIter, err := b.L1Client.FilterTxBatchAppendedEvents(&opts)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	err = b.processTxBatchAppendedEvents(ctx, eventsIter)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	log.Info(
		"Synced L1->L2",
		"l1 start", start,
		"l1 end", l1BlockHead,
		"l2 size", b.Eth.BlockChain().CurrentBlock().Number(),
	)
	return l1BlockHead, nil
}

func (b *BaseService) SyncLoop(ctx context.Context, start uint64, newBatchCh chan<- struct{}) {
	defer b.Wg.Done()
	// Sync to current L1 block head. Can't use WatchOpts.Start,
	// due to github.com/ethereum/go-ethereum/issues/15063
	endBlock, err := b.SyncL2ChainToL1Head(ctx, start)
	if err != nil {
		log.Crit("Failed initial inbox sync", "err", err)
	}
	// Start watching for new TxBatchAppended events.
	batchEventCh := make(chan *bindings.ISequencerInboxTxBatchAppended)
	opts := bind.WatchOpts{Context: ctx}
	batchEventSub, err := b.L1Client.WatchTxBatchAppended(&opts, batchEventCh)
	if err != nil {
		log.Crit("Failed to watch rollup event", "err", err)
	}
	defer batchEventSub.Unsubscribe()
	// Sync again to ensure we don't miss events.
	endBlock, err = b.SyncL2ChainToL1Head(ctx, endBlock+1)
	if err != nil {
		log.Crit("Failed initial inbox sync", "err", err)
	}
	// Process TxBatchAppended events.
	for {
		select {
		case ev := <-batchEventCh:
			// Avoid processing duplicate events.
			if ev.Raw.BlockNumber <= endBlock {
				log.Warn("Ignoring duplicate event", "l1Block", ev.Raw.BlockNumber)
				continue
			}
			log.Info("Processing `TxBatchAppended` event", "l1Block", ev.Raw.BlockNumber)
			err := b.processTxBatchAppendedEvent(ctx, ev)
			if err != nil {
				log.Crit("Failed to process event", "err", err)
			}
			if newBatchCh != nil {
				newBatchCh <- struct{}{}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *BaseService) processTxBatchAppendedEvents(
	ctx context.Context,
	eventsIter *bindings.ISequencerInboxTxBatchAppendedIterator,
) error {
	for eventsIter.Next() {
		log.Info("processTxBatchAppendedEvents l1blocknumber is", "l1blocknumber", eventsIter.Event.Raw.BlockNumber)
		err := b.processTxBatchAppendedEvent(ctx, eventsIter.Event)
		if err != nil {
			return fmt.Errorf("Failed to process event, err: %w", err)
		}
	}
	if err := eventsIter.Error(); err != nil {
		return fmt.Errorf("Failed to iterate through events, err: %w", err)
	}
	return nil
}

// Reads tx data associated with batch event and commits as blocks on L2.
func (b *BaseService) processTxBatchAppendedEvent(
	ctx context.Context,
	ev *bindings.ISequencerInboxTxBatchAppended,
) error {
	tx, _, err := b.L1Client.TransactionByHash(ctx, ev.Raw.TxHash)
	if err != nil {
		return fmt.Errorf("Failed to get transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Decode input to appendTxBatch transaction.
	decoded, err := b.L1Client.DecodeAppendTxBatchInput(tx)
	if err != nil {
		return fmt.Errorf("Failed to decode transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Construct batch. TODO: decode into blocks directly.
	batch, err := rollupTypes.TxBatchFromDecoded(decoded)
	if err != nil {
		return fmt.Errorf("Failed to split AppendTxBatch input into batches, err: %w", err)
	}
	log.Info("Decoded batch", "#txs", len(batch.Txs))
	blocks := batch.SplitToBlocks()
	log.Info("Batch split into blocks", "#blocks", len(blocks))
	//Compare L2 Block Number of the Event to Current L2 Block Number
	if len(blocks) > 0 && blocks[0].BlockNumber > b.Eth.BlockChain().CurrentBlock().Number().Uint64() {
		// If it's greater means we can commit on L2
		b.commitBlocks(blocks)
	}
	return nil
}

// TODO: clean up.
func (b *BaseService) setL2BlockBoundaries(
	assertion *rollupTypes.Assertion,
	parentAssertionCreatedEvent *bindings.IRollupAssertionCreated,
) error {
	numBlocks := b.Eth.BlockChain().CurrentBlock().Number().Uint64()
	if numBlocks == 0 {
		log.Info("Zero-initializing assertion block boundaries.")
		assertion.StartBlock = 0
		assertion.EndBlock = 0
		return nil
	}
	startFound := false
	// Note: by convention defined in Rollup.sol, the parent VmHash is the
	// same as the child only when the assertion is the genesis assertion.
	// This is a hack to avoid mis-setting `assertion.StartBlock`.
	if assertion.ID == parentAssertionCreatedEvent.AssertionID {
		parentAssertionCreatedEvent.VmHash = common.Hash{}
		startFound = true
	}
	log.Info("Searching for start and end blocks for assertion.", "id", assertion.ID)
	// Find start and end blocks using L2 chain (assumes it's synced at least up to the assertion).
	for i := uint64(0); i <= numBlocks; i++ {
		root := b.Eth.BlockChain().GetBlockByNumber(i).Root()
		if root == parentAssertionCreatedEvent.VmHash {
			log.Info("Found start block", "l2 block#", i+1)
			assertion.StartBlock = i + 1
			startFound = true
		} else if root == assertion.VmHash {
			log.Info("Found end block", "l2 block#", i)
			assertion.EndBlock = i
			if !startFound {
				return fmt.Errorf("Found end block before start block for assertion with hash %d", assertion.VmHash)
			}
			return nil
		}
	}
	return fmt.Errorf("Could not find start or end block for assertion with hash %s", assertion.VmHash)
}

// commitBlocks executes and commits sequenced blocks to local blockchain
// TODO: this function shares a lot of codes with Batcher
// TODO: use StateProcessor::Process() instead
func (b *BaseService) commitBlocks(blocks []*rollupTypes.SequenceBlock) error {
	if len(blocks) == 0 {
		return nil
	}
	chainConfig := b.Chain().Config()
	parent := b.Chain().CurrentBlock()
	if parent == nil {
		return fmt.Errorf("missing parent")
	}
	num := parent.Number()
	if num.Uint64() != blocks[0].BlockNumber-1 {
		return fmt.Errorf("rollup services unsynced")
	}
	state, err := b.Chain().StateAt(parent.Root())
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
			receipt, err := core.ApplyTransaction(
				chainConfig, b.Chain(), &b.Config.SequencerAddr, gasPool, state, header, tx, &header.GasUsed, *b.Chain().GetVMConfig())
			if err != nil {
				return err
			}
			receipts = append(receipts, receipt)
		}
		// Finalize header
		header.Root = state.IntermediateRoot(b.Chain().Config().IsEIP158(header.Number))
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
		_, err := b.Chain().WriteBlockAndSetHead(block, receipts, logs, state, true)
		if err != nil {
			return err
		}
	}
	return nil
}
