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
	buffer := make([]byte, 4*n.Width*n.Height)
	_, err := rand.Read(buffer)
	if err != nil {
		panic(err)
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
		Filter: texture.FilterLinear,
		Wrap:   texture.WrapRepeat,
	}
}
