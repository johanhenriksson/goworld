package font

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Glyph struct {
	key     string
	Size    vec2.T
	Bearing vec2.T
	Advance float32
	Mask    *image.Data
}

var _ texture.Ref = &Glyph{}

func (r *Glyph) Key() string  { return r.key }
func (r *Glyph) Version() int { return 1 }

func (r *Glyph) ImageData() *image.Data {
	return r.Mask
}

func (r *Glyph) TextureArgs() texture.Args {
	return texture.Args{
		Wrap:   core1_0.SamplerAddressModeClampToBorder,
		Border: core1_0.BorderColorFloatTransparentBlack,
		Filter: core1_0.FilterNearest,
	}
}
