package physics

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
)

type RigidBody struct {
	object.Component

	world    *World
	handle   rigidbodyHandle
	mass     float32
	tfparent transform.T

	Shape Shape
}

func NewRigidBody(mass float32) *RigidBody {
	body := object.NewComponent(&RigidBody{
		mass: mass,
	})
	runtime.SetFinalizer(body, func(b *RigidBody) {
		b.destroy()
	})
	return body
}

func (b *RigidBody) pullState() {
	if b.handle == nil {
		return
	}
	if b.Kinematic() {
		return
	}
	state := rigidbody_state_pull(b.handle)
	b.Transform().SetWorldPosition(state.position)
	b.Transform().SetWorldRotation(state.rotation)
}

func (b *RigidBody) pushState() {
	if b.handle == nil {
		return
	}
	rigidbody_state_push(b.handle,
		b.Transform().WorldPosition(),
		b.Transform().WorldRotation(),
	)
}

func (b *RigidBody) OnEnable() {
	if b.Shape == nil {
		b.Shape = object.Get[Shape](b)
		if b.Shape == nil {
			log.Println("Rigidbody", b.Parent().Name(), ": no shape in siblings")
			return
		}

		b.Shape.OnChange().Subscribe(b, func(s Shape) {
			if b.handle == nil {
				panic("rigidbody shape set to nil")
			}
			rigidbody_shape_set(b.handle, s.shape())
		})
	}

	if b.handle == nil {
		b.handle = rigidbody_new(unsafe.Pointer(b), b.mass, b.Shape.shape())
	}

	// update physics transforms
	b.pushState()

	b.world = object.GetInParents[*World](b)
	if b.world != nil {
		b.world.addRigidBody(b)

		// detach object transform from parent
		if !b.Kinematic() {
			b.tfparent = b.Transform().Parent()
			b.Transform().SetParent(nil)
		}
	} else {
		log.Println("Rigidbody", b.Name(), ": no physics world in parents")
	}
}

func (b *RigidBody) OnDisable() {
	b.detach()
	if b.tfparent != nil {
		// re-attach transform to parent
		b.Transform().SetParent(b.tfparent)
	}
}

func (b *RigidBody) detach() {
	if b.world != nil {
		b.world.removeRigidBody(b)
		b.world = nil
	}
}

func (b *RigidBody) destroy() {
	b.detach()
	if b.Shape != nil {
		b.Shape.OnChange().Unsubscribe(b)
		b.Shape = nil
	}
	if b.handle != nil {
		rigidbody_delete(&b.handle)
	}
}

func (b *RigidBody) Kinematic() bool {
	return b.mass <= 0
}
