package test

import (
	"testing"

	"github.com/specularl2/specular/clients/geth/specular/utils"
	"github.com/stretchr/testify/assert"
)

func TestIndexOfLessThanOrEq(t *testing.T) {
	type inputs struct {
		sorted []int
		target int
	}
	type output struct{ index int }
	var testCases = map[string]testCase[inputs, output]{
		"empty slice":                       {inputs: inputs{[]int{}, 1}, expected: output{-1}},
		"exact match":                       {inputs: inputs{[]int{1, 2, 3, 4, 5}, 3}, expected: output{2}},
		"target greater than some => match": {inputs: inputs{[]int{1, 3, 5, 7, 9}, 6}, expected: output{2}},
		"target greater than all => match":  {inputs: inputs{[]int{5, 6, 7, 8}, 9}, expected: output{3}},
		"target less than all => no match":  {inputs: inputs{[]int{5, 6, 7, 8}, 3}, expected: output{-1}},
	}
	for name, tc := range testCases {
		t.Run(
			name,
			func(t *testing.T) {
				actual := utils.IndexOfLEq(tc.inputs.sorted, tc.inputs.target)
				assert.Equal(t, tc.expected.index, actual)
			},
		)
	}
}
