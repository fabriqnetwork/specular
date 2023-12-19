package validator

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/specularL2/specular/services/sidecar/utils/log"
)

var (
	ErrInvalidStateCommitment        = errors.New("invalid state commitment")
	ErrInvalidStateCommitmentVersion = errors.New("invalid state commitment version")

	StateCommitmentVersionV0 = Bytes32{}
)

type Bytes32 = [32]byte

type VersionedStateCommitment interface {
	// Version returns the version of the L2 state commitment
	Version() Bytes32

	// Marshal a L2 state commitment into a byte slice for hashing
	Marshal() []byte
}

type StateCommitmentV0 struct {
	l2BlockHash common.Hash
	l2StateRoot common.Hash
}

func (o *StateCommitmentV0) Version() Bytes32 {
	return StateCommitmentVersionV0
}

func (o *StateCommitmentV0) Marshal() []byte {
	var buf [96]byte
	version := o.Version()
	copy(buf[:32], version[:])
	copy(buf[32:64], o.l2BlockHash[:])
	copy(buf[64:], o.l2StateRoot[:])
	return buf[:]
}

// StateCommitment returns the keccak256 hash of the marshaled L2 state commitment
func StateCommitment(stateCommitment VersionedStateCommitment) common.Hash {
	marshaled := stateCommitment.Marshal()
	log.Info("hashed state commitment", "hex", common.Bytes2Hex(crypto.Keccak256(marshaled)))
	return crypto.Keccak256Hash(marshaled)
}

func UnmarshalStateCommitment(data []byte) (VersionedStateCommitment, error) {
	if len(data) < 32 {
		return nil, ErrInvalidStateCommitment
	}
	var ver Bytes32
	copy(ver[:], data[:32])
	switch ver {
	case StateCommitmentVersionV0:
		return unmarshalStateCommitmentV0(data)
	default:
		return nil, ErrInvalidStateCommitmentVersion
	}
}

func unmarshalStateCommitmentV0(data []byte) (*StateCommitmentV0, error) {
	if len(data) != 96 {
		return nil, ErrInvalidStateCommitment
	}
	var l2State StateCommitmentV0
	// data[:32] is the version
	copy(l2State.l2BlockHash[:], data[32:64])
	copy(l2State.l2StateRoot[:], data[64:])
	return &l2State, nil
}
