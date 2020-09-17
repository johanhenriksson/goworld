package ui

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/render"
)

type Text struct {
	*Image
	Text  string
	Font  *render.Font
	Style Style
}

func (t *Text) Set(text string) {
	if text == t.Text {
		return
	}

	width, height := t.Font.Measure(text)
	color := t.Style.Color("color", render.White)

	t.Font.Render(t.Texture, text, color)
	t.Text = text
	t.Resize(Size{width, height})
}

func NewText(text string, style Style) *Text {
	// create font
	size := style.Float("size", 16.0)
	spacing := style.Float("spacing", 1.0)
	font := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", size, spacing)

	// create opengl texture
	width, height := font.Measure(text)
	texture := render.CreateTexture(int32(width), int32(height))

	element := &Text{
		Image: NewImage(texture, float32(width), float32(height), true, style),
		Font:  font,
		Style: style,
	}
	element.Set(text)
	return element
}

func (t *Text) Flow(size Size) Size {
	dw, dh := t.Font.Measure(t.Text)
	desired := Size{dw, dh}
	desired.Width = math.Min(size.Width, desired.Width)
	desired.Height = math.Min(size.Height, desired.Height)
	return desired
}
