package predeploys

import "github.com/ethereum/go-ethereum/common"

// TODO: This list should be generated from a configuration source
const (
	UUPSPlaceholder      = "0x2A00000000000000000000000000000000000000"
	L1Oracle             = "0x2A00000000000000000000000000000000000010"
	L2Portal             = "0x2A00000000000000000000000000000000000011"
	L2StandardBridge     = "0x2A00000000000000000000000000000000000012"
	L1FeeVault           = "0x2A00000000000000000000000000000000000020"
	L2BaseFeeVault       = "0x2A00000000000000000000000000000000000021"
	MintableERC20Factory = "0x2A000000000000000000000000000000000000f0"
)

var (
	UUPSPlaceholderAddr      = common.HexToAddress(UUPSPlaceholder)
	L1OracleAddr             = common.HexToAddress(L1Oracle)
	L2PortalAddr             = common.HexToAddress(L2Portal)
	L2StandardBridgeAddr     = common.HexToAddress(L2StandardBridge)
	L1FeeVaultAddr           = common.HexToAddress(L1FeeVault)
	L2BaseFeeVaultAddr       = common.HexToAddress(L2BaseFeeVault)
	MintableERC20FactoryAddr = common.HexToAddress(MintableERC20Factory)

	Predeploys = make(map[string]*common.Address)
)

func init() {
	Predeploys["UUPSPlaceholder"] = &UUPSPlaceholderAddr
	Predeploys["L1Oracle"] = &L1OracleAddr
	Predeploys["L2Portal"] = &L2PortalAddr
	Predeploys["L2StandardBridge"] = &L2StandardBridgeAddr
	Predeploys["L1FeeVault"] = &L1FeeVaultAddr
	Predeploys["L2BaseFeeVault"] = &L2BaseFeeVaultAddr
	Predeploys["MintableERC20Factory"] = &MintableERC20FactoryAddr
}
