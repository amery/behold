package behold

import "darvaza.org/core"

// CompFunc is a generic comparison function that takes two values of type T and returns an integer.
// The return value follows the standard comparison convention:
// - Negative value if a < b
// - Zero if a == b
// - Positive value if a > b
type CompFunc[T any] func(a, b T) int

// CondFunc is a generic condition function that takes two values of type T and returns a boolean.
// The return value indicates whether the condition is true or false for the given pair of values.
type CondFunc[T any] func(a, b T) bool

// AsLess converts a CompFunc into a less-than condition function.
// It returns a function that returns true if the first argument is less than the second argument.
// It panics if the provided comparison function is nil.
func AsLess[T any](cmp CompFunc[T]) CondFunc[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}

	return func(a, b T) bool {
		return cmp(a, b) < 0
	}
}

// AsEqual converts a CompFunc into an equality condition function.
// It returns a function that returns true if the first argument is equal to the second argument.
// It panics if the provided comparison function is nil.
func AsEqual[T any](cmp CompFunc[T]) CondFunc[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}

	return func(a, b T) bool {
		return cmp(a, b) == 0
	}
}

// Reverse returns a new CompFunc that inverts the comparison result of the given CompFunc.
// It returns a function that negates the original comparison, effectively reversing the order.
// It panics if the provided comparison function is nil.
func Reverse[T any](cmp CompFunc[T]) CompFunc[T] {
	if cmp == nil {
		panic(newNilCompFuncErr())
	}

	return func(a, b T) int {
		return -cmp(a, b)
	}
}

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

// EqFn2 returns true if a is equal to b using a custom equality function.
// It panics if the provided equality function is nil.
func EqFn2[T any](a, b T, eq CondFunc[T]) bool {
	if eq == nil {
		panic(newNilCondFuncErr())
	}
	return eq(a, b)
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

// NotEqFn2 returns true if a is not equal to b using a custom equality function.
// It panics if the provided equality function is nil.
func NotEqFn2[T any](a, b T, eq CondFunc[T]) bool {
	if eq == nil {
		panic(newNilCondFuncErr())
	}
	return !eq(a, b)
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

// GtEqFn2 returns true if a is greater than or equal to b using a custom less-than condition function.
// It panics if the provided less-than condition function is nil.
func GtEqFn2[T any](a, b T, less CondFunc[T]) bool {
	if less == nil {
		panic(newNilCondFuncErr())
	}

	return !less(a, b)
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

// LtFn2 returns true if a is less than b using a custom less-than condition function.
// It panics if the provided less-than condition function is nil.
func LtFn2[T any](a, b T, less CondFunc[T]) bool {
	if less == nil {
		panic(newNilCondFuncErr())
	}
	return less(a, b)
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

func newNilCondFuncErr() error {
	return core.NewPanicError(2, "nil condition function")
}
