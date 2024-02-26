// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const IChallengeStorageLayoutJSON = "{\"storage\":null,\"types\":{}}"

var IChallengeStorageLayout = new(solc.StorageLayout)

var IChallengeDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(IChallengeStorageLayoutJSON), IChallengeStorageLayout); err != nil {
		panic(err)
	}

	layouts["IChallenge"] = IChallengeStorageLayout
	deployedBytecodes["IChallenge"] = IChallengeDeployedBin
}
