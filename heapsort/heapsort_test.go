package heapsort_test

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"testing"

	"github.com/jian-maker-ap/hello-github/heapsort"
)

func TestSortInts(t *testing.T) {
	tests := []struct {
		name string
		in   []int
	}{
		{"nil", nil},
		{"empty", []int{}},
		{"single", []int{1}},
		{"two ordered", []int{1, 2}},
		{"two reversed", []int{2, 1}},
		{"already sorted", []int{-5, -1, 0, 2, 3, 7}},
		{"reverse sorted", []int{7, 3, 2, 0, -1, -5}},
		{"duplicates", []int{4, 2, 4, 1, 2, 4, 1}},
		{"all equal", []int{5, 5, 5, 5, 5, 5}},
		{"negative numbers", []int{-3, -10, 0, -1, 5, 2}},
		{"extremes", []int{math.MaxInt, 0, math.MinInt, -1, 1}},
		{"many duplicates", []int{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := slices.Clone(tt.in)
			want := slices.Clone(tt.in)
			slices.Sort(want)

			heapsort.Sort(tt.in)

			if !slices.Equal(tt.in, want) {
				t.Fatalf("Sort(%v) = %v, want %v", orig, tt.in, want)
			}
		})
	}
}

func TestSortStrings(t *testing.T) {
	in := []string{"pear", "apple", "banana", "", "Apple", "apple"}
	orig := slices.Clone(in)
	want := slices.Clone(in)
	slices.Sort(want)

	heapsort.Sort(in)

	if !slices.Equal(in, want) {
		t.Fatalf("Sort(%v) = %v, want %v", orig, in, want)
	}
}

func TestSortFloats(t *testing.T) {
	in := []float64{3.5, -1.25, 0, 2.5, -1.25, 10.125}
	orig := slices.Clone(in)
	want := slices.Clone(in)
	slices.Sort(want)

	heapsort.Sort(in)

	if !slices.Equal(in, want) {
		t.Fatalf("Sort(%v) = %v, want %v", orig, in, want)
	}
}

func TestSortInPlace(t *testing.T) {
	s := []int{3, 1, 2}
	alias := s // same backing array

	heapsort.Sort(s)

	if !slices.Equal(alias, []int{1, 2, 3}) {
		t.Fatalf("Sort did not sort in place: alias = %v", alias)
	}
}

func TestSortLargeSorted(t *testing.T) {
	s := make([]int, 100000)
	for i := range s {
		s[i] = i
	}

	heapsort.Sort(s)

	if !slices.IsSorted(s) {
		t.Fatal("already-sorted input was not sorted")
	}
}

func TestSortLargeReversed(t *testing.T) {
	s := make([]int, 100000)
	for i := range s {
		s[i] = len(s) - i
	}

	heapsort.Sort(s)

	if !slices.IsSorted(s) {
		t.Fatal("reversed input was not sorted")
	}
}

func TestSortOrganPipe(t *testing.T) {
	const n = 5000
	s := make([]int, 0, 2*n)
	for i := 1; i <= n; i++ {
		s = append(s, i)
	}
	for i := n; i >= 1; i-- {
		s = append(s, i)
	}

	heapsort.Sort(s)

	if !slices.IsSorted(s) {
		t.Fatal("organ-pipe input was not sorted")
	}
}

func TestSortRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 200; i++ {
		n := rng.Intn(1000)
		in := make([]int, n)
		for j := range in {
			in[j] = rng.Intn(100) - 50 // a small range produces duplicates
		}
		orig := slices.Clone(in)
		want := slices.Clone(in)
		slices.Sort(want)

		heapsort.Sort(in)

		if !slices.Equal(in, want) {
			t.Fatalf("case %d: Sort(%v) = %v, want %v", i, orig, in, want)
		}
	}
}

func ExampleSort() {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6}
	heapsort.Sort(s)
	fmt.Println(s)
	// Output: [1 1 2 3 4 5 6 9]
}

func FuzzSort(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{2, 1})
	f.Add([]byte{4, 2, 4, 1, 2, 4, 1})
	f.Fuzz(func(t *testing.T, in []byte) {
		got := slices.Clone(in)
		heapsort.Sort(got)

		want := slices.Clone(in)
		slices.Sort(want)

		if !slices.Equal(got, want) {
			t.Fatalf("Sort(%v) = %v, want %v", in, got, want)
		}
	})
}
