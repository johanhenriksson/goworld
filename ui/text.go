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

	size := t.Font.Measure(text)

	img := t.Font.Render(text, color.White)
	t.Texture.BufferImage(img)

	t.Text = text
	t.Resize(size)
}

func NewText(text string, style Style) *Text {
	// create font
	size := style.Float("size", 16.0)
	spacing := style.Float("spacing", 1.0)
	font := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", size, spacing)

	// create opengl texture
	bounds := font.Measure(text)
	texture := gltex.New(int(bounds.X), int(bounds.Y))

	element := &Text{
		Image: NewImage(texture, bounds, false, style),
		Font:  font,
		Style: style,
	}
	element.Set(text)
	return element
}

func (t *Text) Flow(size vec2.T) vec2.T {
	desired := t.Font.Measure(t.Text)
	desired.X = math.Min(size.X, desired.X)
	desired.Y = math.Min(size.Y, desired.Y)
	return desired
}
