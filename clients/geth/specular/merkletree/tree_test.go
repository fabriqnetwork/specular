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

	"github.com/holiman/uint256"
)

func makeTree() (*MerkleTree, []*uint256.Int) {
	elements := make([]*uint256.Int, 3)
	for i := 0; i < 3; i++ {
		elements[i] = uint256.NewInt(uint64(i + 1))
	}
	tree := New(elements)
	return tree, elements
}

func TestMerkleTreeConstruction(t *testing.T) {
	tree, elements := makeTree()
	if tree.ElementCount() != 4 {
		t.Errorf("got %d, want %d", tree.ElementCount(), 4)
	}
	for i := uint64(0); i < 3; i++ {
		if tree.GetElement(i) != elements[i] {
			t.Errorf("got %d, want %d", tree.GetElement(i), elements[i])
		}
	}
	if *tree.GetElement(3) != (uint256.Int{}) {
		t.Errorf("got %d, want %d", tree.GetElement(3), uint256.Int{})
	}
	for i := uint64(1); i < 4; i++ {
		expected := hashNode(tree.tree[i<<1], tree.tree[i<<1+1])
		if tree.tree[i] != expected {
			t.Errorf("got %d, want %d", tree.tree[i], expected)
		}
	}
}

func TestProof(t *testing.T) {
	tree, elements := makeTree()
	root := tree.GetRoot()
	elementNum := tree.ElementCount()

	proof := tree.GetProof(0)
	result := VerifyProof(root, elementNum, 0, elements[0:1], proof)
	if !result {
		t.Error("Verification Failed for range 0->0")
	}

	proof = tree.GetRangeProof(0, 2)
	result = VerifyProof(root, elementNum, 0, elements[0:2], proof)
	if !result {
		t.Error("Verification Failed for range 0->1")
	}

	proof = tree.GetRangeProof(1, 2)
	result = VerifyProof(root, elementNum, 1, elements[1:3], proof)
	if !result {
		t.Error("Verification Failed for range 1->2")
	}
}

func makeLargeTree() (*MerkleTree, []*uint256.Int) {
	elements := make([]*uint256.Int, 100)
	for i := 0; i < 100; i++ {
		elements[i] = uint256.NewInt(uint64(i * i))
	}
	tree := New(elements)
	return tree, elements
}

func TestProofLarge(t *testing.T) {
	tree, elements := makeLargeTree()
	root := tree.GetRoot()
	elementNum := tree.ElementCount()

	proof := tree.GetRangeProof(30, 20)
	result := VerifyProof(root, elementNum, 30, elements[30:50], proof)
	if !result {
		t.Error("Verification Failed for range 30->49")
	}

	proof = tree.GetRangeProof(0, 100)
	result = VerifyProof(root, elementNum, 0, elements[0:100], proof)
	if !result {
		t.Error("Verification Failed for range 0->99")
	}
}

func makeHugeTree(b *testing.B) (*MerkleTree, []*uint256.Int) {
	b.StopTimer()
	elementNum := uint64(1<<15 + 1<<14)
	elements := make([]*uint256.Int, elementNum)
	for i := uint64(0); i < elementNum; i++ {
		elements[i] = uint256.NewInt(i)
	}
	b.StartTimer()
	tree := New(elements)
	return tree, elements
}

func BenchmarkMerkleTreeConstructionCold(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = makeHugeTree(b)
	}
}

func BenchmarkMerkleTreeConstructionWithRebuild(b *testing.B) {
	b.StopTimer()
	tree, elements := makeHugeTree(b)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		tree.Rebuild(elements)
	}
}

func BenchmarkMerkleTreeProof(b *testing.B) {
	b.StopTimer()
	tree, _ := makeHugeTree(b)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		tree.GetRangeProof(3*(1<<12), 1<<12)
	}
}

func BenchmarkMerkleTreeVerification(b *testing.B) {
	b.StopTimer()
	tree, elements := makeHugeTree(b)
	root := tree.GetRoot()
	elementNum := tree.ElementCount()
	proof := tree.GetRangeProof(3*(1<<12), 1<<12)
	proofElements := elements[3*(1<<12) : 4*(1<<12)]
	result := VerifyProof(root, elementNum, 3*(1<<12), proofElements, proof)
	if !result {
		b.Error("Verification Failed for range 30->49")
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		VerifyProof(root, elementNum, 3*(1<<12), proofElements, proof)
	}
}
