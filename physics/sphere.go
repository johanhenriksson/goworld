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

func NewSphere() *Sphere {
	sphere := &Sphere{
		kind:   SphereShape,
		Radius: object.NewProperty[float32](1),
	}
	sphere.Collider = newCollider(sphere, true)
	return sphere
}

func (s *Sphere) colliderCreate() shapeHandle {
	return shape_new_sphere(unsafe.Pointer(s), s.Radius.Get())
}

func (s *Sphere) colliderRefresh() {}
func (s *Sphere) colliderDestroy() {}
