package player

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type ArcballCamera struct {
	object.T
	Camera   camera.T
	Distance float32

	mouselook bool
}

func NewEye() *ArcballCamera {
	distance := float32(10)
	return object.Builder(object.New(&ArcballCamera{
		Camera: object.Builder(camera.New(60, 0.1, 100, color.Black)).
			Position(vec3.New(0, 0, -distance)).
			Create(),
		Distance: distance,
	})).
		Rotation(quat.Euler(30, 45, 0)).
		Create()
}

func (p *ArcballCamera) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && e.Button() == mouse.Button2 {
		p.mouselook = true
		mouse.Lock()
		e.Consume()
	}
	if e.Action() == mouse.Release && e.Button() == mouse.Button2 {
		p.mouselook = false
		mouse.Show()
		e.Consume()
	}

	// orbit
	if e.Action() == mouse.Move && p.mouselook {
		sensitivity := vec2.New(0.045, 0.04)
		delta := e.Delta().Mul(sensitivity)

		eye := p.Transform().Rotation().Euler()

		xrot := eye.X + delta.Y
		yrot := eye.Y + delta.X

		// camera angle limits
		xrot = math.Clamp(xrot, -89.9, 89.9)
		yrot = math.Mod(yrot, 360)
		rot := quat.Euler(xrot, yrot, 0)

		p.Transform().SetRotation(rot)

		e.Consume()
	}

	// zoom
	if e.Action() == mouse.Scroll {
		p.Distance += e.Scroll().Y
		if p.Distance < 1 {
			p.Distance = 1
		}
		p.Camera.Transform().SetPosition(vec3.New(0, 0, -p.Distance))
	}
}
