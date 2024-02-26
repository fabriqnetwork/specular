// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const IAsymChallengeStorageLayoutJSON = "{\"storage\":null,\"types\":{}}"

var IAsymChallengeStorageLayout = new(solc.StorageLayout)

var IAsymChallengeDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(IAsymChallengeStorageLayoutJSON), IAsymChallengeStorageLayout); err != nil {
		panic(err)
	}

	layouts["IAsymChallenge"] = IAsymChallengeStorageLayout
	deployedBytecodes["IAsymChallenge"] = IAsymChallengeDeployedBin
}
