# mutex

Package mutex provides interfaces and utilities for mutual exclusion and synchronization primitives.

## Overview

This package defines standardized interfaces for mutex operations and utilities for working with multiple locks simultaneously. It's designed to work with standard library synchronization primitives while providing additional functionality.

## Features

- Standardized `Mutex` and `RWMutex` interfaces
- Utilities for read-only mutexes
- Functions for operating on multiple locks simultaneously
- Safe lock/unlock operations with proper error handling
- Collections of mutexes that can be operated on as a group

## Interfaces

### Mutex

The core interface for mutual exclusion operations:

```go
type Mutex interface {
    // Lock acquires the mutex, blocking until it is available.
    Lock()

    // TryLock attempts to acquire the mutex without blocking.
    // Returns true if the lock was acquired, false otherwise.
    TryLock() bool

    // Unlock releases the mutex.
    // Calling Unlock on an unlocked mutex is expected to panic.
    Unlock()
}
```

### RWMutex

Extends the Mutex interface with read-locking capabilities:

```go
type RWMutex interface {
    Mutex

    // RLock acquires a read lock on the mutex, blocking until it is available
    // if necessary.
    RLock()

    // RUnlock releases a read lock on the mutex.
    // Calling RUnlock on a mutex not holding a read lock is expected to panic.
    RUnlock()

    // TryRLock attempts to acquire a read lock without blocking.
    // Returns true if the lock was acquired, false otherwise.
    TryRLock() bool
}
```

## Functions

### Lock Operations

- `Lock(locks ...Mutex)`: Acquires multiple locks in order
- `RLock(locks ...Mutex)`: Acquires multiple read locks when possible
- `TryLock(locks ...Mutex) bool`: Non-blocking attempt to acquire multiple locks
- `TryRLock(locks ...Mutex) bool`: Non-blocking attempt to acquire multiple read locks

### Unlock Operations

- `Unlock(locks ...Mutex)`: Releases multiple locks
- `RUnlock(locks ...Mutex)`: Releases multiple read locks

### Mutex Conversions

- `ROMutex(m RWMutex) Mutex`: Converts an RWMutex to a read-only Mutex

## Collections

### Mutexes

`Mutexes` is a slice-based collection of `Mutex` objects that can be locked and unlocked together:

```go
type Mutexes []Mutex
```

It implements the `Mutex` interface, allowing a collection of mutexes to be treated as a single mutex:

- `Lock()`: Acquires all locks in the collection in a deadlock-safe way
- `Unlock()`: Releases all locks in the collection
- `TryLock() bool`: Attempts to acquire all locks without blocking, returning true only if all locks were acquired

### RWMutexes

`RWMutexes` is a slice-based collection of `RWMutex` objects that can be locked and unlocked together:

```go
type RWMutexes []RWMutex
```

It implements the `RWMutex` interface, providing both exclusive and shared locking operations on the collection:

- `Lock()`: Acquires exclusive locks on all mutexes in the collection
- `Unlock()`: Releases exclusive locks on all mutexes in the collection
- `TryLock() bool`: Attempts to acquire exclusive locks without blocking
- `RLock()`: Acquires shared (read) locks on all mutexes in the collection
- `RUnlock()`: Releases shared (read) locks on all mutexes in the collection
- `TryRLock() bool`: Attempts to acquire shared locks without blocking

Both collections ensure interface compliance with the `Mutex`, `RWMutex`, and `sync.Locker` interfaces.

## Examples

### Basic Usage

```go
import (
    "sync"
    "github.com/amery/behold/mutex"
)

func example() {
    var mu sync.Mutex
    
    // Standard operations
    mu.Lock()
    // Critical section
    mu.Unlock()
    
    // Try lock without blocking
    if mu.TryLock() {
        // Got the lock
        mu.Unlock()
    }
}
```

### Multiple Locks

```go
func multiLockExample(mu1, mu2, mu3 mutex.Mutex) {
    // Acquire multiple locks safely
    mutex.Lock(mu1, mu2, mu3)
    // Critical section with all locks held
    mutex.Unlock(mu3, mu2, mu1) // Can release in any order
}
```

### Using Mutex Collections

```go
func mutexCollectionExample() {
    // Create a collection of mutexes
    locks := mutex.Mutexes{&sync.Mutex{}, &sync.Mutex{}, &sync.Mutex{}}
    
    // Lock all mutexes in the collection at once
    locks.Lock()
    // Critical section with all locks held
    locks.Unlock()
    
    // Create a collection of RW mutexes
    rwlocks := mutex.RWMutexes{&sync.RWMutex{}, &sync.RWMutex{}}
    
    // Acquire read locks on all mutexes
    rwlocks.RLock()
    // Read operations
    rwlocks.RUnlock()
}
```

### Read-Only Mutex

```go
func readOnlyExample() {
    var rwMutex sync.RWMutex
    
    // Convert to read-only mutex
    readOnly := mutex.ROMutex(&rwMutex)
    
    // Now it can only be used for read locking
    readOnly.Lock()   // Uses RLock internally
    // Read-only operations
    readOnly.Unlock() // Uses RUnlock internally
}
```

> **Important Note**: When using custom RWMutex implementations with ROMutex, the behavior may differ from standard library mutexes. Custom implementations can impose additional constraints or provide enhanced capabilities that affect what operations are possible on the protected data. Unlike standard library mutexes which only manage concurrency, custom implementations might implement domain-specific access control or validation, so the "read-only" nature is defined by the implementation rather than being a universally consistent guarantee. Be cautious when using ROMutex with custom mutex implementations and don't rely solely on type constraints for security or data integrity.

## Compatibility

This package works with standard library types like `sync.Mutex` and `sync.RWMutex` which implement the interfaces defined here.

## License

MIT License
