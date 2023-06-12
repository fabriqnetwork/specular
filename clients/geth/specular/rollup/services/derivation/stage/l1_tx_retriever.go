package stage

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

// Retrieves transactions from L1 block associated with block ID returned from previous stage,
// according to filter query.
type L1TxRetriever struct {
	l1Client L1Client
	filterFn func(tx *ethTypes.Transaction) bool
	// Result queue.
	queue []filteredBlock
}

type filteredBlock struct {
	blockID types.BlockID
	txs     []*ethTypes.Transaction
}

type FilterQueryParams struct {
	Addresses []common.Address // contract addresses
	EventIDs  []common.Hash    // events emitted by the target transactions
}

type TxFilterParam struct {
	ContractAddress common.Address
	FunctionName    string
}

func NewL1TxRetriever(l1Client L1Client, filterFn func(tx *ethTypes.Transaction) bool) *L1TxRetriever {
	return &L1TxRetriever{l1Client: l1Client, filterFn: filterFn}
}

func (s *L1TxRetriever) hasNext() bool {
	return len(s.queue) > 0
}

func (s *L1TxRetriever) next() filteredBlock {
	txs := s.queue[0]
	s.queue = s.queue[1:]
	return txs
}

// Enqueues relevant transactions from L1 block associated with block ID returned from previous stage.
// Does not fetch new transactions unless queue is empty.
// Dequeues and returns one transaction batch.
func (s *L1TxRetriever) ingest(ctx context.Context, l1BlockID types.BlockID) error {
	// Retrieve the entire block.
	l1Block, err := s.l1Client.BlockByHash(ctx, l1BlockID.GetHash())
	if err != nil {
		return RetryableError{fmt.Errorf("Failed to get L1 block %s by hash: %w", l1BlockID, err)}
	}
	// Filter transactions in block by `filterFn`.
	s.queue = append(s.queue, filteredBlock{l1BlockID, filterTransactions(l1Block, s.filterFn)})
	return nil
}

func (s *L1TxRetriever) recover(ctx context.Context, l1BlockID types.BlockID) error {
	s.queue = nil
	return nil
}

func filterTransactions(block *ethTypes.Block, filterFn func(tx *ethTypes.Transaction) bool) []*ethTypes.Transaction {
	var txs []*ethTypes.Transaction
	for _, tx := range block.Transactions() {
		if filterFn(tx) {
			txs = append(txs, tx)
		}
	}
	return txs
}

// TODO: this is a hack for hardhat.
// func (s *L1TxRetriever) retrieveBlock(ctx context.Context, l1BlockID types.BlockID) (*ethTypes.Block, error) {
// 	block, err := s.l1Client.BlockByHash(ctx, l1BlockID.GetHash())
// 	if err == nil {
// 		return block, nil
// 	}
// 	log.Warn("Failed to get L1 block by hash; trying by number.", "number", l1BlockID.GetNumber(), "hash", l1BlockID.GetHash())
// 	block, err = s.l1Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(l1BlockID.GetNumber()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if block.Hash() != l1BlockID.GetHash() {
// 		return nil, fmt.Errorf("L1 block hash mismatch: got %s, expected %s", block.Hash(), l1BlockID.GetHash())
// 	}
// 	return block, nil
// }

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
// 	topics, err := abi.MakeTopics([]any{params.EventIDs})
// 	if err != nil {
// 		return ethereum.FilterQuery{}, err
// 	}
// 	return ethereum.FilterQuery{BlockHash: &blockHash, Addresses: params.Addresses, Topics: topics}, nil
// }

///////////

// type ExtractionStage struct {
// 	prev     Stage[any, []types.Transaction]
// 	l1Client EthClient
// }

// func (s *ExtractionStage) Step(ctx context.Context, _ any) (*types.Transaction, error) {
// 	// s.prev.Step(ctx, _)
// 	return nil, nil
// }

// func (s *ExtractionStage) Recover(ctx context.Context, l1BlockNumber uint64) error {
// 	return s.prev.Recover(ctx, l1BlockNumber)
// }
