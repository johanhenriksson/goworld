package vertex

import (
	"github.com/johanhenriksson/goworld/render/backend/types"
)

type Pointer struct {
	Name        string
	Binding     int
	Source      types.Type
	Destination types.Type
	Elements    int
	Stride      int
	Offset      int
	Normalize   bool
}

func (p *Pointer) Bind(binding int, kind types.Type) {
	p.Binding = binding
	p.Destination = kind
}
