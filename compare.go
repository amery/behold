package behold

import "darvaza.org/core"

// CompFunc is a generic comparison function that takes two values of type T and returns an integer.
// The return value follows the standard comparison convention:
// - Negative value if a < b
// - Zero if a == b
// - Positive value if a > b
type CompFunc[T any] func(a, b T) int

// Eq returns true if a is equal to b for comparable types.
func Eq[T comparable](a, b T) bool {
	return a == b
}

// EqFn returns true if a is equal to b using a custom comparison function.
// It panics if the provided comparison function is nil.
func EqFn[T any](a, b T, cmp CompFunc[T]) bool {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return cmp(a, b) == 0
}

// NotEq returns true if a is not equal to b for comparable types.
func NotEq[T comparable](a, b T) bool {
	return a != b
}

// NotEqFn returns true if a is not equal to b using a custom comparison function.
// It panics if the provided comparison function is nil.
func NotEqFn[T any](a, b T, cmp CompFunc[T]) bool {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return cmp(a, b) != 0
}

// Gt returns true if a is greater than b for ordered types.
func Gt[T core.Ordered](a, b T) bool {
	return a > b
}

// GtFn returns true if a is greater than b using a custom comparison function.
// It panics if the provided comparison function is nil.
func GtFn[T any](a, b T, cmp CompFunc[T]) bool {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return cmp(a, b) > 0
}

// GtEq returns true if a is greater than or equal to b for ordered types.
func GtEq[T core.Ordered](a, b T) bool {
	return a >= b
}

// GtEqFn returns true if a is greater than or equal to b using a custom comparison function.
// It panics if the provided comparison function is nil.
func GtEqFn[T any](a, b T, cmp CompFunc[T]) bool {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return cmp(a, b) >= 0
}

// Lt returns true if a is less than b for ordered types.
func Lt[T core.Ordered](a, b T) bool {
	return a < b
}

// LtFn returns true if a is less than b using a custom comparison function.
// It panics if the provided comparison function is nil.
func LtFn[T any](a, b T, cmp CompFunc[T]) bool {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return cmp(a, b) < 0
}

// LtEq returns true if a is less than or equal to b for ordered types.
func LtEq[T core.Ordered](a, b T) bool {
	return a <= b
}

// LtEqFn returns true if a is less than or equal to b using a custom comparison function.
// It panics if the provided comparison function is nil.
func LtEqFn[T any](a, b T, cmp CompFunc[T]) bool {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}
	return cmp(a, b) <= 0
}

func newNilCompFuncErr() error {
	return core.NewPanicError(2, "nil comparison function")
}
