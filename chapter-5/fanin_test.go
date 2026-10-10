package ch05

import (
	"runtime"
	"testing"
	"time"
)

func TestFirstReturnsTheFastestAnswer(t *testing.T) {
	sources := []string{"alpha", "be", "gamma-source"}
	delays := []time.Duration{
		80 * time.Millisecond,
		10 * time.Millisecond,
		120 * time.Millisecond,
	}

	got := First(sources, delays)
	if got.Source != "be" {
		t.Fatalf("First returned %q, want %q (the fastest)",
			got.Source, "be")
	}
}

// The test the naive version cannot pass. First returns the right
// answer either way; what changes is whether the losers can finish.
func TestFirstLeavesNoGoroutinesBehind(t *testing.T) {
	sources := []string{"a", "bb", "ccc", "dddd", "eeeee"}
	delays := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}

	settle()
	before := runtime.NumGoroutine()

	First(sources, delays)

	// Well past the slowest source: any goroutine still alive here is
	// parked on a send, not merely slow.
	time.Sleep(200 * time.Millisecond)
	settle()
	leaked := runtime.NumGoroutine() - before

	if leaked > 0 {
		t.Fatalf("%d goroutines leaked (want 0).\nWith %d senders and "+
			"one receiver, the %d that lose the race are still parked "+
			"on a send. See §5.1's N-1 rule.",
			leaked, len(sources), len(sources)-1)
	}
}

func settle() {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		time.Sleep(5 * time.Millisecond)
	}
}
