package draw

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/command"
)

type Pass interface {
	Name() string
	Record(command.Recorder, Args, object.Component)
	Destroy()
}
