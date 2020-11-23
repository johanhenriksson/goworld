package object

import "github.com/johanhenriksson/goworld/math/vec3"

type builder struct {
	name     string
	position vec3.T
	rotation vec3.T
	scale    vec3.T
	children []Component
}

func Builder(name string) *builder {
	return &builder{
		name:     name,
		position: vec3.Zero,
		rotation: vec3.Zero,
		scale:    vec3.One,
	}
}

func (b *builder) Attach(c Component) *builder {
	b.children = append(b.children, c)
	return b
}

func (b *builder) Position(p vec3.T) *builder {
	b.position = p
	return b
}

func (b *builder) Rotation(r vec3.T) *builder {
	b.rotation = r
	return b
}

func (b *builder) Scale(s vec3.T) *builder {
	b.scale = s
	return b
}

func (b *builder) Create() *T {
	obj := New(b.name, b.children...)
	obj.SetPosition(b.position)
	obj.SetRotation(b.rotation)
	obj.SetScale(b.scale)
	return obj
}

func testbuild() {
	obj := Builder("NewObject").
		Position(vec3.New(1, 2, 3)).
		Create()
	if obj.enabled {
	}
}
