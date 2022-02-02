package texture

import (
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var ErrInvalidTextureUnit = errors.New("invalid texture unit")

type Filter int32
type WrapMode int32
type Format int32

const (
	LinearFilter  = Filter(gl.LINEAR)
	NearestFilter = Filter(gl.NEAREST)
)

const (
	ClampWrap  = WrapMode(gl.CLAMP_TO_EDGE)
	RepeatWrap = WrapMode(gl.REPEAT)
)

const (
	RGB  = Format(gl.RGB)
	RGBA = Format(gl.RGBA)
	Red  = Format(gl.RED)
)
