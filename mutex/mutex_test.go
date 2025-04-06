package mutex

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func doSomething() {
	time.Sleep(10 * time.Millisecond)
}

// TestROMutex tests the read-only mutex adapter functionality
func TestROMutex(t *testing.T) {
	// Test nil case
	rom := ROMutex(nil)
	assert.Nil(t, rom, "ROMutex(nil) should return nil")

	// Test normal RWMutex conversion
	rwm := &sync.RWMutex{}
	rom = ROMutex(rwm)
	assert.NotNil(t, rom, "ROMutex should not return nil for valid input")

	// Test that the ROMutex properly works as a read-only lock
	acquired := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	// First acquire a read lock through the read-only adapter
	rom.Lock()

	// Verify we can acquire another read lock on the same underlying RWMutex
	rwm.RLock()
	doSomething()
	rwm.RUnlock()

	// Try to acquire a write lock in a goroutine (should block)
	go func() {
		defer wg.Done()

		// Try to acquire write lock (should block)
		rwm.Lock()
		acquired <- true
		rwm.Unlock()
	}()

	// Give the goroutine time to potentially acquire the lock (should not succeed)
	time.Sleep(50 * time.Millisecond)

	select {
	case <-acquired:
		assert.Fail(t, "Write lock should not be acquired while read lock is held")
	default:
		// This is the expected case - the channel should be empty
	}

	// Release the read lock
	rom.Unlock()

	// Now the write lock should be acquirable
	wg.Wait()

	select {
	case <-acquired:
		// This is the expected case - the write lock was acquired
	default:
		assert.Fail(t, "Write lock should have been acquired after read lock was released")
	}
}

// TestTryLock tests the TryLock function with multiple mutexes
func TestTryLock(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}
	m3 := &sync.Mutex{}

	// Test successful acquisition of all locks
	assert.True(t, TryLock(m1, m2, m3), "TryLock should succeed on unlocked mutexes")

	// Locks should be acquired, so release them
	Unlock(m1, m2, m3)

	// Test with one mutex already locked
	m1.Lock()
	assert.False(t, TryLock(m1, m2, m3), "TryLock should fail when one mutex is already locked")

	// Verify m2 and m3 are not locked
	assert.True(t, m2.TryLock(), "m2 should not be locked after TryLock failure")
	m2.Unlock()

	assert.True(t, m3.TryLock(), "m3 should not be locked after TryLock failure")
	m3.Unlock()

	m1.Unlock()

	// Test with nil mutex (should panic)
	assert.Panics(t, func() {
		TryLock(m1, nil, m3)
	}, "TryLock should panic with nil mutex")
}

// TestTryRLock tests the TryRLock function with multiple mutexes
func TestTryRLock(t *testing.T) {
	m1 := &sync.RWMutex{}
	m2 := &sync.RWMutex{}
	regularMutex := &sync.Mutex{}

	// Test successful acquisition of all locks
	assert.True(t, TryRLock(m1, m2, regularMutex), "TryRLock should succeed on unlocked mutexes")

	// Locks should be acquired, so release them
	RUnlock(m1, m2, regularMutex)

	// Test with regular mutex already locked
	regularMutex.Lock()
	assert.False(t, TryRLock(m1, m2, regularMutex), "TryRLock should fail when regular mutex is already locked")
	regularMutex.Unlock()

	// Test with RWMutex write-locked
	m1.Lock()
	assert.False(t, TryRLock(m1, m2, regularMutex), "TryRLock should fail when RWMutex is write-locked")
	m1.Unlock()

	// Test concurrent read locks on same RWMutex
	m1.RLock()
	assert.True(t, TryRLock(m1, m2, regularMutex), "TryRLock should succeed even if read lock already held")
	RUnlock(m1, m2, regularMutex)
	m1.RUnlock()

	// Test with nil mutex (should panic)
	assert.Panics(t, func() {
		TryRLock(m1, nil, regularMutex)
	}, "TryRLock should panic with nil mutex")
}

// TestLock tests the Lock function with multiple mutexes
func TestLock(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}

	// Test basic acquisition of locks
	lockAcquired := make(chan bool)
	go func() {
		Lock(m1, m2)
		lockAcquired <- true
	}()

	// Should acquire locks quickly
	select {
	case <-lockAcquired:
		// Expected case - locks were acquired
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Lock should have acquired locks quickly")
	}

	// Clean up
	Unlock(m1, m2)

	// Test blocking behavior when one lock is held
	m1.Lock()

	lockBlocked := make(chan bool)
	lockReleased := make(chan bool)

	go func() {
		lockBlocked <- true
		Lock(m1, m2) // This should block until m1 is released
		lockReleased <- true
	}()

	// Wait for goroutine to start attempting to lock
	<-lockBlocked

	// Should not acquire locks while m1 is held
	select {
	case <-lockReleased:
		assert.Fail(t, "Lock should block when one mutex is already locked")
	case <-time.After(50 * time.Millisecond):
		// Expected case - lock is blocked
	}

	// Release m1 and the Lock call should complete
	m1.Unlock()

	select {
	case <-lockReleased:
		// Expected case - locks were acquired after m1 was released
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Lock should have acquired locks after m1 was released")
	}

	// Clean up
	Unlock(m1, m2)

	// Test with nil mutex (should panic)
	assert.Panics(t, func() {
		Lock(m1, nil)
	}, "Lock should panic with nil mutex")
}

