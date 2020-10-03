package render

import "github.com/go-gl/gl/v4.1-core/gl"

func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func ClearDepth() {
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}
