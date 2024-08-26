package physics

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
)

func init() {
	object.Register[*Sphere](object.Type{
		Name: "Sphere Collider",
		Path: []string{"Physics"},
		Create: func(ctx object.Pool) (object.Component, error) {
			return NewSphere(ctx, 1), nil
		},
	})
}

type Sphere struct {
	kind ShapeType
	*Collider

	Radius object.Property[float32]
}

var _ = checkShape(NewSphere(object.GlobalPool, 1))

func NewSphere(pool object.Pool, radius float32) *Sphere {
	sphere := &Sphere{
		kind:   SphereShape,
		Radius: object.NewProperty[float32](radius),
	}
	sphere.Collider = newCollider(pool, sphere, true)

	// resize shape when radius is modified
	sphere.Radius.OnChange.Subscribe(func(t float32) {
		sphere.refresh()
	})

	return sphere
}

func (s *Sphere) colliderCreate() shapeHandle {
	return shape_new_sphere(unsafe.Pointer(s), s.Radius.Get())
}

func (s *Sphere) colliderIsCompound() bool { return false }

func (s *Sphere) colliderRefresh() {}
func (s *Sphere) colliderDestroy() {}
