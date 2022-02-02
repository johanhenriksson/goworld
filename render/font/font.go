package font

import (
	"image"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"

	fontlib "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type T interface {
	Measure(string, Args) vec2.T
	Render(string, Args) *image.RGBA
	Size() float32
}

type Args struct {
	Color      color.T
	LineHeight float32
}

type font struct {
	file   string
	size   float32
	fnt    *truetype.Font
	drawer *fontlib.Drawer
}

func (f *font) setup() {
	f.drawer = &fontlib.Drawer{
		Face: truetype.NewFace(f.fnt, &truetype.Options{
			Size:    float64(f.size),
			Hinting: fontlib.HintingFull,
		}),
	}
}

func (f *font) Size() float32 { return f.size }

func (f *font) Measure(text string, args Args) vec2.T {
	if args.LineHeight == 0 {
		args.LineHeight = 1
	}

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

	lineHeight := int(math.Ceil(f.size * args.LineHeight))
	height := lineHeight*lines + (lineHeight / 2)
	return vec2.NewI(width, height)
}

func (f *font) Render(text string, args Args) *image.RGBA {
	if args.LineHeight == 0 {
		args.LineHeight = 1
	}

	f.drawer.Src = image.NewUniform(args.Color.RGBA())

	size := f.Measure(text, args)

	// todo: its probably not a great idea to allocate an image on every draw
	// perhaps textures should always have a backing image ?
	output := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(size.X)), int(math.Ceil(size.Y))))
	f.drawer.Dst = output

	s := 0
	line := 1
	lineHeight := int(math.Ceil(f.size * args.LineHeight))

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

	return output
}

// Load a truetype font
func Load(filename string, size int) T {
	fontBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	fnt := &font{
		file: filename,
		size: float32(size),
		fnt:  f,
	}
	fnt.setup()
	return fnt
}
