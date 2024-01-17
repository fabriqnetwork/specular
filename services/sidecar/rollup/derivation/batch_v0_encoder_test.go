package derivation

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/require"
)

type config struct{}

func (c config) GetTargetBatchSize() uint64 { return 10 }

func TestProcessEmptyBlock(t *testing.T) {
	var (
		enc          = NewBatchV0Encoder(config{})
		baselineSize = enc.size()
		block        = types.NewBlock(&types.Header{Number: big.NewInt(0)}, nil, nil, nil, trie.NewStackTrie(nil))
	)
	err := enc.ProcessBlock(block, false)
	require.Equal(t, baselineSize, enc.size())
	require.Nil(t, err, err)
}

func TestProcessBlock(t *testing.T) {
	var (
		enc          = NewBatchV0Encoder(config{})
		baselineSize = enc.size()
		txs          = []*types.Transaction{types.NewTx(&types.LegacyTx{})}
		block        = types.NewBlock(&types.Header{Number: big.NewInt(0)}, txs, nil, nil, trie.NewStackTrie(nil))
	)
	// Test successful case
	err := enc.ProcessBlock(block, false)
	require.Greater(t, enc.size(), baselineSize)
	require.Nil(t, err, err)
	// Test failure case
	err = enc.ProcessBlock(block, false)
	require.Equal(t, err, errBatchFull)
}

func TestFlush(t *testing.T) {
	var (
		enc          = NewBatchV0Encoder(config{})
		baselineSize = enc.size()
		txs          = []*types.Transaction{types.NewTx(&types.LegacyTx{})}
		block        = types.NewBlock(&types.Header{Number: big.NewInt(0)}, txs, nil, nil, trie.NewStackTrie(nil))
	)
	// Flush an empty batch (unforced)
	_, err := enc.Flush(false)
	require.Equal(t, errBatchTooSmall, err)
	// Flush an empty batch (forced)
	batch, err := enc.Flush(true)
	require.Nil(t, err)
	require.NotNil(t, batch)
	require.Equal(t, baselineSize, enc.size())
	// Process a block
	err = enc.ProcessBlock(block, false)
	require.Greater(t, enc.size(), baselineSize)
	require.Equal(t, len(enc.subBatches), 1)
	require.Nil(t, err, err)
	// Flush a non-empty batch
	batch, err = enc.Flush(true)
	require.Nil(t, err, err)
	require.NotNil(t, batch)
	require.Equal(t, baselineSize, enc.size())
}

func TestIsEmpty(t *testing.T) {
	var (
		enc   = NewBatchV0Encoder(config{})
		txs   = []*types.Transaction{types.NewTx(&types.LegacyTx{})}
		block = types.NewBlock(&types.Header{Number: big.NewInt(0)}, txs, nil, nil, trie.NewStackTrie(nil))
	)
	require.True(t, enc.IsEmpty())
	err := enc.ProcessBlock(block, false)
	require.Nil(t, err, err)
	require.False(t, enc.IsEmpty())
}
