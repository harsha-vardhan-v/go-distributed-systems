package ch04

import (
	"slices"
	"testing"
	"time"
)

// collect drains out until it closes. A correct Merge closes; one that
// leaves a closed channel enabled spins on it, so we stop early and say
// so rather than printing a thousand zero values.
func collect(t *testing.T, out <-chan int) []int {
	t.Helper()
	type result struct {
		got     []int
		runaway bool
	}
	done := make(chan result, 1)

	go func() {
		var got []int
		for v := range out {
			got = append(got, v)
			if len(got) > 200 {
				done <- result{got[:8], true}
				return
			}
		}
		done <- result{got, false}
	}()

	select {
	case r := <-done:
		if r.runaway {
			t.Fatalf("Merge produced 200+ values and kept going; the "+
				"first few were %v.\nA closed channel is permanently "+
				"ready and yields the zero value every time, so the "+
				"select is spinning on it. See §4.6.", r.got)
		}
		return r.got
	case <-time.After(2 * time.Second):
		t.Fatal("Merge never closed out. See §4.6.")
		return nil
	}
}

func TestMergeForwardsEverythingThenCloses(t *testing.T) {
	a, b := make(chan int), make(chan int)
	go func() {
		for _, v := range []int{1, 2, 3} {
			a <- v
		}
		close(a)
	}()
	go func() {
		for _, v := range []int{10, 20} {
			b <- v
		}
		close(b)
	}()

	got := collect(t, Merge(a, b))
	slices.Sort(got)
	want := []int{1, 2, 3, 10, 20}
	if !slices.Equal(got, want) {
		t.Fatalf("Merge gave %v, want %v", got, want)
	}
}

// The case that most cleanly separates a real fix from a lucky one: one
// input is already closed before Merge ever looks at it.
func TestMergeHandlesAnAlreadyClosedInput(t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	close(a)
	go func() {
		b <- 7
		close(b)
	}()

	got := collect(t, Merge(a, b))
	if !slices.Equal(got, []int{7}) {
		t.Fatalf("Merge gave %v, want [7]", got)
	}
}
