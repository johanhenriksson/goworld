package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Tool interface {
	Use(*Editor, vec3.T, vec3.T)
	Update(*Editor, float32, vec3.T, vec3.T)
	Draw(*Editor, engine.DrawArgs)
}
