[![Dependabot Updates](https://github.com/fgrzl/enumerators/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/enumerators/actions/workflows/dependabot/dependabot-updates)
[![ci](https://github.com/fgrzl/enumerators/actions/workflows/ci.yaml/badge.svg)](https://github.com/fgrzl/enumerators/actions/workflows/ci.yaml)

# enumerators

A small utility library for defining and consuming enumerators in Go. It provides a consistent interface to iterate over slices, channels, and custom generators.

## Installation

```bash
go get github.com/fgrzl/enumerators
```

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

If the context passed to `Channel` is canceled, iteration stops and `Err()` returns the context error. `Publish` returns `false` after completion, disposal, or cancellation.

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

## Grouping Example

`Group` emits adjacent groups, so each group should be consumed before moving to the next one.

```go
grouped := enumerators.Group(
  enumerators.Slice([]string{"ant", "ape", "bat", "bee"}),
  func(item string) (byte, error) { return item[0], nil },
)

groups, err := enumerators.CollectGroupingSlices(grouped)
// groups: [{Group: 'a', Items: ["ant", "ape"]}, {Group: 'b', Items: ["bat", "bee"]}]
```

## Performance Considerations

### Memory Usage

**Lazy Operations** (recommended for large datasets):
- `Map`, `Filter`, `Take`, `Skip`, `FlatMap` - Process elements one by one, minimal memory usage

**Eager Operations** (materialize entire sequence in memory):
- `Sort`, `Reverse`, `Union`, `Intersect`, `Except`, `Sample` - Load all elements into memory
- **Warning**: These operations can cause out-of-memory errors for large datasets. Consider processing in batches or using lazy alternatives.

**Terminal Operations** (consume the enumerator):
- `ToSlice`, `ToMap`, `Count`, `Sum`, `Min`, `Max`, `First`, `Last` - May load data into memory depending on input size

### Error Handling

Operations that expect at least one element (`First`, `Last`, `Single`, `Aggregate`) return `ErrEmptySequence` for empty enumerators. Operations that work with empty sequences (`Min`, `Max`, `Sum`) return zero values.

## Available Operations

### Transformation
- `Map[T, U](enumerator, mapper)` - Transform each element
- `FlatMap[T, U](enumerator, mapper)` - Transform and flatten
- `Filter[T](enumerator, predicate)` - Keep elements matching predicate
- `FilterMap[TIn, TOut](enumerator, apply)` - Filter and transform in one step
- `Sort[T](enumerator, less)` - Sort elements
- `Reverse[T](enumerator)` - Reverse element order

### Selection
- `Take[T](enumerator, n)` - Take first N elements
- `TakeWhile[T](enumerator, condition)` - Take while condition holds
- `Skip[T](enumerator, n)` - Skip first N elements
- `SkipIf[T](enumerator, condition)` - Skip while condition holds
- `Distinct[T](enumerator)` - Remove duplicates
- `Sample[T](enumerator, n)` - Randomly sample N elements
- `ElementAt[T](enumerator, index)` - Get element at index
- `Single[T](enumerator)` - Get single element or error

### Aggregation
- `Sum[T, TSum](enumerator, selector)` - Sum elements
- `Count[T](enumerator)` - Count elements
- `Min[T](enumerator)` - Find minimum
- `Max[T](enumerator)` - Find maximum
- `First[T](enumerator)` - Get first element
- `Last[T](enumerator)` - Get last element
- `Any[T](enumerator, predicate)` - Check if any element matches
- `All[T](enumerator, predicate)` - Check if all elements match
- `Contains[T](enumerator, value)` - Check if contains value
- `Aggregate[TSource, TAccumulate](enumerator, seed, accumulator)` - Custom accumulation

### Comparison
- `SequenceEqual[T](first, second)` - Compare two sequences for equality

### Default Values
- `DefaultIfEmpty[T](enumerator, defaultValue)` - Provide default value for empty sequences

### Context-Aware Operations
- `MapWithContext[T, U](ctx, enumerator, mapper)` - Transform with context cancellation
- `FilterWithContext[T](ctx, enumerator, predicate)` - Filter with context cancellation

### Grouping and Chunking
- `Group[T, G](enumerator, keySelector)` - Group by key
- `Chunk[T, TSize](enumerator, sizeSelector)` - Chunk by size
- `ChunkByKey[V, K](in, keySelector)` - Chunk by key
- `ChunkWhen[V, K](in, splitFunc)` - Chunk when condition
- `Partition[T](enumerator, predicate)` - Partition into two slices

### Combination
- `Chain[T](enumerators...)` - Concatenate enumerators
- `Interleave[T, TOrdered](enumerators, keySelector)` - Interleave by key
- `Zip[T, U, V](left, right, zipper)` - Zip two enumerators
- `Union[T](enumerators...)` - Union of enumerators
- `Intersect[T](enumerators...)` - Intersection of enumerators
- `Except[T](left, right)` - Set difference

### Conversion
- `ToSlice[T](enumerator)` - Convert to slice
- `ToMap[T, K, V](enumerator, keyFn, valFn)` - Convert to map
- `ToMapWithKey[K, V](enumerator, keyFn)` - Convert to map with identity values

### Utilities
- `ForEach[T](enumerator, action)` - Execute action for each element
- `Consume[T](enumerator)` - Consume without processing
- `Peekable[T](enumerator)` - Add peek functionality
- `Cleanup[T](enumerator, cleanupFunc)` - Add cleanup function

