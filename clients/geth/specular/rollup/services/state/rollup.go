package state

// // Should only be called as goroutine.
// func (r *Rollup) deriveRollupState(ctx context.Context) {
// 	rollupHeaderCh := make(chan *types.Header)
// 	r.SubscribeHead(ctx, rollupHeaderCh)
// 	for {
// 		select {
// 		case header := <-rollupHeaderCh:
// 			r.SysState.L2State.UpdateHead(header)
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// 	// r.SysState.L2State.UpdateSafe()
// }
