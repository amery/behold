# behold

[![Go Reference][godoc_badge]][godoc_link]
[![Go Report Card][goreportcard_badge]][goreportcard_link]

`behold` is a generic key-value store with versioning and transactional capabilities. It's designed to offer a flexible interface for managing data with type safety using Go's generics.

[godoc_badge]: https://pkg.go.dev/badge/github.com/amery/behold.svg
[godoc_link]: https://pkg.go.dev/github.com/amery/behold
[goreportcard_badge]: https://goreportcard.com/badge/github.com/amery/behold
[goreportcard_link]: https://goreportcard.com/report/github.com/amery/behold

## Overview

The `behold` package provides:

- Generic key-value storage with strong type safety
- Transaction support for both read-only and read-write operations
- Built-in versioning of data
- Comprehensive query and filtering capabilities
- Flexible synchronization mechanisms

## Installation

```bash
go get github.com/amery/behold
```

## Key Features

### Type-Safe Generic Store

The package uses Go generics to provide a type-safe key-value store interface that works with any comparable key type and any value type.

### Transactional Operations

`behold` offers both read-only (`View`) and read-write (`Update`) transactions, with optional mutex locking for concurrent access control.

### Query System

A powerful query system allows filtering data with logical operations:

- Combine predicates with `AND` and `OR` operations
- Match values against complex conditions
- Create custom query functions for specific filtering needs

### Comparison Utilities

Built-in comparison functions that work with both standard comparable types and custom types with user-defined comparison logic.

## API Reference

### Core Interfaces

#### Store

```go
type Store[K comparable, V any] interface {
    Version() uint64
    Now() time.Time
    View(ctx context.Context, fn func(Tx[K, V]) error, locks ...Mutex) error
    Update(ctx context.Context, fn func(Tx[K, V]) error, locks ...Mutex) error
    Close() error
}
```

The main interface for the key-value store that handles transactions and versioning.

#### Transaction

```go
type Tx[K comparable, V any] interface {
    Context() context.Context
    Version() uint64
    Now() time.Time
    ForEach(fn func(key K, value V) bool, ors ...Query[any]) error
    Get(key K) (value V, err error)
    Set(key K, value V) error
    Append(key K, value V) error
    Delete(key K) error
    Commit() error
    Close() error
}
```

Interface for operations within a transaction context.

#### Query

```go
type Query[T any] interface {
    And(...Query[T]) Query[T]
    Or(...Query[T]) Query[T]
    Match(T) bool
}
```

Generic interface for filtering data with logical operations.

### Synchronization

`behold` provides flexible synchronization mechanisms:

```go
type Mutex interface {
    Lock()
    TryLock() bool
    Unlock()
}

type RWMutex interface {
    Mutex
    RLock()
    RUnlock()
    TryRLock() bool
}
```

Standard Go mutex types (`sync.Mutex` and `sync.RWMutex`) implement these interfaces.

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/amery/behold"
    // Assuming an implementation exists:
    "github.com/amery/behold/implementation"
)

func main() {
    // Create a new store with string keys and int values
    store := implementation.NewStore[string, int]()
    defer store.Close()
    
    // Write data in a transaction
    err := store.Update(context.Background(), func(tx behold.Tx[string, int]) error {
        if err := tx.Set("one", 1); err != nil {
            return err
        }
        if err := tx.Set("two", 2); err != nil {
            return err
        }
        if err := tx.Set("three", 3); err != nil {
            return err
        }

        return tx.Commit()
    })
    
    if err != nil {
        log.Fatalf("Failed to update: %v", err)
    }
    
    // Create a query that matches values greater than 1
    query := behold.GtQuery(1)
    
    // Read data in a read-only transaction
    err = store.View(context.Background(), func(tx behold.Tx[string, int]) error {
        // Print all key-value pairs where value > 1
        return tx.ForEach(func(key string, value int) bool {
            fmt.Printf("%s: %d\n", key, value)
            return true // continue iteration
        }, query)
    })
    
    if err != nil {
        log.Fatalf("Failed to view: %v", err)
    }
}
```

## See others

- [`darvaza.org/core`][darvaza_core] - A collection of utility packages for Go.
- [`darvaza.org/cache`][darvaza_cache] - A generic cache interface and implementations following
similar principles as this package.
- [`github.com/timshannon/bolthold`][bolthold] - A low-level key-value store for Go,
built on BoltDB, that inspired this package.
- [`github.com/timshannon/badgerhold`][badgerhold] - BadgerDB based sibling of [`bolthold`][bolthold], offering

[darvaza_core]: https://pkg.go.dev/darvaza.org/core
[darvaza_cache]: https://pkg.go.dev/darvaza.org/cache
[bolthold]: https://pkg.go.dev/github.com/timshannon/bolthold
[badgerhold]: https://pkg.go.dev/github.com/timshannon/badgerhold

## License

This project is licensed under the MIT License.
