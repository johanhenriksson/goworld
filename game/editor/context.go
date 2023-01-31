package editor

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/engine/renderer"
)

type Context struct {
	Render renderer.T
	Camera camera.T
}
