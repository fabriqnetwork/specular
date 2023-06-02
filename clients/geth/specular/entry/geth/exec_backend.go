package geth

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

// Implements api.ExecutionBackend.
// Assumes exclusive control of underlying blockchain, i.e.
// mining and blockchain insertion can not happen.
// TODO: cleanup -- use Engine API
// TODO: support Berlin+London fork
type ExecutionBackend struct {
	coinbase common.Address

	eth   GethBackend
	state *state.StateDB // apply state changes here

	// pending block
	header   *types.Header
	gasPool  *core.GasPool // available gas used to pack transactions
	tcount   int
	txs      []*types.Transaction
	receipts []*types.Receipt
}

type GethBackend interface {
	ChainDb() ethdb.Database
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
}

var _ api.ExecutionBackend = (*ExecutionBackend)(nil)

func NewExecutionBackend(eth GethBackend, coinbase common.Address) (*ExecutionBackend, error) {
	b := &ExecutionBackend{coinbase: coinbase, eth: eth}
	err := b.startNewBlock()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *ExecutionBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.eth.TxPool().SubscribeNewTxsEvent(ch)
}

// See [go-ethereum/eth/catalyst/api.go]::ConsensusAPI::forkchoiceUpdated
// Only updates fork-choice, does not build payloads.
func (b *ExecutionBackend) ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error) {
	if update.HeadBlockHash == (common.Hash{}) {
		log.Warn("Forkchoice requested update to zero hash")
		return &STATUS_INVALID, nil
	}
	// Check whether we have the block yet in our database or not.
	block := b.eth.BlockChain().GetBlockByHash(update.HeadBlockHash)
	if block == nil {
		return &STATUS_INVALID, nil
	}
	valid := ForkChoiceResponse{
		PayloadStatus: PayloadStatus{Status: VALID, LatestValidHash: &update.HeadBlockHash},
		PayloadID:     nil,
	}
	if rawdb.ReadCanonicalHash(b.eth.ChainDb(), block.NumberU64()) != update.HeadBlockHash {
		// Block is not canonical, set head.
		if latestValid, err := b.eth.BlockChain().SetCanonical(block); err != nil {
			return &ForkChoiceResponse{PayloadStatus: PayloadStatus{Status: INVALID, LatestValidHash: &latestValid}}, err
		}
	} else if b.eth.BlockChain().CurrentBlock().Hash() == update.HeadBlockHash {
		// If the specified head matches with our local head, do nothing.
		// It's a special corner case that a few slots are missing and we are requested to generate the payload in slot.
	} else {
		// If the head block is already in our canonical chain, the beacon client is
		// probably resyncing. Ignore the update.
		log.Info(
			"Ignoring beacon update to old head",
			"number", block.NumberU64(),
			"hash", update.HeadBlockHash,
			"age", common.PrettyAge(time.Unix(int64(block.Time()), 0)),
			"have", b.eth.BlockChain().CurrentBlock().Number,
		)
		return &valid, nil
	}
	// If the finalized block is not in our canonical tree, somethings wrong
	if update.FinalizedBlockHash != (common.Hash{}) {
		finalBlock := b.eth.BlockChain().GetBlockByHash(update.FinalizedBlockHash)
		if finalBlock == nil {
			log.Warn("Final block not available in database", "hash", update.FinalizedBlockHash)
			return &STATUS_INVALID, InvalidForkChoiceState.With(errors.New("final block not available in database"))
		} else if rawdb.ReadCanonicalHash(b.eth.ChainDb(), finalBlock.NumberU64()) != update.FinalizedBlockHash {
			log.Warn("Final block not in canonical chain", "number", block.NumberU64(), "hash", update.HeadBlockHash)
			return &STATUS_INVALID, InvalidForkChoiceState.With(errors.New("final block not in canonical chain"))
		}
		b.eth.BlockChain().SetFinalized(finalBlock)
	}
	// Check if the safe block hash is in our canonical tree, if not somethings wrong
	if update.SafeBlockHash != (common.Hash{}) {
		safeBlock := b.eth.BlockChain().GetBlockByHash(update.SafeBlockHash)
		if safeBlock == nil {
			log.Warn("Safe block not found", "hash", update.SafeBlockHash)
			return &STATUS_INVALID, InvalidForkChoiceState.With(errors.New("safe block not available in database"))
		}
		if rawdb.ReadCanonicalHash(b.eth.ChainDb(), safeBlock.NumberU64()) != safeBlock.Hash() {
			log.Warn("Safe block not in canonical chain")
			return &STATUS_INVALID, InvalidForkChoiceState.With(errors.New("safe block not in canonical chain"))

		}
		b.eth.BlockChain().SetSafe(safeBlock)
	}
	return &valid, nil
}

