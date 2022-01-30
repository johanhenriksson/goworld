package object

import "sort"

type ComponentPredicate func(Component) bool

// ComponentSorter should return true if a is "less than" b
type ComponentSorter func(a, b Component) bool

type Query struct {
	results []Component
	filters []ComponentPredicate
	sorter  ComponentSorter
}

// NewQuery returns a new component query.
func NewQuery() *Query {
	return &Query{
		filters: make([]ComponentPredicate, 0, 8),
		results: make([]Component, 0, 128),
	}
}

func (q *Query) Where(predicate ComponentPredicate) *Query {
	q.filters = append(q.filters, predicate)
	return q
}

func (q *Query) Sort(sorter ComponentSorter) *Query {
	q.sorter = sorter
	return q
}

// Match returns true if the passed component matches the query
func (q *Query) match(c Component) bool {
	for _, filter := range q.filters {
		if !filter(c) {
			return false
		}
	}
	return true
}

// Append a component to the query results.
func (q *Query) append(result Component) {
	q.results = append(q.results, result)
}

// Clear the query results, without freeing the memory.
func (q *Query) clear() {
	// clear slice, but keep the memory
	q.results = q.results[:0]
}

func (q *Query) Collect(root T) []Component {
	q.clear()

	// collect all matches
	q.collect(root)

	// sort if required
	if q.sorter != nil {
		sort.Slice(q.results, func(i, j int) bool {
			return q.sorter(q.results[i], q.results[j])
		})
	}

	return q.results
}

func (q *Query) collect(object T) {
	for _, component := range object.Components() {
		if !component.Active() {
			continue
		}
		if q.match(component) {
			q.append(component)
		}
	}
	for _, child := range object.Children() {
		if !child.Active() {
			continue
		}
		q.collect(child)
	}
}
