package node

import "github.com/johanhenriksson/goworld/util"

func If(condition bool, n T) T {
	if condition {
		return n
	}
	return nil
}

func Map[K any](items []K, transform func(K) T) []T {
	return util.Map(items, transform)
}
