package gl

import (
	"errors"

	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var ErrInvalidTextureUnit = errors.New("invalid texture unit")

type TextureSlot int

func ActiveTexture(slot TextureSlot) error {
	texture := gl.TEXTURE0 + uint32(slot)
	gl.ActiveTexture(texture)
	if gl.GetError() == gl.INVALID_ENUM {
		return ErrInvalidTextureUnit
	}
	return nil
}

func SetTexture2DFilter(min, mag texture.Filter) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(min))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(mag))
}

func SetTexture2DWrapMode(s, t texture.WrapMode) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(s))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(t))
}
