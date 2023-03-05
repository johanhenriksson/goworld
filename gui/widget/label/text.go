package label

import (
	"github.com/johanhenriksson/goworld/math"
	"golang.org/x/exp/utf8string"
)

const CursorBlinkInterval = float32(0.5)

type Text struct {
	text     *utf8string.String
	cursor   int // cursor position
	selstart int // selection start offset
	blink    bool
	blinkDt  float32
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

	// todo: try to preserve selection
	t.Deselect()
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
	start := math.Min(t.cursor, t.selstart)
	end := math.Max(t.cursor, t.selstart)

	t.text = utf8string.NewString(t.text.Slice(0, start) + string(char) + t.text.Slice(end, t.text.RuneCount()))
	t.SetCursor(start + 1)
	t.Deselect()
}

func (t *Text) DeleteBackward() bool {
	if t.HasSelection() {
		// delete selection
		start := math.Min(t.cursor, t.selstart)
		end := math.Max(t.cursor, t.selstart)
		t.text = utf8string.NewString(t.text.Slice(0, start) + t.text.Slice(end, t.text.RuneCount()))
		t.SetCursor(start)
		t.Deselect()
		return true
	} else {
		// no selection - normal backspace behavior
		if t.cursor > 0 {
			t.SetCursor(t.cursor - 1)
			t.text = utf8string.NewString(t.text.Slice(0, t.cursor) + t.text.Slice(t.cursor+1, t.text.RuneCount()))
			t.Deselect()
			return true
		}
	}
	return false
}

func (t *Text) DeleteForward() bool {
	if t.HasSelection() {
		// delete selection
		start := math.Min(t.cursor, t.selstart)
		end := math.Max(t.cursor, t.selstart)
		t.text = utf8string.NewString(t.text.Slice(0, start) + t.text.Slice(end, t.text.RuneCount()))
		t.SetCursor(start)
		t.Deselect()
		return true
	} else {
		// no selection - normal forward delete behavior
		if t.cursor < t.text.RuneCount() {
			t.text = utf8string.NewString(t.text.Slice(0, t.cursor) + t.text.Slice(t.cursor+1, t.text.RuneCount()))
			t.Deselect()
			return true
		}
	}
	return false
}

func (t *Text) CursorLeft() {
	if t.HasSelection() {
		// move to the lower end of the selection
		t.SetCursor(math.Min(t.cursor, t.selstart))
	} else {
		// if nothing is selected, move cursor 1 step left
		t.SetCursor(t.cursor - 1)
	}
	t.Deselect()
}

func (t *Text) CursorRight() {
	if t.HasSelection() {
		// move to the higher end of the selection
		t.SetCursor(math.Max(t.cursor, t.selstart))
	} else {
		// if nothing is selected, move cursor 1 step right
		t.SetCursor(t.cursor + 1)
	}
	t.Deselect()
}

func (t *Text) SelectLeft() {
	t.SetCursor(t.cursor - 1)
}

func (t *Text) SelectRight() {
	t.SetCursor(t.cursor + 1)
}

func (t *Text) Clear() bool {
	changed := t.text.String() != ""
	t.SetText("")
	return changed
}

func (t *Text) HasSelection() bool {
	return t.cursor != t.selstart
}

func (t *Text) Deselect() {
	t.selstart = t.cursor
}
