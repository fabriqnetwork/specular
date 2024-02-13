// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const IRollupStorageLayoutJSON = "{\"storage\":null,\"types\":{}}"

var IRollupStorageLayout = new(solc.StorageLayout)

var IRollupDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(IRollupStorageLayoutJSON), IRollupStorageLayout); err != nil {
		panic(err)
	}

	layouts["IRollup"] = IRollupStorageLayout
	deployedBytecodes["IRollup"] = IRollupDeployedBin
}
