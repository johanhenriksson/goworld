package ui

import (
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
	t.Font.Render(t.Texture, text, t.Style.Color("color", render.White))
	t.Text = text
	t.SetSize(float32(width), float32(height))
}

func NewText(text string, style Style) *Text {
	// create font
	dpi := 1.0
	size := style.Float("size", 16.0)
	spacing := style.Float("spacing", 1.0)
	fnt := render.LoadFont("assets/fonts/SourceCodeProRegular.ttf",
		float64(size), dpi, float64(spacing))

	// create opengl texture
	width, height := fnt.Measure(text)
	texture := render.CreateTexture(int32(width), int32(height))

	element := &Text{
		Image: NewImage(texture, float32(width), float32(height), true),
		Font:  fnt,
		Style: style,
	}
	element.Set(text)
	return element
}

func (t *Text) DesiredSize(w, h float32) (float32, float32) {
	dwi, dhi := t.Font.Measure(t.Text)
	dw := float32(dwi)
	dh := float32(dhi)
	if dw > w {
		dw = w
	}
	if dh > h {
		dh = h
	}
	t.SetSize(dw, dh)
	return dw, dh
}
