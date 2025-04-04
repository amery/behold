package behold

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"equal values", 5, 5, true},
		{"different values", 5, 10, false},
		{"negative values equal", -5, -5, true},
		{"negative values different", -5, -10, false},
		{"zero and non-zero", 0, 5, false},
		{"zero and zero", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Eq(tt.a, tt.b))
		})
	}
}

func TestEqString(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{"equal strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"empty strings", "", "", true},
		{"empty and non-empty", "", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Eq(tt.a, tt.b))
		})
	}
}

func TestEqFn(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"equal values", 5, 5, true},
		{"different values", 5, 10, false},
		{"negative values equal", -5, -5, true},
		{"negative values different", -5, -10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, EqFn(tt.a, tt.b, cmp))
		})
	}
}

func TestEqFnPanic(t *testing.T) {
	assert.Panics(t, func() {
		EqFn(1, 2, nil)
	})
}

func TestNotEq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"equal values", 5, 5, false},
		{"different values", 5, 10, true},
		{"zero and non-zero", 0, 5, true},
		{"zero and zero", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, NotEq(tt.a, tt.b))
		})
	}
}

func TestNotEqFn(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"equal values", 5, 5, false},
		{"different values", 5, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, NotEqFn(tt.a, tt.b, cmp))
		})
	}
}

func TestNotEqFnPanic(t *testing.T) {
	assert.Panics(t, func() {
		NotEqFn(1, 2, nil)
	})
}

func TestGt(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"greater", 10, 5, true},
		{"less", 5, 10, false},
		{"equal", 5, 5, false},
		{"negative and positive", -5, 5, false},
		{"negative values", -5, -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Gt(tt.a, tt.b))
		})
	}
}

func TestGtFn(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"greater", 10, 5, true},
		{"less", 5, 10, false},
		{"equal", 5, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GtFn(tt.a, tt.b, cmp))
		})
	}
}

func TestGtFnPanic(t *testing.T) {
	assert.Panics(t, func() {
		GtFn(1, 2, nil)
	})
}

func TestGtEq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"greater", 10, 5, true},
		{"less", 5, 10, false},
		{"equal", 5, 5, true},
		{"negative and positive", -5, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GtEq(tt.a, tt.b))
		})
	}
}

func TestGtEqFn(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"greater", 10, 5, true},
		{"less", 5, 10, false},
		{"equal", 5, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GtEqFn(tt.a, tt.b, cmp))
		})
	}
}

func TestGtEqFnPanic(t *testing.T) {
	assert.Panics(t, func() {
		GtEqFn(1, 2, nil)
	})
}

func TestLt(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"less", 5, 10, true},
		{"greater", 10, 5, false},
		{"equal", 5, 5, false},
		{"negative and positive", -5, 5, true},
		{"negative values", -10, -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Lt(tt.a, tt.b))
		})
	}
}

func TestLtFn(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"less", 5, 10, true},
		{"greater", 10, 5, false},
		{"equal", 5, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, LtFn(tt.a, tt.b, cmp))
		})
	}
}

func TestLtFnPanic(t *testing.T) {
	assert.Panics(t, func() {
		LtFn(1, 2, nil)
	})
}

func TestLtEq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"less", 5, 10, true},
		{"greater", 10, 5, false},
		{"equal", 5, 5, true},
		{"negative and positive", -5, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, LtEq(tt.a, tt.b))
		})
	}
}

func TestLtEqFn(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"less", 5, 10, true},
		{"greater", 10, 5, false},
		{"equal", 5, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, LtEqFn(tt.a, tt.b, cmp))
		})
	}
}

func TestLtEqFnPanic(t *testing.T) {
	assert.Panics(t, func() {
		LtEqFn(1, 2, nil)
	})
}

// Custom type tests
type customType struct {
	value int
}

func TestCustomTypeComparison(t *testing.T) {
	cmp := func(a, b customType) int {
		return a.value - b.value
	}

	a := customType{value: 5}
	b := customType{value: 10}
	c := customType{value: 5}

	assert.True(t, EqFn(a, c, cmp))
	assert.False(t, EqFn(a, b, cmp))
	assert.True(t, NotEqFn(a, b, cmp))
	assert.True(t, LtFn(a, b, cmp))
	assert.True(t, LtEqFn(a, b, cmp))
	assert.True(t, LtEqFn(a, c, cmp))
	assert.False(t, GtFn(a, b, cmp))
	assert.True(t, GtFn(b, a, cmp))
	assert.True(t, GtEqFn(a, c, cmp))
}

func TestAsLess(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"first less than second", 3, 7, true},
		{"first equal to second", 5, 5, false},
		{"first greater than second", 8, 4, false},
		{"comparing with zero", 0, 1, true},
		{"comparing negative numbers", -3, -2, true},
		{"comparing positive and negative", -1, 1, true},
		{"comparing large numbers", 1000000, 1000001, true},
		{"comparing same large numbers", 1000000, 1000000, false},
	}

	lessFn := AsLess(cmp)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, lessFn(tt.a, tt.b))
		})
	}
}

func TestAsLessWithCustomStruct(t *testing.T) {
	type version struct {
		major, minor, patch int
	}

	cmp := func(a, b version) int {
		if a.major != b.major {
			return a.major - b.major
		}
		if a.minor != b.minor {
			return a.minor - b.minor
		}
		return a.patch - b.patch
	}

	tests := []struct {
		name     string
		a, b     version
		expected bool
	}{
		{"lower major version", version{1, 0, 0}, version{2, 0, 0}, true},
		{"same major different minor", version{1, 2, 0}, version{1, 3, 0}, true},
		{"same major and minor different patch", version{1, 2, 3}, version{1, 2, 4}, true},
		{"identical versions", version{1, 2, 3}, version{1, 2, 3}, false},
		{"higher version", version{2, 0, 0}, version{1, 9, 9}, false},
	}

	lessFn := AsLess(cmp)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, lessFn(tt.a, tt.b))
		})
	}
}

func TestAsLessPanic(t *testing.T) {
	assert.Panics(t, func() {
		AsLess[int](nil)
	})
}

//revive:disable-next-line:cognitive-complexity
func TestReverse(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"reverse less than", 5, 10, 1},
		{"reverse greater than", 10, 5, -1},
		{"reverse equal", 5, 5, 0},
		{"reverse with zero", 0, 1, 1},
		{"reverse negative numbers", -5, -3, 1},
		{"reverse mixed signs", -1, 1, 1},
	}

	reversedCmp := Reverse(cmp)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, reversedCmp(tt.a, tt.b))
		})
	}

	// Test with custom type
	type score struct {
		value float64
	}
	scoreCmp := func(a, b score) int {
		if a.value < b.value {
			return -1
		}
		if a.value > b.value {
			return 1
		}
		return 0
	}

	reversedScoreCmp := Reverse(scoreCmp)
	s1 := score{3.14}
	s2 := score{2.71}

	assert.Equal(t, -1, reversedScoreCmp(s1, s2))
	assert.Equal(t, 1, reversedScoreCmp(s2, s1))
	assert.Equal(t, 0, reversedScoreCmp(s1, s1))
}

func TestReversePanic(t *testing.T) {
	assert.Panics(t, func() {
		Reverse[int](nil)
	})
}

func TestReverseChained(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	// Test double reverse returns to original ordering
	doubleReversed := Reverse(Reverse(cmp))

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"double reverse less than", 5, 10, -1},
		{"double reverse greater than", 10, 5, 1},
		{"double reverse equal", 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, doubleReversed(tt.a, tt.b))
			assert.Equal(t, tt.expected, cmp(tt.a, tt.b))
		})
	}
}
