package noise

import (
	"fmt"
	"math/rand"

	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type WhiteNoise struct {
	Width  int
	Height int

	key string
}

func NewWhiteNoise(width, height int) *WhiteNoise {
	return &WhiteNoise{
		key:    fmt.Sprintf("noise-white-%dx%d", width, height),
		Width:  width,
		Height: height,
	}
}

func (n *WhiteNoise) Key() string  { return n.key }
func (n *WhiteNoise) Version() int { return 1 }

func (n *WhiteNoise) ImageData() *image.Data {
	size := n.Width * n.Height
	buffer := make([]byte, size)
	for i := 0; i < len(buffer); i++ {
		v := uint8(rand.Intn(100))
		buffer[i+0] = 100 + v
	}
	return &image.Data{
		Width:  n.Width,
		Height: n.Height,
		Format: core1_0.FormatR8UnsignedNormalized,
		Buffer: buffer,
	}
}

func (n *WhiteNoise) TextureArgs() texture.Args {
	return texture.Args{
		Filter: texture.FilterNearest,
		Wrap:   texture.WrapRepeat,
	}
}
