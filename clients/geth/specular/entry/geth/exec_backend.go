package geth

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
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
type ExecutionBackend struct{ eth GethBackend }

type GethBackend interface {
	ChainDb() ethdb.Database
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
}

var _ api.ExecutionBackend = (*ExecutionBackend)(nil)

func NewExecutionBackend(eth GethBackend, coinbase common.Address) *ExecutionBackend {
	return &ExecutionBackend{eth: eth}
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

func (b *ExecutionBackend) BuildPayload(attrs api.BuildPayloadAttributes) error {
	parent := b.eth.BlockChain().CurrentHeader()
	state, err := b.eth.BlockChain().StateAt(parent.Root)
	if err != nil {
		return err
	}
	state.StartPrefetcher("chain")
	defer state.StopPrefetcher()

	var (
		coinbase = attrs.SuggestedFeeRecipient()
		// TODO: this may cause problem if the gas limit generated on sequencer side mismatches with this one
		header = &types.Header{
			ParentHash: parent.Hash(),
			Number:     common.Big0.Add(parent.Number, common.Big1),
			GasLimit:   core.CalcGasLimit(parent.GasLimit, ethconfig.Defaults.Miner.GasCeil),
			Time:       attrs.Timestamp(),
			Coinbase:   coinbase,
			Difficulty: common.Big1, // Fake difficulty. Avoid use 0 here because it means the merge happened
		}
		gasPool = new(core.GasPool).AddGas(header.GasLimit)
		env     = &environment{
			signer:   types.MakeSigner(b.eth.BlockChain().Config(), header.Number),
			state:    state,
			tcount:   0,
			gasPool:  gasPool,
			coinbase: coinbase,
			header:   header,
		}
	)
	// Decode and force-include transactions specified by payloadAttrs, if any.
	if attrs.Txs() != nil {
		txs, err := decodeRLP(attrs.Txs())
		if err != nil {
			return fmt.Errorf("Failed to decode transactions: %w", err)
		}
		for _, tx := range txs {
			state.Prepare(tx.Hash(), env.tcount)
			if err := b.commitTransaction(env, tx); err != nil {
				return fmt.Errorf("Failed to commit force-included transaction: %w", err)
			}
		}
	}
	// Process pending transactions in the tx pool.
	if !attrs.NoTxPool() {
		pending := b.eth.TxPool().Pending(true) // should be false for us?
		pendingTxs := types.NewTransactionsByPriceAndNonce(env.signer, pending, env.header.BaseFee)
		err = b.commitTransactions(env, pendingTxs)
	}
	if env.tcount == 0 {
		log.Info("No transactions to include in payload")
		return nil
	}
	// Finalize header
	header.Root = state.IntermediateRoot(b.eth.BlockChain().Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	// Assemble block
	var (
		block = types.NewBlock(header, env.txs, nil, env.receipts, trie.NewStackTrie(nil))
		hash  = block.Hash()
		logs  []*types.Log
	)
	// Finalize receipts and logs
	for i, receipt := range env.receipts {
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
	_, err = b.eth.BlockChain().WriteBlockAndSetHead(block, env.receipts, logs, state, true)
	log.Info("Built payload", "block#", block.NumberU64(), "#txs", env.tcount)
	return err
}

// Commits transactions in the given transaction queue.
func (b *ExecutionBackend) commitTransactions(env *environment, txs *types.TransactionsByPriceAndNonce) error {
	var (
		chainCfg      = b.eth.BlockChain().Config()
		isBlockEIP155 = chainCfg.IsEIP155(env.header.Number)
	)
	for {
		if env.gasPool.Gas() < params.TxGas {
			log.Trace("Not enough gas for further transactions", "have", env.gasPool, "want", params.TxGas)
			break
		}
		var tx = txs.Peek()
		if tx == nil {
			break
		}
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !isBlockEIP155 {
			log.Trace(
				"Ignoring replay protected transaction",
				"hash", tx.Hash(), "eip155", chainCfg.EIP155Block,
			)
			txs.Pop()
			continue
		}
		env.state.Prepare(tx.Hash(), env.tcount)
		var err = b.commitTransaction(env, tx)
		switch {
		case errors.Is(err, core.ErrNonceTooLow):
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "nonce", tx.Nonce())
			txs.Shift()
		case errors.Is(err, nil):
			txs.Shift()
		default:
			// Transaction is regarded as invalid, drop all consecutive transactions from
			// the same sender because of `nonce-too-high` clause.
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Pop()
		}
	}
	return nil
}

func (b *ExecutionBackend) commitTransaction(env *environment, tx *types.Transaction) error {
	var (
		snap     = env.state.Snapshot()
		chain    = b.eth.BlockChain()
		chainCfg = chain.Config()
		vmCfg    = chain.GetVMConfig()
	)
	receipt, err := core.ApplyTransaction(chainCfg, chain, &env.coinbase, env.gasPool, env.state, env.header, tx, &env.header.GasUsed, *vmCfg)
	if err != nil {
		env.state.RevertToSnapshot(snap)
		// TODO: revert gas
		return err
	}
	env.tcount++
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)
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

type environment struct {
	signer   types.Signer
	state    *state.StateDB // apply state changes here
	tcount   int            // tx count in cycle
	gasPool  *core.GasPool  // available gas used to pack transactions
	coinbase common.Address

	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
}
