package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
)

type Pass interface {
	Name() string
	Record(command.Recorder, render.Args, object.T)
	Destroy()
}
