package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/texture"
)

// ColorBuffer is a framebuffer with a color component
type Color interface {
	T
	Texture() texture.T
}
