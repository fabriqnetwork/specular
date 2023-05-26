package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
)

// TODO: tags
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

// TODO: unused
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

// TODO: unused
type L2BlockRef struct {
	BlockRef
	l1Origin BlockID
}

func (ref L2BlockRef) L1Origin() BlockID { return ref.l1Origin }

// Mapping from L2 block to the L1 block it corresponds to.
type BlockRelation struct {
	L1BlockID BlockID `json:"l1_block_id"`
	L2BlockID BlockID `json:"l2_block_id"`
}

type BlockRelations []BlockRelation

func (r BlockRelations) Append(relation BlockRelation) error {
	// if r[len(r) - 1].L1BlockRef.Number >= relation.L1BlockRef.Number {
	// }
	r = append(r, relation)
	return nil
}

// Marks the latest L2 block corresponding to l1BlockNumber as safe.
// l1BlockNumber is any valid L1 block number (unsafe/safe/finalized).
// If multiple L2 blocks belong to the same L1 block, use the latest L2 block.
func (r BlockRelations) MarkSafe(l1BlockNumber uint64) BlockID {
	idx := indexOfLastBlockRelationOrPrior(r, l1BlockNumber)
	if idx == -1 {
		// TODO: handle
		return BlockID{}
	}
	return r[idx].L2BlockID
}

// This clears any previous L2 blocks from the slice as well, since they must also now be final.
func (r BlockRelations) MarkFinal(finalL1BlockNumber uint64) BlockID {
	idx := indexOfLastBlockRelationOrPrior(r, finalL1BlockNumber)
	if idx == -1 {
		// TODO: handle
		return BlockID{}
	}
	latestFinalL2BlockID := r[idx].L2BlockID
	// Remove from slice since it's final.
	r = r[idx+1:]
	return latestFinalL2BlockID
}

func (r BlockRelations) MarkReorgedOut(existingL1BlockNumber uint64) {
	idx := indexOfLastBlockRelationOrPrior(r, existingL1BlockNumber)
	r = r[idx+1:]
}

func indexOfLastBlockRelationOrPrior(relations []BlockRelation, targetL1BlockNumber uint64) int {
	return utils.IndexOfMappedLEq(
		relations, targetL1BlockNumber, func(relation BlockRelation) uint64 { return relation.L1BlockID.GetNumber() },
	)
}
