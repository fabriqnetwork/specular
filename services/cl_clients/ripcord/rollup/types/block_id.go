package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/services/cl_clients/ripcord/utils/fmt"
)

var (
	EmptyBlockID  = BlockID{}
	EmptyBlockRef = BlockRef{}
)

type BlockID struct {
	Number uint64      `json:"number"`
	Hash   common.Hash `json:"hash"`
}

func NewBlockID(number uint64, hash common.Hash) BlockID { return BlockID{number, hash} }

func NewBlockIDFromHeader(header *types.Header) BlockID {
	return NewBlockID(header.Number.Uint64(), header.Hash())
}

func (id BlockID) GetNumber() uint64    { return id.Number }
func (id BlockID) GetHash() common.Hash { return id.Hash }
func (id BlockID) String() string       { return fmt.Sprintf("(#=%d|h=%s)", id.Number, id.Hash) }

type BlockRef struct {
	BlockID
	ParentHash common.Hash `json:"parent_hash"`
}

func NewBlockRef(number uint64, hash common.Hash, parentHash common.Hash) BlockRef {
	return BlockRef{NewBlockID(number, hash), parentHash}
}

func NewBlockRefFromHeader(header *types.Header) BlockRef {
	return BlockRef{NewBlockIDFromHeader(header), header.ParentHash}
}

func (ref BlockRef) GetParentHash() common.Hash { return ref.ParentHash }
