package object

type ComponentPredicate func(Component) bool

type Query struct {
	Match   ComponentPredicate
	Results []Component
}

// NewQuery returns a new component query.
func NewQuery(predicate ComponentPredicate) Query {
	return Query{
		Match:   predicate,
		Results: make([]Component, 0, 128),
	}
}

// Append a component to the query results.
func (q *Query) Append(result Component) {
	q.Results = append(q.Results, result)
}

// Clear the query results, without freeing the memory.
func (q *Query) Clear() {
	// clear slice, but keep the memory
	q.Results = q.Results[:0]
}
