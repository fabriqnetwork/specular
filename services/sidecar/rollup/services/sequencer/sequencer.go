package sequencer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/proof"
	"github.com/specularL2/specular/services/sidecar/rollup/client"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/specularL2/specular/services/sidecar/rollup/services/api"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

const timeInterval = 3 * time.Second

// TODO: delete this implementation.
// Current Sequencer assumes no Berlin+London fork on L2
type Sequencer struct {
	*services.BaseService
}

func New(eth api.ExecutionBackend, proofBackend proof.Backend, l1Client client.L1BridgeClient, cfg services.BaseConfig) (*Sequencer, error) {
	base, err := services.NewBaseService(eth, proofBackend, l1Client, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to create base service, err: %w", err)
	}
	return &Sequencer{BaseService: base}, nil
}

// Appends tx to batch if not already exists in batch or on chain
func (s *Sequencer) modifyTxnsInBatch(ctx context.Context, batchTxs []*types.Transaction, tx *types.Transaction) ([]*types.Transaction, error) {
	// Check if tx in batch
	for i := len(batchTxs) - 1; i >= 0; i-- {
		if batchTxs[i].Hash() == tx.Hash() {
			return batchTxs, nil
		}
	}
	// Check if tx exists on chain
	prevTx, _, _, _, err := s.ProofBackend.GetTransaction(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("Checking GetTransaction, err: %w", err)
	}
	if prevTx == nil {
		batchTxs = append(batchTxs, tx)
	}
	return batchTxs, nil
}

// Add sorted txs to batch and commit txs
func (s *Sequencer) addTxsToBatchAndCommit(
	ctx context.Context,
	batcher *Batcher,
	txs *types.TransactionsByPriceAndNonce,
	batchTxs []*types.Transaction,
	signer types.Signer,
) ([]*types.Transaction, error) {
	if txs != nil {
		for {
			tx := txs.Peek()
			if tx == nil {
				break
			}
			var err error
			batchTxs, err = s.modifyTxnsInBatch(ctx, batchTxs, tx)
			if err != nil {
				return nil, fmt.Errorf("Modifying batch failed, err: %w", err)
			}
			txs.Pop()
		}
	}
	if len(batchTxs) == 0 {
		return batchTxs, nil
	}
	err := batcher.CommitTransactions(batchTxs)
	if err != nil {
		return nil, fmt.Errorf("Failed to commit transactions, err: %w", err)
	}
	log.Info("Committed tx batch", "batch size", len(batchTxs))
	return batchTxs, nil
}

// This goroutine fetches txs from txpool and batches them
func (s *Sequencer) batchingLoop(ctx context.Context) {
	defer s.Wg.Done()

	// Ticker
	var ticker = time.NewTicker(timeInterval)
	defer ticker.Stop()

	// Watch transactions in TxPool
	txsCh := make(chan core.NewTxsEvent, 4096)
	txsSub := s.Eth.TxPool().SubscribeNewTxsEvent(txsCh)
	defer txsSub.Unsubscribe()

	// Process txns via batcher
	batcher, err := NewBatcher(s.Config.GetAccountAddr(), s.Eth)
	if err != nil {
		log.Crit("Failed to start batcher", "err", err)
	}

	var batchTxs []*types.Transaction

	// Loop over txns
	for {
		select {
		case <-ticker.C:
			// Get pending txs - locals and remotes, sorted by price
			var txs []*types.Transaction
			signer := types.MakeSigner(batcher.chainConfig, batcher.header.Number)

			pending := s.Eth.TxPool().Pending(true)
			localTxs, remoteTxs := make(map[common.Address]types.Transactions), pending
			for _, account := range s.Eth.TxPool().Locals() {
				if txs = remoteTxs[account]; len(txs) > 0 {
					delete(remoteTxs, account)
					localTxs[account] = txs
				}
			}
			if len(localTxs) > 0 {
				sortedTxs := types.NewTransactionsByPriceAndNonce(signer, localTxs, batcher.header.BaseFee)
				batchTxs, err = s.addTxsToBatchAndCommit(ctx, batcher, sortedTxs, batchTxs, signer)
				if err != nil {
					log.Crit("Failed to process local txs", "err", err)
				}
			}
			if len(remoteTxs) > 0 {
				sortedTxs := types.NewTransactionsByPriceAndNonce(signer, remoteTxs, batcher.header.BaseFee)
				batchTxs, err = s.addTxsToBatchAndCommit(ctx, batcher, sortedTxs, batchTxs, signer)
				if err != nil {
					log.Crit("Failed to process remote txs", "err", err)
				}
			}
			if len(batchTxs) > 0 {
				_, err := batcher.Batch()
				if err != nil {
					log.Crit("Failed to send transaction to batch", "err", err)
				}
			}
			batchTxs = nil
		case ev := <-txsCh:
			// Batch txs in case of txEvent
			log.Info("Received txsCh event", "txs", len(ev.Txs))
			txs := make(map[common.Address]types.Transactions)
			signer := types.MakeSigner(batcher.chainConfig, batcher.header.Number)
			for _, tx := range ev.Txs {
				acc, _ := types.Sender(signer, tx)
				txs[acc] = append(txs[acc], tx)
			}
			sortedTxs := types.NewTransactionsByPriceAndNonce(signer, txs, batcher.header.BaseFee)
			batchTxs, err = s.addTxsToBatchAndCommit(ctx, batcher, sortedTxs, batchTxs, signer)
			if err != nil {
				log.Crit("Failed to process txsCh event ", "err", err)
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}

func (s *Sequencer) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting sequencer...")
	err := s.BaseService.Start(ctx, eg)
	if err != nil {
		return fmt.Errorf("Failed to start sequencer: %w", err)
	}
	_, err = s.SyncL2ChainToL1Head(ctx, s.Config.GetRollupGenesisBlock())
	if err != nil {
		return fmt.Errorf("Failed to start sequencer: %w", err)
	}
	// We assume a single sequencer (us) for now, so we don't
	// need to sync transactions sequenced up.
	s.Wg.Add(1)
	go s.batchingLoop(ctx)
	log.Info("Sequencer started")
	return nil
}
