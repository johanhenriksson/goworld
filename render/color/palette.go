package color

import (
	"github.com/johanhenriksson/goworld/math/random"
)

// A Palette is a list of colors
type Palette []T

// RawPalette creates a palette from a list of hex integers
func RawPalette(colors ...int) Palette {
	palette := make(Palette, len(colors))
	for i, clr := range colors {
		palette[i] = RGBA(
			float32((clr>>16)&0xFF)/255.0,
			float32((clr>>8)&0xFF)/255.0,
			float32((clr>>0)&0xFF)/255.0,
			1.0)
	}
	return palette
}

// DefaultPalette https://lospec.com/palette-list/broken-facility
var DefaultPalette = RawPalette(
	0x24211e, 0x898377, 0xada99e, 0xcccac4, 0xf9f8f7,
	0x563735, 0x835748, 0xa37254, 0xb59669, 0xcab880,
	0x4d1c2d, 0x98191e, 0xd12424, 0xdd4b63, 0xf379e2,
	0xc86826, 0xd8993f, 0xe8c04f, 0xf2db89, 0xf8f1c6,
	0x17601f, 0x488c36, 0x7abd40, 0xa4cf41, 0xcdde5e,
	0x5044ba, 0x5e9ccc, 0x7fc6ce, 0x9de2df, 0xcaf1ea,
	0x202c56, 0x3f2d6d, 0x772673, 0xb9284f, 0xcb5135,
	0xeda7d8, 0xf3bedd, 0xdbebeb, 0xe9dde8, 0xd5c4df)

func Random() T {
	return random.Choice(DefaultPalette)
}
