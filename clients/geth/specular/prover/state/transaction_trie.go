// Copyright 2022, Specular contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

type TransactionTrie struct {
	transactions prover_types.Transactions
	trie         *trie.Trie
}

func NewTransactionTrie(txs prover_types.Transactions) *TransactionTrie {
	t := trie.NewEmpty(trie.NewDatabase(memorydb.New()))
	valueBuf := new(bytes.Buffer)
	var indexBuf []byte
	for i := 0; i < txs.Len(); i++ {
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		valueBuf.Reset()
		txs.EncodeIndex(i, valueBuf)
		t.Update(indexBuf, common.CopyBytes(valueBuf.Bytes()))
	}
	return &TransactionTrie{
		transactions: txs,
		trie:         t,
	}
}

func (t *TransactionTrie) Root() common.Hash {
	return t.trie.Hash()
}

func (t *TransactionTrie) EncodeState() []byte {
	return t.Root().Bytes()
}

func (t *TransactionTrie) GetTransaction(index int) prover_types.Transaction {
	return t.transactions[index]
}

func (t *TransactionTrie) Prove(index int) ([][]byte, error) {
	var proof mptProofList
	var indexBuf []byte
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(index))
	err := t.trie.Prove(indexBuf, 0, &proof)
	return proof, err
}
