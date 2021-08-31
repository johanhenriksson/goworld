package object

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// T is the basic building block of the scene graph
type T struct {
	name       string
	enabled    bool
	parent     *T
	components []Component

	local    mat4.T
	world    mat4.T
	position vec3.T
	rotation vec3.T
	scale    vec3.T
	forward  vec3.T
	right    vec3.T
	up       vec3.T
}

// New instantiates a new game object
func New(name string, components ...Component) *T {
	obj := &T{
		enabled: true,
		name:    name,
		parent:  nil,
		world:   mat4.Ident(),
		local:   mat4.Ident(),
		scale:   vec3.One,
		forward: vec3.UnitZN,
		right:   vec3.UnitX,
		up:      vec3.UnitY,
	}
	obj.Attach(components...)
	return obj
}

func (o *T) String() string { return o.name }

// Parent returns a pointer to the parent object (or nil)
func (o *T) Parent() *T { return o.parent }

// SetParent sets the parent object pointer
func (o *T) SetParent(p *T) Component {
	o.parent = p
	o.updateTransform()
	return o
}

// SetActive sets the objects active state
func (o *T) SetActive(active bool) Component {
	o.enabled = active
	return o
}

// Active indicates whether the object is currently enabled
func (o *T) Active() bool { return o.enabled }

// Collect performs a query against this objects child components
func (o *T) Collect(query *Query) {
	for _, component := range o.components {
		if !component.Active() {
			continue
		}
		if query.Match(component) {
			query.Append(component)
		}
		component.Collect(query)
	}
}

// Attach a component to this object
func (o *T) Attach(components ...Component) {
	for _, component := range components {
		// find the ancestor component
		// we always attach the whole object tree
		for component.Parent() != nil {
			component = component.Parent()
		}

		// attach it
		o.components = append(o.components, component)
		component.SetParent(o)
	}
}

// Update this object and its child components
func (o *T) Update(dt float32) {
	// update transformation matrix
	// by placing the call here we can avoid the problem of recursively updating
	// child transforms when the parent changes. we also guarantee that the transform
	// is only recalculated once per frame, instead of every call to SetPosition/Rot etc
	o.updateTransform()

	// update components
	for _, component := range o.components {
		if !component.Active() {
			continue
		}
		component.Update(dt)
	}
}

func (o *T) updateTransform() {
	// Update local transform
	o.local = mat4.Transform(o.position, o.rotation, o.scale)

	// Update local -> world matrix
	// if we have a parent, apply its transforms to our local-to-world matrix
	// otherwise, the world transform is equal to the local transform
	if o.parent != nil {
		o.world = o.parent.world.Mul(&o.local)
	} else {
		o.world = o.local
	}

	// Grab axis vectors from transformation matrix
	o.up = o.world.Up()
	o.right = o.world.Right()
	o.forward = o.world.Forward()
}

// TransformPoint transforms a world point into this coordinate system
func (o *T) TransformPoint(point vec3.T) vec3.T {
	return o.world.TransformPoint(point)
}

// TransformDir transforms a world direction vector into this coordinate system
func (o *T) TransformDir(dir vec3.T) vec3.T {
	return o.world.TransformDir(dir)
}

// Forward returns the objects forward vector in world space
func (o *T) Forward() vec3.T { return o.forward }

// Right returns the objects right vector in world space
func (o *T) Right() vec3.T { return o.right }

// Up returns the objects up vector in world space
func (o *T) Up() vec3.T { return o.up }

// Position returns the objects position relative to its parent
func (o *T) Position() vec3.T { return o.position }

// Rotation returns the objects rotation relative to its parent
func (o *T) Rotation() vec3.T { return o.rotation }

// Scale returns the objects scale relative to its parent
func (o *T) Scale() vec3.T { return o.scale }

// SetPosition sets the objects position.
func (o *T) SetPosition(p vec3.T) Component {
	o.position = p
	return o
}

// SetRotation sets the objects rotation.
func (o *T) SetRotation(r vec3.T) Component {
	o.rotation = r
	return o
}

// SetScale sets the objects scale.
func (o *T) SetScale(s vec3.T) Component {
	o.scale = s
	return o
}

// Transform returns the affine transformation matrix of the object.
func (o *T) Transform() mat4.T {
	return o.world
}
