package client

import (
	"time"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/player"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type LocalController struct {
	object.Object
	Target *Entity

	Character *physics.Character
	Camera    *player.ArcballCamera

	Speed    float32
	TurnRate float32

	mgr      *Manager
	keys     keys.State
	mouse    mouse.State
	velocity vec3.T
	tickRate time.Duration
	lastTick time.Time
	moving   bool
	lastPos  vec3.T
	rotating bool
	lastRot  float32
}

func NewLocalController() *LocalController {
	updatesPerSecond := 20
	return object.New("LocalController", &LocalController{
		Character: physics.NewCharacter(1, 0.5, 0.1),
		Camera:    player.NewEye(),

		Speed:    7,
		TurnRate: 40,
		keys:     keys.NewState(),
		mouse:    mouse.NewState(),
		tickRate: time.Second / time.Duration(updatesPerSecond),
		lastTick: time.Now(),
	})
}

func (p *LocalController) Observe(entity *Entity) {
	p.Target = entity
	p.velocity = vec3.Zero

	if entity != nil {
		p.Transform().SetPosition(entity.Transform().Position())
		object.Enable(p)
	} else {
		object.Disable(p)
	}
}

func (p *LocalController) Update(scene object.Component, dt float32) {
	p.Object.Update(scene, dt)

	if p.mgr == nil {
		var ok bool
		p.mgr, ok = p.Parent().(*Manager)
		if !ok {
			return
		}
	}

	p.updateMovement(dt)
}

func (p *LocalController) Shoot() {
	// shoot projectile

	dir := p.Target.Transform().Forward().Scaled(8)
	proj := object.Builder(NewProjectile(dir)).
		Position(p.Target.Transform().Position()).
		Create()
	object.Attach(p.mgr, proj)
}

func (p *LocalController) updateMovement(dt float32) {
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

	pos := p.Transform().Position()
	rotY := p.Camera.Transform().Rotation().Euler().Y

	// update target
	if p.Target != nil {
		p.Target.Transform().SetPosition(pos)
		p.Target.Transform().SetRotation(quat.Euler(0, rotY, 0))

		moving := vec3.Distance(pos, p.lastPos) > 0.01
		p.Target.Sprite.animating = moving

		sinceUpdate := time.Now().Sub(p.lastTick)
		if sinceUpdate >= p.tickRate {
			rotating := math.Abs(rotY-p.lastPos.Y) > 0.001
			stopped := (!moving && p.moving) || (!rotating && p.rotating)
			p.moving = moving
			p.rotating = rotating
			p.lastPos = pos
			p.lastRot = rotY
			p.mgr.Client.SendMove(p.Target.EntityID(), pos, rotY, stopped, float32(p.tickRate.Seconds()))
			p.lastTick = time.Now().Add(p.tickRate - sinceUpdate)
		}
	}
}

func (p *LocalController) KeyEvent(e keys.Event) {
	p.Object.KeyEvent(e)
	p.keys.KeyEvent(e)

	if p.Character.Grounded() && e.Code() == keys.Space {
		p.Character.Jump()
	}

	if e.Action() == keys.Press && e.Code() == keys.Key1 {
		p.Shoot()
	}
}

func (p *LocalController) MouseEvent(e mouse.Event) {
	p.Object.MouseEvent(e)
	p.mouse.MouseEvent(e)
}
