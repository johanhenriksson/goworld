package gui

import (
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/rect"
)

func Render(Root func() rect.T) rect.T {
	hooks.Reset()
	root := Root()
	return root
}
