package sequencer

import (
	"context"
	errors "errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

const timeInterval = 3 * time.Second

// TODO: get rid of batcher; use engine API
// TODO: modularize challenge

type challengeCtx struct {
	challengeAddr common.Address
	assertion     *rollupTypes.Assertion
}

// Current Sequencer assumes no Berlin+London fork on L2
type Sequencer struct {
	*services.BaseService

	blockCh              chan types.Blocks
	pendingAssertionCh   chan *rollupTypes.Assertion
	confirmedIDCh        chan *big.Int
	challengeCh          chan *challengeCtx
	challengeResoutionCh chan struct{}
}

func New(eth services.Backend, proofBackend proof.Backend, l1Client client.L1BridgeClient, cfg *services.Config) (*Sequencer, error) {
	base, err := services.NewBaseService(eth, proofBackend, l1Client, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to create base service, err: %w", err)
	}
	return &Sequencer{
		BaseService:          base,
		blockCh:              make(chan types.Blocks, 4096),
		pendingAssertionCh:   make(chan *rollupTypes.Assertion, 4096),
		confirmedIDCh:        make(chan *big.Int, 4096),
		challengeCh:          make(chan *challengeCtx),
		challengeResoutionCh: make(chan struct{}),
	}, nil
}

// Get tx index in batch
func getTxIndexInBatch(slice []*types.Transaction, elem *types.Transaction) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i].Hash() == elem.Hash() {
			return i
		}
	}
	return -1
}

