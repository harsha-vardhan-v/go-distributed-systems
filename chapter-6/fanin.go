package ch06

import "sync"

// TODO(reader): FanIn merges several input channels into one. Two
// things are wrong with it, and this chapter fixes both.
//
//  1. Every goroutine closes `out` when its own input dries up.
//     With more than one input that is a double close, and the
//     panic you see is a race: usually `send on closed channel`
//     (another worker was mid-send), sometimes `close of closed
//     channel`. §6.4 Pattern 2 says why: no single producer can
//     know it is the last, so the close belongs to a coordinator
//     that waits for all of them.
//
//  2. The signature hands the caller a bidirectional `chan int`, so
//     nothing stops them sending into it or closing it early. §6.5
//     says what to return instead — and the type system will then
//     enforce it for you.
//
// Do not change the parameter type, and do not buffer your way out of
// the panic: a buffer would only postpone it.
func FanIn(inputs []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Go(func() {
			for v := range in {
				out <- v
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
