package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render/color"
)

type CullMode int

// Culling modes
const (
	CullNone CullMode = iota
	CullBack
	CullFront
)

type BlendValue uint32

// Blend Functions
const (
	One              = BlendValue(gl.ONE)
	Zero             = BlendValue(gl.ZERO)
	SrcAlpha         = BlendValue(gl.SRC_ALPHA)
	OneMinusSrcAlpha = BlendValue(gl.ONE_MINUS_SRC_ALPHA)
)

type State struct {
	Blend       bool
	BlendSrc    BlendValue
	BlendDst    BlendValue
	DepthTest   bool
	DepthOutput bool
	ClearColor  color.T
	CullMode    CullMode

	// Viewport dimensions
	Y      int
	X      int
	Width  int
	Height int
}

var state = State{
	ClearColor: color.Black,
}

func (s *State) Enable() {
	// blending
	Blend(s.Blend)
	if state.Blend {
		BlendFunc(s.BlendSrc, s.BlendDst)
	}

	// depth
	DepthTest(s.DepthTest)
	DepthOutput(s.DepthOutput)

	ClearColor(s.ClearColor)
	SetViewport(s.X, s.Y, s.Width, s.Height)
	CullFace(s.CullMode)
}

type RenderFunc func()

func (s *State) Use(scope RenderFunc) {
	// copy current state
	previous := state
	s.Enable()

	// run render
	scope()

	// reset render state
	previous.Enable()
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

func BlendFunc(src, dst BlendValue) {
	if state.BlendSrc != src || state.BlendDst != dst {
		gl.BlendFunc(uint32(src), uint32(dst))
		state.BlendSrc = src
		state.BlendDst = dst
	}
}

func BlendAdditive() {
	Blend(true)
	BlendFunc(One, One)
}

func BlendMultiply() {
	BlendFunc(SrcAlpha, OneMinusSrcAlpha)
}

func DepthOutput(enabled bool) {
	if state.DepthOutput == enabled {
		return
	}
	gl.DepthMask(enabled)
	state.DepthOutput = enabled
}

func DepthTest(enabled bool) {
	if state.DepthTest == enabled {
		return
	}
	if enabled {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	state.DepthTest = enabled
}

func ClearColor(color color.T) {
	color = color.WithAlpha(1)
	// if color != state.ClearColor {
	gl.ClearColor(color.R, color.G, color.B, 1)
	state.ClearColor = color
	// }
}

func SetViewport(x, y, w, h int) {
	if state.Width != w || state.Height != h || state.Y != y || state.X != x {
		gl.Viewport(int32(x), int32(y), int32(w), int32(h))
		state.Width = w
		state.Height = h
		state.X = x
		state.Y = y
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
