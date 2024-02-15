// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const ISequencerInboxStorageLayoutJSON = "{\"storage\":null,\"types\":{}}"

var ISequencerInboxStorageLayout = new(solc.StorageLayout)

var ISequencerInboxDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(ISequencerInboxStorageLayoutJSON), ISequencerInboxStorageLayout); err != nil {
		panic(err)
	}

	layouts["ISequencerInbox"] = ISequencerInboxStorageLayout
	deployedBytecodes["ISequencerInbox"] = ISequencerInboxDeployedBin
}
