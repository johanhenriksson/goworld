package object

import (
	"sort"

	"github.com/samber/lo"
)

type Query[K Component] struct {
	results []K
	filters []func(b K) bool
	sorter  func(a, b K) bool
}

// NewQuery returns a new query for the given component type
func NewQuery[K Component]() *Query[K] {
	return &Query[K]{
		filters: make([]func(K) bool, 0, 8),
		results: make([]K, 0, 128),
	}
}

// Where applies a filter predicate to the results
func (q *Query[K]) Where(predicate func(K) bool) *Query[K] {
	q.filters = append(q.filters, predicate)
	return q
}

// Sort the result using a compare function.
// The compare function should return true if a is "less than" b
func (q *Query[K]) Sort(sorter func(a, b K) bool) *Query[K] {
	q.sorter = sorter
	return q
}

// Match returns true if the passed component matches the query
func (q *Query[K]) match(component K) bool {
	for _, filter := range q.filters {
		if !filter(component) {
			return false
		}
	}
	return true
}

// Append a component to the query results.
func (q *Query[K]) append(result K) {
	q.results = append(q.results, result)
}

// Clear the query results, without freeing the memory.
func (q *Query[K]) Reset() *Query[K] {
	// clear slice, but keep the memory
	q.results = q.results[:0]
	q.filters = q.filters[:0]
	return q
}

// First returns the first match in a depth-first fashion
func (q *Query[K]) First(root Component) (K, bool) {
	result, hit := q.first(root)
	return result, hit
}

func (q *Query[K]) first(root Component) (K, bool) {
	var empty K
	if !root.Enabled() {
		return empty, false
	}
	if k, ok := root.(K); ok {
		if q.match(k) {
			return k, true
		}
	}
	if group, ok := root.(Object); ok {
		for child := range group.Children() {
			if match, found := q.first(child); found {
				return match, true
			}
		}
	}
	return empty, false
}

// Collect returns all matching components
func (q *Query[K]) Collect(roots ...Component) []K {
	// collect all matches
	for _, root := range roots {
		q.collect(root)
	}

	// sort if required
	if q.sorter != nil {
		sort.Slice(q.results, func(i, j int) bool {
			return q.sorter(q.results[i], q.results[j])
		})
	}

	return q.results
}

func (q *Query[K]) CollectObjects(roots ...Component) []Component {
	return lo.Map(NewQuery[K]().Collect(roots...), func(s K, _ int) Component { return s })
}

func (q *Query[K]) collect(object Component) {
	if !object.Enabled() {
		return
	}
	if k, ok := object.(K); ok {
		if q.match(k) {
			q.append(k)
		}
	}
	if group, ok := object.(Object); ok {
		for child := range group.Children() {
			q.collect(child)
		}
	}
}
