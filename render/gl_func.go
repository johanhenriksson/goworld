package render

import (
	"github.com/johanhenriksson/goworld/render/color"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func Clear() {
	ClearColor(color.Black)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func ClearWith(color color.T) {
	ClearColor(color)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func ClearDepth() {
	gl.ClearDepth(0)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

func BindScreenBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
