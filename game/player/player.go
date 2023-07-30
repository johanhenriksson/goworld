package player

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
)

type T struct {
	object.Object
	Character *physics.Character
	Camera    *ArcballCamera
	Model     object.Object

	Speed    float32
	TurnRate float32

	keys     keys.State
	mouse    mouse.State
	velocity vec3.T
}

func New() *T {
	model := cube.New(cube.Args{
		Size: 1,
		Mat:  material.StandardDeferred(),
	})
	model.SetTexture("diffuse", color.White)

	return object.New("Player", &T{
		Character: physics.NewCharacter(1, 0.5, 0.2),
		Camera:    NewEye(),
		Model: object.Builder(object.Empty("Model")).
			Scale(vec3.New(1, 2, 1)).
			Attach(model).
			Create(),

		Speed:    7,
		TurnRate: 40,
		keys:     keys.NewState(),
		mouse:    mouse.NewState(),
	})
}

func (p *T) Name() string {
	return "Player"
}

func (p *T) Update(scene object.Component, dt float32) {
	p.Object.Update(scene, dt)

	forward, right := float32(0), float32(0)
	mouseMove := p.mouse.Down(mouse.Button1) && p.mouse.Down(mouse.Button2)
	if p.keys.Down(keys.D) {
		right += 1
	}
	if p.keys.Down(keys.A) {
		right -= 1
	}
	if p.keys.Down(keys.W) || mouseMove {
		forward += 1
	}
	if p.keys.Down(keys.S) {
		forward -= 1
	}

	rotate := float32(0)
	if p.keys.Down(keys.LeftArrow) {
		rotate -= 1
	}
	if p.keys.Down(keys.RightArrow) {
		rotate += 1
	}

	// apply keyboard turning
	rot := p.Camera.Transform().Rotation().Euler()
	rot.Y += rotate * p.TurnRate * dt
	p.Camera.Transform().SetRotation(quat.Euler(rot.X, rot.Y, rot.Z))

	// calculate forward & right vectors relative to camera
	camFwd := p.Camera.Transform().Forward()
	camFwd.Y = 0
	camFwd.Normalize()

	camRight := p.Camera.Transform().Right()
	camRight.Y = 0
	camRight.Normalize()

	// compute movement direction
	dirForward := camFwd.Scaled(forward)
	dirRight := camRight.Scaled(right)
	dir := dirForward.Add(dirRight).Normalized()

	if p.Character.Grounded() {
		p.velocity = dir.Scaled(p.Speed)
	} else {
		// the player is allowed some air acceleration
		// ensure the total velocity does not exceed the maximum speed
		p.velocity = p.velocity.Add(dir.Scaled(0.016 * p.Speed))
		if p.velocity.Length() > p.Speed {
			p.velocity = p.velocity.Normalized().Scaled(p.Speed)
		}
	}
	p.Character.Move(p.velocity)
}

func (p *T) KeyEvent(e keys.Event) {
	p.Object.KeyEvent(e)
	p.keys.KeyEvent(e)

	if p.Character.Grounded() && e.Code() == keys.Space {
		p.Character.Jump()
	}
}

func (p *T) MouseEvent(e mouse.Event) {
	p.Object.MouseEvent(e)
	p.mouse.MouseEvent(e)
}
