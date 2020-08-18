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
			DPI:     72 * f.DPI,
			Hinting: font.HintingFull,
		}),
	}
}

func (f *Font) LineHeight() float32 {
	return float32(math.Ceil(f.Size * f.Spacing * f.DPI))
}

func (f *Font) Measure(text string) (int, int) {
	lines := 1
	width := 0
	s := 0
	for i, c := range text {
		if c == '\n' {
			line := text[s:i]
			w := f.drawer.MeasureString(line).Ceil()
			if w > width {
				width = w
			}
			s = i + 1
			lines++
		}
	}
	r := len(text)
	if s < r {
		line := text[s:]
		w := f.drawer.MeasureString(line).Ceil()
		if w > width {
			width = w
		}
	}

	lineHeight := int(f.LineHeight())
	height := lineHeight*lines + (lineHeight / 2)
	return width, height
}

func (f *Font) RenderNew(text string, color Color) *Texture {
	w, h := f.Measure(text)
	texture := CreateTexture(int32(w), int32(h))
	f.Render(texture, text, color)
	return texture
}

func (f *Font) Render(tx *Texture, text string, color Color) {
	f.drawer.Src = image.NewUniform(color.RGBA())

	width, height := f.Measure(text)

	// todo: its probably not a great idea to allocate an image on every draw
	// perhaps textures should always have a backing image ?
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	f.drawer.Dst = rgba

	s := 0
	line := 1
	lineHeight := int(f.LineHeight())
	for i, c := range text {
		if c == '\n' {
			if i == s {
				continue // skip empty rows
			}
			f.drawer.Dot = fixed.P(0, line*int(lineHeight))
			f.drawer.DrawString(text[s:i])
			s = i + 1
			line++
		}
	}
	if s < len(text) {
		f.drawer.Dot = fixed.P(0, line*int(lineHeight))
		f.drawer.DrawString(text[s:])
	}

	tx.Bind()
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
