package vertex

import (
	"strings"

	"log"

	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/util"
)

type Pointers []Pointer

type AttributeResolver interface {
	Attribute(name string) (shader.AttributeDesc, error)
}

func (ps Pointers) BufferString() string {
	names := util.Map(ps, func(i int, p Pointer) string { return p.Name })
	return strings.Join(names, ",")
}

func (ps Pointers) Bind(shader AttributeResolver) {
	for i, ptr := range ps {
		attr, err := shader.Attribute(ptr.Name)
		if err != nil {
			log.Printf("no attribute in shader %s\n", ptr.Name)
		}
		ptr.Bind(attr.Bind, attr.Type)
		ps[i] = ptr
	}
}

func (ps Pointers) Stride() int {
	if len(ps) == 0 {
		return 0
	}
	return ps[0].Stride
}
