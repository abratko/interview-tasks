// Package debounce delays running a function until a quiet period has passed.
//
// Debounce: each new call within the wait window resets the timer, so the
// function runs only after calls stop for the given wait time.
//
// Typical use: search input, window resize, button clicks — cases where
// you want one action after rapid events, not one action per event.
//
// Contrast with throttle: throttle runs at most once per interval while
// events keep coming; debounce waits for the stream to settle, then runs once.
package debounce

import (
	"sync"
	"time"
)

func Debounce[T any](
	fn func(T),
	wait time.Duration,
) (func(T), func()) {
	var (
		timer *time.Timer
		mu sync.Mutex
	)

	return func(arg T) {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(wait, func() { fn(arg) })
	}, func() {
		mu.Lock()
		defer mu.Unlock()
		
		if (timer == nil) {
			return 
		}

		timer.Stop()
		timer = nil
	}
}
