# Interview practice

Short implementations for **Go** and **JavaScript** career interview tasks.

Each topic lives in its own folder (e.g. `debounce/`).  
Source files hold a short description of the problem; implementations and tests are added as you practice.

Tests for a topic are under `<topic>/tests/`.

## How to run tests

Run commands from the **project root** (`interview/`).

### JavaScript

Uses Node’s built-in test runner (`node:test`). No extra packages needed.

```bash
# one topic
node --test debounce/tests/debounce.test.js

# all JS tests under a topic (or the whole tree)
node --test debounce/tests/**/*.test.js
node --test **/tests/**/*.test.js
```

### Go

Uses the standard `testing` package. Requires Go installed and `go.mod` at the project root.

```bash
# one topic
go test ./debounce/tests/ -v

# all test packages
go test ./... -v
```

`-v` prints each test name; omit it for a short pass/fail summary.
