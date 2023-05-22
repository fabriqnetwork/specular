// Copyright 2022, Specular contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package merkletree

import (
	"testing"
)

func TestRoundUpToPowerOf2(t *testing.T) {
	testCases := []struct {
		input uint64
		want  uint64
	}{
		{0, 1},
		{1, 1},
		{5, 8},
		{8, 8},
	}
	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			got := roundUpToPowerOf2(testCase.input)
			if got != testCase.want {
				t.Errorf("got %d, want %d", got, testCase.want)
			}
		})
	}
}
