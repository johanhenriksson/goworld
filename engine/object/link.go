package object

import "github.com/johanhenriksson/goworld/math/vec3"

// Link behaviours to objects
type Link struct {
	parent  *T
	name    string
	enabled bool
}

func NewLink(name string) *Link {
	return &Link{
		name:    name,
		enabled: true,
	}
}

func (b *Link) String() string { return b.name }

// Parent refers to the parent object
func (b *Link) Parent() *T { return b.parent }

// SetParent sets the parent object
func (b *Link) SetParent(o *T) { b.parent = o }

// Active returns true if the object is active
func (b *Link) Active() bool { return b.enabled }

// SetActive sets the active state of the object
func (b *Link) SetActive(active bool) { b.enabled = active }

// Update the object
func (b *Link) Update(dt float32) {}

// Collect query results on an object
func (b *Link) Collect(q *Query) {}

// Forward returns the objects forward vector in world space
func (o *Link) Forward() vec3.T { return o.parent.Forward() }

// Right returns the objects right vector in world space
func (o *Link) Right() vec3.T { return o.parent.Right() }

// Up returns the objects up vector in world space
func (o *Link) Up() vec3.T { return o.parent.Up() }

// Position returns the objects position relative to its parent
func (o *Link) Position() vec3.T { return o.parent.Position() }

// Rotation returns the objects rotation relative to its parent
func (o *Link) Rotation() vec3.T { return o.parent.Rotation() }

// Scale returns the objects scale relative to its parent
func (o *Link) Scale() vec3.T { return o.parent.Scale() }

func (o *Link) SetPosition(p vec3.T) {
	o.parent.SetPosition(p)
}

func (o *Link) SetRotation(r vec3.T) {
	o.parent.SetRotation(r)
}

func (o *Link) SetScale(s vec3.T) {
	o.parent.SetScale(s)
}