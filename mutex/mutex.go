// Package mutex provides interfaces and utilities for mutual exclusion and synchronization primitives.
package mutex

import (
	"errors"
	"sync"

	"darvaza.org/core"
)

// Mutex defines a standard interface for mutual exclusion locking mechanisms
// that support basic locking, unlocking, and non-blocking lock attempts.
//
// This interface is implemented by standard library types like sync.Mutex
// and sync.RWMutex.
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

// RWMutex extends the Mutex interface with read-locking capabilities,
// allowing multiple readers or a single writer to access a shared resource.
//
// This interface is implemented by standard library types like sync.RWMutex.
// When a write lock is held, all read lock attempts will block until the write
// lock is released. Multiple read locks can be held simultaneously.
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

// ROMutex converts an RWMutex to a read-only Mutex, allowing only read locking operations.
// If the input mutex is nil, it returns nil.
//
// The returned Mutex adapts Lock/Unlock/TryLock calls to use the underlying
// RWMutex's RLock/RUnlock/TryRLock methods respectively. This is useful when you want
// to restrict access to read-only operations for certain code paths while maintaining
// the Mutex interface contract.
// Example usage:
//
//	var rwm sync.RWMutex
//	readOnlyAccess := ROMutex(&rwm)
//	// Now readOnlyAccess can only perform read locking operations
func ROMutex(m RWMutex) Mutex {
	if m == nil {
		return nil
	}

	return &readOnlyMutex{m: m}
}

// readOnlyMutex uses a RWMutex as a Mutex in read-only mode
type readOnlyMutex struct {
	m RWMutex
}

// Lock acquires a read lock on the underlying RWMutex.
func (m readOnlyMutex) Lock() { m.m.RLock() }

// Unlock releases a read lock on the underlying RWMutex.
func (m readOnlyMutex) Unlock() { m.m.RUnlock() }

// TryLock attempts to acquire a read lock on the underlying RWMutex.
// Returns true if the lock was acquired, false otherwise.
func (m readOnlyMutex) TryLock() bool { return m.m.TryRLock() }

// TryLock attempts to acquire locks on multiple mutexes simultaneously without blocking.
// Returns true if all locks were successfully acquired, false otherwise.
// Panics if an error occurs during the locking process.
//
// This function will panic in two cases:
// 1. If any of the provided mutexes are nil
// 2. If any mutex panics during lock/unlock operations (after releasing locks)
func TryLock(locks ...Mutex) bool {
	ok, err := doTryLock(false, locks)
	if err != nil {
		panic(err)
	}
	return ok
}

// doTryLock attempts to acquire locks on multiple mutexes simultaneously.
// If readOnly is true, it tries to acquire read locks when possible.
// Returns a boolean indicating success and any errors encountered.
//
//revive:disable-next-line:cognitive-complexity
func doTryLock(readOnly bool, locks []Mutex) (bool, error) {
	if len(locks) == 0 {
		return false, errors.New("cannot lock on nothing")
	}

	var i int

	for _, mu := range locks {
		ok, err := doTryLockOne(readOnly, mu)
		if !ok || err != nil {
			var errs core.CompoundError

			if err != nil {
				errs.AppendError(err)
			}

			if err := reverseUnlock(readOnly, locks[:i]...); err != nil {
				errs.AppendError(err)
			}

			return false, errs.AsError()
		}

		i++
	}

	return true, nil
}

// doTryLockOne attempts to acquire a lock on a single mutex.
// If readOnly is true and the mutex implements RWMutex, it tries to acquire a read lock.
// Returns a boolean indicating success and any errors encountered.
func doTryLockOne(readOnly bool, mu Mutex) (bool, error) {
	var ok bool

	if mu == nil {
		// can't lock on nothing
		return false, core.NewPanicError(1, "nil Mutex not allowed")
	}

	fn := makeTryLockFunc(readOnly, mu, &ok)
	if err := core.Catch(fn); err != nil {
		return false, err
	}

	return ok, nil
}

//revive:disable-next-line:flag-parameter
func makeTryLockFunc(readOnly bool, mu Mutex, ok *bool) func() error {
	if readOnly {
		if rw, is := mu.(RWMutex); is {
			// TryRLock
			return func() error {
				*ok = rw.TryRLock()
				return nil
			}
		}
	}

	// TryLock
	return func() error {
		*ok = mu.TryLock()
		return nil
	}
}

// TryRLock attempts to acquire read locks on multiple mutexes simultaneously.
// It returns true if all read locks are successfully acquired, and false otherwise.
// If any read lock fails, it releases any previously acquired read locks in reverse order.
//
// For RWMutex instances, it uses the TryRLock method instead of TryLock.
// For regular Mutex instances, it falls back to TryLock.
//
// This function will panic in two cases:
// 1. If any of the provided mutexes are nil
// 2. If any mutex panics during lock/unlock operations (after releasing locks)
func TryRLock(locks ...Mutex) bool {
	ok, err := doTryLock(true, locks)
	if err != nil {
		panic(err)
	}
	return ok
}

// Unlock releases multiple mutexes simultaneously.
// It attempts to unlock all provided mutexes even if some operations fail.
// If any unlock operation panics, errors will be aggregated and raised at the end.
//
// This function will panic:
// 1. If any of the provided mutexes are nil
// 2. If any mutex panics during unlock operations
func Unlock(locks ...Mutex) {
	if err := doUnlock(false, locks); err != nil {
		panic(err)
	}
}

