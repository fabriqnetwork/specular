package l2types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockID struct {
	number uint64
	hash   common.Hash
}

func NewBlockID(number uint64, hash common.Hash) BlockID { return BlockID{number, hash} }

func NewBlockIDFromHeader(header *types.Header) BlockID {
	return NewBlockID(header.Number.Uint64(), header.Hash())
}

func (id BlockID) Number() uint64    { return id.number }
func (id BlockID) Hash() common.Hash { return id.hash }

// Mapping from L2 block to the L1 block it corresponds to.
type BlockRelation struct {
	L1BlockID BlockID
	L2BlockID BlockID
}

type BlockRelations []BlockRelation

func (r BlockRelations) Append(relation BlockRelation) error {
	// if r[len(r) - 1].L1BlockRef.Number >= relation.L1BlockRef.Number {
	// }
	r = append(r, relation)
	return nil
}

// If multiple L2 blocks belong to the same L1 block, use the latest L2 block.
func (r BlockRelations) MarkSafe(safeL1BlockNumber uint64) BlockID {
	idx := findLastBlockRelationOrPriorIndex(r, safeL1BlockNumber)
	if idx == -1 {
		// TODO: handle
		return BlockID{}
	}
	return r[idx].L2BlockID
}

func (r BlockRelations) MarkFinal(finalL1BlockNumber uint64) BlockID {
	idx := findLastBlockRelationOrPriorIndex(r, finalL1BlockNumber)
	if idx == -1 {
		// TODO: handle
		return BlockID{}
	}
	latestSafeL2BlockRef := r[idx].L2BlockID
	// Remove from slice since it's final.
	r = r[idx+1:]
	return latestSafeL2BlockRef
}

func (r BlockRelations) MarkReorgedOut(existingL1BlockNumber uint64) {
	idx := findLastBlockRelationOrPriorIndex(r, existingL1BlockNumber)
	r = r[idx+1:]
}

// Assumes array sorted by strictly increasing L1 block numbers.
func findLastBlockRelationOrPriorIndex(relations []BlockRelation, targetL1BlockNumber uint64) int {
	start := 0
	end := len(relations)
	var mid int
	midL1BlockRef := relations[mid].L1BlockID
	for start <= end {
		mid = (start + end) / 2
		if targetL1BlockNumber < midL1BlockRef.Number() {
			end = mid - 1
		} else if targetL1BlockNumber > midL1BlockRef.Number() {
			start = mid + 1
		} else {
			return mid
		}
	}
	// Return index of last block# < l1BlockNumber.
	if midL1BlockRef.Number() < targetL1BlockNumber {
		return mid
	} else {
		return mid - 1
	}
}
