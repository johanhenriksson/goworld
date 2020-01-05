package ui

import (
	"fmt"
	"unicode/utf8"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
)

type Textbox struct {
	*Image
	Size  float32
	Text  string
	Font  *render.Font
	Color render.Color

	focused bool
}

func (t *Textbox) Set(text string) {
	t.Text = text
	if t.focused {
		text += "_"
	}
	t.Font.RenderOn(t.Texture, text, float32(t.Texture.Width), float32(t.Texture.Height), t.Color)
}

func (m *Manager) NewTextbox(text string, color render.Color, x, y, z float32) *Textbox {
	/* TODO: calculate size of text */
	w, h := float32(290.0), float32(25.0)
	fnt := render.LoadFont("assets/fonts/SourceCodeProRegular.ttf", 12.0, 100.0, 1.5)
	texture := fnt.Render(text, w, h, color)
	img := m.NewImage(texture, x, y, w, h, z)
	img.Quad.FlipY()

	t := &Textbox{
		Image: img,
		Font:  fnt,
		Text:  text,
		Color: color,
		Size:  h,
	}
	t.OnClick(func(ev MouseEvent) {
		fmt.Println("caught input focus")
		ev.UI.Focus(t)
	})
	return t
}

func (t *Textbox) Append(text string) {
	t.Set(t.Text + text)
}

func (t *Textbox) backspace() {
	r, size := utf8.DecodeLastRuneInString(t.Text)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		return
	}
	if len(t.Text) >= size {
		newText := t.Text[:len(t.Text)-size]
		t.Set(newText)
	}
}

func (t *Textbox) HandleInput(char rune) {
	t.Append(string(char))
}

func (t *Textbox) HandleKey(event KeyEvent) {
	if !event.Press {
		return
	}

	// backspace
	if event.Key == engine.KeyBackspace {
		t.backspace()
	}

	// drop focus on esc
	if event.Key == engine.KeyEscape {
		event.UI.Focus(nil)
	}
}

func (t *Textbox) Focus() {
	t.focused = true
	t.Set(t.Text)
}

func (t *Textbox) Blur() {
	t.focused = false
	t.Set(t.Text)
}
