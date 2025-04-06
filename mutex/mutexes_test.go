package mutex

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMutexesLock tests the Lock method of Mutexes
func TestMutexesLock(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}
	mutexes := Mutexes{m1, m2}

	// Test basic acquisition of locks
	lockAcquired := make(chan bool, 1)
	go func() {
		mutexes.Lock()
		lockAcquired <- true
	}()

	// Should acquire locks quickly
	select {
	case <-lockAcquired:
		// Expected case - locks were acquired
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Mutexes.Lock should have acquired locks quickly")
	}

	// Verify mutexes are locked by trying to lock them directly
	assert.False(t, m1.TryLock(), "m1 should be locked after Mutexes.Lock")
	assert.False(t, m2.TryLock(), "m2 should be locked after Mutexes.Lock")

	// Clean up
	mutexes.Unlock()
}

// TestMutexesLockBlocking tests the blocking behavior of Mutexes.Lock
func TestMutexesLockBlocking(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}
	mutexes := Mutexes{m1, m2}

	// Test with one mutex already locked
	m1.Lock()

	mutexBlocked := make(chan bool, 1)
	mutexReleased := make(chan bool, 1)

	go func() {
		mutexBlocked <- true
		mutexes.Lock() // This should block until m1 is released
		mutexReleased <- true
	}()

	// Wait for goroutine to start attempting to lock
	<-mutexBlocked

	// Should not acquire locks while m1 is held
	select {
	case <-mutexReleased:
		assert.Fail(t, "Mutexes.Lock should block when one mutex is already locked")
	case <-time.After(50 * time.Millisecond):
		// Expected case - lock is blocked
	}

	// Release m1 and the Lock call should complete
	m1.Unlock()

	select {
	case <-mutexReleased:
		// Expected case - locks were acquired after m1 was released
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Mutexes.Lock should have acquired locks after m1 was released")
	}

	// Clean up
	mutexes.Unlock()
}

// TestMutexesTryLock tests the TryLock method of Mutexes
func TestMutexesTryLock(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}
	mutexes := Mutexes{m1, m2}

	// Test successful acquisition of all locks
	assert.True(t, mutexes.TryLock(), "Mutexes.TryLock should succeed on unlocked mutexes")

	// Locks should be acquired, so verify by trying to lock them directly
	assert.False(t, m1.TryLock(), "m1 should be locked after Mutexes.TryLock")
	assert.False(t, m2.TryLock(), "m2 should be locked after Mutexes.TryLock")

	// Release locks
	mutexes.Unlock()

	// Test with one mutex already locked
	m1.Lock()
	assert.False(t, mutexes.TryLock(), "Mutexes.TryLock should fail when one mutex is already locked")

	// Verify m2 is not locked
	assert.True(t, m2.TryLock(), "m2 should not be locked after Mutexes.TryLock failure")
	m2.Unlock()

	m1.Unlock()
}

// TestMutexesUnlock tests the Unlock method of Mutexes
func TestMutexesUnlock(t *testing.T) {
	m1 := &sync.Mutex{}
	m2 := &sync.Mutex{}
	mutexes := Mutexes{m1, m2}

	// Lock mutexes first
	mutexes.Lock()

	doSomething()

	// Test unlock
	mutexes.Unlock()

	// Verify they are unlocked by trying to lock them again
	assert.True(t, m1.TryLock(), "m1 should be unlocked after Mutexes.Unlock")
	m1.Unlock()

	assert.True(t, m2.TryLock(), "m2 should be unlocked after Mutexes.Unlock")
	m2.Unlock()

	// Test with empty Mutexes
	emptyMutexes := Mutexes{}
	emptyMutexes.Lock() // Should not block or panic
	doSomething()
	emptyMutexes.Unlock() // Should not panic
}

// TestMutexesEmptyAndNil tests edge cases for Mutexes
func TestMutexesEmptyAndNil(t *testing.T) {
	// Empty slice should not panic
	emptyMutexes := Mutexes{}
	assert.NotPanics(t, func() {
		emptyMutexes.Lock()
		doSomething()
		emptyMutexes.Unlock()
	})

	assert.True(t, emptyMutexes.TryLock(), "TryLock on empty Mutexes should return true")
	emptyMutexes.Unlock()

	// Nil slice should behave like empty slice and not panic
	var nilMutexes Mutexes = nil
	assert.NotPanics(t, func() {
		nilMutexes.Lock()
		doSomething()
		nilMutexes.Unlock()
	})

	assert.True(t, nilMutexes.TryLock(), "TryLock on nil Mutexes should return true")
	nilMutexes.Unlock()
}

