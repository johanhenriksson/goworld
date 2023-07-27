package editor

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
)

type Context struct {
	Camera *camera.Camera
	Root   object.Component
	Scene  object.Object
}
