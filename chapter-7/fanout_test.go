package ch07

import (
	"slices"
	"testing"
	"time"
)

func feed(values ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

// A fan-out that never closes its output does not crash. It hangs, so
// the deadline is the test: without it this would sit here until the
// package timeout, with no useful failure message.
func TestFanOutClosesItsOutput(t *testing.T) {
	in := feed(1, 2, 3, 4, 5, 6, 7, 8)

	done := make(chan []int, 1)
	go func() {
		var got []int
		for v := range FanOut(in, 3, func(n int) int { return n * n }) {
			got = append(got, v)
		}
		done <- got
	}()

	select {
	case got := <-done:
		slices.Sort(got)
		want := []int{1, 4, 9, 16, 25, 36, 49, 64}
		if !slices.Equal(got, want) {
			t.Fatalf("FanOut gave %v, want %v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("FanOut never closed its output channel: the " +
			"consumer's range loop is still waiting. See §7.2.")
	}
}

// Fan-out divides work; it does not duplicate it. Every input must
// come out exactly once no matter how the workers interleave.
func TestFanOutDeliversEachValueOnce(t *testing.T) {
	const n = 200
	values := make([]int, n)
	for i := range values {
		values[i] = i
	}

	seen := make(map[int]int, n)
	done := make(chan map[int]int, 1)
	identity := func(n int) int { return n }
	go func() {
		for v := range FanOut(feed(values...), 8, identity) {
			seen[v]++
		}
		done <- seen
	}()

	select {
	case got := <-done:
		for _, v := range values {
			if got[v] != 1 {
				t.Fatalf("value %d came out %d times; want exactly 1",
					v, got[v])
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatal("FanOut never closed its output channel")
	}
}
