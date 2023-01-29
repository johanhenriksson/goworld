package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/sync"
)

type Pass interface {
	Name() string
	Record(command.Recorder, render.Args, object.T)
	Draw(args render.Args, scene object.T)
	Completed() sync.Semaphore
	Destroy()
}
