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

package proof

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/specularl2/specular/clients/geth/specular/proof/state"
)

func GetAccountWithProof(
	db vm.StateDB,
	address common.Address,
) (*state.Account, [][]byte, error) {
	balance, _ := uint256.FromBig(db.GetBalance(address))
	account := &state.Account{
		Nonce:       db.GetNonce(address),
		Balance:     *balance,
		StorageRoot: db.GetStateRootForProof(address),
		CodeHash:    db.GetCodeHash(address),
	}
	accountProof, err := db.GetProof(address)
	if err != nil {
		return nil, nil, err
	}
	return account, accountProof, nil
}

func GetStorageAtWithProof(
	db vm.StateDB,
	address common.Address,
	key common.Hash,
) (*uint256.Int, *state.Account, [][]byte, [][]byte, error) {
	balance, _ := uint256.FromBig(db.GetBalance(address))
	account := &state.Account{
		Nonce:       db.GetNonce(address),
		Balance:     *balance,
		StorageRoot: db.GetStateRootForProof(address),
		CodeHash:    db.GetCodeHash(address),
	}
	value := new(uint256.Int).SetBytes(db.GetState(address, key).Bytes())
	accountProof, err := db.GetProof(address)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	storageProof, err := db.GetStorageProof(address, key)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return value, account, accountProof, storageProof, nil
}
