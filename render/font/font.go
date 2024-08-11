package font

import (
	"errors"
	"fmt"
	"sync"

	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/util"

	fontlib "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var ErrNoGlyph = errors.New("no glyph for rune")

type Args struct {
	Color      color.T
	LineHeight float32
}

type Font struct {
	size   float32
	scale  float32
	name   string
	face   fontlib.Face
	drawer *fontlib.Drawer
	mutex  *sync.Mutex
	glyphs *util.SyncMap[rune, *Glyph]
	kern   *util.SyncMap[runepair, float32]
}

type runepair struct {
	a, b rune
}

func (f *Font) Name() string  { return f.name }
func (f *Font) Size() float32 { return f.size }

func (f *Font) Glyph(r rune) (*Glyph, error) {
	if cached, exists := f.glyphs.Load(r); exists {
		return cached, nil
	}

	// grab the font lock
	f.mutex.Lock()
	defer f.mutex.Unlock()

	bounds, advance, ok := f.face.GlyphBounds(r)
	if !ok {
		return nil, ErrNoGlyph
	}

	// calculate bearing
	bearing := vec2.New(FixToFloat(bounds.Min.X), FixToFloat(bounds.Min.Y)).Floor()

	// texture size
	size := vec2.New(FixToFloat(bounds.Max.X), FixToFloat(bounds.Max.Y)).Ceil().Sub(bearing)

	// glyph texture
	_, mask, offset, _, _ := f.face.Glyph(fixed.Point26_6{
		X: FloatToFix(-bearing.X),
		Y: FloatToFix(-bearing.Y),
	}, r)
	width, height := int(size.X), int(size.Y)

	img := &image.Data{
		Width:  width,
		Height: height,
		Buffer: make([]byte, 4*width*height),
		Format: image.FormatRGBA8Unorm,
	}
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// grab alpha value as 16-bit integer
			_, _, _, alpha := mask.At(offset.X+x, offset.Y+y).RGBA()
			img.Buffer[i+0] = 0xFF // red
			img.Buffer[i+1] = 0xFF // green
			img.Buffer[i+2] = 0xFF // blue
			img.Buffer[i+3] = uint8(alpha >> 8)
			i += 4
		}
	}

	scaleFactor := 1 / f.scale
	glyph := &Glyph{
		key:     fmt.Sprintf("glyph:%s:%dx%.2f:%c", f.Name(), int(f.size), f.scale, r),
		Size:    size.Scaled(scaleFactor),
		Bearing: bearing.Scaled(scaleFactor),
		Advance: FixToFloat(advance) * scaleFactor,
		Mask:    img,
	}
	f.glyphs.Store(r, glyph)

	return glyph, nil
}

func (f *Font) Kern(a, b rune) float32 {
	pair := runepair{a, b}
	if k, exists := f.kern.Load(pair); exists {
		return k
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()
	k := FixToFloat(f.face.Kern(a, b))
	f.kern.Store(pair, k)
	return k
}

func (f *Font) MeasureLine(text string) vec2.T {
	size := vec2.Zero
	var prev rune
	for i, r := range text {
		g, err := f.Glyph(r)
		if err != nil {
			panic("no such glyph")
		}
		if i > 0 {
			size.X += f.Kern(prev, r)
		}
		if i < len(text)-1 {
			size.X += g.Advance
		} else {
			size.X += g.Bearing.X + g.Size.X
		}
		size.Y = math.Max(size.Y, g.Size.Y)
		prev = r
	}
	return size.Scaled(f.scale)
}

func (f *Font) Measure(text string, args Args) vec2.T {
	if args.LineHeight == 0 {
		args.LineHeight = 1
	}

	lines := 1
	width := float32(0)
	s := 0
	for i, c := range text {
		if c == '\n' {
			line := text[s:i]
			// w := f.drawer.MeasureString(line).Ceil()
			w := f.MeasureLine(line)
			if w.X > width {
				width = w.X
			}
			s = i + 1
			lines++
		}
	}
	r := len(text)
	if s < r {
		line := text[s:]
		// w := f.drawer.MeasureString(line).Ceil()
		w := f.MeasureLine(line)
		if w.X > width {
			width = w.X
		}
	}

	lineHeight := int(math.Ceil(f.size * f.scale * args.LineHeight))
	height := lineHeight*lines + (lineHeight/2)*(lines-1)
	return vec2.New(width, float32(height)).Scaled(1 / f.scale).Ceil()
}

func FixToFloat(v fixed.Int26_6) float32 {
	const scalar = 1 / float32(1<<6)
	return float32(v) * scalar
}

func FloatToFix(f float32) fixed.Int26_6 {
	const scalar = float32(1 << 6)
	return fixed.Int26_6(f * scalar)
}
