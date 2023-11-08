package player

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

type ArcballCamera struct {
	object.Object
	Camera   *camera.Object
	Distance float32

	mouselook bool
}

func NewEye() *ArcballCamera {
	distance := float32(10)
	return object.Builder(object.New("Arcball", &ArcballCamera{
		Camera: object.Builder(camera.NewObject(camera.Args{
			Fov:   58,
			Near:  0.1,
			Far:   1000,
			Clear: color.Black,
		})).
			Position(vec3.New(0, 0, -distance)).
			Attach(NewSkydome()).
			Create(),

		Distance: distance,
	})).
		Position(vec3.New(0, 0.8, 0)).
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

	world := object.GetInParents[*physics.World](p)
	if world != nil {
		hit, exists := world.Raycast(p.Transform().WorldPosition(), p.Camera.Transform().WorldPosition(), physics.All)
		if exists {
			maxDist := vec3.Distance(p.Transform().WorldPosition(), hit.Point)
			if p.Distance > maxDist {
				p.Distance = maxDist
				p.Camera.Transform().SetPosition(vec3.New(0, 0, -p.Distance))
			}
		}
	}

	// orbit
	if e.Action() == mouse.Move && p.mouselook {
		sensitivity := vec2.New(0.1, 0.1)
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
		if p.Distance < 0 {
			p.Distance = 0
		}
		p.Camera.Transform().SetPosition(vec3.New(0, 0, -p.Distance))
	}
}
