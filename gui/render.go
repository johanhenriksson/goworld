package gui

import (
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/widget"
)

func Render(Root func() widget.T) widget.T {
	hooks.Reset()
	root := Root()
	return root
}
