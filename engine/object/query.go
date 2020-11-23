package object

type ComponentPredicate func(Component) bool

type Query struct {
	Match   ComponentPredicate
	Results []Component
}

func NewQuery(predicate ComponentPredicate) Query {
	return Query{
		Match:   predicate,
		Results: make([]Component, 0, 128),
	}
}

func (q *Query) Append(result Component) {
	q.Results = append(q.Results, result)
}

func (q *Query) Clear() {
	// clear slice, but keep the memory
	q.Results = q.Results[:0]
}