// Appends tx to batch if not already exists in batch or on chain
func (s *Sequencer) modifyTxnsInBatch(ctx context.Context, batchTxs []*types.Transaction, tx *types.Transaction) ([]*types.Transaction, error) {
	// Check if tx in batch
	txIndex := getTxIndexInBatch(batchTxs, tx)
	if txIndex < 0 {
		// Check if tx exists on chain
		prevTx, _, _, _, err := s.ProofBackend.GetTransaction(ctx, tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("Checking GetTransaction, err: %w", err)
		}
		if prevTx == nil {
			batchTxs = append(batchTxs, tx)
		}
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
	defer close(s.blockCh)

	// Ticker
	var ticker = time.NewTicker(timeInterval)
	defer ticker.Stop()

	// Watch transactions in TxPool
	txsCh := make(chan core.NewTxsEvent, 4096)
	txsSub := s.Eth.TxPool().SubscribeNewTxsEvent(txsCh)
	defer txsSub.Unsubscribe()

	// Process txns via batcher
	batcher, err := NewBatcher(s.Config.Coinbase, s.Eth)
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
				blocks, err := batcher.Batch()
				if err != nil {
					log.Crit("Failed to send transaction to batch", "err", err)
				}
				s.blockCh <- blocks
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

func (s *Sequencer) sequencingLoop(ctx context.Context) {
	defer s.Wg.Done()

	// Ticker
	var ticker = time.NewTicker(timeInterval)
	defer ticker.Stop()

	// Watch AssertionCreated event
	createdCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionCreated](
		ctx,
		s.L1Syncer.LatestHeaderBroker,
		s.L1Client.FilterAssertionCreated,
		s.L1Syncer.Latest.Number.Uint64(),
	)

	// Last validated assertion, initalize it to genesis
	// TODO: change name to lastValidatedAssertion since "confirmed" may imply L1-confirmed.
	confirmedAssertion, err := s.GetLastValidatedAssertion(ctx)
	if err != nil {
		log.Crit("Failed to get last validated assertion", "err", err)
	}
	// Assertion created and pending for confirmation
	var pendingAssertion *rollupTypes.Assertion
	// Assertion to be created on L1 Rollup
	queuedAssertion := confirmedAssertion.Copy()

	// Create assertion on L1 Rollup
	commitAssertion := func() {
		pendingAssertion = queuedAssertion.Copy()
		queuedAssertion.StartBlock = queuedAssertion.EndBlock + 1
		_, err = s.L1Client.CreateAssertion(pendingAssertion.VmHash, pendingAssertion.InboxSize)
		if errors.Is(err, core.ErrInsufficientFunds) {
			log.Crit("Insufficient Funds to send Tx", "error", err)
		}
		if err != nil {
			log.Error("Can not create DA", "error", err)
		}
		log.Info(
			"Created assertion",
			"id", pendingAssertion.ID,
			"vmHash", pendingAssertion.VmHash,
			"start block", pendingAssertion.StartBlock,
			"end block", pendingAssertion.EndBlock,
		)
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
			contexts, txLengths, firstL2BlockNumber, txs, err := batch.SerializeToArgs()
			if err != nil {
				log.Error("Can not serialize batch", "error", err)
				continue
			}
			log.Info(
				"Serialized new Tx Batch",
				"#txs", len(batch.Txs),
				"#numBlocks", len(contexts) / 2,
				"#firsBlockNumber", firstL2BlockNumber,
			)

			_, err = s.L1Client.AppendTxBatch(contexts, txLengths, firstL2BlockNumber, txs)
			if errors.Is(err, core.ErrInsufficientFunds) {
				log.Crit("Insufficient Funds to send Tx", "error", err)
			}
			if err != nil {
				log.Error("Can not sequence batch", "error", err)
				continue
			}
			log.Info("Sequenced batch", "batch size", len(batch.Txs))
			// Update queued assertion to latest batch
			// queuedAssertion.ID.Add(queuedAssertion.ID, big.NewInt(1))
			queuedAssertion.VmHash = batch.LastBlockRoot()
			queuedAssertion.InboxSize.Add(queuedAssertion.InboxSize, batch.InboxSize())
			queuedAssertion.EndBlock = batch.LastBlockNumber()
			// If no assertion is pending, commit it
			if pendingAssertion == nil {
				commitAssertion()
			}
			batchBlocks = nil
		case blocks := <-s.blockCh:
			// Add blocks
			batchBlocks = append(batchBlocks, blocks...)
		case ev := <-createdCh:
			// New assertion created on L1 Rollup
			log.Info("Received `AssertionCreated` event.", "assertion id", ev.AssertionID)
			if common.Address(ev.AsserterAddr) == s.Config.Coinbase {
				if ev.VmHash == pendingAssertion.VmHash {
					// If assertion is created by us, get ID and deadline
					pendingAssertion.ID = ev.AssertionID
					assertionFromRollup, err := s.L1Client.GetAssertion(ev.AssertionID)
					if err != nil {
						log.Error("Could not get DA", "error", err)
						continue
					}
					pendingAssertion.Deadline = assertionFromRollup.Deadline
					// Send to confirmation goroutine to confirm it
					s.pendingAssertionCh <- pendingAssertion
				}
			}
		case id := <-s.confirmedIDCh:
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
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}

// This goroutine tries to confirm created assertions
func (s *Sequencer) confirmationLoop(ctx context.Context) {
	defer s.Wg.Done()

	headCh := s.L1Syncer.LatestHeaderBroker.Subscribe()
	confirmedCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionConfirmed](
		ctx, s.L1Syncer.LatestHeaderBroker, s.L1Client.FilterAssertionConfirmed, s.L1Syncer.Latest.Number.Uint64(),
	)
	challengedCh := client.SubscribeHeaderMapped[*bindings.IRollupAssertionChallenged](
		ctx, s.L1Syncer.LatestHeaderBroker, s.L1Client.FilterAssertionChallenged, s.L1Syncer.Latest.Number.Uint64(),
	)

	// Current pending assertion from sequencing goroutine
	// TODO: watch multiple pending assertions
	var pendingAssertion *rollupTypes.Assertion
	pendingConfirmationSent := true
	pendingConfirmed := true

	for {
		select {
		case header := <-headCh:
			// New block mined on L1
			log.Info("Received new header", "number", header.Number.Uint64())
			if !pendingConfirmationSent && !pendingConfirmed {
				if header.Number.Uint64() >= pendingAssertion.Deadline.Uint64() {
					log.Info("We can now confirm", "pending assertion", pendingAssertion.Deadline.Uint64())
					// Confirmation period has past, confirm it
					_, err := s.L1Client.ConfirmFirstUnresolvedAssertion()
					if errors.Is(err, core.ErrInsufficientFunds) {
						log.Crit("Insufficient Funds to send Tx", "error", err)
					}
					if err != nil {
						// log.Error("Failed to confirm DA", "error", err)
						log.Crit("Failed to confirm DA", "err", err)
						// TODO: wait some time before retry
					}
					pendingConfirmationSent = true
				}
			}
		case ev := <-confirmedCh:
			log.Info("Received `AssertionConfirmed` event ", "assertion id", ev.AssertionID)
			// New confirmed assertion
			if ev.AssertionID.Cmp(pendingAssertion.ID) == 0 {
				// Notify sequencing goroutine
				s.confirmedIDCh <- pendingAssertion.ID
				pendingConfirmed = true
			}
		case newPendingAssertion := <-s.pendingAssertionCh:
			log.Info("Received pending assertion")
			// New assertion created by sequencing goroutine
			if !pendingConfirmed {
				// TODO: support multiple pending assertion
				log.Error("Got another DA request before current is confirmed")
				continue
			}
			pendingAssertion = newPendingAssertion.Copy()
			pendingConfirmationSent = false
			pendingConfirmed = false
		case ev := <-challengedCh:
			// New challenge raised
			log.Info("Received `AssertionChallenged` event ", "assertion id", ev.AssertionID)
			if ev.AssertionID.Cmp(pendingAssertion.ID) == 0 {
				s.challengeCh <- &challengeCtx{
					ev.ChallengeAddr,
					pendingAssertion,
				}
				wait(ctx, s.challengeResoutionCh, "challenge resolution")
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}

func wait(ctx context.Context, ch <-chan struct{}, taskName string) {
	log.Info("Waiting for %s...", taskName)
	select {
	case <-ch:
		log.Info("%s notification received", taskName)
		return
	case <-ctx.Done():
		log.Info("Aborting wait for %s", taskName)
		return
	}
}

func (s *Sequencer) challengeLoop(ctx context.Context) {
	defer s.Wg.Done()
	// Watch L1 blockchain for challenge timeout
	headCh := make(chan *types.Header, 4096)
	headSub, err := s.L1Client.ResubscribeErrNewHead(ctx, headCh)
	if err != nil {
		log.Crit("Failed to watch l1 chain head", "err", err)
	}
	defer headSub.Unsubscribe()
	for {
		select {
		case chalCtx := <-s.challengeCh:
			err := s.handleChallenge(ctx, chalCtx, headCh)
			if err != nil {
				log.Crit("Failed to handle challenge", "err", err)
			}
		case <-headCh:
			continue // consume channel values
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}

func (s *Sequencer) handleChallenge(
	ctx context.Context,
	chalCtx *challengeCtx,
	headCh chan *types.Header,
) error {
	err := s.L1Client.InitNewChallengeSession(ctx, chalCtx.challengeAddr)
	if err != nil {
		return fmt.Errorf("Failed to access ongoing challenge (address=%s), err: %w", chalCtx.challengeAddr, err)
	}
	states, err := proof.GenerateStates(
		s.ProofBackend,
		ctx,
		chalCtx.assertion.StartBlock,
		chalCtx.assertion.EndBlock+1,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to generate states, err: %w", err)
	}
	_, err = s.L1Client.InitializeChallengeLength(new(big.Int).SetUint64(uint64(len(states)) - 1))
	if err != nil {
		return fmt.Errorf("Failed to initialize challenge, err: %w", err)
	}

	subCtx, subCancel := context.WithCancel(ctx)
	defer subCancel()
	bisectedCh := client.SubscribeHeaderMapped[*bindings.ISymChallengeBisected](
		subCtx, s.L1Syncer.LatestHeaderBroker, s.L1Client.FilterBisected, s.L1Syncer.Latest.Number.Uint64(),
	)
	challengeCompletedCh := client.SubscribeHeaderMapped[*bindings.ISymChallengeCompleted](
		subCtx, s.L1Syncer.LatestHeaderBroker, s.L1Client.FilterChallengeCompleted, s.L1Syncer.Latest.Number.Uint64(),
	)
	log.Info("to generate state from", "start", chalCtx.assertion.StartBlock, "to", chalCtx.assertion.EndBlock)
	log.Info("backend", "start", chalCtx.assertion.StartBlock, "to", chalCtx.assertion.EndBlock)
	var opponentTimeoutBlock uint64
	for {
		select {
		case ev := <-bisectedCh:
			// case get bisection, if is our turn
			//   if in single step, submit proof
			//   if multiple step, track current segment, update
			responder, err := s.L1Client.CurrentChallengeResponder()
			if err != nil {
				// TODO: error handling
				log.Error("Cannot get current responder", "error", err)
				continue
			}
			if responder == common.Address(s.Config.Coinbase) {
				// If it's our turn
				err := services.RespondBisection(ctx, s.ProofBackend, s.L1Client, ev, states, common.Hash{}, false)
				if err != nil {
					// TODO: error handling
					log.Error("Cannot respond to bisection", "error", err)
					continue
				}
			} else {
				opponentTimeLeft, err := s.L1Client.CurrentChallengeResponderTimeLeft()
				if err != nil {
					// TODO: error handling
					log.Error("Cannot get current responder left time", "error", err)
					continue
				}
				log.Info("Opponent time left", "time", opponentTimeLeft)
				opponentTimeoutBlock = ev.Raw.BlockNumber + opponentTimeLeft.Uint64()
			}
		case header := <-headCh:
			if opponentTimeoutBlock == 0 {
				continue
			}
			// TODO: can we use >= here?
			if header.Number.Uint64() > opponentTimeoutBlock {
				_, err := s.L1Client.TimeoutChallenge()
				if err != nil {
					log.Error("Cannot timeout opponent", "error", err)
					continue
					// TODO: wait some time before retry
					// TODO: fix race condition
				}
			}
		case ev := <-challengeCompletedCh:
			// TODO: handle if we are not winner --> state corrupted
			log.Info("Challenge completed", "winner", ev.Winner)
			states = []*proof.ExecutionState{}
			s.challengeResoutionCh <- struct{}{}
			return nil
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

func (s *Sequencer) Start() error {
	log.Info("Starting sequencer...")
	ctx, err := s.BaseService.Start()
	if err != nil {
		return fmt.Errorf("Failed to start sequencer: %w", err)
	}
	if err := s.Stake(ctx); err != nil {
		return fmt.Errorf("Failed to start sequencer: %w", err)
	}
	_, err = s.SyncL2ChainToL1Head(ctx, s.Config.L1RollupGenesisBlock)
	if err != nil {
		return fmt.Errorf("Failed to start sequencer: %w", err)
	}
	// We assume a single sequencer (us) for now, so we don't
	// need to sync transactions sequenced up.
	s.Wg.Add(4)
	go s.batchingLoop(ctx)
	go s.sequencingLoop(ctx)
	go s.confirmationLoop(ctx)
	go s.challengeLoop(ctx)
	log.Info("Sequencer started")
	return nil
}

func (s *Sequencer) APIs() []rpc.API {
	// TODO: sequencer APIs
	return []rpc.API{}
}
