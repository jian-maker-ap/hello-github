// Package quicksort implements an in-place quicksort for ordered types.
package quicksort

import "cmp"

// insertionSortCutoff is the largest slice that Sort sorts with insertion
// sort instead of further partitioning. Insertion sort is faster than
// quicksort for very small inputs and avoids deep recursion on them.
const insertionSortCutoff = 12

// Sort sorts s in place in ascending order.
//
// Sort is a quicksort with median-of-three pivot selection and three-way
// partitioning, so runs of equal elements keep the algorithm linear in their
// size. It compares elements with < and therefore works for any
// cmp.Ordered type. Like every quicksort its average cost is O(n log n)
// comparisons and its worst case O(n²) for adversarial inputs; common inputs,
// including already-sorted ones with distinct elements, are close to
// O(n log n). It sorts in place, allocates no heap memory and uses at most
// O(log n) stack space.
func Sort[T cmp.Ordered](s []T) {
	for len(s) > insertionSortCutoff {
		// Median-of-three pivot selection: order the first, middle and
		// last elements so that the middle one is a better pivot than an
		// arbitrary choice.
		mid := len(s) / 2
		last := len(s) - 1
		if s[mid] < s[0] {
			s[0], s[mid] = s[mid], s[0]
		}
		if s[last] < s[0] {
			s[0], s[last] = s[last], s[0]
		}
		if s[last] < s[mid] {
			s[mid], s[last] = s[last], s[mid]
		}

		lt, gt := partition(s, s[mid])

		// Recurse into the smaller side and loop on the larger one, which
		// bounds the recursion depth to O(log n) even for bad pivots.
		if lt < len(s)-gt {
			Sort(s[:lt])
			s = s[gt:]
		} else {
			Sort(s[gt:])
			s = s[:lt]
		}
	}
	insertionSort(s)
}

// partition partitions s into elements smaller than pivot, elements equal to
// pivot, and elements greater than pivot. It returns the bounds of the equal
// middle part: s[:lt] < pivot, s[lt:gt] == pivot, and s[gt:] > pivot.
func partition[T cmp.Ordered](s []T, pivot T) (lt, gt int) {
	lt, gt = 0, len(s)
	for i := 0; i < gt; {
		switch {
		case s[i] < pivot:
			s[i], s[lt] = s[lt], s[i]
			i++
			lt++
		case pivot < s[i]:
			gt--
			s[i], s[gt] = s[gt], s[i]
		default: // s[i] == pivot
			i++
		}
	}
	return lt, gt
}

// insertionSort sorts s in place with insertion sort, which is efficient for
// the short slices produced by partitioning.
func insertionSort[T cmp.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
