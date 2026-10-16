package ch07

import "sync"

// TODO(reader): FanOut spreads `in` across `workers` goroutines and
// merges their results onto one output channel. The work is already
// correct — every value is processed exactly once. What is missing is
// the end.
//
// Nobody closes `out`. Each worker just returns when `in` dries up,
// and the caller's `for range` waits forever for a close that never
// comes. §7.2 names the fix in one line: no single worker can know it
// is the last, so the close belongs to a coordinator that waits for
// all of them.
//
// This one does not panic. It hangs — which is why the test below
// carries its own deadline instead of trusting the runtime to notice.
//
// Do not change the signature, and do not close `out` from inside a
// worker: with more than one worker that is a double close.
func FanOut(
	in <-chan int,
	workers int,
	work func(int) int,
) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Go(func() {
			for n := range in {
				out <- work(n)
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
