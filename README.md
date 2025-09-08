[![Dependabot Updates](https://github.com/fgrzl/enumerators/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/enumerators/actions/workflows/dependabot/dependabot-updates)
[![ci](https://github.com/fgrzl/enumerators/actions/workflows/ci.yaml/badge.svg)](https://github.com/fgrzl/enumerators/actions/workflows/ci.yaml)

# enumerators

A small utility library for defining and consuming enumerators in Go. It provides a consistent interface to iterate over slices, channels, and custom generators.

## Installation

```bash
go get github.com/fgrzl/enumerators
````

## Overview

This library defines a generic `Enumerator[T]` interface with concrete implementations for common data sources:

* Slices
* Channels
* Callback-based generators

## Interface

All enumerators implement the `Enumerator[T]` interface and extend the `Disposable` interface:

```go
type Enumerator[T any] interface {
  Disposable
  MoveNext() bool
  Current() (T, error)
  Err() error
}

type Disposable interface {
  Dispose()
}
```

**Important**: All enumerators must be disposed when no longer needed to ensure proper resource cleanup. Most functions in this library automatically dispose enumerators for you, but when manually iterating, always call `Dispose()` or use `defer` to ensure cleanup.

## Built-in Enumerators

### SliceEnumerator

Iterates over a Go slice.

```go
e := enumerators.Slice([]int{1, 2, 3})
defer e.Dispose() // Ensure cleanup
for e.MoveNext() {
  v, _ := e.Current()
  fmt.Println(v)
}
```
```

### ChannelEnumerator

Iterates over values published to a channel-based enumerator.

```go
e := enumerators.Channel[string](context.Background(), 0)
defer e.Dispose() // Ensure cleanup

// Publish data in a separate goroutine
go func() {
  e.Publish("hello")
  e.Publish("world")
  e.Complete() // Signal completion
}()

for e.MoveNext() {
  v, _ := e.Current()
  fmt.Println(v)
}
```

## Automatic Disposal

Many functions in this library automatically dispose enumerators for you:

```go
// ToSlice automatically disposes the enumerator
slice, err := enumerators.ToSlice(enumerators.Slice([]int{1, 2, 3}))

// ForEach automatically disposes the enumerator
err := enumerators.ForEach(enumerators.Slice([]string{"a", "b"}), func(s string) error {
  fmt.Println(s)
  return nil
})

// Consume automatically disposes the enumerator
err := enumerators.Consume(enumerators.Slice([]int{1, 2, 3}))
```

## Chaining Operations

Enumerators can be chained together for complex data processing:

```go
result, err := enumerators.ToSlice(
  enumerators.Map(
    enumerators.Filter(
      enumerators.Slice([]int{1, 2, 3, 4, 5}),
      func(x int) bool { return x%2 == 0 }, // Keep even numbers
    ),
    func(x int) (string, error) { return fmt.Sprintf("num_%d", x), nil }, // Convert to string
  ),
)
// result: []string{"num_2", "num_4"}
```

