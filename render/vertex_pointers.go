package render

import (
	"strings"
)

type Pointers []Pointer

func (ps Pointers) BufferName() string {
	names := make([]string, len(ps))
	for i, p := range ps {
		names[i] = p.Name
	}
	return strings.Join(names, ",")
}

func (ps Pointers) Enable() {
	for _, p := range ps {
		p.Enable()
	}
}

func (ps Pointers) Disable() {
	for _, p := range ps {
		p.Disable()
	}
}
