package render

import (
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/swapchain"
)

// Args holds the arguments used to perform a draw pass.
// Includes the various transformation matrices and position of the camera.
type Args struct {
	VP         mat4.T
	MVP        mat4.T
	Projection mat4.T
	View       mat4.T
	Transform  mat4.T
	Position   vec3.T
	Forward    vec3.T
	Viewport   Screen
	Clear      color.T
	Context    swapchain.Context
}

type Screen struct {
	Width  int
	Height int
	Scale  float32
}

func (s Screen) Size() vec2.T {
	return vec2.NewI(s.Width, s.Height)
}

func (s Screen) NormalizeCursor(cursor vec2.T) vec2.T {
	return cursor.Div(s.Size()).Sub(vec2.New(0.5, 0.5)).Scaled(2)
}

// Apply the effects of a transform
func (d Args) Apply(t mat4.T) Args {
	d.Transform = d.Transform.Mul(&t)
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}

func (d Args) Set(t transform.T) Args {
	d.Transform = t.World()
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}