// TestRLock tests the RLock function with multiple mutexes
func TestRLock(t *testing.T) {
	m1 := &sync.RWMutex{}
	m2 := &sync.RWMutex{}
	regularMutex := &sync.Mutex{}

	// Test basic acquisition
	RLock(m1, m2, regularMutex)

	// Should be able to acquire another read lock on the same RWMutexes
	m1.RLock()
	doSomething()
	m1.RUnlock()

	// But should not be able to acquire write lock
	assert.False(t, m1.TryLock(), "Should not be able to acquire write lock while read lock is held")

	// Clean up
	RUnlock(m1, m2, regularMutex)

	// Test blocking behavior with write lock
	m1.Lock()

	rlockBlocked := make(chan bool)
	rlockReleased := make(chan bool)

	go func() {
		rlockBlocked <- true
		RLock(m1, m2, regularMutex) // Should block until m1 write lock is released
		rlockReleased <- true
	}()

	// Wait for goroutine to start attempting to lock
	<-rlockBlocked

	// Should not acquire locks while m1 write lock is held
	select {
	case <-rlockReleased:
		assert.Fail(t, "RLock should block when write lock is held")
	case <-time.After(50 * time.Millisecond):
		// Expected case - rlock is blocked
	}

	// Release m1 and the RLock call should complete
	m1.Unlock()

	select {
	case <-rlockReleased:
		// Expected case - locks were acquired after m1 was released
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "RLock should have acquired locks after m1 was released")
	}

	// Clean up
	RUnlock(m1, m2, regularMutex)

	// Test with nil mutex (should panic)
	assert.Panics(t, func() {
		RLock(m1, nil, regularMutex)
	}, "RLock should panic with nil mutex")
}

// TestUnlock tests the Unlock function with multiple mutexes
func TestUnlock(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}

	// Lock mutexes first
	m1.Lock()
	m2.Lock()

	// Test unlock
	Unlock(m1, m2)

	// Verify they are unlocked by trying to lock them again
	assert.True(t, m1.TryLock(), "m1 should be unlocked after Unlock")
	m1.Unlock()

	assert.True(t, m2.TryLock(), "m2 should be unlocked after Unlock")
	m2.Unlock()

	// Create a custom mutex that panics on Unlock
	panicMutex := &panicOnUnlockMutex{}

	// Test with mutex that panics on unlock
	assert.Panics(t, func() {
		Unlock(panicMutex)
	}, "Unlock should panic when a mutex panics during unlock")
}

// panicOnUnlockMutex is a custom mutex that panics when Unlock is called
type panicOnUnlockMutex struct{}

func (*panicOnUnlockMutex) Lock()         {}
func (*panicOnUnlockMutex) TryLock() bool { return true }
func (*panicOnUnlockMutex) Unlock()       { panic("intentional panic on unlock") }

// TestRUnlock tests the RUnlock function with multiple mutexes
func TestRUnlock(t *testing.T) {
	m1 := &sync.RWMutex{}
	m2 := &sync.RWMutex{}
	regularMutex := &sync.Mutex{}

	// Lock mutexes first
	m1.RLock()
	m2.RLock()
	regularMutex.Lock()

	// Test RUnlock
	RUnlock(m1, m2, regularMutex)

	// Verify they are unlocked by trying to write-lock the RWMutexes
	assert.True(t, m1.TryLock(), "m1 should be read-unlocked after RUnlock")
	m1.Unlock()

	assert.True(t, m2.TryLock(), "m2 should be read-unlocked after RUnlock")
	m2.Unlock()

	// Verify regular mutex is unlocked
	assert.True(t, regularMutex.TryLock(), "regularMutex should be unlocked after RUnlock")
	regularMutex.Unlock()

	// Create a custom mutex that panics on Unlock
	panicMutex := &panicOnRUnlockMutex{}

	// Test with mutex that panics on unlock
	assert.Panics(t, func() {
		RUnlock(panicMutex)
	}, "RUnlock should panic when a mutex panics during unlock")
}

// panicOnRUnlockMutex is needed for the test but wasn't defined in the original code
type panicOnRUnlockMutex struct{}

func (*panicOnRUnlockMutex) Lock()         {}
func (*panicOnRUnlockMutex) TryLock() bool { return true }
func (*panicOnRUnlockMutex) Unlock()       { panic("intentional panic on unlock") }

