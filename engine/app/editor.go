package app

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
)

func RunEditor(args Args, scenefunc object.SceneFunc) {
	Run(args, editor.Scene(scenefunc))
}
