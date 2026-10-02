package tests

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/abratko/interview/debounce"
)

// Race tests stress concurrent call/stop on one Debounce instance.
// Run with the race detector:
//
//	go test ./debounce/tests/ -race -run Race -count=1

func TestDebounce_RaceConcurrentCalls(t *testing.T) {
	var calls atomic.Int32
	d, stop := debounce.Debounce(func(Args) {
		calls.Add(1)
	}, 5*time.Millisecond)
	defer stop()

	const goroutines = 32
	const perGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				d(Args{ID: id*perGoroutine + i, Msg: "race-call"})
			}
		}(g)
	}
	wg.Wait()

	time.Sleep(20 * time.Millisecond)
	// At most one invocation after the burst settles; under -race the
	// important check is absence of a data race on timer.
	if n := calls.Load(); n > 1 {
		t.Fatalf("calls = %d, want at most 1 after concurrent burst", n)
	}
}

func TestDebounce_RaceConcurrentCallAndStop(t *testing.T) {
	var calls atomic.Int32
	d, stop := debounce.Debounce(func(Args) {
		calls.Add(1)
	}, 5*time.Millisecond)

	const goroutines = 32
	const perGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)
	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				d(Args{ID: id*perGoroutine + i, Msg: "race-mix"})
			}
		}(g)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				stop()
			}
		}()
	}
	wg.Wait()
	stop()

	time.Sleep(20 * time.Millisecond)
	// stop may cancel everything; any count is fine if -race is clean.
	_ = calls.Load()
}

func TestDebounce_RaceConcurrentStops(t *testing.T) {
	d, stop := debounce.Debounce(func(Args) {}, 5*time.Millisecond)
	d(Args{ID: 1, Msg: "pending"})

	const goroutines = 32
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				stop()
			}
		}()
	}
	wg.Wait()
}
