package ch06

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func feed(values ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

func TestFanInReturnsEveryValue(t *testing.T) {
	in := []<-chan int{feed(1, 2, 3), feed(10, 20), feed(100)}

	done := make(chan []int, 1)
	go func() {
		var got []int
		for v := range FanIn(in) {
			got = append(got, v)
		}
		done <- got
	}()

	select {
	case got := <-done:
		slices.Sort(got)
		want := []int{1, 2, 3, 10, 20, 100}
		if !slices.Equal(got, want) {
			t.Fatalf("FanIn gave %v, want %v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("FanIn never closed its output channel")
	}
}

// This chapter's own thesis, as a test: a function that hands a stream
// to a caller should return <-chan T, so the caller cannot send into
// it or close it. §6.5.
func TestFanInReturnsAReceiveOnlyChannel(t *testing.T) {
	out := reflect.TypeOf(FanIn).Out(0)

	if out.Kind() != reflect.Chan {
		t.Fatalf("FanIn returns %v, want a channel", out)
	}
	if out.ChanDir() != reflect.RecvDir {
		t.Fatalf("FanIn returns %v; want <-chan int so callers "+
			"cannot send into it or close it. See §6.5.", out)
	}
}
