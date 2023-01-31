package object

import "github.com/johanhenriksson/goworld/math/vec3"

// Builder API for game objects
type Builder[K T] struct {
	name     string
	position vec3.T
	rotation vec3.T
	scale    vec3.T
	active   bool

	parent   T
	children []T
}

// Build instantiates a new object builder.
func Build[K T](name string) *Builder[K] {
	return &Builder[K]{
		name:     name,
		position: vec3.Zero,
		rotation: vec3.Zero,
		scale:    vec3.One,
		active:   true,
	}
}

func (b *Builder[K]) Attach(child T) *Builder[K] {
	b.children = append(b.children, child)
	return b
}

func (b *Builder[K]) Parent(parent T) *Builder[K] {
	b.parent = parent
	return b
}

// Position sets the intial position of the object.
func (b *Builder[K]) Position(p vec3.T) *Builder[K] {
	b.position = p
	return b
}

// Rotation sets the intial rotation of the object.
func (b *Builder[K]) Rotation(r vec3.T) *Builder[K] {
	b.rotation = r
	return b
}

// Scale sets the intial scale of the object.
func (b *Builder[K]) Scale(s vec3.T) *Builder[K] {
	b.scale = s
	return b
}

// Active sets the objects active flag.
func (b *Builder[K]) Active(active bool) *Builder[K] {
	b.active = active
	return b
}

func (b *Builder[K]) Name(name string) *Builder[K] {
	b.name = name
	return b
}

// Create instantiates a new object with the current builder settings.
func (b *Builder[K]) Create(obj K) K {
	obj = New(obj)
	obj.setName(b.name)
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
