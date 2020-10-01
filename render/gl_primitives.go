package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type GLPrimitive uint32

const (
	Triangles = GLPrimitive(gl.TRIANGLES)
	Points    = GLPrimitive(gl.POINTS)
	Lines     = GLPrimitive(gl.LINES)
)
