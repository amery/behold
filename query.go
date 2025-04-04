package behold

// Query is a generic interface for filtering and combining predicates of type T.
// It allows logical AND and OR operations between query conditions, and matching against a value.
type Query[T any] interface {
	// And combines this query with others using logical AND.
	// All conditions must match for the combined query to match.
	And(...Query[T]) Query[T]

	// Or combines this query with others using logical OR.
	// At least one condition must match for the combined query to match.
	Or(...Query[T]) Query[T]

	// Match tests if the given value satisfies this query's conditions.
	Match(T) bool
}

// QueryFunc is a function type that implements the Query interface.
// It allows simple functions to be used as queries.
type QueryFunc[T any] func(T) bool

// And combines this query function with others using logical AND.
func (fn QueryFunc[T]) And(others ...Query[T]) Query[T] {
	return ands[T](qJoin(fn, others))
}

// Or combines this query function with others using logical OR.
func (fn QueryFunc[T]) Or(others ...Query[T]) Query[T] {
	return ors[T](qJoin(fn, others))
}

// Match calls the query function with the provided value.
// If the function is nil, it returns true (matches everything).
func (fn QueryFunc[T]) Match(value T) bool {
	if fn == nil {
		return true
	}
	return fn(value)
}

// MatchAny returns a query that matches if any of the provided queries match.
// If no queries are provided, the result will match nothing (return false).
// Nil queries in the provided list are ignored during matching.
func MatchAny[T any](queries ...Query[T]) Query[T] {
	return ors[T](queries)
}

// MatchAll returns a query that matches if all of the provided queries match.
// If no queries are provided, the result will match everything (return true).
// Nil queries in the provided list are ignored during matching.
func MatchAll[T any](queries ...Query[T]) Query[T] {
	return ands[T](queries)
}

// ands is a slice of queries that implements the Query interface with AND logic.
type ands[T any] []Query[T]

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
func (c ands[T]) And(others ...Query[T]) Query[T] {
	return append(c, others...)
}

// Or combines this AND query with others using logical OR.
func (c ands[T]) Or(others ...Query[T]) Query[T] {
	return ors[T](qJoin(c, others))
}

// ors is a slice of queries that implements the Query interface with OR logic.
type ors[T any] []Query[T]

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
func (c ors[T]) And(others ...Query[T]) Query[T] {
	return ands[T](qJoin(c, others))
}

// Or combines this OR query with others using logical OR.
func (c ors[T]) Or(others ...Query[T]) Query[T] {
	return append(c, others...)
}

// qJoin combines a query with a slice of other queries into a single slice.
// If the first query is nil, it simply returns the others slice.
func qJoin[T any](fn Query[T], others []Query[T]) []Query[T] {
	if fn == nil {
		return others
	}
	return append([]Query[T]{fn}, others...)
}
