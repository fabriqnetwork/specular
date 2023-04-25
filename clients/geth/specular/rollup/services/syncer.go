package services

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/data"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

type Syncer struct {
	backend ExecutionBackend
	// TODO: reduce client scope
	L1Client EthBridgeClient
	L1Syncer *client.EthSyncer
}

func NewSyncer(
	backend ExecutionBackend,
	l1Client EthBridgeClient,
	l1Syncer *client.EthSyncer,
) *Syncer {
	return &Syncer{backend, l1Client, l1Syncer}
}

func (s *Syncer) SyncLoop(ctx context.Context, start uint64, newBatchCh chan<- struct{}) {
	// Start watching for new TxBatchAppended events.
	subCtx, cancel := context.WithCancel(ctx)
	batchEventCh := client.SubscribeHeaderMapped[*bindings.ISequencerInboxTxBatchAppended](
		subCtx, s.L1Syncer.LatestHeaderBroker, s.L1Client.FilterTxBatchAppendedEvents, start,
	)
	defer cancel()
	// Process TxBatchAppended events.
	for {
		select {
		case ev := <-batchEventCh:
			log.Info("Processing `TxBatchAppended` event", "l1Block", ev.Raw.BlockNumber)
			err := s.processTxBatchAppendedEvent(ctx, ev)
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

// Sync to current L1 block head and commit blocks.
// `start` is the block number to start syncing from.
// Returns the last synced block number (inclusive).
func (s *Syncer) SyncL2ChainToL1Head(ctx context.Context, start uint64) (uint64, error) {
	l1BlockHead, err := s.L1Client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	opts := bind.FilterOpts{Start: start, End: &l1BlockHead, Context: ctx}
	eventsIter, err := s.L1Client.FilterTxBatchAppendedEvents(&opts)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	err = s.processTxBatchAppendedEvents(ctx, eventsIter)
	if err != nil {
		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
	}
	log.Info(
		"Synced L1->L2",
		"l1 start", start,
		"l1 end", l1BlockHead,
	)
	return l1BlockHead, nil
}

func (s *Syncer) processTxBatchAppendedEvents(
	ctx context.Context,
	eventsIter *bindings.ISequencerInboxTxBatchAppendedIterator,
) error {
	for eventsIter.Next() {
		err := s.processTxBatchAppendedEvent(ctx, eventsIter.Event)
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
func (s *Syncer) processTxBatchAppendedEvent(
	ctx context.Context,
	ev *bindings.ISequencerInboxTxBatchAppended,
) error {
	tx, _, err := s.L1Client.TransactionByHash(ctx, ev.Raw.TxHash)
	if err != nil {
		return fmt.Errorf("Failed to get transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Decode input to appendTxBatch transaction.
	decoded, err := client.UnpackAppendTxBatchInput(tx)
	if err != nil {
		return fmt.Errorf("Failed to decode transaction associated with TxBatchAppended event, err: %w", err)
	}
	// Construct batch.
	blocks, err := data.BlocksFromDecoded(decoded)
	if err != nil {
		return fmt.Errorf("Failed to split AppendTxBatch input into batches, err: %w", err)
	}
	log.Info("Decoded batch", "#blocks", len(blocks))
	for _, block := range blocks {
		s.backend.CommitPayload(block)
	}
	return nil
}
