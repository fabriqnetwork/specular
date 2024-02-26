// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const IChallengeResultReceiverStorageLayoutJSON = "{\"storage\":null,\"types\":{}}"

var IChallengeResultReceiverStorageLayout = new(solc.StorageLayout)

var IChallengeResultReceiverDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(IChallengeResultReceiverStorageLayoutJSON), IChallengeResultReceiverStorageLayout); err != nil {
		panic(err)
	}

	layouts["IChallengeResultReceiver"] = IChallengeResultReceiverStorageLayout
	deployedBytecodes["IChallengeResultReceiver"] = IChallengeResultReceiverDeployedBin
}
