package behold

import (
	"context"
	"time"
)

// Store represents a generic key-value store with versioning and transaction support.
// It provides a clean interface for working with data in a consistent manner.
//
// Type Parameters:
//   - K comparable: The key type, which must be comparable (supports == and != operators)
//   - V any: The value type, which can be any type
type Store[K comparable, V any] interface {
	// Version returns the current version of the data in the Store.
	Version() uint64

	// Now returns the store's current time reference
	Now() time.Time

	// View executes a read-only transaction with optional mutex locks
	// The provided function will be called with a transaction object that
	// can be used to access data in the store
	View(ctx context.Context, fn func(Tx[K, V]) error, locks ...Mutex) error

	// Update executes a read-write transaction with optional mutex locks
	// The provided function will be called with a transaction object that
	// can be used to access and modify data in the store
	Update(ctx context.Context, fn func(Tx[K, V]) error, locks ...Mutex) error

	// Close closes the store and releases its resources
	Close() error
}

// Tx represents a transaction within the key-value store.
// It provides methods for reading, writing, and managing data within a transactional context.
//
// Type Parameters:
//   - K comparable: The key type, matching the store's key type
//   - V any: The value type, matching the store's value type
type Tx[K comparable, V any] interface {
	// Context returns the transaction's context
	Context() context.Context

	// Version returns the data version accessed by this transaction.
	Version() uint64

	// Now returns the transaction's time reference
	Now() time.Time

	// ForEach iterates through key-value pairs matching optional queries
	// The provided function is called for entries matching any of the queries,
	// or for every entry if no queries are provided.
	// Iteration can be ended early by returning false from the callback.
	ForEach(fn func(key K, value V) bool, ors ...Query[any]) error

	// Get retrieves a value by key
	Get(key K) (value V, err error)

	// Set associates a value with a key
	Set(key K, value V) error

	// Append adds a value to an existing key (implementation depends on value type)
	Append(key K, value V) error

	// Delete removes a key-value pair
	Delete(key K) error

	// Commit persists changes made within the transaction
	Commit() error

	// Close aborts the transaction if not already committed
	Close() error
}
