package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/renderpass"
)

const MainSubpass = renderpass.Name("main")

type Pass interface {
	Name() string
	Record(command.Recorder, render.Args, object.Component)
	Destroy()
}
