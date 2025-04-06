// Package mutex provides utilities for working with collections of mutexes.
package mutex

import "sync"

// Mutexes is a collection of Mutex objects that can be locked and unlocked together.
type Mutexes []Mutex

// Lock acquires all locks in the collection in a deadlock-safe way.
func (ms Mutexes) Lock() { Lock(ms...) }

// Unlock releases all locks in the collection.
func (ms Mutexes) Unlock() { Unlock(ms...) }

// TryLock attempts to acquire all locks without blocking.
// Returns true if all locks were acquired, false otherwise.
func (ms Mutexes) TryLock() bool { return TryLock(ms...) }

// RWMutexes is a collection of RWMutex objects that can be locked and unlocked together.
type RWMutexes []RWMutex

// export converts an RWMutexes slice to a Mutexes slice.
func (ms RWMutexes) export() []Mutex {
	out := make(Mutexes, len(ms))
	for i, lock := range ms {
		out[i] = lock
	}
	return out
}

// Lock acquires exclusive locks on all RWMutexes in the collection.
func (ms RWMutexes) Lock() { Lock(ms.export()...) }

// Unlock releases exclusive locks on all RWMutexes in the collection.
func (ms RWMutexes) Unlock() { Unlock(ms.export()...) }

// TryLock attempts to acquire exclusive locks on all RWMutexes without blocking.
// Returns true if all locks were acquired, false otherwise.
func (ms RWMutexes) TryLock() bool { return TryLock(ms.export()...) }

// RLock acquires shared (read) locks on all RWMutexes in the collection.
func (ms RWMutexes) RLock() { RLock(ms.export()...) }

// RUnlock releases shared (read) locks on all RWMutexes in the collection.
func (ms RWMutexes) RUnlock() { RUnlock(ms.export()...) }

// TryRLock attempts to acquire shared (read) locks on all RWMutexes without blocking.
// Returns true if all locks were acquired, false otherwise.
func (ms RWMutexes) TryRLock() bool { return TryRLock(ms.export()...) }

// Type assertions to ensure interface compliance
var _ Mutex = Mutexes{}
var _ Mutex = RWMutexes{}
var _ sync.Locker = Mutexes{}
var _ sync.Locker = RWMutexes{}
var _ RWMutex = RWMutexes{}
