package font

import (
	"github.com/johanhenriksson/goworld/assets/fs"
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

//
// assets.Texture implementation
//

func (r *Glyph) Key() string  { return r.key }
func (r *Glyph) Version() int { return 1 }

func (r *Glyph) LoadTexture(fs.Filesystem) *texture.Data {
	return &texture.Data{
		Image: r.Mask,
		Args: texture.Args{
			Filter:  texture.FilterLinear,
			Wrap:    texture.WrapBorder,
			Border:  core1_0.BorderColorFloatTransparentBlack,
			Mipmaps: false,
		},
	}
}
