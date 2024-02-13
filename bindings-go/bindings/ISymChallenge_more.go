// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const ISymChallengeStorageLayoutJSON = "{\"storage\":null,\"types\":{}}"

var ISymChallengeStorageLayout = new(solc.StorageLayout)

var ISymChallengeDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(ISymChallengeStorageLayoutJSON), ISymChallengeStorageLayout); err != nil {
		panic(err)
	}

	layouts["ISymChallenge"] = ISymChallengeStorageLayout
	deployedBytecodes["ISymChallenge"] = ISymChallengeDeployedBin
}
