package tests

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/abratko/interview/debounce"
)

// Args packs call arguments into one value for Debounce[T].
type Args struct {
	ID  int
	Msg string
}

func TestDebounce_DoesNotCallBeforeWait(t *testing.T) {
	var calls atomic.Int32
	var mu sync.Mutex
	var got Args

	d, _ := debounce.Debounce(func(a Args) {
		calls.Add(1)
		mu.Lock()
		got = a
		mu.Unlock()
	}, 40*time.Millisecond)

	d(Args{ID: 1, Msg: "first"})
	time.Sleep(20 * time.Millisecond)
	if n := calls.Load(); n != 0 {
		t.Fatalf("calls = %d, want 0 before wait", n)
	}

	time.Sleep(30 * time.Millisecond)
	if n := calls.Load(); n != 1 {
		t.Fatalf("calls = %d, want 1 after wait", n)
	}
	mu.Lock()
	defer mu.Unlock()
	if got != (Args{ID: 1, Msg: "first"}) {
		t.Fatalf("got = %+v, want {ID:1 Msg:first}", got)
	}
}

func TestDebounce_ResetsTimerOnEachCall(t *testing.T) {
	var calls atomic.Int32
	var mu sync.Mutex
	var got Args

	d, _ := debounce.Debounce(func(a Args) {
		calls.Add(1)
		mu.Lock()
		got = a
		mu.Unlock()
	}, 40*time.Millisecond)

	d(Args{ID: 1, Msg: "a"})
	time.Sleep(25 * time.Millisecond)
	d(Args{ID: 2, Msg: "b"})
	time.Sleep(25 * time.Millisecond)
	d(Args{ID: 3, Msg: "c"})
	time.Sleep(25 * time.Millisecond)

	if n := calls.Load(); n != 0 {
		t.Fatalf("calls = %d, want 0 while calls keep resetting", n)
	}

	time.Sleep(25 * time.Millisecond)
	if n := calls.Load(); n != 1 {
		t.Fatalf("calls = %d, want 1 after quiet period", n)
	}
	mu.Lock()
	defer mu.Unlock()
	if got != (Args{ID: 3, Msg: "c"}) {
		t.Fatalf("got = %+v, want latest Args {ID:3 Msg:c}", got)
	}
}

func TestDebounce_OnlyLatestArgsRunOnce(t *testing.T) {
	var calls atomic.Int32
	var mu sync.Mutex
	var got Args

	d, _ := debounce.Debounce(func(a Args) {
		calls.Add(1)
		mu.Lock()
		got = a
		mu.Unlock()
	}, 20*time.Millisecond)

	for i := 0; i < 10; i++ {
		d(Args{ID: i, Msg: "burst"})
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(30 * time.Millisecond)
	if n := calls.Load(); n != 1 {
		t.Fatalf("calls = %d, want 1", n)
	}
	mu.Lock()
	defer mu.Unlock()
	if got != (Args{ID: 9, Msg: "burst"}) {
		t.Fatalf("got = %+v, want latest Args {ID:9 Msg:burst}", got)
	}
}

func TestDebounce_StopCancelsPendingCall(t *testing.T) {
	var calls atomic.Int32

	d, stop := debounce.Debounce(func(Args) {
		calls.Add(1)
	}, 40*time.Millisecond)

	d(Args{ID: 1, Msg: "pending"})
	time.Sleep(10 * time.Millisecond)
	stop()

	time.Sleep(50 * time.Millisecond)
	if n := calls.Load(); n != 0 {
		t.Fatalf("calls = %d, want 0 after stop", n)
	}
}

func TestDebounce_StopIsSafeWhenNothingPending(t *testing.T) {
	var calls atomic.Int32

	_, stop := debounce.Debounce(func(Args) {
		calls.Add(1)
	}, 40*time.Millisecond)

	stop()
	stop()

	if n := calls.Load(); n != 0 {
		t.Fatalf("calls = %d, want 0", n)
	}
}

func TestDebounce_CallWorksAgainAfterStop(t *testing.T) {
	var calls atomic.Int32
	var mu sync.Mutex
	var got Args

	d, stop := debounce.Debounce(func(a Args) {
		calls.Add(1)
		mu.Lock()
		got = a
		mu.Unlock()
	}, 30*time.Millisecond)

	d(Args{ID: 1, Msg: "canceled"})
	stop()

	d(Args{ID: 2, Msg: "after-stop"})
	time.Sleep(50 * time.Millisecond)

	if n := calls.Load(); n != 1 {
		t.Fatalf("calls = %d, want 1 after reschedule", n)
	}
	mu.Lock()
	defer mu.Unlock()
	if got != (Args{ID: 2, Msg: "after-stop"}) {
		t.Fatalf("got = %+v, want {ID:2 Msg:after-stop}", got)
	}
}
