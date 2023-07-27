package physics

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
)

type Sphere struct {
	kind ShapeType
	*Collider

	Radius object.Property[float32]
}

func init() {
	checkShape(NewSphere(1))
}

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
