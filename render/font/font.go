package font

import (
	"errors"
	"fmt"
	"image"
	imgcolor "image/color"
	"io"
	"io/fs"
	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"

	fontlib "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var ErrNoGlyph = errors.New("no glyph for rune")

type T interface {
	Name() string
	Glyph(rune) (*Glyph, error)
	Measure(string, Args) vec2.T
	Render(string, Args) *image.RGBA
	Size() float32
}

type Args struct {
	Color      color.T
	LineHeight float32
}

type font struct {
	size   float32
	fnt    *truetype.Font
	face   fontlib.Face
	drawer *fontlib.Drawer
	mutex  *sync.Mutex
	glyphs map[rune]*Glyph
}

func (f *font) Name() string {
	return f.fnt.Name(truetype.NameIDFontFullName)
}

func (f *font) Size() float32 { return f.size }

func (f *font) Glyph(r rune) (*Glyph, error) {
	if cached, exists := f.glyphs[r]; exists {
		return cached, nil
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	bounds, advance, ok := f.face.GlyphBounds(r)
	if !ok {
		return nil, ErrNoGlyph
	}

	// calculate size and bearing
	bearing := vec2.New(
		FixToFloat(bounds.Min.X),
		FixToFloat(bounds.Min.Y))

	dr, mask, offset, _, _ := f.face.Glyph(fixed.Point26_6{X: 0, Y: 0}, r)

	// texture size
	size := vec2.New(
		float32(dr.Max.X-dr.Min.X),
		float32(dr.Max.Y-dr.Min.Y))

	// copy image mask
	img := image.NewRGBA(image.Rect(0, 0, int(size.X), int(size.Y)))
	for y := 0; y < int(size.Y); y++ {
		for x := 0; x < int(size.X); x++ {
			// grab alpha value as 16-bit integer
			_, _, _, alpha := mask.At(offset.X+x, offset.Y+y).RGBA()
			// create a white texture using the mask as alpha
			c := imgcolor.RGBA{
				R: 0xFF,
				G: 0xFF,
				B: 0xFF,
				A: uint8(alpha >> 8),
			}
			img.Set(x, y, c)
		}
	}

	glyph := &Glyph{
		key:     fmt.Sprintf("glyph:%s:%d:%c", f.Name(), int(f.size), r),
		Size:    size,
		Bearing: bearing,
		Advance: FixToFloat(advance),
		Mask:    img,
	}
	f.glyphs[r] = glyph
	return glyph, nil
}

func (f *font) Measure(text string, args Args) vec2.T {
	f.mutex.Lock()
	defer f.mutex.Unlock()

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
	height := lineHeight*lines + (lineHeight/2)*(lines-1)
	return vec2.NewI(width, height)
}

func (f *font) Render(text string, args Args) *image.RGBA {
	size := f.Measure(text, args)

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if args.LineHeight == 0 {
		args.LineHeight = 1
	}

	f.drawer.Src = image.NewUniform(args.Color.RGBA())

	// todo: its probably not a great idea to allocate an image on every draw
	// perhaps textures should always have a backing image ?
	output := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(size.X)), int(math.Ceil(size.Y))))
	f.drawer.Dst = output

	// debug outline
	// for y := 0; y < output.Bounds().Size().Y; y++ {
	// 	for x := 0; x < output.Bounds().Size().X; x++ {
	// 		if x == 0 || y == 0 || x == output.Bounds().Size().X-1 || y == output.Bounds().Size().Y-1 {
	// 			output.Set(x, y, color.Red.RGBA())
	// 		}
	// 	}
	// }

	s := 0
	line := 1
	lineHeight := int(math.Ceil(f.size * args.LineHeight))
	yOffset := int(math.Ceil(f.size * 0.20))

	for i, c := range text {
		if c == '\n' {
			if i == s {
				continue // skip empty rows
			}
			f.drawer.Dot = fixed.P(0, line*int(lineHeight)-yOffset)
			f.drawer.DrawString(text[s:i])
			s = i + 1
			line++
		}
	}
	if s < len(text) {
		f.drawer.Dot = fixed.P(0, line*int(lineHeight)-yOffset)
		f.drawer.DrawString(text[s:])
	}

	return output
}

// Load a truetype font
func Load(file fs.File, size int) (T, error) {
	// todo: read & parse could be cached, instead of doing it for each size of the font
	fontBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(fnt, &truetype.Options{
		Size:    float64(size),
		Hinting: fontlib.HintingFull,
	})

	return &font{
		size:   float32(size),
		fnt:    fnt,
		face:   face,
		glyphs: make(map[rune]*Glyph, 128),

		drawer: &fontlib.Drawer{
			Face: face,
		},
		mutex: &sync.Mutex{},
	}, nil
}

func FixToFloat(v fixed.Int26_6) float32 {
	div := 1 / float32(1<<6)
	return float32(v) * div
}

func PointToVec(p fixed.Point26_6) vec2.T {
	return vec2.T{
		X: FixToFloat(p.X),
		Y: FixToFloat(p.Y),
	}
}
