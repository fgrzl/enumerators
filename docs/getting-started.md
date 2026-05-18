# Getting started

## Install

```bash
go get github.com/fgrzl/enumerators
```

## Iterate a slice

```go
e := enumerators.Slice([]int{1, 2, 3})
defer e.Dispose()

for e.MoveNext() {
    v, _ := e.Current()
    fmt.Println(v)
}
```

## Channel enumerator

```go
ctx := context.Background()
e := enumerators.Channel[string](ctx, 0)
defer e.Dispose()

go func() {
    e.Publish("hello")
    e.Publish("world")
    e.Complete()
}()

for e.MoveNext() {
    v, _ := e.Current()
    fmt.Println(v)
}
```

Canceling the context stops iteration; `Err()` returns the context error.

## Pipeline

```go
nums := enumerators.Slice([]int{1, 2, 3, 4, 5})
evens := enumerators.Filter(nums, func(n int) bool { return n%2 == 0 })
doubled := enumerators.Map(evens, func(n int) int { return n * 2 })

defer doubled.Dispose()
for doubled.MoveNext() {
    v, _ := doubled.Current()
    fmt.Println(v)
}
```

Prefer helpers like `enumerators.ForEach` when they dispose the chain for you.

## Next steps

- [Operations](operations.md) — full list of combinators
- [Overview](overview.md) — disposal rules and error handling
