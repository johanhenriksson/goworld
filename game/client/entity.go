package client

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type Entity struct {
	object.Object
	Sprite    *CharacterSprite
	Nameplate gui.Fragment
	Height    float32
	Health    float32

	id  server.Identity
	cam *camera.Camera

	moveFrom   vec3.T
	moveTo     vec3.T
	moveVel    vec3.T
	rotFrom    float32
	rotTo      float32
	rotVel     float32
	lastUpdate time.Time
	duration   time.Duration
	animating  bool
	rotating   bool
	stopAfter  bool
}

func NewEntity(id server.Identity, pos vec3.T, rot float32) *Entity {
	height := float32(2)
	spriteIndex := int(id % 365)
	sprite := NewCharacterSprite(spriteIndex, height)

	entity := object.New(fmt.Sprintf("Entity %x", id), &Entity{
		Sprite: sprite,
		Height: height,

		id:         id,
		animating:  false,
		moveFrom:   pos,
		rotFrom:    rot,
		lastUpdate: time.Now(),
	})

	entity.Nameplate = gui.NewFragment(gui.FragmentArgs{
		Slot:   "plates",
		Render: entity.renderNameplate,
	})
	object.Attach(entity, entity.Nameplate)

	return entity
}

func (e *Entity) EntityID() server.Identity {
	return e.id
}

func (e *Entity) Move(ev EntityMoveEvent) {
	e.rotFrom = e.Transform().Rotation().Euler().Y
	e.moveFrom = e.Transform().Position()
	e.rotTo = ev.Rotation
	e.moveTo = ev.Position
	e.stopAfter = ev.Stopped
	e.duration = time.Duration(ev.Delta * float32(time.Second))
	e.lastUpdate = time.Now()

	e.animating = true // vec3.Distance(e.moveFrom, e.moveTo) > 0.01 || math.Abs(e.rotFrom-ev.Rotation) > 0.001
	e.moveVel = e.moveTo.Sub(e.moveFrom).Scaled(1 / float32(e.duration.Seconds()))
	e.rotVel = (e.rotTo - e.rotFrom) / float32(e.duration.Seconds())
}

func (e *Entity) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.cam = object.GetInChildren[*camera.Camera](scene)

	e.Health -= 0.25 * dt
	if e.Health < 0 {
		e.Health = 1
	}

	// update sprite animation state
	e.Sprite.animating = e.moveVel.LengthSqr() > 0.3

	if e.animating {
		elapsed := time.Now().Sub(e.lastUpdate)
		f := float32(elapsed.Seconds() / e.duration.Seconds())

		e.Transform().SetPosition(e.Transform().Position().Add(e.moveVel.Scaled(dt)))
		e.Transform().SetRotation(quat.Euler(0, e.Transform().Rotation().Euler().Y+e.rotVel*dt, 0))

		// if we reached the end of the move, stop
		done := f >= 1
		if e.stopAfter && done {
			e.animating = false
		}
	}
}

func (e *Entity) renderNameplate() node.T {
	cam := object.GetInChildren[*camera.Camera](object.Root(e))
	if cam == nil {
		return nil
	}

	viewPortPos := cam.Project(e.Transform().WorldPosition().Add(vec3.New(0, e.Height+0.25, 0)))

	const plateWidth = 120
	return rect.New(fmt.Sprintf("plate_%x", e.EntityID()), rect.Props{
		Style: rect.Style{
			Position: style.Absolute{
				Left: style.Px(viewPortPos.X - plateWidth/2),
				Top:  style.Px(viewPortPos.Y),
			},
			Width: style.Px(plateWidth),
			Color: color.Black,
			Border: style.Border{
				Width: style.Px(2),
				Color: color.Black,
			},
		},
		Children: []node.T{
			rect.New("hp", rect.Props{
				Style: rect.Style{
					Position: style.Relative{},
					Color:    color.Lerp(color.Red, color.Green, e.Health),
					Height:   style.Pct(100),
					Width:    style.Pct(100 * e.Health),
				},
			}),
			rect.New("name", rect.Props{
				Style: rect.Style{
					AlignItems: style.AlignCenter,
				},
				Children: []node.T{
					label.New("name", label.Props{
						Text: fmt.Sprintf("%x", e.EntityID()),
						Style: label.Style{
							Color: color.White,
						},
					}),
				},
			}),
		},
	})
}
