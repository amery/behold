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
