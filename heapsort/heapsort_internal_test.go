package heapsort

import (
	"slices"
	"testing"
)

// isMaxHeap reports whether s satisfies the max-heap property: no element is
// greater than its parent.
func isMaxHeap(s []int) bool {
	for i := range s {
		if left := 2*i + 1; left < len(s) && s[i] < s[left] {
			return false
		}
		if right := 2*i + 2; right < len(s) && s[i] < s[right] {
			return false
		}
	}
	return true
}

func TestSiftDown(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		root int
		n    int
	}{
		{"empty", nil, 0, 0},
		{"single", []int{1}, 0, 1},
		{"already heap", []int{9, 8, 7}, 0, 3},
		{"root smaller than left", []int{1, 8, 7}, 0, 3},
		{"root smaller than right", []int{1, 7, 8}, 0, 3},
		{"root smaller than both", []int{1, 8, 9}, 0, 3},
		{"inner node", []int{9, 1, 7, 2, 3}, 1, 5},
		{"deep violation", []int{1, 8, 7, 6, 5, 4, 3}, 0, 7},
		{"tied values", []int{1, 1, 1, 1}, 0, 4},
		{"limited heap", []int{1, 9, 9, 9}, 0, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := slices.Clone(tt.in)
			siftDown(s, tt.root, tt.n)

			if !isMaxHeap(s[:tt.n]) {
				t.Fatalf("siftDown(%v, %d, %d) = %v, want a max-heap", tt.in, tt.root, tt.n, s)
			}

			got := slices.Clone(s)
			slices.Sort(got)
			want := slices.Clone(tt.in)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Fatalf("siftDown changed the elements: got %v, want %v", got, want)
			}
		})
	}
}
