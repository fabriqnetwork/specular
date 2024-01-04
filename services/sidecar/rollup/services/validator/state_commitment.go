package validator

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/specularL2/specular/services/sidecar/utils/log"
)

const V0BufSize = 96
const V0VersionSize = 32
const V0BlockHashOffset = 32
const V0BlockHashSize = 32
const V0StateRootOffset = 64

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
	var buf [V0BufSize]byte
	version := o.Version()
	copy(buf[:V0VersionSize], version[:])
	copy(buf[V0BlockHashOffset:V0BlockHashSize], o.l2BlockHash[:])
	copy(buf[V0StateRootOffset:], o.l2StateRoot[:])
	return buf[:]
}

// StateCommitment returns the keccak256 hash of the marshaled L2 state commitment
func StateCommitment(stateCommitment VersionedStateCommitment) common.Hash {
	marshaled := stateCommitment.Marshal()
	log.Info("hashed state commitment", "hex", common.Bytes2Hex(crypto.Keccak256(marshaled)))
	return crypto.Keccak256Hash(marshaled)
}

func UnmarshalStateCommitment(data []byte) (VersionedStateCommitment, error) {
	if len(data) < V0VersionSize {
		return nil, ErrInvalidStateCommitment
	}
	var ver Bytes32
	copy(ver[:], data[:V0VersionSize])
	switch ver {
	case StateCommitmentVersionV0:
		return unmarshalStateCommitmentV0(data)
	default:
		return nil, ErrInvalidStateCommitmentVersion
	}
}

func unmarshalStateCommitmentV0(data []byte) (*StateCommitmentV0, error) {
	if len(data) != V0BufSize {
		return nil, ErrInvalidStateCommitment
	}
	var l2State StateCommitmentV0
	// data[:32] is the version
	copy(l2State.l2BlockHash[:], data[V0BlockHashOffset:V0BlockHashSize])
	copy(l2State.l2StateRoot[:], data[V0StateRootOffset:])
	return &l2State, nil
}