// TestReadOnlyMutexImplementation tests the readOnlyMutex directly
func TestReadOnlyMutexImplementation(t *testing.T) {
	rwm := &sync.RWMutex{}
	rom := &readOnlyMutex{m: rwm}

	// Test locking
	rom.Lock()

	// Test that we have a read lock, not a write lock
	rwm.RLock() // Should not block if we have a read lock
	doSomething()
	rwm.RUnlock()

	assert.False(t, rwm.TryLock(), "Should not be able to acquire write lock while read lock is held")

	// Test unlocking
	rom.Unlock()

	// Verify we can now acquire a write lock
	assert.True(t, rwm.TryLock(), "Should be able to acquire write lock after read lock is released")
	rwm.Unlock()

	// Test TryLock when unlocked
	assert.True(t, rom.TryLock(), "TryLock should succeed on unlocked mutex")
	rom.Unlock()

	// Test TryLock when write-locked
	rwm.Lock()
	assert.False(t, rom.TryLock(), "TryLock should fail when mutex is write-locked")
	rwm.Unlock()
}

// TestConcurrentLocking tests concurrent locking behavior
func TestConcurrentLocking(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}
	const goroutines = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			// Try to lock both mutexes
			if TryLock(m1, m2) {
				// Hold the locks briefly
				time.Sleep(1 * time.Millisecond)
				// Release the locks
				Unlock(m1, m2)
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify both mutexes are unlocked at the end
	assert.True(t, m1.TryLock(), "m1 should be unlocked after concurrent operations")
	m1.Unlock()

	assert.True(t, m2.TryLock(), "m2 should be unlocked after concurrent operations")
	m2.Unlock()
}

// TestErrorHandling tests error cases when locking/unlocking mutexes
func TestErrorHandling(t *testing.T) {
	// Test internal functions directly

	// Test doTryLockOne with nil mutex
	_, err := doTryLockOne(false, nil)
	assert.Error(t, err, "doTryLockOne should return error with nil mutex")

	// Test doUnlockOne with nil mutex
	err = doUnlockOne(false, nil)
	assert.Error(t, err, "doUnlockOne should return error with nil mutex")

	// Test doLockOne with nil mutex
	err = doLockOne(false, nil)
	assert.Error(t, err, "doLockOne should return error with nil mutex")
}

// TestReverseUnlock tests that locks are released in reverse order on failure
func TestReverseUnlock(t *testing.T) {
	// Create a channel to track unlock order
	unlockOrder := make(chan int, 3)

	// Create custom mutex implementations to track unlock order
	customM1 := &testMutex{id: 1, unlockOrder: unlockOrder}
	customM2 := &testMutex{id: 2, unlockOrder: unlockOrder}
	customM3 := &testMutex{id: 3, unlockOrder: unlockOrder}

	// Set up m3 to fail when TryLock is called
	customM3.shouldFailTryLock = true

	// Try to lock all three (should fail on m3)
	assert.False(t, TryLock(customM1, customM2, customM3), "TryLock should fail when one mutex is configured to fail")

	// Check the unlock order (should be 2, 1 - reverse order)
	// We should have exactly 2 unlocks (not 3, since m3's lock failed)
	assert.Equal(t, 2, len(unlockOrder), "Expected 2 unlocks")

	order1 := <-unlockOrder
	assert.Equal(t, 2, order1, "Expected first unlock to be mutex 2")

	order2 := <-unlockOrder
	assert.Equal(t, 1, order2, "Expected second unlock to be mutex 1")
}

// testMutex is a custom implementation of Mutex for testing unlock order
type testMutex struct {
	id                int
	unlockOrder       chan<- int
	locked            bool
	shouldFailTryLock bool
}

func (m *testMutex) Lock() {
	m.locked = true
}

func (m *testMutex) Unlock() {
	if !m.locked {
		panic("unlock of unlocked mutex")
	}
	m.unlockOrder <- m.id
	m.locked = false
}

func (m *testMutex) TryLock() bool {
	if m.shouldFailTryLock {
		return false
	}
	if m.locked {
		return false
	}
	m.locked = true
	return true
}

// TestNilImplementations tests interface implementations for nil values
func TestNilImplementations(t *testing.T) {
	// Verify various nil behaviors

	// Nil ROMutex
	assert.Nil(t, ROMutex(nil), "ROMutex(nil) should return nil")

	// Test panic recovery for nil mutexes
	assert.Panics(t, func() {
		Lock(nil)
	}, "Lock(nil) should panic")

	assert.Panics(t, func() {
		RLock(nil)
	}, "RLock(nil) should panic")

	assert.Panics(t, func() {
		Unlock(nil)
	}, "Unlock(nil) should panic")

	assert.Panics(t, func() {
		RUnlock(nil)
	}, "RUnlock(nil) should panic")
}
