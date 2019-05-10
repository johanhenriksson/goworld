package render

import (
	"image"
	"io/ioutil"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	File    string
	Size    float64
	DPI     float64
	Spacing float64
	Color   Color

	fnt    *truetype.Font
	drawer *font.Drawer
}

func (f *Font) setup() {
	f.drawer = &font.Drawer{
		Face: truetype.NewFace(f.fnt, &truetype.Options{
			Size:    f.Size,
			DPI:     f.DPI,
			Hinting: font.HintingFull,
		}),
	}
}

func (f *Font) Render(text string, width, height float32, color Color) *Texture {
	tx := CreateTexture(int32(width), int32(height))
	f.RenderOn(tx, text, width, height, color)
	return tx
}

func (f *Font) RenderOn(tx *Texture, text string, width, height float32, color Color) {
	/* Set color */
	f.drawer.Src = image.NewUniform(color.RGBA())

	line := math.Ceil(f.Size * f.DPI / 72)
	//height := int(f.Spacing * line)
	//width := int(float64(f.drawer.MeasureString(text)) / f.Size)

	/* Create and attach destination image */
	rgba := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	f.drawer.Dst = rgba

	/* Draw text */
	f.drawer.Dot = fixed.P(0, int(line))
	f.drawer.DrawString(text)

	tx.Buffer(rgba)
}

/** Load a truetype font */
func LoadFont(filename string, size, dpi, spacing float64) *Font {
	fontBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	fnt := &Font{
		Size:    size,
		DPI:     dpi,
		Spacing: spacing,
		fnt:     f,
	}
	fnt.setup()
	return fnt
}
