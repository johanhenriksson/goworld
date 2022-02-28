package vertex

import (
	"strings"

	"log"

	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/util"
)

type Pointers []Pointer

func (ps Pointers) BufferString() string {
	names := util.Map(ps, func(i int, p Pointer) string { return p.Name })
	return strings.Join(names, ",")
}

func (ps Pointers) Bind(shader shader.T) {
	for i, ptr := range ps {
		attr, err := shader.Attribute(ptr.Name)
		if err != nil {
			log.Printf("no attribute in shader %s: %s\n", shader.Name(), ptr.Name)
		}
		ptr.Bind(attr.Bind, attr.Type)
		ps[i] = ptr
	}
}
