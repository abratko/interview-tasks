const { describe, it, mock, afterEach } = require('node:test');
const assert = require('node:assert/strict');
const { debounce } = require('../debounce.js');

describe('debounce', () => {
  afterEach(() => {
    mock.timers.reset();
  });

  it('does not call the function before wait elapses', () => {
    mock.timers.enable({ apis: ['setTimeout'] });

    let calls = 0;
    const d = debounce(() => {
      calls += 1;
    }, 100);

    d();
    assert.equal(calls, 0);

    mock.timers.tick(99);
    assert.equal(calls, 0);

    mock.timers.tick(1);
    assert.equal(calls, 1);
  });

  it('resets the timer on each call within the wait window', () => {
    mock.timers.enable({ apis: ['setTimeout'] });

    let calls = 0;
    const d = debounce(() => {
      calls += 1;
    }, 100);

    d();
    mock.timers.tick(50);
    d();
    mock.timers.tick(50);
    d();
    mock.timers.tick(99);
    assert.equal(calls, 0);

    mock.timers.tick(1);
    assert.equal(calls, 1);
  });

  it('passes the latest arguments to the function', () => {
    mock.timers.enable({ apis: ['setTimeout'] });

    /** @type {unknown[]} */
    const received = [];
    const d = debounce((...args) => {
      received.push(...args);
    }, 50);

    d('a');
    d('b', 2);
    mock.timers.tick(50);

    assert.deepEqual(received, ['b', 2]);
  });

  it('preserves `this` when called as a method', () => {
    mock.timers.enable({ apis: ['setTimeout'] });

    /** @type {unknown} */
    let seen = null;
    const obj = {
      value: 42,
      run: debounce(function run() {
        seen = this.value;
      }, 30),
    };

    obj.run();
    mock.timers.tick(30);

    assert.equal(seen, 42);
  });
});