// TestRWMutexesExport tests the export method of RWMutexes
func TestRWMutexesExport(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	exported := rwMutexes.export()

	assert.Equal(t, 2, len(exported), "export should return slice of same length")
	assert.Same(t, rwm1, exported[0], "first element should be same object")
	assert.Same(t, rwm2, exported[1], "second element should be same object")
}

// TestRWMutexesLock tests the Lock method of RWMutexes
func TestRWMutexesLock(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Test basic acquisition of locks
	lockAcquired := make(chan bool, 1)
	go func() {
		rwMutexes.Lock()
		lockAcquired <- true
	}()

	// Should acquire locks quickly
	select {
	case <-lockAcquired:
		// Expected case - locks were acquired
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "RWMutexes.Lock should have acquired locks quickly")
	}

	// Verify mutexes are write-locked by trying to read-lock them directly
	assert.False(t, rwm1.TryRLock(), "rwm1 should be write-locked after RWMutexes.Lock")
	assert.False(t, rwm2.TryRLock(), "rwm2 should be write-locked after RWMutexes.Lock")

	// Clean up
	rwMutexes.Unlock()
}

// TestRWMutexesLockBlocking tests the blocking behavior of RWMutexes.Lock
func TestRWMutexesLockBlocking(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Test with one mutex already write-locked
	rwm1.Lock()

	rwmBlocked := make(chan bool, 1)
	rwmReleased := make(chan bool, 1)

	go func() {
		rwmBlocked <- true
		rwMutexes.Lock() // This should block until rwm1 is released
		rwmReleased <- true
	}()

	// Wait for goroutine to start attempting to lock
	<-rwmBlocked

	// Should not acquire locks while rwm1 is write-locked
	select {
	case <-rwmReleased:
		assert.Fail(t, "RWMutexes.Lock should block when one mutex is already write-locked")
	case <-time.After(50 * time.Millisecond):
		// Expected case - lock is blocked
	}

	// Release rwm1 and the Lock call should complete
	rwm1.Unlock()

	select {
	case <-rwmReleased:
		// Expected case - locks were acquired after rwm1 was released
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "RWMutexes.Lock should have acquired locks after rwm1 was released")
	}

	// Clean up
	rwMutexes.Unlock()
}

// TestRWMutexesTryLock tests the TryLock method of RWMutexes
func TestRWMutexesTryLock(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Test successful acquisition of all locks
	assert.True(t, rwMutexes.TryLock(), "RWMutexes.TryLock should succeed on unlocked mutexes")

	// Locks should be acquired, so verify by trying to lock them directly
	assert.False(t, rwm1.TryLock(), "rwm1 should be write-locked after RWMutexes.TryLock")
	assert.False(t, rwm2.TryLock(), "rwm2 should be write-locked after RWMutexes.TryLock")

	// Release locks
	rwMutexes.Unlock()

	// Test with one mutex already write-locked
	rwm1.Lock()
	assert.False(t, rwMutexes.TryLock(), "RWMutexes.TryLock should fail when one mutex is already write-locked")

	// Verify rwm2 is not locked
	assert.True(t, rwm2.TryLock(), "rwm2 should not be locked after RWMutexes.TryLock failure")
	rwm2.Unlock()

	rwm1.Unlock()
}

// TestRWMutexesUnlock tests the Unlock method of RWMutexes
func TestRWMutexesUnlock(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Lock mutexes first
	rwMutexes.Lock()

	doSomething()

	// Test unlock
	rwMutexes.Unlock()

	// Verify they are unlocked by trying to write-lock them again
	assert.True(t, rwm1.TryLock(), "rwm1 should be unlocked after RWMutexes.Unlock")
	rwm1.Unlock()

	assert.True(t, rwm2.TryLock(), "rwm2 should be unlocked after RWMutexes.Unlock")
	rwm2.Unlock()

	// Test with empty RWMutexes
	emptyRWMutexes := RWMutexes{}
	emptyRWMutexes.Lock() // Should not block or panic
	doSomething()
	emptyRWMutexes.Unlock() // Should not panic
}

