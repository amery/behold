package behold

import "darvaza.org/core"

// ComposeQuery creates a new Query by applying an accessor function to transform input values
// before matching against an existing query. It allows composing queries on different types
// by first extracting a specific field or transforming the input. Panics if the accessor
// function or the base query is nil.
func ComposeQuery[T any, V any](fn func(T) V, query Query[V]) Query[T] {
	if fn == nil {
		panic(core.NewPanicError(1, "nil accessor function"))
	}

	if query == nil {
		panic(core.NewPanicError(1, "nil value query"))
	}

	return QueryFunc[T](func(x T) bool {
		return query.Match(fn(x))
	})
}

// EqQuery creates a Query that checks for equality with the given value.
// It returns a function that returns true if the input is equal to the specified value.
func EqQuery[T comparable](v T) Query[T] {
	return QueryFunc[T](func(v0 T) bool {
		return Eq(v0, v)
	})
}

// EqQueryFn creates a Query that checks for equality using a custom comparison function.
// It returns a function that returns true if the input is equal to the specified value
// according to the provided comparison function. Panics if the comparison function is nil.
func EqQueryFn[T any](v T, cmp CompFunc[T]) Query[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return EqFn(v0, v, cmp)
	})
}

// EqQueryFn2 creates a Query that checks for equality using a custom equality function.
// It returns a function that returns true if the input is equal to the specified value
// according to the provided equality function. Panics if the equality function is nil.
func EqQueryFn2[T any](v T, eq CondFunc[T]) Query[T] {
	if eq == nil {
		panic(newNilCondFuncErr())
	}

	return QueryFunc[T](func(v0 T) bool {
		return EqFn2(v0, v, eq)
	})
}

// NotEqQuery creates a Query that checks for inequality with the given value.
// It returns a function that returns true if the input is not equal to the specified value.
func NotEqQuery[T comparable](v T) Query[T] {
	return QueryFunc[T](func(v0 T) bool {
		return NotEq(v0, v)
	})
}

// NotEqQueryFn creates a Query that checks for inequality using a custom comparison function.
// It returns a function that returns true if the input is not equal to the specified value
// according to the provided comparison function. Panics if the comparison function is nil.
func NotEqQueryFn[T any](v T, cmp CompFunc[T]) Query[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return NotEqFn(v0, v, cmp)
	})
}

// NotEqQueryFn2 creates a Query that checks for inequality using a custom equality function.
// It returns a function that returns true if the input is not equal to the specified value
// according to the provided equality function. Panics if the equality function is nil.
func NotEqQueryFn2[T any](v T, eq CondFunc[T]) Query[T] {
	if eq == nil {
		panic(newNilCondFuncErr())
	}

	return QueryFunc[T](func(v0 T) bool {
		return NotEqFn2(v0, v, eq)
	})
}

// GtQuery creates a Query that checks if a value is strictly greater than the given value.
// It returns a function that returns true if the input is greater than the specified value.
func GtQuery[T core.Ordered](v T) Query[T] {
	return QueryFunc[T](func(v0 T) bool {
		return Gt(v0, v)
	})
}

// GtQueryFn creates a Query that checks if a value is strictly greater than the given value
// using a custom comparison function. It returns a function that returns true if the input
// is greater than the specified value according to the provided comparison function.
// Panics if the comparison function is nil.
func GtQueryFn[T core.Ordered](v T, cmp CompFunc[T]) Query[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return GtFn(v0, v, cmp)
	})
}

// GtEqQuery creates a Query that checks if a value is greater than or equal to the given value.
// It returns a function that returns true if the input is greater than or equal to the specified value.
func GtEqQuery[T core.Ordered](v T) Query[T] {
	return QueryFunc[T](func(v0 T) bool {
		return GtEq(v0, v)
	})
}

// GtEqQueryFn creates a Query that checks if a value is greater than or equal to the given value
// using a custom comparison function. It returns a function that returns true if the input
// is greater than or equal to the specified value according to the provided comparison function.
// Panics if the comparison function is nil.
func GtEqQueryFn[T core.Ordered](v T, cmp CompFunc[T]) Query[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return GtEqFn(v0, v, cmp)
	})
}

// GtEqQueryFn2 creates a Query that checks if a value is greater than or equal to the given value
// using a custom condition function. It returns a function that returns true if the input
// is greater than or equal to the specified value according to the provided condition function.
// Panics if the condition function is nil.
func GtEqQueryFn2[T core.Ordered](v T, less CondFunc[T]) Query[T] {
	if less == nil {
		panic(newNilCondFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return GtEqFn2(v0, v, less)
	})
}

// LtQuery creates a Query that checks if a value is strictly less than the given value.
// It returns a function that returns true if the input is less than the specified value.
func LtQuery[T core.Ordered](v T) Query[T] {
	return QueryFunc[T](func(v0 T) bool {
		return Lt(v0, v)
	})
}

// LtQueryFn creates a Query that checks if a value is strictly less than the given value
// using a custom comparison function. It returns a function that returns true if the input
// is less than the specified value according to the provided comparison function.
// Panics if the comparison function is nil.
func LtQueryFn[T core.Ordered](v T, cmp CompFunc[T]) Query[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return LtFn(v0, v, cmp)
	})
}

// LtQueryFn2 creates a Query that checks if a value is strictly less than the given value
// using a custom condition function. It returns a function that returns true if the input
// is less than the specified value according to the provided condition function.
// Panics if the condition function is nil.
func LtQueryFn2[T core.Ordered](v T, less CondFunc[T]) Query[T] {
	if less == nil {
		panic(newNilCondFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return LtFn2(v0, v, less)
	})
}

// LtEqQuery creates a Query that checks if a value is less than or equal to the given value.
// It returns a function that returns true if the input is less than or equal to the specified value.
func LtEqQuery[T core.Ordered](v T) Query[T] {
	return QueryFunc[T](func(v0 T) bool {
		return LtEq(v0, v)
	})
}

// LtEqQueryFn creates a Query that checks if a value is less than or equal to the given value
// using a custom comparison function. It returns a function that returns true if the input
// is less than or equal to the specified value according to the provided comparison function.
// Panics if the comparison function is nil.
func LtEqQueryFn[T core.Ordered](v T, cmp CompFunc[T]) Query[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return QueryFunc[T](func(v0 T) bool {
		return LtEqFn(v0, v, cmp)
	})
}
