package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

func NewAssertionFrom(
	assertion *bindings.IRollupAssertion,
	event *bindings.IRollupAssertionCreated,
) *rollupTypes.Assertion {
	// TODO: set StartBlock, EndBlock if necessary (or remove from this struct).
	return &rollupTypes.Assertion{
		ID:        event.AssertionID,
		VmHash:    event.VmHash,
		InboxSize: assertion.InboxSize,
		Deadline:  assertion.Deadline,
	}
}

// For debugging purposes.
func LogBlockChainInfo(backend Backend, start, end uint64) {
	for i := start; i < end; i++ {
		block := backend.BlockChain().GetBlockByNumber(i)
		log.Info("Block", "number", i, "hash", block.Hash(), "root", block.Root(), "num txs", len(block.Transactions()))
	}
}

func SubmitOneStepProof(
	ctx context.Context,
	proofBackend proof.Backend,
	l1Client client.L1BridgeClient,
	state *proof.ExecutionState,
	challengedStepIndex *big.Int,
	prevBisection [][32]byte,
	prevChallengedSegmentStart *big.Int,
	prevChallengedSegmentLength *big.Int,
) error {
	osp, err := proof.GenerateProof(proofBackend, ctx, state, nil)
	if err != nil {
		log.Crit("UNHANDLED: osp generation failed", "err", err)
	}
	_, err = l1Client.VerifyOneStepProof(
		uint8(osp.VerifierType),
		osp.Encode(),
		[]byte{},                                 // TODO: fix
		bindings.VerificationContextRawContext{}, // TODO: fix
		challengedStepIndex,
		prevBisection,
		prevChallengedSegmentStart,
		prevChallengedSegmentLength,
	)
	log.Info("OSP submitted")
	if err != nil {
		log.Error("OSP verification failed")
	}
	return err
}

func RespondBisection(
	ctx context.Context,
	proofBackend proof.Backend,
	l1Client client.L1BridgeClient,
	ev *bindings.ISymChallengeBisected,
	states []*proof.ExecutionState,
	opponentEndStateHash common.Hash,
	isDefender bool,
) error {
	// Get bisection info from event
	segStart := ev.ChallengedSegmentStart.Uint64()
	segLen := ev.ChallengedSegmentLength.Uint64()
	// Get previous bisections from call data
	tx, _, err := l1Client.TransactionByHash(ctx, ev.Raw.TxHash)
	if err != nil {
		// TODO: error handling
		log.Error("Failed to get challenge data", "error", err)
		return nil
	}
	decoded, err := l1Client.DecodeBisectExecutionInput(tx)
	if err != nil {
		if isDefender {
			// Defender always starts first
			log.Error("Failed to decode bisection data", "error", err)
			return nil
		}
		// We are in the first round when the defender calls initializeChallengeLength
		// Get initialized challenge length from event
		steps := segLen
		if steps != uint64(len(states)-1) {
			log.Crit("UNHANDLED: currently not support diverge on steps")
		}
		prevBisection := [][32]byte{
			states[0].Hash(),
			opponentEndStateHash,
		}
		if segLen == 1 {
			// This assertion only has one step
			err = SubmitOneStepProof(
				ctx,
				proofBackend,
				l1Client,
				states[0],
				common.Big1,
				prevBisection,
				ev.ChallengedSegmentStart,
				ev.ChallengedSegmentLength,
			)
			if err != nil {
				log.Crit("UNHANDLED: osp failed")
			}
		} else {
			// This assertion has multiple steps
			startState := states[0].Hash()
			midState := states[steps/2+steps%2].Hash()
			endState := states[steps].Hash()
			bisection := [][32]byte{
				startState,
				midState,
				endState,
			}
			_, err := l1Client.BisectExecution(
				bisection,
				common.Big1,
				prevBisection,
				ev.ChallengedSegmentStart,
				ev.ChallengedSegmentLength,
			)
			log.Info("BisectExecution", "bisection", bisection, "cidx", common.Big1, "psegStart", segStart, "psegLen", segLen, "prev", prevBisection)
			if err != nil {
				log.Crit("UNHANDLED: bisection excution failed", "err", err)
			}
		}
		return nil
	}
	prevBisection := decoded[0].([][32]byte)
	startState := states[segStart].Hash()
	midState := states[segStart+segLen/2+segLen%2].Hash()
	endState := states[segStart+segLen].Hash()
	if segLen == 1 {
		// We've reached one step
		err = SubmitOneStepProof(
			ctx,
			proofBackend,
			l1Client,
			states[segStart],
			common.Big1,
			prevBisection,
			ev.ChallengedSegmentStart,
			ev.ChallengedSegmentLength,
		)
		if err != nil {
			log.Crit("UNHANDLED: osp failed")
		}
	} else {
		challengeIdx := uint64(1)
		if prevBisection[1] == midState {
			challengeIdx = 2
		}
		if segLen == 2 || (segLen == 3 && challengeIdx == 2) {
			// The next challenge segment is a single step
			stateIndex := segStart
			if challengeIdx != 1 {
				stateIndex = segStart + segLen/2
			}
			err = SubmitOneStepProof(
				ctx,
				proofBackend,
				l1Client,
				states[stateIndex],
				new(big.Int).SetUint64(challengeIdx),
				prevBisection,
				ev.ChallengedSegmentStart,
				ev.ChallengedSegmentLength,
			)
			if err != nil {
				log.Crit("UNHANDLED: osp failed")
			}
		} else {
			var newLen uint64 // New segment length
			var bisection [][32]byte
			if challengeIdx == 1 {
				newLen = segLen/2 + segLen%2
				bisection = [][32]byte{
					startState,
					states[segStart+newLen/2+newLen%2].Hash(),
					midState,
				}
			} else {
				newLen = segLen / 2
				bisection = [][32]byte{
					midState,
					states[segStart+segLen/2+segLen%2+newLen/2+newLen%2].Hash(),
					endState,
				}
			}
			_, err := l1Client.BisectExecution(
				bisection,
				new(big.Int).SetUint64(challengeIdx),
				prevBisection,
				ev.ChallengedSegmentStart,
				ev.ChallengedSegmentLength,
			)
			log.Info("BisectExecution", "bisection", bisection, "cidx", challengeIdx, "psegStart", segStart, "psegLen", segLen, "prev", prevBisection)
			if err != nil {
				log.Crit("UNHANDLED: bisection excution failed", "err", err)
			}
		}
	}
	return nil
}
