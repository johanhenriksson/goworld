package physics

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func init() {
	object.Register[*Box](object.TypeInfo{
		Name:        "Box Collider",
		Path:        []string{"Physics"},
		Deserialize: DeserializeBox,
		Create: func() (object.Component, error) {
			return NewBox(vec3.One), nil
		},
	})
}

type Box struct {
	kind ShapeType
	*Collider

	Extents object.Property[vec3.T]
}

var _ = checkShape(NewBox(vec3.Zero))

func NewBox(size vec3.T) *Box {
	box := object.NewComponent(&Box{
		kind:    BoxShape,
		Extents: object.NewProperty(size),
	})
	box.Collider = newCollider(box, true)

	// resize shape when extents are modified
	box.Extents.OnChange.Subscribe(func(t vec3.T) {
		box.refresh()
	})

	return box
}

func (b *Box) Name() string {
	return "BoxShape"
}

func (b *Box) String() string {
	return fmt.Sprintf("Box[Size=%s]", b.Extents.Get())
}

func (b *Box) colliderCreate() shapeHandle {
	return shape_new_box(unsafe.Pointer(b), b.Extents.Get().Scaled(0.5))
}

func (b *Box) colliderRefresh() {}
func (b *Box) colliderDestroy() {}

type boxState struct {
	Extents vec3.T
}

func (s *Box) Serialize(enc object.Encoder) error {
	return enc.Encode(boxState{
		Extents: s.Extents.Get(),
	})
}

func DeserializeBox(dec object.Decoder) (object.Component, error) {
	var state boxState
	if err := dec.Decode(&state); err != nil {
		return nil, err
	}
	return NewBox(state.Extents), nil
}
