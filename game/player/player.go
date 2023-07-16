package player

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type T struct {
	object.G
	Character *physics.Character
	Camera    *ArcballCamera
	Speed     float32
	TurnRate  float32

	keys     keys.State
	mouse    mouse.State
	velocity vec3.T
}

func New() *T {
	return object.Group("Player", &T{
		Character: physics.NewCharacter(1.8, 0.5, 0.2),
		Camera:    NewEye(),
		Speed:     7,
		TurnRate:  40,
		keys:      keys.NewState(),
		mouse:     mouse.NewState(),
	})
}

func (p *T) Name() string {
	return "Player"
}

func (p *T) Update(scene object.Component, dt float32) {
	p.G.Update(scene, dt)

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
	}
	p.Character.Move(p.velocity)
}

func (p *T) KeyEvent(e keys.Event) {
	p.G.KeyEvent(e)
	p.keys.KeyEvent(e)

	if p.Character.Grounded() && e.Code() == keys.Space {
		p.Character.Jump()
	}
}

func (p *T) MouseEvent(e mouse.Event) {
	p.G.MouseEvent(e)
	p.mouse.MouseEvent(e)
}
