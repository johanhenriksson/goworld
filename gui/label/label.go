package label

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Label interface {
	widget.T

	Text() string
	SetText(string)
}

type TextMeasurer interface {
	Measure(string) vec2.T
}

type label struct {
	widget.T

	text string
}

func NewLabel() Label {
	return &label{
		// Widget: NewWidget(),
	}
}

func (l *label) Text() string        { return l.text }
func (l *label) SetText(text string) { l.text = text }
