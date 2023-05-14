package stage

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
)

// Retrieves transactions from L1 block associated with block ID returned from previous stage,
// according to filter query.
type L1TxRetrievalStage struct {
	prev     Stage[l2types.BlockID]
	l1Client EthClient
	filterFn func(tx *types.Transaction) bool
	// Result queue.
	queue []filteredBlock
}

type filteredBlock struct {
	blockID l2types.BlockID
	txs     []*types.Transaction
}

type FilterQueryParams struct {
	Addresses []common.Address // contract addresses
	EventIDs  []common.Hash    // events emitted by the target transactions
}

type TxFilterParam struct {
	ContractAddress common.Address
	FunctionName    string
}

// Enqueues relevant transactions from L1 block associated with block ID returned from previous stage.
// Does not fetch new transactions unless queue is empty.
// Dequeues and returns one transaction batch.
func (s *L1TxRetrievalStage) Step(ctx context.Context) (filteredBlock, error) {
	if len(s.queue) > 0 {
		return s.dequeue(), nil
	}
	// Get the next block ID.
	l1BlockID, err := s.prev.Step(ctx)
	if err != nil {
		return filteredBlock{}, err
	}
	// Retrieve the entire block.
	l1Block, err := s.l1Client.BlockByHash(ctx, l1BlockID.Hash())
	if err != nil {
		return filteredBlock{}, &RetryableError{fmt.Errorf("Failed to get L1 block by hash: %w", err)}
	}
	// Filter transactions in block by `filterFn`.
	txs := filterTransactions(l1Block, s.filterFn)
	s.queue = append(s.queue, filteredBlock{l1BlockID, txs})
	return s.dequeue(), nil
}

func (s *L1TxRetrievalStage) Recover(ctx context.Context, l1BlockID l2types.BlockID) error {
	s.queue = nil
	return s.prev.Recover(ctx, l1BlockID)
}

func (s *L1TxRetrievalStage) dequeue() filteredBlock {
	if len(s.queue) == 0 {
		return filteredBlock{}
	}
	txs := s.queue[0]
	s.queue = s.queue[1:]
	return txs
}

func filterTransactions(block *types.Block, filterFn func(tx *types.Transaction) bool) []*types.Transaction {
	var txs []*types.Transaction
	for _, tx := range block.Transactions() {
		if filterFn(tx) {
			txs = append(txs, tx)
		}
	}
	return txs
}

// .....
// query, err := createFilterQuery(l1BlockID.Hash(), s.params)
// if err != nil {
// 	return nil, err
// }
// err = s.filterAndEnqueue(ctx, query)
// if err != nil {
// 	return nil, err
// }

// func (s *L1TxRetrievalStage) filterAndEnqueue(ctx context.Context, query ethereum.FilterQuery) error {
// 	logs, err := s.l1Client.FilterLogs(ctx, query)
// 	if err != nil {
// 		return &RetryableError{fmt.Errorf("Failed to filter logs: %w", err)}
// 	}
// 	if len(logs) == 0 {
// 		return nil
// 	}
// 	var batch transactions
// 	for _, log := range logs {
// 		tx, _, err := s.l1Client.TransactionByHash(ctx, log.TxHash)
// 		if err != nil {
// 			return &RetryableError{fmt.Errorf("failed to get tx corresponding to log by hash: %w", err)}
// 		}
// 		batch = append(batch, tx)
// 	}
// 	s.queue = append(s.queue, batch)
// 	return nil
// }

// func createFilterQuery(blockHash common.Hash, params FilterQueryParams) (ethereum.FilterQuery, error) {
// 	topics, err := abi.MakeTopics([]interface{}{params.EventIDs})
// 	if err != nil {
// 		return ethereum.FilterQuery{}, err
// 	}
// 	return ethereum.FilterQuery{BlockHash: &blockHash, Addresses: params.Addresses, Topics: topics}, nil
// }

///////////

// type ExtractionStage struct {
// 	prev     Stage[interface{}, []types.Transaction]
// 	l1Client EthClient
// }

// func (s *ExtractionStage) Step(ctx context.Context, _ interface{}) (*types.Transaction, error) {
// 	// s.prev.Step(ctx, _)
// 	return nil, nil
// }

// func (s *ExtractionStage) Recover(ctx context.Context, l1BlockNumber uint64) error {
// 	return s.prev.Recover(ctx, l1BlockNumber)
// }
