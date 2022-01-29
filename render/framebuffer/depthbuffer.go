package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/texture"
)

// Depth is a framebuffer with only a depth component
type Depth interface {
	T
	Depth() texture.T
}
