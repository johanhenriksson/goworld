package pass

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/swapchain"
)

type Pass interface {
	Name() string
	Record(command.Recorder, render.Args, object.Component)
	Destroy()
}

type Args struct {
	Camera    *camera.T
	Transform mat4.T
	MVP       mat4.T
	Viewport  render.Screen
	Context   *swapchain.Context
}
