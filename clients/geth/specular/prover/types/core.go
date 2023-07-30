package types

import (
	"bytes"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Transaction interface {
	AccessList() types.AccessList
	AsMessage(s interface{}, baseFee *big.Int, rules interface{}) (interface{}, error)
	ChainId() *big.Int
	Data() []byte
	EncodeRLP(w io.Writer) error
	Gas() uint64
	GasPrice() *big.Int
	Hash() common.Hash
	// MarshalBinary() ([]byte, error)
	// MarshalJSON() ([]byte, error)
	Nonce() uint64
	RawSignatureValues() (v *big.Int, r *big.Int, s *big.Int)
	To() *common.Address
	Type() uint8
	// UnmarshalBinary(b []byte) error
	// UnmarshalJSON(input []byte) error
	Value() *big.Int
}

type Transactions []Transaction

func (txs Transactions) EncodeIndex(i int, w *bytes.Buffer) {
	rlp.Encode(w, txs[i])
}

func (txs Transactions) Len() int {
	return len(txs)
}

type Block interface {
	BaseFee() *big.Int
	Bloom() types.Bloom
	Coinbase() common.Address
	Difficulty() *big.Int
	EncodeRLP(w io.Writer) error
	Extra() []byte
	GasLimit() uint64
	GasUsed() uint64
	Hash() common.Hash
	Header() *types.Header
	MixDigest() common.Hash
	Nonce() uint64
	Number() *big.Int
	NumberU64() uint64
	ParentHash() common.Hash
	ReceiptHash() common.Hash
	Root() common.Hash
	Time() uint64
	Transaction(hash common.Hash) Transaction
	Transactions() Transactions
	TxHash() common.Hash
}
