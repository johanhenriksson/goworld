package vertex

import (
	"strings"

	"github.com/johanhenriksson/goworld/util"
)

type Pointers []Pointer

func (ps Pointers) BufferString() string {
	names := util.Map(ps, func(p Pointer) string { return p.Name })
	return strings.Join(names, ",")
}

func (ps Pointers) Stride() int {
	if len(ps) == 0 {
		return 0
	}
	return ps[0].Stride
}
