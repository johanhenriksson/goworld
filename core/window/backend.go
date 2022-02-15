package window

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GlfwBackend interface {
	GlfwHints(Args) []GlfwHint
	GlfwSetup(*glfw.Window, Args) error
	Resize(int, int)
}

type GlfwHint struct {
	Hint  glfw.Hint
	Value int
}

type OpenGLBackend struct {
	window *glfw.Window
}

func (b *OpenGLBackend) GlfwHints(args Args) []GlfwHint {
	hints := []GlfwHint{
		{glfw.ContextVersionMajor, 4},
		{glfw.ContextVersionMinor, 1},
		{glfw.OpenGLProfile, glfw.OpenGLCoreProfile},
		{glfw.OpenGLForwardCompatible, glfw.True},
		{glfw.Samples, 1},
	}

	if args.Debug {
		hints = append(hints, GlfwHint{glfw.OpenGLDebugContext, glfw.True})
	}

	return hints
}

func (b *OpenGLBackend) GlfwSetup(w *glfw.Window, args Args) error {
	b.window = w

	if args.Vsync {
		glfw.SwapInterval(1)
	}

	// activate OpenGL context
	w.MakeContextCurrent()

	// initialize OpenGL. context must be active first
	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	// set up debugging
	if args.Debug {
		var flags int32
		gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
		if flags&gl.CONTEXT_FLAG_DEBUG_BIT == gl.CONTEXT_FLAG_DEBUG_BIT {
			gl.Enable(gl.DEBUG_OUTPUT)
			gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
			gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
			gl.DebugMessageCallback(b.onDebugMessage, nil)
		} else {
			fmt.Println("warning: failed to enable opengl debugging")
		}
	}

	return nil
}

func (b *OpenGLBackend) Resize(width, height int) {

}

func (b *OpenGLBackend) Aquire() {

}

func (b *OpenGLBackend) Present() {
	b.window.SwapBuffers()
}

func (b *OpenGLBackend) onDebugMessage(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {
	// todo: proper messages
	// see https://learnopengl.com/In-Practice/Debugging
	fmt.Println("GL Debug:", message)
}
