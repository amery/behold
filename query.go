package behold

import (
	"reflect"
)

type Store interface{}
type Transaction interface{}
type iterBookmark struct{}

// Query is a chained collection of criteria of which an object needs to match to
// be returned. An Empty query matches against all records.
type Query struct {
	index         string
	currentField  string
	fieldCriteria map[string][]*Criterion
	ors           []*Query

	badIndex bool
	dataType reflect.Type
	tx       Transaction
	writable bool
	subquery bool
	bookmark *iterBookmark

	limit   int
	skip    int
	sort    []string
	reverse bool
}

func (q *Query) IsEmpty() bool {
	if q.index != "" {
		return false
	}

	if len(q.fieldCriteria) != 0 {
		return false
	}

	if len(q.ors) != 0 {
		return false
	}

	return true
}