// CommitTransactions will try fill transactions into blocks, and insert
// full blocks into the blockchain
// TODO: recover from failed commitBlock, rewind blockchain
func (b *ExecutionBackend) CommitTransactions(txs []*types.Transaction) error {
	// Update timestamp if we do not have any transactions in current block
	if len(b.txs) == 0 {
		b.header.Time = uint64(time.Now().Unix())
	}

	for i := 0; i < len(txs); i++ {
		// If we don't have enough gas for any further transactions, commit and
		// start a new block
		if b.gasPool.Gas() < params.TxGas {
			err := b.commitStoredBlock()
			if err != nil {
				return err
			}
		}
		// Retrieve the next transaction
		tx := txs[i]
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !b.eth.BlockChain().Config().IsEIP155(b.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", b.eth.BlockChain().Config().EIP155Block)
			continue
		}
		// Start executing the transaction
		b.state.Prepare(tx.Hash(), b.tcount)
		snap := b.state.Snapshot()
		receipt, err := core.ApplyTransaction(b.eth.BlockChain().Config(), b.eth.BlockChain(), &b.coinbase, b.gasPool, b.state, b.header, tx, &b.header.GasUsed, *b.eth.BlockChain().GetVMConfig())
		if err != nil {
			b.state.RevertToSnapshot(snap)
		}
		switch {
		case errors.Is(err, core.ErrGasLimitReached):
			// Commit block and retry transaction in new block
			err = b.commitStoredBlock()
			if err != nil {
				return err
			}
			i--

		case errors.Is(err, nil):
			// Everything ok, collect the tx and receipt
			b.txs = append(b.txs, tx)
			b.receipts = append(b.receipts, receipt)
			b.tcount++
		}
	}

	return nil
}

// BuildPayload executes and commits a block to local blockchain *deterministically*
// TODO: dedup with CommitTransactions & commitStoredBlock
// TODO: use StateProcessor::Process() instead
func (b *ExecutionBackend) BuildPayload(payload api.ExecutionPayload) error {
	parent := b.eth.BlockChain().CurrentBlock()
	if parent == nil {
		return fmt.Errorf("missing parent")
	}
	num := parent.Number()
	if num.Uint64() != payload.BlockNumber()-1 {
		return fmt.Errorf("rollup services unsynced")
	}
	state, err := b.eth.BlockChain().StateAt(parent.Root())
	if err != nil {
		return err
	}
	state.StartPrefetcher("chain")
	defer state.StopPrefetcher()

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).SetUint64(payload.BlockNumber()),
		GasLimit:   core.CalcGasLimit(parent.GasLimit(), ethconfig.Defaults.Miner.GasCeil), // TODO: this may cause problem if the gas limit generated on sequencer side mismatch with this one
		Time:       payload.Timestamp(),
		Coinbase:   b.coinbase,
		Difficulty: common.Big1, // Fake difficulty. Avoid use 0 here because it means the merge happened
	}
	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	var receipts []*types.Receipt

	txs, err := decodeRLP(payload.Txs())
	if err != nil {
		return fmt.Errorf("Failed to decode batch, err: %w", err)
	}

	for idx, tx := range txs {
		state.Prepare(tx.Hash(), idx)
		receipt, err := core.ApplyTransaction(
			b.eth.BlockChain().Config(), b.eth.BlockChain(), &b.coinbase, gasPool, state, header, tx, &header.GasUsed, *b.eth.BlockChain().GetVMConfig())
		if err != nil {
			return err
		}
		receipts = append(receipts, receipt)
	}
	// Finalize header
	header.Root = state.IntermediateRoot(b.eth.BlockChain().Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	// Assemble block
	block := types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil))
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
	_, err = b.eth.BlockChain().WriteBlockAndSetHead(block, receipts, logs, state, true)
	if err != nil {
		return err
	}
	return nil
}

