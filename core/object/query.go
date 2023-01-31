package object

import (
	"sort"
)

type query[K T] struct {
	results []K
	filters []func(b K) bool
	sorter  func(a, b K) bool
}

// Any returns a query for generic components
func Any() *query[T] {
	return Query[T]()
}

// Query returns a new query for the given component type
func Query[K T]() *query[K] {
	return &query[K]{
		filters: make([]func(K) bool, 0, 8),
		results: make([]K, 0, 128),
	}
}

// Where applies a filter predicate to the results
func (q *query[K]) Where(predicate func(K) bool) *query[K] {
	q.filters = append(q.filters, predicate)
	return q
}

// Sort the result using a compare function.
// The compare function should return true if a is "less than" b
func (q *query[K]) Sort(sorter func(a, b K) bool) *query[K] {
	q.sorter = sorter
	return q
}

// Match returns true if the passed component matches the query
func (q *query[K]) match(component K) bool {
	for _, filter := range q.filters {
		if !filter(component) {
			return false
		}
	}
	return true
}

// Append a component to the query results.
func (q *query[K]) append(result K) {
	q.results = append(q.results, result)
}

// Clear the query results, without freeing the memory.
func (q *query[K]) clear() {
	// clear slice, but keep the memory
	q.results = q.results[:0]
}

// First returns the first match
func (q *query[K]) First(root T) K {
	result, _ := q.first(root)
	return result
}

func (q *query[K]) first(root T) (K, bool) {
	var empty K
	if !root.Active() {
		return empty, false
	}
	if k, ok := root.(K); ok {
		if q.match(k) {
			return k, true
		}
	}
	for _, child := range root.Children() {
		if match, found := q.first(child); found {
			return match, true
		}
	}
	return empty, false
}

// Collect returns all matching components
func (q *query[K]) Collect(root T) []K {
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

func (q *query[K]) collect(object T) {
	if !object.Active() {
		return
	}
	if k, ok := object.(K); ok {
		if q.match(k) {
			q.append(k)
		}
	}
	for _, child := range object.Children() {
		q.collect(child)
	}
}
