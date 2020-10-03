package render

import "github.com/go-gl/gl/v4.1-core/gl"

type CullMode int

const (
	CullNone CullMode = iota
	CullBack
	CullFront
)

type State struct {
	Blend      bool
	BlendSrc   uint32
	BlendDst   uint32
	DepthTest  bool
	DepthMask  bool
	ClearColor Color
	CullMode   CullMode

	// Viewport dimensions
	Y      int
	X      int
	Width  int
	Height int
}

var state = State{}

func (s *State) Enable() {
	// blending
	Blend(s.Blend)
	if state.Blend {
		BlendFunc(s.BlendSrc, s.BlendDst)
	}

	// depth
	DepthTest(s.DepthTest)
	DepthMask(s.DepthMask)

	ClearColor(s.ClearColor)
	Viewport(s.X, s.Y, s.Width, s.Height)
	CullFace(s.CullMode)
}

func Blend(enabled bool) {
	if state.Blend == enabled {
		return
	}

	if enabled {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	state.Blend = enabled
}

func BlendFunc(src, dst uint32) {
	if state.BlendSrc != src || state.BlendDst != dst {
		gl.BlendFunc(src, dst)
		state.BlendSrc = src
		state.BlendDst = dst
	}
}

func DepthMask(enabled bool) {
	if state.DepthMask == enabled {
		return
	}
	gl.DepthMask(enabled)
	state.DepthMask = enabled
}

func DepthTest(enabled bool) {
	if state.DepthMask == enabled {
		return
	}
	if enabled {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func ClearColor(color Color) {
	color = color.WithAlpha(1)
	if color != state.ClearColor {
		gl.ClearColor(color.R, color.G, color.B, 1)
	}
}

func Viewport(x, y, w, h int) {
	if state.Width != w || state.Height != h || state.Y != y || state.X != x {
		state.Width = w
		state.Height = h
		state.X = x
		state.Y = y
		gl.Viewport(int32(x), int32(y), int32(w), int32(h))
	}
}

func CullFace(mode CullMode) {
	if state.CullMode == mode {
		return
	}
	switch mode {
	case CullNone:
		gl.Disable(gl.CULL_FACE)
	case CullBack:
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	case CullFront:
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.FRONT)
	}
	state.CullMode = mode
}
