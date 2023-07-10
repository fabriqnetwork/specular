package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
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

var EmptyBlockID = BlockID{}

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

type L2BlockRef struct {
	BlockRef
	L1Origin BlockID
	Time     uint64
}

func NewL2BlockRef(number uint64, hash common.Hash, parentHash common.Hash, l1Origin BlockID, time uint64) L2BlockRef {
	return L2BlockRef{NewBlockRef(number, hash, parentHash), l1Origin, time}
}

func NewL2BlockRefFromHeader(header *types.Header, l1Origin BlockID) L2BlockRef {
	return L2BlockRef{NewBlockRefFromHeader(header), l1Origin, header.Time}
}

func (ref L2BlockRef) GetL1Origin() BlockID { return ref.L1Origin }

func (ref L2BlockRef) GetTime() uint64 { return ref.Time }
