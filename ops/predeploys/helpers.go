package predeploys

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// codeNamespace represents the namespace of implementations of predeploys
	codeNamespace = common.HexToAddress("0xc0D3C0d3C0d3C0D3c0d3C0d3c0D3C0d3c0d30000")
	// l2PredeployNamespace represents the namespace of L2 predeploys
	l2PredeployNamespace = common.HexToAddress("0x2A00000000000000000000000000000000000000")
	// bigL2PredeployNamespace represents the predeploy namespace as a big.Int
	BigL2PredeployNamespace = new(big.Int).SetBytes(l2PredeployNamespace.Bytes())
	// bigCodeNamespace represents the predeploy namespace as a big.Int
	bigCodeNameSpace = new(big.Int).SetBytes(codeNamespace.Bytes())
	// implementationSlot represents the EIP 1967 implementation storage slot
	ImplementationSlot = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")
	// adminSlot represents the EIP 1967 admin storage slot
	AdminSlot = common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103")
	// predeployProxyCount represents the number of predeploy proxies in the namespace
	PredeployProxyCount uint64 = 2048
)

// AddressToCodeNamespace takes a predeploy address and computes
// the implementation address that the implementation should be deployed at
func AddressToCodeNamespace(addr common.Address) (common.Address, error) {
	if !IsL2Predeploy(addr) {
		return common.Address{}, fmt.Errorf("cannot handle non predeploy: %s", addr)
	}
	bigAddress := new(big.Int).SetBytes(addr[18:])
	num := new(big.Int).Or(bigCodeNameSpace, bigAddress)
	return common.BigToAddress(num), nil
}

func IsL2Predeploy(addr common.Address) bool {
	return bytes.Equal(addr[0:2], []byte{0x2A, 0x00})
}
