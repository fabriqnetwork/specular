package utils

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Returns the index of the last value in slice that is less than or equal to the target.
// Assumes slice is sorted in ascending order and strictly increasing.
func IndexOfLEq[T constraints.Ordered](sorted []T, target T) int {
	return IndexOfMappedLEq(sorted, target, func(in T) T { return in })
}

// Returns the index of the last value in slice (mapped) that is less than or equal to the target.
// Assumes mapped values of slice are in ascending order and strictly increasing.
func IndexOfMappedLEq[T any, U constraints.Ordered](sorted []T, target U, mapFn func(T) U) int {
	var (
		start  = 0
		end    = len(sorted) - 1
		mid    int
		midVal U
	)
	if end == -1 {
		return -1
	}
	// Do a binary search for exact match.
	for start <= end {
		mid = (start + end) / 2
		midVal = mapFn(sorted[mid])
		fmt.Println("mid", start, mid, end)
		if midVal > target {
			end = mid - 1
		} else if midVal < target {
			start = mid + 1
		} else {
			return mid
		}
	}
	// Couldn't find exact match. Return index of the last value less than the target.
	// We assume strictly increasing values, so mid - 1 is guaranteed to be less than target.
	if midVal < target {
		return mid
	} else {
		return mid - 1
	}
}
