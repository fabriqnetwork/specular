package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	rollupTypes "github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

func NewAssertionFrom(
	assertion *bindings.IRollupAssertion,
	event *bindings.RollupBaseAssertionCreated,
) *rollupTypes.Assertion {
	// TODO: set StartBlock, EndBlock if necessary (or remove from this struct).
	return &rollupTypes.Assertion{
		ID:        event.AssertionID,
		VmHash:    event.VmHash,
		InboxSize: assertion.InboxSize,
		Deadline:  assertion.Deadline,
	}
}

type challengeClient interface {
	VerifyOneStepProof(
		proof []byte,
		txInclusionProof []byte,
		verificationRawCtx bindings.VerificationContextLibRawContext,
		challengedStepIndex *big.Int,
		prevBisection [][32]byte,
		prevChallengedSegmentStart *big.Int,
		prevChallengedSegmentLength *big.Int,
	) (*types.Transaction, error)
}

type challengeClient interface {
	VerifyOneStepProof(
		proof []byte,
		txInclusionProof []byte,
		verificationRawCtx bindings.VerificationContextLibRawContext,
		challengedStepIndex *big.Int,
		prevBisection [][32]byte,
		prevChallengedSegmentStart *big.Int,
		prevChallengedSegmentLength *big.Int,
	) (*types.Transaction, error)
}

func SubmitOneStepProof(
	ctx context.Context,
	proofBackend proof.Backend,
	l1Client challengeClient,
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
		osp.Encode(),
		[]byte{}, // TODO: fix
		bindings.VerificationContextLibRawContext{}, // TODO: fix
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
