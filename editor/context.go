package editor

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
)

type Context struct {
	Render renderer.T
	Camera *camera.T
	Root   object.T
	Scene  object.G
}
