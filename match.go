package behold

import "darvaza.org/core"

// Matcher is a generic interface for filtering and combining predicates of type T.
// It allows logical AND and OR operations between conditions, and matching against a value.
type Matcher[T any] interface {
	// And combines this query with others using logical AND.
	// All conditions must match for the combined query to match.
	And(...Matcher[T]) Matcher[T]

	// Or combines this query with others using logical OR.
	// At least one condition must match for the combined query to match.
	Or(...Matcher[T]) Matcher[T]

	// Match tests if the given value satisfies this query's conditions.
	Match(T) bool
}

// MatchFunc is a function type that implements the Matcher interface.
// It allows simple functions to be used as matchers.
type MatchFunc[T any] func(T) bool

// And combines this query function with others using logical AND.
func (fn MatchFunc[T]) And(others ...Matcher[T]) Matcher[T] {
	return ands[T](qJoin(fn, others))
}

// Or combines this query function with others using logical OR.
func (fn MatchFunc[T]) Or(others ...Matcher[T]) Matcher[T] {
	return ors[T](qJoin(fn, others))
}

// Match calls the query function with the provided value.
// If the function is nil, it returns true (matches everything).
func (fn MatchFunc[T]) Match(value T) bool {
	if fn == nil {
		return true
	}
	return fn(value)
}

// MatchAny returns a query that matches if any of the provided queries match.
// If no queries are provided, the result will match nothing (return false).
// Nil queries in the provided list are ignored during matching.
func MatchAny[T any](queries ...Matcher[T]) Matcher[T] {
	return ors[T](queries)
}

// MatchAll returns a query that matches if all of the provided queries match.
// If no queries are provided, the result will match everything (return true).
// Nil queries in the provided list are ignored during matching.
func MatchAll[T any](queries ...Matcher[T]) Matcher[T] {
	return ands[T](queries)
}

// ComposeMatch creates a new Matcher by applying an accessor function to transform input values
// before matching against an existing matcher. It allows composing matchers on different types
// by first extracting a specific field or transforming the input. Panics if the accessor
// function or the base query is nil.
func ComposeMatch[T any, V any](fn func(T) (V, bool), match Matcher[V]) Matcher[T] {
	if fn == nil {
		panic(core.NewPanicError(1, "nil accessor function"))
	}

	if match == nil {
		panic(core.NewPanicError(1, "no match condition"))
	}

	return MatchFunc[T](func(x T) bool {
		if v, ok := fn(x); ok {
			return match.Match(v)
		}
		return false
	})
}

// ands is a slice of matchers that implements the Matcher interface with AND logic.
type ands[T any] []Matcher[T]

// Match returns true if all non-nil queries in the slice match the provided value.
// An empty slice matches everything (returns true).
func (c ands[T]) Match(value T) bool {
	for _, q := range c {
		if q != nil && !q.Match(value) {
			return false
		}
	}
	return true
}

// And combines this AND query with others using logical AND.
func (c ands[T]) And(others ...Matcher[T]) Matcher[T] {
	return append(c, others...)
}

// Or combines this AND query with others using logical OR.
func (c ands[T]) Or(others ...Matcher[T]) Matcher[T] {
	return ors[T](qJoin(c, others))
}

// ors is a slice of matchers that implements the Matcher interface with OR logic.
type ors[T any] []Matcher[T]

// Match returns true if any non-nil query in the slice matches the provided value.
// An empty slice matches nothing (returns false).
func (c ors[T]) Match(value T) bool {
	for _, q := range c {
		if q != nil && q.Match(value) {
			return true
		}
	}
	return false
}

// And combines this OR query with others using logical AND.
func (c ors[T]) And(others ...Matcher[T]) Matcher[T] {
	return ands[T](qJoin(c, others))
}

// Or combines this OR query with others using logical OR.
func (c ors[T]) Or(others ...Matcher[T]) Matcher[T] {
	return append(c, others...)
}

// qJoin combines a query with a slice of other queries into a single slice.
// If the first query is nil, it simply returns the others slice.
func qJoin[T any](fn Matcher[T], others []Matcher[T]) []Matcher[T] {
	if fn == nil {
		return others
	}
	return append([]Matcher[T]{fn}, others...)
}
