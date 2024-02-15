// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/specularL2/specular/bindings-go/solc"
)

const L2PortalDeterministicStorageStorageLayoutJSON = "{\"storage\":[{\"astId\":1000,\"contract\":\"src/bridge/L2Portal.sol:L2PortalDeterministicStorage\",\"label\":\"initiatedWithdrawals\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_mapping(t_bytes32,t_bool)\"}],\"types\":{\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_bytes32\":{\"encoding\":\"inplace\",\"label\":\"bytes32\",\"numberOfBytes\":\"32\"},\"t_mapping(t_bytes32,t_bool)\":{\"encoding\":\"mapping\",\"label\":\"mapping(bytes32 =\u003e bool)\",\"numberOfBytes\":\"32\",\"key\":\"t_bytes32\",\"value\":\"t_bool\"}}}"

var L2PortalDeterministicStorageStorageLayout = new(solc.StorageLayout)

var L2PortalDeterministicStorageDeployedBin = "0x"
func init() {
	if err := json.Unmarshal([]byte(L2PortalDeterministicStorageStorageLayoutJSON), L2PortalDeterministicStorageStorageLayout); err != nil {
		panic(err)
	}

	layouts["L2PortalDeterministicStorage"] = L2PortalDeterministicStorageStorageLayout
	deployedBytecodes["L2PortalDeterministicStorage"] = L2PortalDeterministicStorageDeployedBin
}
