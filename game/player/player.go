package player

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type T struct {
	object.T
	Character *physics.Character
	Camera    *ArcballCamera
	Speed     float32

	keys     keys.State
	velocity vec3.T
}

func New() *T {
	return object.New(&T{
		Character: physics.NewCharacter(1.8, 0.5, 0.2),
		Camera:    NewEye(),
		Speed:     0.05,
		keys:      keys.NewState(),
	})
}

func (p *T) Update(scene object.T, dt float32) {
	p.T.Update(scene, dt)
	if !p.Character.Grounded() {
		p.velocity.Scale(1 - 0.9*dt)
	}
	p.Character.Move(p.velocity)
}

func (p *T) Move(dir vec3.T) {
	if p.Character.Grounded() {
		p.velocity = dir
	}
}

func (p *T) KeyEvent(e keys.Event) {
	p.keys.KeyEvent(e)

	if p.Character.Grounded() && e.Code() == keys.Space {
		p.Character.Jump()
	}

	forward, right := float32(0), float32(0)
	if p.keys.Down(keys.D) {
		right += 1
	}
	if p.keys.Down(keys.A) {
		right -= 1
	}
	if p.keys.Down(keys.W) {
		forward += 1
	}
	if p.keys.Down(keys.S) {
		forward -= 1
	}

	camFwd := p.Camera.Transform().Forward()
	camFwd.Y = 0
	camFwd.Normalize()

	camRight := p.Camera.Transform().Right()
	camRight.Y = 0
	camRight.Normalize()

	dirForward := camFwd.Scaled(forward)
	dirRight := camRight.Scaled(right)
	dir := dirForward.Add(dirRight)
	p.Move(dir.Scaled(p.Speed))
}
