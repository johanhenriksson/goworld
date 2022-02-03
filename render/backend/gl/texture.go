package gl

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func ActiveTexture(slot texture.Slot) error {
	tex := gl.TEXTURE0 + uint32(slot)
	gl.ActiveTexture(tex)
	return GetError()
}

func SetTexture2DFilter(min, mag texture.Filter) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(min))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(mag))
}

func SetTexture2DWrapMode(s, t texture.WrapMode) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(s))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(t))
}

func GenTexture() (texture.ID, error) {
	var id uint32
	gl.GenTextures(1, &id)
	return texture.ID(id), GetError()
}

func BindTexture2D(id texture.ID) error {
	gl.BindTexture(gl.TEXTURE_2D, uint32(id))

	err := GetError()
	switch err {
	case ErrInvalidEnum:
		return fmt.Errorf("%w: texture target is not one of the allowable values", err)
	case ErrInvalidValue:
		return fmt.Errorf("%w: texture is not a name returned from a previous call to glGenTextures", err)
	case ErrInvalidOperation:
		return fmt.Errorf("%w: texture was previously created with a target that doesn't match that of target.", err)
	default:
		return err
	}
}
