package physics

import (
	"runtime"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Object interface {
	pullState()
	pushState()
}

type World struct {
	object.Component
	handle worldHandle
	debug  bool

	objects []Object
}

func NewWorld() *World {
	handle := world_new()
	world := object.NewComponent(&World{
		handle: handle,
	})
	runtime.SetFinalizer(world, func(w *World) {
		world_delete(&world.handle)
	})
	return world
}

func (w *World) Update(scene object.Component, dt float32) {
	// todo: optimize using bullet's MotionState
	for _, obj := range w.objects {
		obj.pushState()
	}
	world_step_simulation(w.handle, dt)
	for _, obj := range w.objects {
		obj.pullState()
	}
}

func (w *World) OnEnable() {
}

func (w *World) SetGravity(gravity vec3.T) {
	world_gravity_set(w.handle, gravity)
}

func (w *World) addObject(obj Object) bool {
	w.objects = append(w.objects, obj)
	return true
}

func (w *World) removeObject(obj Object) bool {
	for i, o := range w.objects {
		if o == obj {
			w.objects = append(w.objects[:i], w.objects[i+1:]...)
			return true
		}
	}
	return false
}

func (w *World) addRigidBody(body *RigidBody) {
	if w.addObject(body) {
		world_add_rigidbody(w.handle, body.handle)
	}
}

func (w *World) removeRigidBody(body *RigidBody) {
	if w.removeObject(body) {
		world_remove_rigidbody(w.handle, body.handle)
	}
}

func (w *World) AddCharacter(character *Character) {
	if w.addObject(character) {
		world_add_character(w.handle, character.handle)
	}
}

func (w *World) RemoveCharacter(character *Character) {
	if w.removeObject(character) {
		world_remove_character(w.handle, character.handle)
	}
}

func (w *World) Debug(enabled bool) {
	if enabled {
		world_debug_enable(w)
	} else {
		world_debug_disable(w)
	}
	w.debug = enabled
}

func (w *World) DebugDraw() {
	if w.handle != nil {
		world_debug_draw(w.handle)
	}
}
