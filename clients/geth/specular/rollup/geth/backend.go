package geth

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

// TODO: cleanup; use Engine API

// Assumes exclusive control of underlying blockchain, i.e.
// mining and blockchain insertion can not happen.
// TODO: support Berlin+London fork
type ExecutionBackend struct {
	coinbase    common.Address
	chainConfig *params.ChainConfig

	chain  *core.BlockChain
	txPool *core.TxPool
	state  *state.StateDB // apply state changes here

	// pending block
	header   *types.Header
	gasPool  *core.GasPool // available gas used to pack transactions
	tcount   int
	txs      []*types.Transaction
	receipts []*types.Receipt
}

type GethBackend interface {
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
}

func NewExecutionBackend(eth GethBackend, coinbase common.Address) (*ExecutionBackend, error) {
	b := &ExecutionBackend{coinbase: coinbase, chainConfig: eth.BlockChain().Config(), chain: eth.BlockChain(), txPool: eth.TxPool()}
	err := b.startNewBlock()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *ExecutionBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.txPool.SubscribeNewTxsEvent(ch)
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
		if tx.Protected() && !b.chainConfig.IsEIP155(b.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", b.chainConfig.EIP155Block)
			continue
		}
		// Start executing the transaction
		b.state.Prepare(tx.Hash(), b.tcount)
		snap := b.state.Snapshot()
		receipt, err := core.ApplyTransaction(b.chainConfig, b.chain, &b.coinbase, b.gasPool, b.state, b.header, tx, &b.header.GasUsed, *b.chain.GetVMConfig())
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

// Sorts transactions to be committed.
func (b *ExecutionBackend) Prepare(txs []*types.Transaction) services.TransactionsByPriceAndNonce {
	sortedTxs := make(map[common.Address]types.Transactions)
	signer := types.MakeSigner(b.chainConfig, b.header.Number)
	for _, tx := range txs {
		acc, _ := types.Sender(signer, tx)
		sortedTxs[acc] = append(sortedTxs[acc], tx)
	}
	return types.NewTransactionsByPriceAndNonce(signer, sortedTxs, b.header.BaseFee)
}

// CommitBlock executes and commits a block to local blockchain *deterministically*
// TODO: dedup with CommitTransactions & commitStoredBlock
// TODO: use StateProcessor::Process() instead
func (b *ExecutionBackend) CommitPayload(payload services.ExecutionPayload) error {
	parent := b.chain.CurrentBlock()
	if parent == nil {
		return fmt.Errorf("missing parent")
	}
	num := parent.Number()
	if num.Uint64() != payload.BlockNumber()-1 {
		return fmt.Errorf("rollup services unsynced")
	}
	state, err := b.chain.StateAt(parent.Root())
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

	txs, err := DecodeRLP(payload.Txs())
	if err != nil {
		return fmt.Errorf("Failed to decode batch, err: %w", err)
	}

	for idx, tx := range txs {
		state.Prepare(tx.Hash(), idx)
		receipt, err := core.ApplyTransaction(
			b.chainConfig, b.chain, &b.coinbase, gasPool, state, header, tx, &header.GasUsed, *b.chain.GetVMConfig())
		if err != nil {
			return err
		}
		receipts = append(receipts, receipt)
	}
	// Finalize header
	header.Root = state.IntermediateRoot(b.chain.Config().IsEIP158(header.Number))
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
	_, err = b.chain.WriteBlockAndSetHead(block, receipts, logs, state, true)
	if err != nil {
		return err
	}
	return nil
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
	b.header.Root = b.state.IntermediateRoot(b.chain.Config().IsEIP158(b.header.Number))
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
	_, err := b.chain.WriteBlockAndSetHead(block, b.receipts, logs, b.state, true)
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
	parent := b.chain.CurrentBlock()
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
	state, err := b.chain.StateAt(parent.Root())
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
