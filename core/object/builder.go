package object

import (
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Builder API for game objects
type builder[K G] struct {
	object K

	position vec3.T
	rotation quat.T
	scale    vec3.T
	active   bool

	parent   G
	children []Component
}

// Builder instantiates a new group builder.
func Builder[K G](object K) *builder[K] {
	return &builder[K]{
		object:   object,
		position: vec3.Zero,
		rotation: quat.Ident(),
		scale:    vec3.One,
		active:   true,
	}
}

// Attach a child component
func (b *builder[K]) Attach(child Component) *builder[K] {
	b.children = append(b.children, child)
	return b
}

// Set the parent of the object
func (b *builder[K]) Parent(parent G) *builder[K] {
	b.parent = parent
	return b
}

// Position sets the intial position of the object.
func (b *builder[K]) Position(p vec3.T) *builder[K] {
	b.position = p
	return b
}

// Rotation sets the intial rotation of the object.
func (b *builder[K]) Rotation(r quat.T) *builder[K] {
	b.rotation = r
	return b
}

// Scale sets the intial scale of the object.
func (b *builder[K]) Scale(s vec3.T) *builder[K] {
	b.scale = s
	return b
}

// Active sets the objects active flag.
func (b *builder[K]) Active(active bool) *builder[K] {
	b.active = active
	return b
}

// Create instantiates a new object with the current builder settings.
func (b *builder[K]) Create() K {
	obj := b.object
	obj.Transform().SetPosition(b.position)
	obj.Transform().SetRotation(b.rotation)
	obj.Transform().SetScale(b.scale)
	obj.SetActive(b.active)
	if b.parent != nil {
		Attach(b.parent, obj)
	}
	for _, child := range b.children {
		Attach(obj, child)
	}
	return obj
}
