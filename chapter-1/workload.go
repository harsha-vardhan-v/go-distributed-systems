package chapter1

import (
	"io"
	"net/http"
	"sync"
)

// fetchOne retrieves a single URL and discards the body.
// Provided for you — the exercise is about structure, not HTTP.
func fetchOne(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}
func FetchSequential(urls []string) {
	for _, url := range urls {
		fetchOne(url)
	}
}

// TODO(reader): make these fetches run concurrently and wait for
// all of them before returning.
func FetchConcurrent(urls []string) {
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Go(func() {
			fetchOne(url)
		})
	}
	wg.Wait()
}