// Sorts transactions to be committed. Does not modify any state.
func (b *ExecutionBackend) Order(txs []*types.Transaction) api.TransactionQueue {
	sortedTxs := make(map[common.Address]types.Transactions)
	signer := types.MakeSigner(b.eth.BlockChain().Config(), b.header.Number)
	for _, tx := range txs {
		acc, _ := types.Sender(signer, tx)
		sortedTxs[acc] = append(sortedTxs[acc], tx)
	}
	return types.NewTransactionsByPriceAndNonce(signer, sortedTxs, b.header.BaseFee)
}

// commitStoredBlock will assemble the pending block, insert it into the blockchain
// and start a new block
func (b *ExecutionBackend) commitStoredBlock() error {
	// TODO: return if nothing to be committed

	// Stop state prefetcher (see env.discard in worker.go)
	if b.state != nil {
		b.state.StopPrefetcher()
	}
	// Finalize header
	b.header.Root = b.state.IntermediateRoot(b.eth.BlockChain().Config().IsEIP158(b.header.Number))
	b.header.UncleHash = types.CalcUncleHash(nil)
	// Assemble block
	block := types.NewBlock(b.header, b.txs, nil, b.receipts, trie.NewStackTrie(nil))
	hash := block.Hash()
	// Finalize receipts and logs
	var logs []*types.Log
	for i, receipt := range b.receipts {
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
	// Write block to chain
	_, err := b.eth.BlockChain().WriteBlockAndSetHead(block, b.receipts, logs, b.state, true)
	if err != nil {
		return err
	}
	err = b.startNewBlock()
	if err != nil {
		return err
	}
	return nil
}

// startNewBlock should be called when a block is full and inserted into the
// blockchain. It will reset the batcher except batched blocks
func (b *ExecutionBackend) startNewBlock() error {
	parent := b.eth.BlockChain().CurrentBlock()
	if parent == nil {
		return fmt.Errorf("missing parent")
	}
	num := parent.Number()
	b.header = &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   core.CalcGasLimit(parent.GasLimit(), ethconfig.Defaults.Miner.GasCeil),
		Time:       uint64(time.Now().Unix()),
		Coinbase:   b.coinbase,
		Difficulty: common.Big1, // Fake difficulty. Avoid use 0 here because it means the merge happened
	}
	state, err := b.eth.BlockChain().StateAt(parent.Root())
	if err != nil {
		return err
	}
	state.StartPrefetcher("chain")
	b.state = state
	b.gasPool = new(core.GasPool).AddGas(b.header.GasLimit)
	b.tcount = 0
	b.txs = nil
	b.receipts = nil
	return nil
}

func decodeRLP(txs [][]byte) ([]*types.Transaction, error) {
	var decodedTxs []*types.Transaction
	for _, tx := range txs {
		// TODO: use tx.DecodeRLP instead?
		var decodedTx *types.Transaction
		err := rlp.DecodeBytes(tx, decodedTx)
		if err != nil {
			return nil, err
		}
		decodedTxs = append(decodedTxs, decodedTx)
	}
	return decodedTxs, nil
}
