---
description: 'Instructions for writing Go code following idiomatic Go practices and community standards'
applyTo: '**/*.go,**/go.mod,**/go.sum'
---

# Go Development Instructions

Follow idiomatic Go practices and community standards when writing Go code.
These instructions are based on [Effective Go](https://go.dev/doc/effective_go),
[Go Code Review Comments](https://go.dev/wiki/CodeReviewComments),
[Go Test Code Comments](https://go.dev/wiki/TestComments),
and [Google's Go Style Guide](https://google.github.io/styleguide/go/).

## Go Version Awareness

- Always determine the project's Go version from the `go` directive in `go.mod` before suggesting changes.
- Do **not** propose rewrites that replace modern standard-library features (Go 1.21+) with hand-rolled equivalents.
- Do **not** flag, downgrade, or "polyfill" the following modern features when the `go.mod` permits them:
  - Go 1.20: `errors.Join`.
  - Go 1.21: built-ins `min`, `max`, `clear`; stdlib packages `slices`, `maps`, `cmp`; structured logging `log/slog`.
  - Go 1.22: range-over-integer (`for i := range n`); per-iteration loop variable scoping; `math/rand/v2`; enhanced `http.ServeMux` patterns; `cmp.Or`.
  - Go 1.23: range-over-function iterators (`iter.Seq`, `iter.Seq2`); `unique` package; `slices.Sorted`, `slices.Collect`, etc.
  - Go 1.24: generic type aliases; `omitzero` JSON tag; `weak` package; `testing.B.Loop`; `os.Root`.
  - Go 1.25: `testing/synctest`; `runtime.AddCleanup`; container-aware `GOMAXPROCS`; `encoding/json/v2` (experimental).
- When proposing code that requires a feature newer than `go.mod` allows, call this out explicitly rather than silently using it.
- Prefer the newest idiomatic form available for the project's Go version; only suggest the older pattern when the feature is genuinely unavailable.

## General Instructions

- Write simple, clear, and idiomatic Go code
- Favor clarity and simplicity over cleverness
- Follow the principle of least surprise
- Keep the happy path left-aligned (minimize indentation)
- Return early to reduce nesting
- Make the zero value useful
- Document exported types, functions, methods, and packages
- Use Go modules for dependency management

## Naming Conventions

### Packages

- Use lowercase, single-word package names
- Avoid underscores, hyphens, or mixedCaps
- Choose names that describe what the package provides, not what it contains
- Avoid generic names like `util`, `common`, or `base`
- Package names should be singular, not plural

### Variables and Functions

- Use mixedCaps or MixedCaps (camelCase) rather than underscores
- Keep names short but descriptive
- Use single-letter variables only for very short scopes (like loop indices)
- Exported names start with a capital letter
- Unexported names start with a lowercase letter
- Avoid stuttering (e.g., avoid `http.HTTPServer`, prefer `http.Server`)

### Interfaces

- Name interfaces with -er suffix when possible (e.g., `Reader`, `Writer`, `Formatter`)
- Single-method interfaces should be named after the method (e.g., `Read` → `Reader`)
- Keep interfaces small and focused

### Constants

- Use MixedCaps for exported constants
- Use mixedCaps for unexported constants
- Group related constants using `const` blocks
- Consider using typed constants for better type safety

## Code Style and Formatting

### Formatting

- Always use `gofmt` to format code
- Use `goimports` to manage imports automatically
- Keep line length reasonable (no hard limit, but consider readability)
- Add blank lines to separate logical groups of code

### Comments

- Write comments in complete sentences
- Start sentences with the name of the thing being described
- Package comments should start with "Package [name]"
- Use line comments (`//`) for most comments
- Use block comments (`/* */`) only for files in `testdata`
- Document why, not what, unless the what is complex

### Error Handling

- Check errors immediately after the function call
- Don't ignore errors using `_` unless you have a good reason (document why)
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Create custom error types when you need to check for specific errors
- Place error returns as the last return value
- Name error variables `err`
- Keep error messages lowercase and don't end with punctuation

## Architecture and Project Structure

### Package Organization

- Follow standard Go project layout conventions
- Use `internal/` for packages that shouldn't be imported by external projects
- Group related functionality into packages
- Avoid circular dependencies

### Dependency Management

- Use Go modules (`go.mod` and `go.sum`)
- Keep dependencies minimal
- Regularly update dependencies for security patches
- Use `go mod tidy` to clean up unused dependencies

## Type Safety and Language Features

### Type Definitions

- Define types to add meaning and type safety
- Use struct tags for JSON, XML, database mappings
- Prefer explicit type conversions
- Use type assertions carefully and check the second return value

### Pointers vs Values

- Use pointers for large structs or when you need to modify the receiver
- Use values for small structs and when immutability is desired
- Be consistent within a type's method set
- Consider the zero value when choosing pointer vs value receivers

### Interfaces and Composition

- Accept interfaces, return concrete types
- Keep interfaces small (1-3 methods is ideal)
- Use embedding for composition
- Define interfaces close to where they're used, not where they're implemented
- Don't export interfaces unless necessary

## Concurrency

### Goroutines

- Don't create goroutines in libraries; let the caller control concurrency
- Always know how a goroutine will exit
- Use `sync.WaitGroup` or channels to wait for goroutines
- Avoid goroutine leaks by ensuring cleanup

### Channels

- Use channels to communicate between goroutines
- Don't communicate by sharing memory; share memory by communicating
- Close channels from the sender side, not the receiver
- Use buffered channels when you know the capacity
- Use `select` for non-blocking operations

### Synchronization

- Use `sync.Mutex` for protecting shared state
- Keep critical sections small
- Use `sync.RWMutex` when you have many readers
- Prefer channels over mutexes when possible
- Use `sync.Once` for one-time initialization

## Error Handling Patterns

### Creating Errors

- Use `errors.New` for simple static errors
- Use `fmt.Errorf` for dynamic errors
- Create custom error types for domain-specific errors
- Export error variables for sentinel errors
- Use `errors.Is` and `errors.As` for error checking

### Error Propagation

- Add context when propagating errors up the stack
- Don't log and return errors (choose one)
- Handle errors at the appropriate level
- Consider using structured errors for better debugging

## Performance Optimization

### Memory Management

- Minimize allocations in hot paths
- Reuse objects when possible (consider `sync.Pool`)
- Use value receivers for small structs
- Preallocate slices when size is known
- Avoid unnecessary string conversions

### Profiling

- Use built-in profiling tools (`pprof`)
- Benchmark critical code paths
- Profile before optimizing
- Focus on algorithmic improvements first
- Consider using `testing.B` for benchmarks

## Testing

### Test Organization

- Keep tests in the same package (white-box testing)
- Use `_test` package suffix for black-box testing
- Name test files with `_test.go` suffix
- Place test files next to the code they test

### Writing Tests

- Use table-driven tests for multiple test cases
- Name tests descriptively using `TestStructName_FunctionName`
- Use subtests with `t.Run` for better organization
- Test both success and error cases
- Avoid the use of `assert` libraries for clearer test failures
- When you use `t.Run` to create a subtest, the first argument is used as a descriptive name for the test
- Test outputs should output the actual value that the function returned before printing the value that was expected

### Test Helpers

- Mark helper functions with `t.Helper()`
- Create test fixtures for complex setup
- Clean up resources using `t.Cleanup()`

## Documentation

### Code Documentation

- Document all exported symbols
- Start documentation with the symbol name
- Keep documentation close to code
- Update documentation when code changes

## Tools and Development Workflow

### Essential Tools

- `go fmt`: Format code
- `go vet`: Find suspicious constructs
- `go test`: Run tests
- `go mod`: Manage dependencies
- `revive` or `golangci-lint`: Additional linting

### Development Practices

- Run tests before committing
- Keep commits focused and atomic
- Write meaningful commit messages
- Review diffs before committing

## Modern Go Features (Go 1.21+) — Do / Don't

Use these examples when the `go.mod` version permits the feature. Do not suggest reverting any of the "Do" forms below to the "Don't" forms.

### Built-in `min`, `max`, `clear` (Go 1.21)

Don't:

```go
func maxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}
n := maxInt(x, y)

for k := range m {
    delete(m, k)
}
```

Do:

```go
n := max(x, y)
clear(m)
```

### `slices` and `maps` packages (Go 1.21)

Don't:

```go
func contains(haystack []string, needle string) bool {
    for _, s := range haystack {
        if s == needle {
            return true
        }
    }
    return false
}

sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
```

Do:

```go
if slices.Contains(haystack, needle) { /* ... */ }
slices.Sort(xs)
keys := slices.Sorted(maps.Keys(m)) // Go 1.23
```

### `errors.Join` (Go 1.20+) and `%w` wrapping

Don't:

```go
type multiErr struct{ errs []error }
func (m *multiErr) Error() string { /* manual join */ }
```

Do:

```go
return errors.Join(err1, err2, err3)
```

### `log/slog` for structured logging (Go 1.21)

Don't:

```go
log.Printf("user %s failed login from %s: %v", user, ip, err)
```

Do:

```go
slog.Error("login failed",
    "user", user,
    "ip", ip,
    "err", err,
)
```

### Range-over-integer (Go 1.22)

Don't:

```go
for i := 0; i < n; i++ {
    work()
}
```

Do:

```go
for range n {
    work()
}
```

### Loop variable scoping (Go 1.22)

Don't (no longer necessary):

```go
for _, v := range items {
    v := v // shadow to capture per iteration
    go func() { handle(v) }()
}
```

Do:

```go
for _, v := range items {
    go func() { handle(v) }() // each iteration has its own v since Go 1.22
}
```

### `cmp.Or` for fallback chains (Go 1.22)

Don't:

```go
name := userName
if name == "" {
    name = envName
}
if name == "" {
    name = "anonymous"
}
```

Do:

```go
name := cmp.Or(userName, envName, "anonymous")
```

### Range-over-function iterators (Go 1.23)

Don't:

```go
ch := make(chan Item)
go func() {
    defer close(ch)
    for _, it := range source {
        ch <- it
    }
}()
for it := range ch { /* ... */ }
```

Do:

```go
func Items[T any](source []T) iter.Seq[T] {
    return func(yield func(T) bool) {
        for _, it := range source {
            if !yield(it) {
                return
            }
        }
    }
}

for it := range Items(source) { /* ... */ }
```

### `omitzero` JSON tag (Go 1.24)

Don't (pointer just to make `omitempty` skip the zero struct):

```go
type Config struct {
    Window *Window `json:"window,omitempty"`
}
```

Do:

```go
type Config struct {
    Window Window `json:"window,omitzero"`
}
```

### Generic type aliases (Go 1.24)

Don't (re-declaring instead of aliasing):

```go
type Vec[T any] []T
```

Do:

```go
type VecAlias[T any] = []T
```

### `testing/synctest` for time-dependent tests (Go 1.25)

Don't:

```go
func TestDebounce(t *testing.T) {
    d := NewDebouncer(100 * time.Millisecond)
    d.Trigger()
    time.Sleep(150 * time.Millisecond) // real wall-clock wait
    // assert fired
}
```

Do:

```go
func TestDebounce(t *testing.T) {
    synctest.Run(func() {
        d := NewDebouncer(100 * time.Millisecond)
        d.Trigger()
        time.Sleep(150 * time.Millisecond) // synthetic time, instant
        // assert fired
    })
}
```

### `testing.B.Loop` for benchmarks (Go 1.24)

Don't:

```go
func BenchmarkParse(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Parse(input)
    }
}
```

Do:

```go
func BenchmarkParse(b *testing.B) {
    for b.Loop() {
        Parse(input)
    }
}
```

### Container-aware `GOMAXPROCS` (Go 1.25)

Don't:

```go
// Manually clamp to cgroup CPU quota at startup.
runtime.GOMAXPROCS(detectCgroupCPUs())
```

Do:

Rely on the default — since Go 1.25 the runtime reads the cgroup CPU quota on Linux automatically. Only override when you have a measured reason.

## Common Pitfalls to Avoid

- Not checking errors
- Ignoring race conditions
- Creating goroutine leaks
- Not using defer for cleanup
- Modifying maps concurrently
- Not understanding nil interfaces vs nil pointers
- Forgetting to close resources (files, connections)
- Using global variables unnecessarily
- Over-using empty interfaces (`interface{}` or `any`)
- Not considering the zero value of types
- Suggesting hand-rolled replacements for standard-library features the project's Go version already provides (see "Modern Go Features" above)
- Re-introducing pre-Go-1.22 loop-variable shadowing (`v := v`) when the module targets Go 1.22 or later
- Replacing `log/slog` with `log.Printf`-style unstructured logging
