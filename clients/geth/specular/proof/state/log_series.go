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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type LogSeries struct {
	Logs  []*types.Log
	Bloom types.Bloom
}

func EmptyLogSeries() *LogSeries {
	return &LogSeries{
		Logs:  make([]*types.Log, 0),
		Bloom: types.Bloom{},
	}
}

func (l *LogSeries) Add(log *types.Log) *LogSeries {
	logs := make([]*types.Log, len(l.Logs)+1)
	copy(logs, l.Logs)
	var bin types.Bloom
	bin.SetBytes(types.LogsBloom(l.Logs))
	return &LogSeries{
		Logs:  append(logs, log),
		Bloom: bin,
	}
}

func (l *LogSeries) Hash() common.Hash {
	return common.BytesToHash(l.EncodeState())
}

func (l *LogSeries) EncodeState() []byte {
	// TODO: should we check rlp encode error here?
	logBytes, _ := rlp.EncodeToBytes(l)
	return crypto.Keccak256(logBytes, l.Bloom.Bytes())
}
