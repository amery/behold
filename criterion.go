package behold

import (
	"regexp"
)

type operator int

const (
	undef operator = iota

	eq    // ==
	ne    // !=
	gt    // >
	lt    // <
	ge    // >=
	le    // <=
	in    // in
	re    // regular expression
	fn    // func
	isnil // tests for nil
	sw    // string starts with
	ew    // string ends with

	hk       // match map keys
	contains // slice only
	any      // slice only
	all      // slice only
)

// Criterion is an operator and a value that a given field needs to match on
type Criterion struct {
	query    *Query
	operator operator
	value    interface{}
	values   []interface{}
}

// Where starts a Query for specifying the criteria that an object needs to match to
// be returned in a Find result
func Where(field string) *Criterion {
	if !exportedField(field) {
		panic(ErrNotExported)
	}

	return &Criterion{
		query: &Query{
			currentField:  field,
			fieldCriteria: make(map[string][]*Criterion),
		},
	}
}

// And creates another set of criterior that needs to apply to a query
func (q *Query) And(field string) *Criterion {
	if !exportedField(field) {
		panic(ErrNotExported)
	}

	q.currentField = field
	return &Criterion{
		query: q,
	}
}

// Skip skips a number of records that match all the rest of the query criteria, and
// aren't included in the result set.
// Setting Skip multiple times or to a negative value will panic.
func (q *Query) Skip(amount int) *Query {
	if amount < 0 {
		panic(ErrNotPositive("Skip"))
	}

	if q.skip != 0 {
		panic(ErrAlreadySet("Skip", q.skip))
	}

	q.skip = amount
	return q
}

// Limit sets the maximum number of records that can be returned by a query.
// Setting Limit multiple times or to a negative value will panic.
func (q *Query) Limit(amount int) *Query {
	if amount < 0 {
		panic(ErrNotPositive("Limit"))
	}

	if q.limit != 0 {
		panic(ErrAlreadySet("Limit", q.limit))
	}

	q.limit = amount
	return q
}

// Reverse will reverse the current result set.
// Useful with SortBy
func (q *Query) Reverse() *Query {
	q.reverse = !q.reverse
	return q
}

// Contains tests if the current field is a slice that contains the passed in value
func (c *Criterion) Contains(value interface{}) *Query {
	return c.op(contains, value)
}

// ContainsAll tests if the current field is a slice that contains all of the passed in values
func (c *Criterion) ContainsAll(values ...interface{}) *Query {
	return c.op(all, values...)
}

// ContainsAny tests if the current field is a slice that contains any of the passed in values
func (c *Criterion) ContainsAny(values ...interface{}) *Query {
	return c.op(any, values...)
}

// In tests if the current field is in the passed in values
func (c *Criterion) In(values ...interface{}) *Query {
	return c.op(in, values...)
}

// Match will test if the current field matches against a regular expression.
// The field value will be converted to string (%s) before testing
func (c *Criterion) Match(expression *regexp.Regexp) *Query {
	return c.op(re, expression)
}

// HasPrefix will test if a field starts with the given string
func (c *Criterion) HasPrefix(prefix string) *Query {
	return c.op(sw, prefix)
}

// HasSuffix will test if a field ends with the given string
func (c *Criterion) HasSuffix(suffix string) *Query {
	return c.op(sw, suffix)
}

// IsNil will test if the current field is equal to nil
func (c *Criterion) IsNil() *Query {
	return c.op(isnil, nil)
}

// HasKey tests if the current field has a map key matching the passed in value
func (c *Criterion) HasKey(value interface{}) *Query {
	return c.op(hk, value)
}

// Eq tests if the current field is Equal To the passed in value
func (c *Criterion) Eq(value interface{}) *Query {
	return c.op(eq, value)
}

// Ne tests if the current field is Not Equal To the passed in value
func (c *Criterion) Ne(value interface{}) *Query {
	return c.op(ne, value)
}

// Gt tests if the current field is Greater Than the passed in value
func (c *Criterion) Gt(value interface{}) *Query {
	return c.op(gt, value)
}

// Lt tests if the current field is Lower Than the passed in value
func (c *Criterion) Lt(value interface{}) *Query {
	return c.op(lt, value)
}

// Ge tests if the current field is Greater Than or Equal To the passed in value
func (c *Criterion) Ge(value interface{}) *Query {
	return c.op(ge, value)
}

// Le tests if the current field is Lower Than or Equal To the passed in value
func (c *Criterion) Le(value interface{}) *Query {
	return c.op(le, value)
}

func (c *Criterion) op(op operator, values ...interface{}) *Query {

	c.operator = op
	if len(values) == 1 {
		c.value = values[0]
	} else {
		c.values = values
	}

	q := c.query
	q.appendCriteria(c)
	return q
}

func (q *Query) appendCriteria(c *Criterion) *Query {
	q.fieldCriteria[q.currentField] = append(q.fieldCriteria[q.currentField], c)
	return q
}
