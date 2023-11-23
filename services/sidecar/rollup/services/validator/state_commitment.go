package validator

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
)

var (
	ErrInvalidStateCommitment        = errors.New("invalid state commitment")
	ErrInvalidStateCommitmentVersion = errors.New("invalid state commitment version")

	StateCommitmentVersionV0 = types.Bytes32{}
)

type VersionedStateCommitment interface {
	// Version returns the version of the L2 state commitment
	Version() types.Bytes32

	// Marshal a L2 state commitment into a byte slice for hashing
	Marshal() []byte
}

type StateCommitmentV0 struct {
	l2VmHash common.Hash
}

func (o *StateCommitmentV0) Version() types.Bytes32 {
	return StateCommitmentVersionV0
}

func (o *StateCommitmentV0) Marshal() []byte {
	var buf [64]byte
	version := o.Version()
	copy(buf[:32], version[:])
	copy(buf[32:], o.l2VmHash[:])
	return buf[:]
}

// StateCommitment returns the keccak256 hash of the marshaled L2 state commitment
func StateCommitment(stateCommitment VersionedStateCommitment) types.Bytes32 {
	marshaled := stateCommitment.Marshal()
	return types.Bytes32(crypto.Keccak256Hash(marshaled))
}

func UnmarshalStateCommitment(data []byte) (VersionedStateCommitment, error) {
	if len(data) < 32 {
		return nil, ErrInvalidStateCommitment
	}
	var ver types.Bytes32
	copy(ver[:], data[:32])
	switch ver {
	case StateCommitmentVersionV0:
		return unmarshalStateCommitmentV0(data)
	default:
		return nil, ErrInvalidStateCommitmentVersion
	}
}

func unmarshalStateCommitmentV0(data []byte) (*StateCommitmentV0, error) {
	if len(data) != 64 {
		return nil, ErrInvalidStateCommitment
	}
	var l2State StateCommitmentV0
	// data[:32] is the version
	copy(l2State.l2VmHash[:], data[32:64])
	return &l2State, nil
}
