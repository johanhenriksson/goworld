package color

import (
	"encoding/gob"
	"fmt"
	"image/color"

	"github.com/johanhenriksson/goworld/math/byte4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/texture"
)

// Predefined Colors
var (
	White       = T{1, 1, 1, 1}
	Black       = T{0, 0, 0, 1}
	Red         = T{1, 0, 0, 1}
	Green       = T{0, 1, 0, 1}
	Blue        = T{0, 0, 1, 1}
	Purple      = T{1, 0, 1, 1}
	Yellow      = T{1, 1, 0, 1}
	Cyan        = T{0, 1, 1, 1}
	Transparent = T{0, 0, 0, 0}
	None        = T{0, 0, 0, 0}

	DarkGrey = T{0.2, 0.2, 0.2, 1}
)

func init() {
	gob.Register(T{})
}

// T holds 32-bit RGBA colors
type T struct {
	R, G, B, A float32
}

var _ texture.Ref = T{}

// Color4 creates a color struct from its RGBA components
func RGBA(r, g, b, a float32) T {
	return T{r, g, b, a}
}

func RGB(r, g, b float32) T {
	return T{r, g, b, 1}
}

func RGBA8(r, g, b, a uint8) T {
	return RGBA(float32(r)/255, float32(g)/255, float32(b)/255, float32(a)/255)
}

func RGB8(r, g, b uint8) T {
	return RGBA8(r, g, b, 255)
}

// RGBA returns an 8-bit RGBA image/color
func (c T) RGBA() color.RGBA {
	return color.RGBA{
		uint8(255.0 * c.R),
		uint8(255.0 * c.G),
		uint8(255.0 * c.B),
		uint8(255.0 * c.A),
	}
}

func FromVec3(v vec3.T) T {
	return RGB(v.X, v.Y, v.Z)
}

// Vec3 returns a vec3 containing the RGB components of the color
func (c T) Vec3() vec3.T {
	return vec3.New(c.R, c.G, c.B)
}

func FromVec4(v vec4.T) T {
	return RGBA(v.X, v.Y, v.Z, v.W)
}

// Vec4 returns a vec4 containing the RGBA components of the color
func (c T) Vec4() vec4.T {
	return vec4.New(c.R, c.G, c.B, c.A)
}

func (c T) Byte4() byte4.T {
	return byte4.New(
		byte(255.0*c.R),
		byte(255.0*c.G),
		byte(255.0*c.B),
		byte(255.0*c.A))
}

func (c T) String() string {
	return fmt.Sprintf("(R:%.2f G:%.2f B:%.2f A:%.2f)", c.R, c.G, c.B, c.A)
}

var hexDigits = []byte("0123456789abcdef")

func (c T) Hex() string {
	rgba := c.Byte4()
	bytes := make([]byte, 9)
	bytes[0] = '#'
	bytes[1] = hexDigits[rgba.X>>4]
	bytes[2] = hexDigits[rgba.X&0x0F]
	bytes[3] = hexDigits[rgba.Y>>4]
	bytes[4] = hexDigits[rgba.Y&0x0F]
	bytes[5] = hexDigits[rgba.Z>>4]
	bytes[6] = hexDigits[rgba.Z&0x0F]
	if c.A < 1 {
		bytes[7] = hexDigits[rgba.W>>4]
		bytes[8] = hexDigits[rgba.W&0x0F]
	} else {
		bytes = bytes[:7]
	}
	return string(bytes)
}

// WithAlpha returns a new color with a modified alpha value
func (c T) WithAlpha(a float32) T {
	c.A = a
	return c
}

func Hex(s string) T {
	if s[0] != '#' {
		panic("invalid color value")
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		panic("invalid color value")
	}

	c := T{A: 1}
	switch len(s) {
	case 9:
		c.A = float32(hexToByte(s[7])<<4+hexToByte(s[6])) / 255
		fallthrough
	case 7:
		c.R = float32(hexToByte(s[1])<<4+hexToByte(s[2])) / 255
		c.G = float32(hexToByte(s[3])<<4+hexToByte(s[4])) / 255
		c.B = float32(hexToByte(s[5])<<4+hexToByte(s[6])) / 255
	case 4:
		c.R = float32(hexToByte(s[1])*17) / 255
		c.G = float32(hexToByte(s[2])*17) / 255
		c.B = float32(hexToByte(s[3])*17) / 255
	default:
		panic("invalid color value")
	}
	return c
}

//
// implement texture reference interface, so that colors may be easily loaded as textures
//

func (c T) Key() string  { return c.Hex() }
func (c T) Version() int { return 1 }

func (c T) ImageData() *image.Data {
	rgba := c.Byte4()
	return &image.Data{
		Width:  1,
		Height: 1,
		Format: image.FormatRGBA8Unorm,
		Buffer: []byte{
			rgba.X, rgba.Y, rgba.Z, rgba.W,
		},
	}
}

func (c T) TextureArgs() texture.Args {
	return texture.Args{
		Filter: texture.FilterNearest,
		Wrap:   texture.WrapClamp,
	}
}
