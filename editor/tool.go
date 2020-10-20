package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Tool interface {
	engine.Component
	Use(*Editor, vec3.T, vec3.T)
	Hover(*Editor, vec3.T, vec3.T)
}
