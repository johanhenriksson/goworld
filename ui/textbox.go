package ui

import (
	"fmt"
	"unicode/utf8"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
)

type Textbox struct {
	*Image
	Text string
	Font *render.Font

	focused bool
}

func (t *Textbox) Set(text string) {
	t.Text = text
	if t.focused {
		text += "_"
	}
	t.Font.Render(t.Texture, text, t.Color("color", render.White))
}

func NewTextbox(text string, style Style) *Textbox {
	size := style.Float("size", 12.0)
	spacing := style.Float("spacing", 1.5)
	font := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", size, spacing)
	w, h := font.Measure(text)
	texture := render.CreateTexture(int32(w), int32(h))

	t := &Textbox{
		Image: NewImage(texture, float32(w), float32(h), true, style),
		Font:  font,
		Text:  text,
	}
	t.OnClick(func(ev MouseEvent) {
		fmt.Println("caught input focus")
		ev.UI.Focus(t)
	})
	t.Set(text)
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
