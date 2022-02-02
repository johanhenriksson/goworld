package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Primitive uint32

const (
	Triangles = Primitive(gl.TRIANGLES)
	Points    = Primitive(gl.POINTS)
	Lines     = Primitive(gl.LINES)
)
