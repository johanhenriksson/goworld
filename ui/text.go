package ui

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	gltex "github.com/johanhenriksson/goworld/render/backend/gl/gl_texture"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
)

type Text struct {
	*Image
	Text  string
	Font  font.T
	Style Style
}

func (t *Text) Set(text string) {
	if text == t.Text {
		return
	}

	args := font.Args{
		LineHeight: t.Spacing(),
		Color:      color.White,
	}

	size := t.Font.Measure(text, args)
	img := t.Font.Render(text, args)
	t.Texture.BufferImage(img)

	t.Text = text
	t.Resize(size.Scaled(0.5))
}

func NewText(text string, style Style) *Text {
	// create font
	size := int(style.Float("size", 16.0))
	spacing := style.Float("spacing", 1.0)
	fnt := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", size*2)

	// create opengl texture
	bounds := fnt.Measure(text, font.Args{
		LineHeight: spacing,
	})
	texture := gltex.New(int(bounds.X), int(bounds.Y))

	element := &Text{
		Image: NewImage(texture, bounds.Scaled(0.5), false, style),
		Font:  fnt,
		Style: style,
	}
	element.Set(text)
	return element
}

func (t *Text) Spacing() float32 {
	return t.Style.Float("spacing", 1.0)
}
func (t *Text) Size() int {
	return int(t.Style.Float("size", 1.0))
}

func (t *Text) Flow(size vec2.T) vec2.T {
	desired := t.Font.Measure(t.Text, font.Args{
		LineHeight: t.Spacing(),
	}).Scaled(0.5)
	desired.X = math.Min(size.X, desired.X)
	desired.Y = math.Min(size.Y, desired.Y)
	return desired
}
