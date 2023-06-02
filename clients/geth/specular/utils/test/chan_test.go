package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/utils"
	"github.com/stretchr/testify/assert"
)

type testCase[T, U any] struct {
	inputs   T
	expected U
}

func TestSubscribeBatched(t *testing.T) {
	type inputs struct {
		minBatchInterval time.Duration
		maxBatchInterval time.Duration
	}
	type output struct{ index int }

	var testCases = map[string]testCase[inputs, output]{
		"no max": {inputs: inputs{10 * time.Millisecond, 0}, expected: output{2}},
		// "no min; 1s max": {inputs: inputs{0, 1}, expected: output{2}},
	}
	for name, tc := range testCases {
		t.Run(
			name,
			func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				inCh := make(chan int)
				fmt.Println("sub")
				outCh := utils.SubscribeBatched(ctx, inCh, tc.inputs.minBatchInterval, tc.inputs.maxBatchInterval)
				fmt.Println("pub")
				// Publish
				sl := [5]int{}
				for i := 0; i < len(sl); i++ {
					sl[i] = i
					inCh <- sl[i]
				}
				fmt.Println("slep")
				time.Sleep(tc.inputs.minBatchInterval)
				fmt.Println("sub")
				batch := <-outCh
				assert.True(t, utils.Equal(sl[:], batch))
				cancel()
			},
		)
	}
}
