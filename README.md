[![Dependabot Updates](https://github.com/fgrzl/enumerators/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/enumerators/actions/workflows/dependabot/dependabot-updates)
[![ci](https://github.com/fgrzl/enumerators/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/enumerators/actions/workflows/ci.yml)

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

```go
type Enumerator[T any] interface {
  MoveNext() bool
  Current() (T, error)
  Err() error
}
```

## Built-in Enumerators

### SliceEnumerator

Iterates over a Go slice.

```go
e := enumerators.NewSliceEnumerator([]int{1, 2, 3})
for e.MoveNext() {
  v, _ := e.Current()
  fmt.Println(v)
}
```

### ChannelEnumerator

Iterates over a channel until it's closed.

```go
ch := make(chan string)
go func() {
  ch <- "hello"
  ch <- "world"
  close(ch)
}()

e := enumerators.NewChannelEnumerator(ch)
for e.MoveNext() {
  v, _ := e.Current()
  fmt.Println(v)
}
```

