package window

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/johanhenriksson/goworld/render/swapchain"
)

type GlfwBackend interface {
	GlfwHints(Args) []GlfwHint
	GlfwSetup(*glfw.Window, Args) error
	Resize(int, int)
	Aquire() (swapchain.Context, error)
	Present()
	Destroy()
}

type GlfwHint struct {
	Hint  glfw.Hint
	Value int
}
