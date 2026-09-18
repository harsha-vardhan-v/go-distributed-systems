package ch02

import (
	"context"
	"time"
)

type Result struct {
	URL  string
	Size int
}

// Provided for you: a slow fetch that honours cancellation.
func fetchOne(ctx context.Context, url string) (Result, error) {
	select {
	case <-time.After(50 * time.Millisecond):
		return Result{URL: url, Size: len(url) * 10}, nil
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}
}

// TODO(reader): this leaks. Every goroutine that loses the race is left
// holding a value for a channel nobody will ever read again. Fix it
// WITHOUT waiting for the slow fetches — the caller wants the first
// answer, fast.
func FetchFirst(ctx context.Context, urls []string) (Result, error) {
	results := make(chan Result, len(urls))

	for _, url := range urls {
		go func(url string) {
			r, err := fetchOne(ctx, url)
			if err != nil {
				return
			}
			results <- r // <- your move
		}(url)
	}

	return <-results, nil
}
