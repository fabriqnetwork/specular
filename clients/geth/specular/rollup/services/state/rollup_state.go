package state

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
)

// Tracks L2 state as a function of synced L1 state.
type RollupState struct {
	// Synced L1 state and derived L2 state.
	L1State *L1State
	L2State *L2State

	L1Syncer *client.EthSyncer
	L2Syncer *client.EthSyncer
}

func NewRollupState() *RollupState {
	l1State := NewL1State()
	l2State := NewL2State()
	return &RollupState{
		L1State:  l1State,
		L2State:  l2State,
		L1Syncer: client.NewEthSyncer(l1State),
		L2Syncer: client.NewEthSyncer(l2State),
	}
}

// Sequencer:      tx -> Execute
// L1State -> L2State ->    |->     Sequence (unsafe l2 blocks)
// Validator:
// L1State -> L2State -> 			CreateAssertion (safe l2 blocks)
//
//	-> read existing assertions -> validate -^
//	-> sync    -> Execute
//	-> L2State -> ChallengeAssertion
func (r *RollupState) StartSync(ctx context.Context, l1Client, l2Client client.EthPollingClient) {
	// Sync headers from L1.
	r.L1Syncer.Start(ctx, l1Client)
	// Sync headers from L2.
	r.L2Syncer.Start(ctx, l2Client)
	// TODO: consider moving
	for r.L1State.Head() == nil {
		log.Info("Waiting for L1 latest header...")
		time.Sleep(100 * time.Millisecond)
	}
	log.Info("Latest header received", "number", r.L1State.Head().Number)
	for r.L1State.Finalized() == nil {
		log.Info("Waiting for L1 finalized header...")
		time.Sleep(100 * time.Millisecond)
	}
	log.Info("Finalized header received", "number", r.L1State.Head().Number)
}

func (r *RollupState) StopSync(ctx context.Context) {
	r.L1Syncer.Stop(ctx)
	r.L2Syncer.Stop(ctx)
}
