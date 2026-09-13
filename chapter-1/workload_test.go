package chapter1

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	urlCount = 10
	latency  = 50 * time.Millisecond
)

// A server that is slow on purpose. Provided for you.
func slowServer(t *testing.T) *httptest.Server {
	t.Helper()
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(latency)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(s.Close)
	return s
}

// n copies of the same URL. Provided for you.
func urls(base string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = base
	}
	return out
}

// The waits overlap, so N requests should cost about as much
// as one — and this holds on a single core, because nothing
// here is CPU-bound.
func TestConcurrentOverlapsWaiting(t *testing.T) {
	s := slowServer(t)
	u := urls(s.URL, urlCount)
	start := time.Now()
	FetchConcurrent(u)
	elapsed := time.Since(start)
	// Sequential would be ~500ms. Anything under 150ms means the
	// waits are genuinely overlapping, not merely faster.
	budget := latency * 3
	if elapsed > budget {
		t.Fatalf("FetchConcurrent took %v for %d URLs of %v "+
			"latency; want < %v.\nSequential would be ~%v. "+
			"The waits are not overlapping yet.",
			elapsed.Round(time.Millisecond), urlCount, latency,
			budget, latency*urlCount)
	}
}
