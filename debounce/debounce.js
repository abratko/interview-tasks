/**
 * Debounce
 *
 * Debounce delays running a function until a quiet period has passed.
 * Each new call within that period resets the timer, so the function
 * runs only after calls stop for the given wait time.
 *
 * Typical use: search input, window resize, button clicks — cases where
 * you want one action after rapid events, not one action per event.
 *
 * Contrast with throttle: throttle runs at most once per interval while
 * events keep coming; debounce waits for the stream to settle, then runs once.
 */

export function debounce(fn, waitMs) {
  let setTimeoutId = null
  const debounce = function (...args) {
    if (setTimeoutId) {
      clearTimeout(setTimeoutId)
    }

    setTimeoutId = setTimeout(
      () => fn.call(this, ...args),
      waitMs
    )
  }

  debounce.stop =function () {
    if (!setTimeoutId) {
      return 
    }

    clearTimeout(setTimeoutId)
    setTimeoutId = null
  }

  return debounce
}
