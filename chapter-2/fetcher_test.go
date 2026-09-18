package ch02

import (
	"context"
	"testing"
	"time"

	"go.uber.org/goleak"
)

// Fails the run if any goroutine outlives the test binary.
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestFetchFirstDoesNotLeak(t *testing.T) {
	urls := []string{
		"https://go.dev",
		"https://pkg.go.dev",
		"https://go.dev/blog",
		"https://go.dev/doc",
	}

	start := time.Now()
	got, err := FetchFirst(context.Background(), urls)
	t.Logf("FetchFirst took %v", time.Since(start))
	if err != nil {
		t.Fatalf("FetchFirst returned %v", err)
	}
	if got.URL == "" {
		t.Fatal("FetchFirst returned an empty result")
	}
	// The losing goroutines are blocked on `results <- r` right now.
	// VerifyTestMain catches them when the binary exits.
}
