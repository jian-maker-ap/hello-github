// Package heapsort implements an in-place heapsort for ordered types.
package heapsort

import "cmp"

// Sort sorts s in place in ascending order.
//
// Sort builds a max-heap over s in O(n) comparisons, then repeatedly moves
// the heap maximum to the end of the unsorted range and sifts it back down.
// It takes O(n log n) comparisons in the worst case, so unlike quicksort its
// running time has no quadratic worst case. It is not stable, sorts in place,
// allocates no heap memory and uses O(1) auxiliary space.
func Sort[T cmp.Ordered](s []T) {
	n := len(s)

	// Build a max-heap bottom-up. Every element past the first parent is a
	// leaf and therefore already a valid heap.
	for root := n/2 - 1; root >= 0; root-- {
		siftDown(s, root, n)
	}

	// Repeatedly swap the maximum to the end of the heap and shrink it.
	for n > 1 {
		n--
		s[0], s[n] = s[n], s[0]
		siftDown(s, 0, n)
	}
}

// siftDown restores the max-heap property of the subtree rooted at root
// within the heap s[:n].
func siftDown[T cmp.Ordered](s []T, root, n int) {
	for {
		left := 2*root + 1
		if left >= n {
			return
		}
		largest := left
		if right := left + 1; right < n && s[left] < s[right] {
			largest = right
		}
		if !(s[root] < s[largest]) {
			return
		}
		s[root], s[largest] = s[largest], s[root]
		root = largest
	}
}
