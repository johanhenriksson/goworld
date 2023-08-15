package physics

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
)

func init() {
	object.Register[*Sphere](DeserializeSphere)
}

type Sphere struct {
	kind ShapeType
	*Collider

	Radius object.Property[float32]
}

var _ = checkShape(NewSphere(1))

func NewSphere(radius float32) *Sphere {
	sphere := &Sphere{
		kind:   SphereShape,
		Radius: object.NewProperty[float32](radius),
	}
	sphere.Collider = newCollider(sphere, true)

	// resize shape when radius is modified
	sphere.Radius.OnChange.Subscribe(func(t float32) {
		sphere.refresh()
	})

	return sphere
}

func (s *Sphere) colliderCreate() shapeHandle {
	return shape_new_sphere(unsafe.Pointer(s), s.Radius.Get())
}

func (s *Sphere) colliderRefresh() {}
func (s *Sphere) colliderDestroy() {}

type sphereState struct {
	Radius float32
}

func (s *Sphere) Serialize(enc object.Encoder) error {
	return enc.Encode(sphereState{
		Radius: s.Radius.Get(),
	})
}

func DeserializeSphere(dec object.Decoder) (object.Component, error) {
	var state sphereState
	if err := dec.Decode(&state); err != nil {
		return nil, err
	}
	return NewSphere(state.Radius), nil
}
