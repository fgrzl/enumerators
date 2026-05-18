# Operations

Pipeline functions take an `Enumerator[T]` and return a new enumerator. Unless noted, callers must still dispose the outermost enumerator (or use `ForEach` / `ToSlice` helpers that handle disposal).

## Transform

| Function | Description |
|----------|-------------|
| `Map` | Project each element |
| `Select` | Alias for `Map` |
| `Cast` | Change element type with a converter |

## Filter

| Function | Description |
|----------|-------------|
| `Filter` | Keep elements matching predicate |
| `Where` | Alias for `Filter` |
| `Take` | First *n* elements |
| `Skip` | Skip first *n* elements |
| `TakeWhile` | Prefix while predicate holds |
| `SkipWhile` | Drop prefix while predicate holds |
| `Distinct` | Unique elements (comparable) |

## Combine

| Function | Description |
|----------|-------------|
| `Concat` | Sequence enumerators back-to-back |
| `Zip` | Pair-wise combine two enumerators |

## Terminal

| Function | Description |
|----------|-------------|
| `ToSlice` | Materialize to `[]T` |
| `ForEach` | Run side effect per element |
| `Any` / `All` | Short-circuit predicates |
| `Count` | Count elements (may enumerate fully) |
| `First` / `Last` | Single element with ok flag |

## Notes

- Passing a disposed enumerator yields undefined behavior — dispose only the head you own
- `Channel` enumerators: call `Complete()` when publishing finishes so `MoveNext` terminates
- Errors from `Current()` do not always stop iteration; check `Err()` after the loop

For signatures and edge cases, see package godoc and tests under the repository root.
