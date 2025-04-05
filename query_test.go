package behold

import (
	"testing"
)

// Name constants to reduce duplication
const (
	nameJohn  = "John"
	nameAlice = "Alice"
	nameBob   = "Bob"
)

// Simple test type
type testPerson struct {
	Name string
	Age  int
}

// Field accessor functions for testPerson
func personAge(p testPerson) int {
	return p.Age
}

func personName(p testPerson) string {
	return p.Name
}

func TestQueryFunc(t *testing.T) {
	isAdult := ComposeQuery(personAge, GtEqQuery(18))
	isNamedJohn := ComposeQuery(personName, EqQuery(nameJohn))

	adult := testPerson{Name: nameAlice, Age: 30}
	child := testPerson{Name: nameBob, Age: 10}
	john := testPerson{Name: nameJohn, Age: 25}

	if !isAdult.Match(adult) {
		t.Error("Adult should match isAdult query")
	}

	if isAdult.Match(child) {
		t.Error("Child should not match isAdult query")
	}

	if !isNamedJohn.Match(john) {
		t.Error("John should match isNamedJohn query")
	}

	if isNamedJohn.Match(adult) {
		t.Error("Alice should not match isNamedJohn query")
	}

	// Test nil function behavior
	var nilFunc QueryFunc[testPerson]
	if !nilFunc.Match(adult) {
		t.Error("Nil function should match everything")
	}
}

func TestMatchAny(t *testing.T) {
	isAdult := ComposeQuery(personAge, GtEqQuery(18))
	isNamedJohn := ComposeQuery(personName, EqQuery(nameJohn))

	// Test with queries
	anyQuery := MatchAny(isAdult, isNamedJohn)

	adult := testPerson{Name: nameAlice, Age: 30}
	child := testPerson{Name: nameBob, Age: 10}
	john := testPerson{Name: nameJohn, Age: 15}

	if !anyQuery.Match(adult) {
		t.Error("Adult should match anyQuery")
	}

	if !anyQuery.Match(john) {
		t.Error("John should match anyQuery")
	}

	if anyQuery.Match(child) {
		t.Error("Child should not match anyQuery")
	}

	// Test with empty queries
	emptyQuery := MatchAny[testPerson]()
	if emptyQuery.Match(adult) {
		t.Error("MatchAny with no queries should match nothing")
	}

	// Test with nil queries
	var nilQuery Query[testPerson]
	anyWithNil := MatchAny(isAdult, nilQuery, isNamedJohn)

	if !anyWithNil.Match(adult) {
		t.Error("Adult should match anyWithNil")
	}

	if anyWithNil.Match(child) {
		t.Error("Child should not match anyWithNil")
	}
}

//revive:disable-next-line:cognitive-complexity
func TestMatchAll(t *testing.T) {
	isAdult := ComposeQuery(personAge, GtEqQuery(18))
	isNamedJohn := ComposeQuery(personName, EqQuery(nameJohn))

	// Test with queries
	allQuery := MatchAll(isAdult, isNamedJohn)

	adult := testPerson{Name: nameAlice, Age: 30}
	child := testPerson{Name: nameBob, Age: 10}
	john := testPerson{Name: nameJohn, Age: 25}
	youngJohn := testPerson{Name: nameJohn, Age: 15}

	if allQuery.Match(adult) {
		t.Error("Adult named Alice should not match allQuery")
	}

	if allQuery.Match(youngJohn) {
		t.Error("Young John should not match allQuery")
	}

	if !allQuery.Match(john) {
		t.Error("Adult John should match allQuery")
	}

	if allQuery.Match(child) {
		t.Error("Child named Bob should not match allQuery")
	}

	// Test with empty queries
	emptyQuery := MatchAll[testPerson]()
	if !emptyQuery.Match(adult) {
		t.Error("MatchAll with no queries should match everything")
	}

	// Test with nil queries
	var nilQuery Query[testPerson]
	allWithNil := MatchAll(isAdult, nilQuery, isNamedJohn)

	if allWithNil.Match(adult) {
		t.Error("Adult named Alice should not match allWithNil")
	}

	if !allWithNil.Match(john) {
		t.Error("Adult John should match allWithNil")
	}

	if allWithNil.Match(child) {
		t.Error("Child named Bob should not match allWithNil")
	}
}

