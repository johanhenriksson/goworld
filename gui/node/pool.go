package node

import (
	"reflect"
)

type Pool struct {
	free map[reflect.Type][]any
}

var globalPool = &Pool{
	free: make(map[reflect.Type][]any),
}

func Alloc[P any](pool *Pool, kind reflect.Type) *node[P] {
	if freelist, exists := pool.free[kind]; exists && len(freelist) > 0 {
		idx := len(freelist) - 1
		n := freelist[idx]
		pool.free[kind] = freelist[:idx]
		return n.(*node[P])
	}
	return &node[P]{
		kind: kind,
	}
}

func Free[P any](pool *Pool, n *node[P]) {
	freelist, exists := pool.free[n.kind]
	if !exists {
		freelist = make([]any, 0, 128)
	}
	pool.free[n.kind] = append(freelist, n)
}