// doUnlock releases multiple locks.
// If readOnly is true, it releases read locks when possible.
// Returns any errors encountered during the unlock operations.
func doUnlock(readOnly bool, locks []Mutex) error {
	var errs core.CompoundError

	for _, mu := range locks {
		if err := doUnlockOne(readOnly, mu); err != nil {
			errs.AppendError(err)
		}
	}

	return errs.AsError()
}

// doUnlockOne releases a lock on a single mutex.
// If readOnly is true and the mutex implements RWMutex, it releases a read lock.
// Returns any errors encountered during the unlock operation.
func doUnlockOne(readOnly bool, mu Mutex) error {
	if mu == nil {
		// can't unlock nothing
		return core.NewPanicError(2, "nil Mutex not allowed")
	}

	fn := makeUnlockFn(readOnly, mu)
	return core.Catch(fn)
}

//revive:disable-next-line:flag-parameter
func makeUnlockFn(readOnly bool, mu Mutex) func() error {
	if readOnly {
		if rw, is := mu.(RWMutex); is {
			// RUnlock
			return func() error {
				rw.RUnlock()
				return nil
			}
		}
	}

	// Unlock
	return func() error {
		mu.Unlock()
		return nil
	}
}

// reverseUnlock releases previously acquired locks in reverse order.
// collecting any possible panic. It's used when a lock request fails
// to prevent deadlocks.
// this internal function won't be called with nil mutexes as they
// would have never been acquired in the first place.
func reverseUnlock(readOnly bool, locks ...Mutex) error {
	var errs core.CompoundError

	for i := len(locks) - 1; i >= 0; i-- {
		if err := doUnlockOne(readOnly, locks[i]); err != nil {
			errs.AppendError(err)
		}
	}

	return errs.AsError()
}

// RUnlock releases multiple mutexes simultaneously using read-unlock operations.
// For RWMutex instances, it uses the RUnlock method instead of Unlock.
// For regular Mutex instances, it falls back to Unlock.
//
// It attempts to unlock all provided mutexes even if some operations fail.
// If any unlock operation panics, errors will be aggregated and raised at the end.
//
// This function will panic:
// 1. If any of the provided mutexes are nil
// 2. If any mutex panics during unlock operations
func RUnlock(locks ...Mutex) {
	if err := doUnlock(true, locks); err != nil {
		panic(err)
	}
}

// Lock acquires multiple mutexes simultaneously.
// It acquires locks in the order provided and blocks until all locks are acquired.
// If any lock operation fails or panics, previously acquired locks are released
// in reverse order to prevent deadlocks.
//
// This function will panic:
// 1. If any of the provided mutexes are nil
// 2. If any mutex panics during lock operations (after releasing locks)
func Lock(locks ...Mutex) {
	if err := doLock(false, locks); err != nil {
		panic(err)
	}
}

// doLock acquires locks on multiple mutexes in sequence.
// If readOnly is true, it acquires read locks when possible.
// If any lock acquisition fails, it releases previously acquired locks.
// Returns any errors encountered during the lock operations.
func doLock(readOnly bool, locks []Mutex) error {
	var i int

	for _, mu := range locks {
		if err := doLockOne(readOnly, mu); err != nil {
			// failed, log and reverse any previously acquired locks
			var errs core.CompoundError

			errs.AppendError(err)
			if err = reverseUnlock(readOnly, locks[:i]...); err != nil {
				errs.AppendError(err)
			}

			return errs.AsError()
		}

		// success, next
		i++
	}

	return nil
}

// doLockOne acquires a lock on a single mutex.
// If readOnly is true and the mutex implements RWMutex, it acquires a read lock.
// Returns any errors encountered during the lock operation.
func doLockOne(readOnly bool, mu Mutex) error {
	if mu == nil {
		return core.NewPanicError(2, "nil Mutex not allowed")
	}

	fn := makeLockFunc(readOnly, mu)
	return core.Catch(fn)
}

//revive:disable-next-line:flag-parameter
func makeLockFunc(readOnly bool, mu Mutex) func() error {
	if readOnly {
		if rw, is := mu.(*sync.RWMutex); is {
			// RLock()
			return func() error {
				rw.RLock()
				return nil
			}
		}
	}

	// Lock()
	return func() error {
		mu.Lock()
		return nil
	}
}

// RLock acquires multiple read locks simultaneously.
// For RWMutex instances, it uses the RLock method instead of Lock.
// For regular Mutex instances, it falls back to Lock.
//
// It acquires locks in the order provided and blocks until all locks are acquired.
// If any lock operation fails or panics, previously acquired locks are released
// in reverse order to prevent deadlocks.
//
// This function will panic:
// 1. If any of the provided mutexes are nil
// 2. If any mutex panics during lock operations (after releasing locks)
//
// Internally uses doLock with readOnly=true to acquire read locks.
func RLock(locks ...Mutex) {
	if err := doLock(true, locks); err != nil {
		panic(err)
	}
}

// interface assertions
var _ Mutex = (*readOnlyMutex)(nil)
var _ Mutex = (*sync.Mutex)(nil)
var _ Mutex = (*sync.RWMutex)(nil)
var _ RWMutex = (*sync.RWMutex)(nil)