//revive:disable-next-line:cognitive-complexity
//revive:disable-next-line:cyclomatic
func TestQueryComposition(t *testing.T) {
	isAdult := ComposeQuery(personAge, GtEqQuery(18))

	isNamedJohn := ComposeQuery(personName, EqQuery(nameJohn))
	isNamedBob := ComposeQuery(personName, EqQuery(nameBob))

	// Test AND composition
	adultNamedJohn := isAdult.And(isNamedJohn)

	john := testPerson{Name: nameJohn, Age: 25}
	youngJohn := testPerson{Name: nameJohn, Age: 15}
	adult := testPerson{Name: nameAlice, Age: 30}

	if !adultNamedJohn.Match(john) {
		t.Error("Adult John should match adultNamedJohn")
	}

	if adultNamedJohn.Match(youngJohn) {
		t.Error("Young John should not match adultNamedJohn")
	}

	if adultNamedJohn.Match(adult) {
		t.Error("Adult Alice should not match adultNamedJohn")
	}

	// Test OR composition
	johnOrBob := isNamedJohn.Or(isNamedBob)

	bob := testPerson{Name: nameBob, Age: 40}

	if !johnOrBob.Match(john) {
		t.Error("John should match johnOrBob")
	}

	if !johnOrBob.Match(bob) {
		t.Error("Bob should match johnOrBob")
	}

	if johnOrBob.Match(adult) {
		t.Error("Alice should not match johnOrBob")
	}

	// Test complex composition
	complexQuery := isAdult.And(isNamedJohn.Or(isNamedBob))

	youngBob := testPerson{Name: nameBob, Age: 15}

	if !complexQuery.Match(john) {
		t.Error("Adult John should match complexQuery")
	}

	if !complexQuery.Match(bob) {
		t.Error("Adult Bob should match complexQuery")
	}

	if complexQuery.Match(youngJohn) {
		t.Error("Young John should not match complexQuery")
	}

	if complexQuery.Match(youngBob) {
		t.Error("Young Bob should not match complexQuery")
	}

	if complexQuery.Match(adult) {
		t.Error("Adult Alice should not match complexQuery")
	}
}

//revive:disable-next-line:cognitive-complexity
func TestAndsAndOrs(t *testing.T) {
	// Test ands type
	a := ands[testPerson]{
		QueryFunc[testPerson](func(p testPerson) bool { return p.Age >= 18 }),
		QueryFunc[testPerson](func(p testPerson) bool { return p.Name == nameJohn }),
	}

	john := testPerson{Name: nameJohn, Age: 25}
	youngJohn := testPerson{Name: nameJohn, Age: 15}

	if !a.Match(john) {
		t.Error("Adult John should match ands query")
	}

	if a.Match(youngJohn) {
		t.Error("Young John should not match ands query")
	}

	// Test ors type
	o := ors[testPerson]{
		QueryFunc[testPerson](func(p testPerson) bool { return p.Age >= 18 }),
		QueryFunc[testPerson](func(p testPerson) bool { return p.Name == nameJohn }),
	}

	if !o.Match(john) {
		t.Error("Adult John should match ors query")
	}

	if !o.Match(youngJohn) {
		t.Error("Young John should match ors query (name matches)")
	}

	// Test with nil values in the slice
	andsWithNil := ands[testPerson]{
		nil,
		QueryFunc[testPerson](func(p testPerson) bool { return p.Age >= 18 }),
	}

	adult := testPerson{Name: nameAlice, Age: 30}
	child := testPerson{Name: nameBob, Age: 10}

	if !andsWithNil.Match(adult) {
		t.Error("Adult should match andsWithNil")
	}

	if andsWithNil.Match(child) {
		t.Error("Child should not match andsWithNil")
	}

	// Empty slices
	emptyAnds := ands[testPerson]{}
	emptyOrs := ors[testPerson]{}

	if !emptyAnds.Match(adult) {
		t.Error("Empty ands should match everything")
	}

	if emptyOrs.Match(adult) {
		t.Error("Empty ors should match nothing")
	}
}

func TestQJoin(t *testing.T) {
	q1 := QueryFunc[testPerson](func(p testPerson) bool { return p.Age >= 18 })
	q2 := QueryFunc[testPerson](func(p testPerson) bool { return p.Name == nameJohn })
	q3 := QueryFunc[testPerson](func(p testPerson) bool { return p.Name == nameBob })

	// Test with non-nil first query
	result := qJoin(q1, []Query[testPerson]{q2, q3})

	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}

	// Test with nil first query
	var nilQuery Query[testPerson]
	result = qJoin(nilQuery, []Query[testPerson]{q2, q3})

	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
}
