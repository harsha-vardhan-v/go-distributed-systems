package ch05

import "time"

type Result struct {
	Source string
	Value  int
}

// Provided for you: one backend answering after a delay.
func query(source string, delay time.Duration) Result {
	time.Sleep(delay)
	return Result{Source: source, Value: len(source)}
}

// TODO(reader): First asks every source at once and returns whichever
// answers first, discarding the rest. It is sized wrong in exactly the
// way §5.1 describes.
//
// `results` is unbuffered, so the one goroutine whose answer is
// received finishes — and the other N-1 block forever on a send that
// nobody will ever take. What makes this easy to miss is that First
// still returns the right answer, quickly: from the caller's side
// nothing looks wrong.
//
// Fix it with §5.1's rule. Two constraints, so you reach for the
// right tool:
//   - do not add a WaitGroup and do not wait for the slow sources,
//   - do not change the signature,
//   - First must still return as soon as the first answer lands.
func First(sources []string, delays []time.Duration) Result {
	results := make(chan Result, len(sources)) // <- your move

	for i, src := range sources {
		go func() {
			results <- query(src, delays[i])
		}()
	}

	return <-results
}
