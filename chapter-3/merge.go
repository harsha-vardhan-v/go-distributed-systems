package ch03

import "sync"

// Provided for you: a producer that emits n values and then returns.
// It never closes the channel — who does that is the exercise.
func produce(id, n int, out chan<- int) {
	for i := 0; i < n; i++ {
		out <- id*100 + i
	}
}

// TODO(reader): Merge fans three producers into one channel and
// collects everything they send. It is broken in exactly the way
// §3.3 warns about: every producer closes the channel when it
// finishes, so whichever one finishes second panics — either
// "close of closed channel" or, if it is still mid-send, "send on
// closed channel".
//
// Fix it so that:
//   - the channel is closed exactly once,
//   - it is closed only after EVERY producer has finished, and
//   - the `for range` below still terminates.
//
// Two constraints, so you reach for the right tool:
//   - do not change the signature, and
//   - do not drain with a fixed count. The caller does not know
//     how many values are coming, which is the whole reason
//     close() exists.
func Merge(counts []int) []int {
	ch := make(chan int)
	var wg sync.WaitGroup

	for id, n := range counts {
		wg.Go(func() {
			produce(id, n, ch)
		})
	}

	//Coordinator
	go func() {
		wg.Wait()
		close(ch)
	}()

	var got []int
	for v := range ch {
		got = append(got, v)
	}
	return got
}
