package derivation

// func (s *Syncer) SyncLoop(ctx context.Context, start uint64, newBatchCh chan<- struct{}) error {
// 	// Start watching for new TxBatchAppended events.
// 	subCtx, cancel := context.WithCancel(ctx)
// 	batchEventCh := client.SubscribeHeaderMapped[*bindings.ISequencerInboxTxBatchAppended](
// 		subCtx, s.L1Syncer.LatestHeaderBroker, s.L1Client.FilterTxBatchAppendedEvents, start,
// 	)
// 	defer cancel()
// 	// Process TxBatchAppended events.
// 	for {
// 		select {
// 		case ev := <-batchEventCh:
// 			log.Info("Processing `TxBatchAppended` event", "l1Block", ev.Raw.BlockNumber)
// 			err := s.processTxBatchAppendedEvent(ctx, ev)
// 			if err != nil {
// 				return fmt.Errorf("Failed to process event: %w", err)
// 			}
// 			if newBatchCh != nil {
// 				newBatchCh <- struct{}{}
// 			}
// 		case <-ctx.Done():
// 			return nil
// 		}
// 	}
// }

// Sync to current L1 block head and commit blocks.
// `start` is the block number to start syncing from.
// Returns the last synced block number (inclusive).
// func (s *Syncer) SyncL2ChainToL1Head(ctx context.Context, start uint64) (uint64, error) {
// 	l1BlockHead, err := s.L1Client.BlockNumber(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
// 	}
// 	opts := bind.FilterOpts{Start: start, End: &l1BlockHead, Context: ctx}
// 	eventsIter, err := s.L1Client.FilterTxBatchAppendedEvents(&opts)
// 	if err != nil {
// 		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
// 	}
// 	err = s.processTxBatchAppendedEvents(ctx, eventsIter)
// 	if err != nil {
// 		return 0, fmt.Errorf("Failed to sync to L1 head, err: %w", err)
// 	}
// 	log.Info(
// 		"Synced L1->L2",
// 		"l1 start", start,
// 		"l1 end", l1BlockHead,
// 	)
// 	return l1BlockHead, nil
// }
