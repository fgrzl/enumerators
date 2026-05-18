# Operations

Pipeline functions take an `Enumerator[T]` and return a new enumerator. Dispose the outermost enumerator when you own the chain (or use helpers like `ForEach` / `ToSlice` that dispose for you).

## Transform

| Function | Description |
|----------|-------------|
| `Map[T,U](enumerator, func(T) (U, error))` | Project each element; mapper errors stop iteration |
| `FlatMap` | Map to inner enumerators and flatten |
| `FilterMap` | Map with optional skip |

## Filter

| Function | Description |
|----------|-------------|
| `Filter` | Keep elements matching predicate |
| `Take` | First *n* elements |
| `Skip` | Skip first *n* elements |
| `TakeWhile` | Prefix while predicate is true |
| `SkipIf` | Skip elements where predicate is true (not LINQ `SkipWhile`) |
| `Distinct` | Unique comparable elements |

## Combine

| Function | Description |
|----------|-------------|
| `Chain[T](enumerators ...Enumerator[T])` | Sequence enumerators back-to-back |
| `Zip` | Pair-wise combine two enumerators into `Enumerator[Pair[T,U]]` |
| `Interleave` | Merge ordered streams by priority |

## Group / chunk

| Function | Description |
|----------|-------------|
| `Group` | Group by key into `Enumerator[Grouping[T,G]]` |
| `Chunk` / `ChunkByCount` | Fixed-size sub-enumerators |
| `Collect` | Materialize `Enumerator[Enumerator[T]]` to `[][]T` |

## Terminal

| Function | Description |
|----------|-------------|
| `ToSlice` | `([]T, error)` |
| `ForEach` | Side effect per element |
| `Any` / `All` | Short-circuit predicates |
| `Count` | Element count |
| `First` / `Last` | `(T, error)` — `ErrEmptySequence` when empty |
| `Min` / `Max` | Ordered types |
| `Sum` | With selector |

## Context variants

`MapWithContext`, `FilterWithContext` pass `context.Context` into callbacks.

## Notes

- `Channel` enumerators: call `Complete()` when publishing finishes
- Check `Err()` after loops; `Current()` can return per-element errors
- See package godoc and tests for full signatures
