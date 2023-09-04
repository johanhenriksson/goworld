package noise

import (
	"fmt"
	"math/rand"

	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/texture"
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
	buffer := make([]byte, 4*size)
	for i := 0; i < len(buffer); i += 4 {
		v := uint8(rand.Intn(255))
		buffer[i+0] = v
		buffer[i+1] = v
		buffer[i+2] = v
		buffer[i+3] = 255
	}
	return &image.Data{
		Width:  n.Width,
		Height: n.Height,
		Format: image.FormatRGBA8Unorm,
		Buffer: buffer,
	}
}

func (n *WhiteNoise) TextureArgs() texture.Args {
	return texture.Args{
		Filter: texture.FilterNearest,
		Wrap:   texture.WrapRepeat,
	}
}
