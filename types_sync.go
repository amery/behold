package behold

import "sync"

// Mutex defines a standard interface for mutual exclusion locking mechanisms
// that support basic locking, unlocking, and non-blocking lock attempts.
type Mutex interface {
	Lock()
	TryLock() bool
	Unlock()
}

// RWMutex extends the Mutex interface with read-locking capabilities,
// allowing multiple readers or a single writer to access a shared resource.
type RWMutex interface {
	Mutex

	RLock()
	RUnlock()
	TryRLock() bool
}

// ROMutex converts an RWMutex to a read-only Mutex, allowing only read locking operations.
// If the input mutex is nil, it returns nil.
func ROMutex(m RWMutex) Mutex {
	if m == nil {
		return nil
	}

	return &readOnlyMutex{m: m}
}

// interface assertions
var _ Mutex = (*readOnlyMutex)(nil)
var _ Mutex = (*sync.Mutex)(nil)
var _ Mutex = (*sync.RWMutex)(nil)
var _ RWMutex = (*sync.RWMutex)(nil)

// readOnlyMutex uses a RWMutex as a Mutex in read-only mode
type readOnlyMutex struct {
	m RWMutex
}

func (m readOnlyMutex) Lock()         { m.m.RLock() }
func (m readOnlyMutex) Unlock()       { m.m.RUnlock() }
func (m readOnlyMutex) TryLock() bool { return m.m.TryRLock() }
