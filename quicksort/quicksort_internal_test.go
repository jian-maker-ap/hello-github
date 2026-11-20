package quicksort

import (
	"slices"
	"testing"
)

func TestPartition(t *testing.T) {
	tests := []struct {
		name  string
		in    []int
		pivot int
	}{
		{"empty", nil, 0},
		{"single equal", []int{5}, 5},
		{"single less", []int{5}, 6},
		{"single greater", []int{5}, 4},
		{"pivot minimum", []int{3, 1, 2, 4}, 1},
		{"pivot maximum", []int{3, 1, 2, 4}, 4},
		{"pivot middle", []int{5, 1, 4, 2, 3}, 3},
		{"all equal", []int{7, 7, 7, 7}, 7},
		{"pivot absent", []int{4, 2, 6, 1}, 5},
		{"duplicates", []int{2, 1, 2, 3, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := slices.Clone(tt.in)
			lt, gt := partition(s, tt.pivot)

			if lt < 0 || gt < lt || gt > len(s) {
				t.Fatalf("partition returned out-of-range bounds lt=%d gt=%d for len(s)=%d", lt, gt, len(s))
			}
			for i, v := range s {
				switch {
				case i < lt:
					if !(v < tt.pivot) {
						t.Fatalf("s[%d] = %d in lower partition, want < %d (result %v)", i, v, tt.pivot, s)
					}
				case i < gt:
					if v != tt.pivot {
						t.Fatalf("s[%d] = %d in equal partition, want == %d (result %v)", i, v, tt.pivot, s)
					}
				default:
					if !(tt.pivot < v) {
						t.Fatalf("s[%d] = %d in upper partition, want > %d (result %v)", i, v, tt.pivot, s)
					}
				}
			}

			got := slices.Clone(s)
			slices.Sort(got)
			want := slices.Clone(tt.in)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Fatalf("partition changed the elements: got %v, want %v", got, want)
			}
		})
	}
}

func TestInsertionSort(t *testing.T) {
	tests := [][]int{
		nil,
		{},
		{1},
		{2, 1},
		{1, 2, 3, 4},
		{4, 3, 2, 1},
		{2, 3, 1, 5, 4},
		{4, 4, 2, 2, 1, 1},
		{-1, -3, 0, 2, -2},
	}
	for _, in := range tests {
		got := slices.Clone(in)
		insertionSort(got)

		want := slices.Clone(in)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("insertionSort(%v) = %v, want %v", in, got, want)
		}
	}
}
