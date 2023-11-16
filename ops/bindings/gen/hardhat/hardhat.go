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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/specularL2/specular/ops/bindings/solc"
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

// DebugFile represents the debug file that contains the path
// to the build info file
type DebugFile struct {
	Format    string `json:"_format"`
	BuildInfo string `json:"buildInfo"`
}

// BuildInfo represents a hardhat build info artifact that is created
// after compilation
type BuildInfo struct {
	Format          string              `json:"_format"`
	Id              string              `json:"id"`
	SolcVersion     string              `json:"solcVersion"`
	SolcLongVersion string              `json:"solcLongVersion"`
	Input           solc.CompilerInput  `json:"input"`
	Output          solc.CompilerOutput `json:"output"`
}

// ReadBuildInfos will read all the build info files in the artifact path
// and return a map of contract name to build info
func ReadBuildInfos(artifactPath string) (map[string]*BuildInfo, error) {
	pathToBuildInfo := make(map[string]*BuildInfo)
	targetToBuildInfo := make(map[string]*BuildInfo)

	fileSystem := os.DirFS(artifactPath)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := filepath.Join(artifactPath, path)

		if !strings.HasSuffix(name, ".dbg.json") {
			return nil
		}

		// Remove ".dbg.json"
		target := filepath.Base(name[:len(name)-9])

		file, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		var debugFile DebugFile
		if err := json.Unmarshal(file, &debugFile); err != nil {
			return err
		}
		relPath := filepath.Join(filepath.Dir(name), debugFile.BuildInfo)
		if err != nil {
			return err
		}
		debugPath, _ := filepath.Abs(relPath)

		// If we have already read the build info file, we can just use it
		if buildInfo, ok := pathToBuildInfo[debugPath]; ok {
			targetToBuildInfo[target] = buildInfo
			return nil
		}

		buildInfoFile, err := os.ReadFile(debugPath)
		if err != nil {
			return err
		}

		var buildInfo BuildInfo
		if err := json.Unmarshal(buildInfoFile, &buildInfo); err != nil {
			return err
		}

		pathToBuildInfo[debugPath] = &buildInfo
		targetToBuildInfo[target] = &buildInfo

		return nil
	})
	if err != nil {
		return nil, err
	}

	return targetToBuildInfo, nil
}

func GetStorageLayout(contractName string, buildInfo *BuildInfo) (*solc.StorageLayout, error) {
	for _, source := range buildInfo.Output.Contracts {
		for name, contract := range source {
			if name == contractName {
				return &contract.StorageLayout, nil
			}
		}
	}

	return nil, fmt.Errorf("contract not found for %s", contractName)
}
