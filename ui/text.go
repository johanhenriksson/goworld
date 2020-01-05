package ui

import (
	"github.com/johanhenriksson/goworld/render"
)

type Text struct {
	*Image
	Size  float32
	Text  string
	Font  *render.Font
	Color render.Color
}

func (t *Text) Set(text string) {
	t.Font.RenderOn(t.Texture, text, float32(t.Texture.Width), float32(t.Texture.Height), t.Color)
}

func (m *Manager) NewText(text string, color render.Color, x, y, z float32) *Text {
	/* TODO: calculate size of text */
	w, h := float32(290.0), float32(25.0)
	fnt := render.LoadFont("assets/fonts/SourceCodeProRegular.ttf", 12.0, 100.0, 1.5)
	texture := fnt.Render(text, w, h, color)
	img := m.NewImage(texture, x, y, w, h, z)

	return &Text{
		Image: img,
		Font:  fnt,
		Text:  text,
		Color: color,
		Size:  h,
	}
}
