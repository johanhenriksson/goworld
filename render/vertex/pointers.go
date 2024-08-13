package vertex

import (
	"strings"

	"github.com/samber/lo"
)

type Pointers []Pointer

func (ps Pointers) BufferString() string {
	names := lo.Map(ps, func(p Pointer, _ int) string { return p.Name })
	return strings.Join(names, ",")
}

func (ps Pointers) Stride() int {
	if len(ps) == 0 {
		return 0
	}
	return ps[0].Stride
}
