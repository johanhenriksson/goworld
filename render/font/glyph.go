package font

import (
	"image"

	"github.com/johanhenriksson/goworld/math/vec2"
)

type Glyph struct {
	key     string
	Size    vec2.T
	Bearing vec2.T
	Advance float32
	Mask    *image.RGBA
}

func (r *Glyph) Key() string  { return r.key }
func (r *Glyph) Version() int { return 1 }

func (r *Glyph) Load() *image.RGBA {
	return r.Mask
}
