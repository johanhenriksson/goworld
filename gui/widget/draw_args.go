package widget

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
)

type Renderer[W T] interface {
	Draw(DrawArgs, W)
}

type DrawArgs struct {
	Time     float32
	Delta    float32
	Commands command.Recorder
	Textures cache.SamplerCache
	Viewport render.Screen
	Position vec3.T
}
