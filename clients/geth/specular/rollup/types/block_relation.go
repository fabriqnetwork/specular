package types

import (
	"github.com/specularl2/specular/clients/geth/specular/utils"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

// Mapping from L2 block to the L1 block it corresponds to.
type BlockRelation struct {
	L1BlockID BlockID `json:"l1_block_id"`
	L2BlockID BlockID `json:"l2_block_id"`
}

var EmptyRelation = BlockRelation{}

func NewBlockRelation(l1BlockID BlockID, l2BlockID BlockID) BlockRelation {
	return BlockRelation{L1BlockID: l1BlockID, L2BlockID: l2BlockID}
}

type BlockRelations []BlockRelation

func (r *BlockRelations) Append(relation BlockRelation) error {
	sl := *r
	if len(sl) > 0 && relation.L1BlockID.Number <= sl[len(sl)-1].L1BlockID.Number {
		return fmt.Errorf("block relation %s is not newer than %s", relation, sl[len(sl)-1])
	}
	*r = append(sl, relation)
	return nil
}

// Marks the latest L2 block corresponding to l1BlockNumber as safe.
// l1BlockNumber is any valid L1 block number (unsafe/safe/finalized).
// If multiple L2 blocks belong to the same L1 block, use the latest L2 block.
func (r *BlockRelations) MarkSafe(l1BlockNumber uint64) BlockID {
	idx := r.indexOfLastBlockRelationOrPrior(l1BlockNumber)
	if idx == -1 {
		// TODO: handle
		return BlockID{}
	}
	return (*r)[idx].L2BlockID
}

// Clears the latest L2 block corresponding to l1BlockNumber (and any preceding it) as finalized.
func (r *BlockRelations) MarkFinal(finalL1BlockNumber uint64) BlockID {
	idx := r.indexOfLastBlockRelationOrPrior(finalL1BlockNumber)
	if idx == -1 {
		// TODO: handle
		return BlockID{}
	}
	latestFinalL2BlockID := (*r)[idx].L2BlockID
	// Remove from slice since it's final.
	*r = (*r)[idx+1:]
	return latestFinalL2BlockID
}

func (r *BlockRelations) MarkReorgedOut(existingL1BlockNumber uint64) {
	idx := r.indexOfLastBlockRelationOrPrior(existingL1BlockNumber)
	*r = (*r)[idx+1:]
}

func (r *BlockRelations) indexOfLastBlockRelationOrPrior(targetL1BlockNumber uint64) int {
	return utils.IndexOfMappedLEq(
		*r, targetL1BlockNumber, func(relation BlockRelation) uint64 { return relation.L1BlockID.GetNumber() },
	)
}
