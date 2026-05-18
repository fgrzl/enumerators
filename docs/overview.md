# Overview

The enumerators package standardizes **lazy iteration** in Go behind one interface, similar to `IEnumerable<T>` or Java streams, without allocating intermediate slices for every step.

## Core interface

```go
type Enumerator[T any] interface {
    Disposable
    MoveNext() bool
    Current() (T, error)
    Err() error
}
```

Always **`Dispose()`** enumerators when done (or use helpers that dispose automatically).

## Built-in sources

| Constructor | Source |
|-------------|--------|
| `Slice([]T)` | In-memory slice |
| `Channel(ctx, buffer)` | Goroutine-fed channel with `Publish` / `Complete` |
| Callback / generator helpers | Custom pull-based sequences |

## Composition

Operations (`Map`, `Filter`, `Take`, `Skip`, `Distinct`, `Chain`, etc.) return new enumerators. Most pipeline helpers dispose upstream enumerators when the pipeline ends.

## Design goals

- **Uniform consumption** — same `for MoveNext()` loop for KV enumeration, channels, and slices
- **Explicit errors** — `Current()` and `Err()` separate value errors from iteration state
- **Resource safety** — `Disposable` makes channel and generator cleanup predictable

## Typical use cases

- **KV and storage scans** — paginate large key ranges without loading everything into memory
- **Stream processing** — chain transforms before materializing with `ToSlice` or `ForEach`
- **Testing** — build enumerators from slices in unit tests, compose like production code
