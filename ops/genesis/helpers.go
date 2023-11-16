package genesis

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func newHexBig(in uint64) *hexutil.Big {
	b := new(big.Int).SetUint64(in)
	hb := hexutil.Big(*b)
	return &hb
}
