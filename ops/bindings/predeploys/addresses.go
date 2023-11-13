package predeploys

import "github.com/ethereum/go-ethereum/common"

// TODO: This list should be generated from a configuration source
const (
	L1Oracle         = "0x2A00000000000000000000000000000000000010"
	L2Portal         = "0x2A00000000000000000000000000000000000011"
	L2StandardBridge = "0x2A00000000000000000000000000000000000012"
)

var (
	L1OracleAddr         = common.HexToAddress(L1Oracle)
	L2PortalAddr         = common.HexToAddress(L2Portal)
	L2StandardBridgeAddr = common.HexToAddress(L2StandardBridge)

	Predeploys = make(map[string]*common.Address)
)

// IsProxied returns true for predeploys that will sit behind a proxy contract
func IsProxied(predeployAddr common.Address) bool {
	// Currently all predeploys are proxied
	return true
}

func init() {
	Predeploys["L1Oracle"] = &L1OracleAddr
	Predeploys["L2Portal"] = &L2PortalAddr
	Predeploys["L2StandardBridge"] = &L2StandardBridgeAddr
}
