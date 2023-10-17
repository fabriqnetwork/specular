package geth

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestScaleBigInt(t *testing.T) {
	largeResult, _ := new(big.Int).SetString("12912720851596686130", 10)

	var tests = []struct {
		num    *big.Int
		scalar float64
		want   *big.Int
	}{
		{big.NewInt(10), 1, big.NewInt(10)},
		{big.NewInt(1), 1.2, big.NewInt(2)},
		{big.NewInt(10), 9.999, big.NewInt(100)},
		{big.NewInt(math.MaxInt64), 1, big.NewInt(math.MaxInt64)},
		{big.NewInt(math.MaxInt64), 1.4, largeResult},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%f", tt.num.String(), tt.scalar)
		t.Run(testname, func(t *testing.T) {
			ans := ScaleBigInt(tt.num, tt.scalar)
			if ans.Cmp(tt.want) != 0 {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
