// Copyright 2023, Specular contributors
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

package hardhat

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Artifact struct {
	ContractName     string        `json:"contractName"`
	SourceName       string        `json:"sourceName"`
	Abi              interface{}   `json:"abi"`
	Bytecode         hexutil.Bytes `json:"bytecode"`
	DeployedBytecode hexutil.Bytes `json:"deployedBytecode"`
}

func ReadArtifact(path string) (*Artifact, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var artifact Artifact
	if err := json.Unmarshal(file, &artifact); err != nil {
		return nil, err
	}
	return &artifact, nil
}

func ReadStorageLayoutFromCacheValidations(path string) (map[string][]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cacheValidationFile map[string]interface{}
	if err := json.Unmarshal(file, &cacheValidationFile); err != nil {
		return nil, err
	}
	layout := make(map[string][]byte)
	logs := cacheValidationFile["log"].([]interface{})
	for _, log := range logs {
		log := log.(map[string]interface{})
		for contract, data := range log {
			if layout[contract] != nil {
				continue
			}
			layout[contract], err = json.Marshal(data.(map[string]interface{})["layout"])
			if err != nil {
				return nil, err
			}
		}
	}
	return layout, nil
}
