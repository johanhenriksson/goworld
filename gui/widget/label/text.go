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

func (t *Text) Len() int {
	return t.text.RuneCount()
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

func (t *Text) Insert(text string) {
	start, end := t.SelectedRange()
	t.text = utf8string.NewString(t.text.Slice(0, start) + text + t.text.Slice(end, t.text.RuneCount()))
	t.SetCursor(start + len(text))
	t.Deselect()
}

func (t *Text) DeleteBackward() bool {
	if !t.HasSelection() {
		t.SelectLeft()
	}
	return t.DeleteSelection()
}

func (t *Text) DeleteForward() bool {
	if !t.HasSelection() {
		t.SelectRight()
	}
	return t.DeleteSelection()
}

func (t *Text) DeleteSelection() bool {
	if !t.HasSelection() {
		return false
	}
	t.Insert("")
	return true
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

func (t *Text) Selection() string {
	start, end := t.SelectedRange()
	return t.Slice(start, end)
}

func (t *Text) SelectedRange() (start, end int) {
	start = math.Min(t.cursor, t.selstart)
	end = math.Max(t.cursor, t.selstart)
	return
}

func (t *Text) Deselect() {
	t.selstart = t.cursor
}
