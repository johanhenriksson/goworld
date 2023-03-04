package label

import (
	"github.com/johanhenriksson/goworld/math"
	"golang.org/x/exp/utf8string"
)

const CursorBlinkInterval = float32(0.5)

type Text struct {
	text    *utf8string.String
	cursor  int
	blink   bool
	blinkDt float32
}

func NewText(text string) *Text {
	t := &Text{}
	t.SetText(text)
	return t
}

func (t *Text) String() string {
	return t.text.String()
}

func (t *Text) Slice(i, j int) string {
	return t.text.Slice(i, j)
}

func (t *Text) SetText(text string) {
	t.text = utf8string.NewString(text)
	t.SetCursor(len(text))
}

func (t *Text) SetCursor(cursor int) {
	t.cursor = math.Clamp(cursor, 0, t.text.RuneCount())
	t.ResetBlink()
}

func (t *Text) Blink() bool {
	return t.blink
}

func (t *Text) UpdateBlink(delta float32) {
	t.blinkDt -= delta
	if t.blinkDt < 0 {
		t.blink = !t.blink
		t.blinkDt = CursorBlinkInterval
	}
}

func (t *Text) ResetBlink() {
	t.blink = true
	t.blinkDt = CursorBlinkInterval
}

func (t *Text) Insert(char rune) {
	t.text = utf8string.NewString(t.text.Slice(0, t.cursor) + string(char) + t.text.Slice(t.cursor, t.text.RuneCount()))
	t.SetCursor(t.cursor + 1)
}

func (t *Text) DeleteBackward() bool {
	if t.cursor > 0 {
		t.SetCursor(t.cursor - 1)
		t.text = utf8string.NewString(t.text.Slice(0, t.cursor) + t.text.Slice(t.cursor+1, t.text.RuneCount()))
		return true
	}
	return false
}

func (t *Text) DeleteForward() bool {
	if t.cursor < t.text.RuneCount() {
		t.text = utf8string.NewString(t.text.Slice(0, t.cursor) + t.text.Slice(t.cursor+1, t.text.RuneCount()))
		return true
	}
	return false
}

func (t *Text) CursorLeft() {
	t.SetCursor(t.cursor - 1)
}

func (t *Text) CursorRight() {
	t.SetCursor(t.cursor + 1)
}

func (t *Text) Clear() bool {
	changed := t.text.String() != ""
	t.SetText("")
	return changed
}
