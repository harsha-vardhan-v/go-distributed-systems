package ch04

// Merge forwards every value from a and b onto a single channel and
// closes it once both inputs are done.
//
// TODO(reader): this is broken in the way §4.6 warns about. A closed
// channel is not "finished" as far as select is concerned — it is
// PERMANENTLY READY, and yields the zero value on every iteration. So
// the moment either input closes, this loop spins at full CPU pushing
// zeros into out, and never reaches close(out).
//
// Fix it with §4.6's technique:
//   - when a receive reports ok == false, take that case out of the
//     running so it stops being selected,
//   - skip the send for that iteration,
//   - and return once both are done, so the defer can close out.
//
// Do not change the signature, and do not count values — Merge has
// no idea how many are coming, which is the whole point.
func Merge(a, b <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()

	return out
}
