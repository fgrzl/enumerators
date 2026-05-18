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
doubled, _ := enumerators.Map(evens, func(n int) (int, error) { return n * 2, nil })

defer doubled.Dispose()
for doubled.MoveNext() {
    v, _ := doubled.Current()
    fmt.Println(v)
}
```

Prefer `enumerators.ForEach` when it disposes the chain for you.

## Next steps

- [Operations](operations.md) — combinators in this package
- [Overview](overview.md) — disposal and error handling
