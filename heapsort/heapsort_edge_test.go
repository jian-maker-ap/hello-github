package heapsort_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/jian-maker-ap/hello-github/heapsort"
)

func TestSortSizes(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	for n := 0; n <= 64; n++ {
		for _, duplicates := range []bool{false, true} {
			in := make([]int, n)
			for i := range in {
				if duplicates {
					in[i] = rng.Intn(5) - 2
				} else {
					in[i] = rng.Intn(1_000_000) - 500_000
				}
			}
			orig := slices.Clone(in)
			want := slices.Clone(in)
			slices.Sort(want)

			heapsort.Sort(in)

			if !slices.Equal(in, want) {
				t.Fatalf("n=%d duplicates=%v: Sort(%v) = %v, want %v", n, duplicates, orig, in, want)
			}
		}
	}
}

func TestSortLargeDuplicates(t *testing.T) {
	s := make([]int, 100000)
	for i := range s {
		s[i] = i % 3
	}

	heapsort.Sort(s)

	if !slices.IsSorted(s) {
		t.Fatal("duplicate-heavy input was not sorted")
	}
}