// TestRWMutexesRLockBlocking tests the blocking behavior of RWMutexes.RLock
func TestRWMutexesRLockBlocking(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Test with one mutex already write-locked
	rwm1.Lock()

	rwmBlocked := make(chan bool, 1)
	rwmReleased := make(chan bool, 1)

	go func() {
		rwmBlocked <- true
		rwMutexes.RLock() // Should block until rwm1 write lock is released
		rwmReleased <- true
	}()

	// Wait for goroutine to start attempting to lock
	<-rwmBlocked

	// Should not acquire read locks while rwm1 write lock is held
	select {
	case <-rwmReleased:
		assert.Fail(t, "RWMutexes.RLock should block when write lock is held")
	case <-time.After(50 * time.Millisecond):
		// Expected case - rlock is blocked
	}

	// Release rwm1 and the RLock call should complete
	rwm1.Unlock()

	select {
	case <-rwmReleased:
		// Expected case - read locks were acquired after rwm1 was released
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "RWMutexes.RLock should have acquired locks after rwm1 was released")
	}

	// Clean up
	rwMutexes.RUnlock()
}

// TestRWMutexesTryRLock tests the TryRLock method of RWMutexes
func TestRWMutexesTryRLock(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Test successful acquisition of all read locks
	assert.True(t, rwMutexes.TryRLock(), "RWMutexes.TryRLock should succeed on unlocked mutexes")

	// Read locks should be acquired, write locks should fail
	assert.False(t, rwm1.TryLock(), "rwm1 should be read-locked after RWMutexes.TryRLock")
	assert.False(t, rwm2.TryLock(), "rwm2 should be read-locked after RWMutexes.TryRLock")

	// But additional read locks should still succeed
	assert.True(t, rwm1.TryRLock(), "Additional read locks should succeed")
	rwm1.RUnlock()

	// Release read locks
	rwMutexes.RUnlock()

	// Test with one mutex already write-locked
	rwm1.Lock()
	assert.False(t, rwMutexes.TryRLock(), "RWMutexes.TryRLock should fail when one mutex is already write-locked")

	rwm1.Unlock()
}

// TestRWMutexesRUnlock tests the RUnlock method of RWMutexes
func TestRWMutexesRUnlock(t *testing.T) {
	rwm1 := &sync.RWMutex{}
	rwm2 := &sync.RWMutex{}
	rwMutexes := RWMutexes{rwm1, rwm2}

	// Lock mutexes first
	rwMutexes.RLock()
	doSomething()

	// Test RUnlock
	rwMutexes.RUnlock()

	// Verify they are unlocked by trying to write-lock them now
	assert.True(t, rwm1.TryLock(), "rwm1 should be read-unlocked after RWMutexes.RUnlock")
	rwm1.Unlock()

	assert.True(t, rwm2.TryLock(), "rwm2 should be read-unlocked after RWMutexes.RUnlock")
	rwm2.Unlock()

	// Test with empty RWMutexes
	emptyRWMutexes := RWMutexes{}
	emptyRWMutexes.RLock() // Should not block or panic
	doSomething()
	emptyRWMutexes.RUnlock() // Should not panic
}

// TestRWMutexesEmptyAndNil tests edge cases for RWMutexes
func TestRWMutexesEmptyAndNil(t *testing.T) {
	// Empty slice should not panic
	emptyRWMutexes := RWMutexes{}
	assert.NotPanics(t, func() {
		emptyRWMutexes.Lock()
		doSomething()
		emptyRWMutexes.Unlock()

		emptyRWMutexes.RLock()
		doSomething()
		emptyRWMutexes.RUnlock()
	})

	assert.True(t, emptyRWMutexes.TryLock(), "TryLock on empty RWMutexes should return true")
	emptyRWMutexes.Unlock()

	assert.True(t, emptyRWMutexes.TryRLock(), "TryRLock on empty RWMutexes should return true")
	emptyRWMutexes.RUnlock()

	// Nil slice should behave like empty slice and not panic
	var nilRWMutexes RWMutexes = nil
	assert.NotPanics(t, func() {
		nilRWMutexes.Lock()
		doSomething()
		nilRWMutexes.Unlock()
		nilRWMutexes.RLock()
		doSomething()
		nilRWMutexes.RUnlock()
	})

	assert.True(t, nilRWMutexes.TryLock(), "TryLock on nil RWMutexes should return true")
	nilRWMutexes.Unlock()

	assert.True(t, nilRWMutexes.TryRLock(), "TryRLock on nil RWMutexes should return true")
	nilRWMutexes.RUnlock()
}
