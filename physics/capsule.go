package physics

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
)

func init() {
	object.Register[*Capsule](object.TypeInfo{
		Name:        "Capsule Collider",
		Path:        []string{"Physics"},
		Deserialize: DeserializeCapsule,
		Create: func(ctx object.Pool) (object.Component, error) {
			return NewCapsule(ctx, 1.8, 0.3), nil
		},
	})
}

type Capsule struct {
	kind ShapeType
	*Collider

	Radius object.Property[float32]
	Height object.Property[float32]
}

var _ = checkShape(NewCapsule(object.GlobalPool, 1, 1))

func NewCapsule(pool object.Pool, height, radius float32) *Capsule {
	capsule := &Capsule{
		kind:   CapsuleShape,
		Radius: object.NewProperty(radius),
		Height: object.NewProperty(height),
	}
	capsule.Collider = newCollider(pool, capsule, true)

	capsule.Radius.OnChange.Subscribe(func(radius float32) {
		capsule.refresh()
	})
	capsule.Height.OnChange.Subscribe(func(height float32) {
		capsule.refresh()
	})

	return object.NewComponent(pool, capsule)
}

func (c *Capsule) Name() string {
	return "CapsuleShape"
}

func (c *Capsule) String() string {
	return fmt.Sprintf("Capsule[Height=%.2f,Radius=%.2f]", c.Height.Get(), c.Radius.Get())
}

func (c *Capsule) colliderCreate() shapeHandle {
	return shape_new_capsule(unsafe.Pointer(c), c.Radius.Get(), c.Height.Get())
}

func (c *Capsule) colliderIsCompound() bool { return false }

func (c *Capsule) colliderRefresh() {}
func (c *Capsule) colliderDestroy() {}

type capsuleState struct {
	Radius float32
	Height float32
}

func (s *Capsule) Serialize(enc object.Encoder) error {
	return enc.Encode(capsuleState{
		Height: s.Height.Get(),
		Radius: s.Radius.Get(),
	})
}

func DeserializeCapsule(ctx object.Pool, dec object.Decoder) (object.Component, error) {
	var state capsuleState
	if err := dec.Decode(&state); err != nil {
		return nil, err
	}
	return NewCapsule(ctx, state.Height, state.Radius), nil
}
