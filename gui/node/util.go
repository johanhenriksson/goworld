package node

import (
	"github.com/samber/lo"
)

func If(condition bool, n T) T {
	if condition {
		return n
	}
	return nil
}

func Map[K any](items []K, transform func(K) T) []T {
	return lo.Map(items, func(it K, _ int) T {
		return transform(it)
	})
}
