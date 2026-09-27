package ch03

import (
	"slices"
	"testing"
	"time"
)

func TestMergeCollectsEveryValue(t *testing.T) {
	got := Merge([]int{3, 4, 5})

	want := []int{0, 1, 2, 100, 101, 102, 103, 200, 201, 202, 203, 204}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Fatalf("Merge returned %v, want %v", got, want)
	}
}

// The edge case that separates a real fix from a lucky one. With no
// producers, a close that lives inside a producer never runs at all,
// so the `for range` waits forever. A coordinator that waits on the
// WaitGroup closes immediately and the range ends at once.
//
// Merge runs on its own goroutine so a wrong answer fails in two
// seconds with a message instead of hanging until the test timeout.
func TestMergeTerminatesWithNoProducers(t *testing.T) {
	done := make(chan []int, 1)
	go func() { done <- Merge(nil) }()

	select {
	case got := <-done:
		if len(got) != 0 {
			t.Fatalf("Merge(nil) gave %d values, want 0", len(got))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Merge(nil) never returned: with no producers " +
			"nothing closed the channel, so for range still waits")
	}

}
