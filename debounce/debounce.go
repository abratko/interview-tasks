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

import "time"

func Debounce[T any](
	fn func(T),
	wait time.Duration,
) func(T) {
	var timer *time.Timer
	return func(arg T) {
		if(timer != nil) {
			timer.Stop()
		}

		timer = time.AfterFunc(wait, func() { fn(arg) })
	}
}