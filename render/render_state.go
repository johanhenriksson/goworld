package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/math/vec2"
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

type rect struct {
	X, Y, W, H int
}

type Viewport struct {
	X, Y          int
	Width, Height int
}

// Blend Functions
const (
	One              = BlendValue(gl.ONE)
	Zero             = BlendValue(gl.ZERO)
	SrcAlpha         = BlendValue(gl.SRC_ALPHA)
	OneMinusSrcAlpha = BlendValue(gl.ONE_MINUS_SRC_ALPHA)
)

type State struct {
	Blend         bool
	BlendSrc      BlendValue
	BlendDst      BlendValue
	DepthTest     bool
	DepthOutput   bool
	ClearColor    color.T
	CullMode      CullMode
	ScissorEnable bool
	Scissor       rect
	Viewport      Viewport
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
	SetViewport(s.Viewport)
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
	Blend(true)
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

func SetViewport(vp Viewport) {
	if state.Viewport != vp {
		gl.Viewport(int32(vp.X), int32(vp.Y), int32(vp.Width), int32(vp.Height))
		state.Viewport = vp
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

func Scissor(pos vec2.T, size vec2.T) {
	area := rect{
		X: int(pos.X), Y: int(pos.Y),
		W: int(size.X), H: int(size.Y),
	}
	if !state.ScissorEnable {
		gl.Enable(gl.SCISSOR_TEST)
		state.ScissorEnable = true
	}
	if state.Scissor != area {
		gl.Scissor(int32(area.X), int32(area.Y), int32(area.W), int32(area.H))
		state.Scissor = area
	}
}

func ScissorDisable() {
	if state.ScissorEnable {
		gl.Disable(gl.SCISSOR_TEST)
		state.ScissorEnable = false
	}
}
