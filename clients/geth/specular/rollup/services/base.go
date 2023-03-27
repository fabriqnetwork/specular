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
	L1Syncer     *client.L1Syncer

	Cancel context.CancelFunc
	Wg     sync.WaitGroup
}

func NewBaseService(eth Backend, proofBackend proof.Backend, l1Client client.L1BridgeClient, cfg *Config) (*BaseService, error) {
	return &BaseService{
		Config:       cfg,
		Eth:          eth,
		ProofBackend: proofBackend,
		L1Client:     l1Client,
		L1Syncer:     nil,
	}, nil
}

// Starts the rollup service.
func (b *BaseService) Start() (context.Context, error) {
	ctx, cancel := context.WithCancel(context.Background())
	b.Cancel = cancel
	b.L1Syncer = client.NewL1Syncer(ctx, b.L1Client)
	b.L1Syncer.Start(ctx)
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
		assertionCreatedIter, err := b.L1Client.FilterAssertionCreated(&opts)
		if err != nil {
			return nil, fmt.Errorf("Failed to get `AssertionCreated` event for last validated assertion, err: %w", err)
		}
		assertionCreatedEvent, err = filterAssertionCreatedWithID(assertionCreatedIter, assertionID)
	}
	// Initialize assertion.
	assertion := NewAssertionFrom(&lastValidatedAssertion, assertionCreatedEvent)
	// Set its boundaries using parent. TODO: move this out. Use local caching.
	opts = bind.FilterOpts{Start: b.Config.L1RollupGenesisBlock, Context: ctx}
	parentAssertionCreatedIter, err := b.L1Client.FilterAssertionCreated(&opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to get `AssertionCreated` event for parent assertion, err: %w", err)
	}
	parentAssertionCreatedEvent, err := filterAssertionCreatedWithID(parentAssertionCreatedIter, lastValidatedAssertion.Parent)
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
	// Start watching for new TxBatchAppended events.
	subCtx, cancel := context.WithCancel(ctx)
	batchEventCh := client.SubscribeHeaderMapped[*bindings.ISequencerInboxTxBatchAppended](
		subCtx, b.L1Syncer.LatestHeaderBroker, b.L1Client.FilterTxBatchAppendedEvents, start,
	)
	defer cancel()
	// Process TxBatchAppended events.
	for {
		select {
		case ev := <-batchEventCh:
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

func filterAssertionCreatedWithID(iter *bindings.IRollupAssertionCreatedIterator, assertionID *big.Int) (*bindings.IRollupAssertionCreated, error) {
	var assertionCreated *bindings.IRollupAssertionCreated
	for iter.Next() {
		// Assumes invariant: only one `AssertionCreated` event per assertion ID.
		if iter.Event.AssertionID.Cmp(assertionID) == 0 {
			assertionCreated = iter.Event
			break
		}
	}
	if iter.Error() != nil {
		return nil, fmt.Errorf("Failed to iterate through `AssertionCreated` events, err: %w", iter.Error())
	}
	if assertionCreated == nil {
		return nil, fmt.Errorf("No `AssertionCreated` event found for %v.", assertionID)
	}
	return assertionCreated, nil
}

func (b *BaseService) processTxBatchAppendedEvents(
	ctx context.Context,
	eventsIter *bindings.ISequencerInboxTxBatchAppendedIterator,
) error {
	for eventsIter.Next() {
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
	// Compare the L2 Block Number of the Event to the Current L2 Block Number and commit blocks ahead of the current L2 chain
	// If it causes any performance issue, then it can optimized at the batch level by sending 2 batches at once. If 1st element of 1st batch is behind and the 2nd batch ahead then only we look at the block level, else look only at the batch level
	for i, block := range blocks {
		if block.BlockNumber > b.Eth.BlockChain().CurrentBlock().Number().Uint64() {
			blocks = blocks[i:]
			b.commitBlocks(blocks)
			break
		}
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
		// TODO: remove assumption of VM hash being the block root.
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
