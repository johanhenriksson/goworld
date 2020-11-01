package object

import "github.com/johanhenriksson/goworld/engine/transform"

// T is the basic building block of the scene graph
type T struct {
	*transform.T
	name       string
	enabled    bool
	parent     *T
	components []Component
}

// NewObject instantiates a new game object
func New(name string, components ...Component) *T {
	obj := &T{
		T:       transform.Identity(),
		enabled: true,
		name:    name,
		parent:  nil,
	}
	obj.Attach(components...)
	return obj
}

func (o *T) String() string { return o.name }

// Parent returns a pointer to the parent object (or nil)
func (o *T) Parent() *T { return o.parent }

// SetParent sets the parent object pointer
func (o *T) SetParent(p *T) {
	o.parent = p
	if p != nil {
		o.T.Update(&p.T.World)
	}
}

// SetActive sets the objects active state
func (o *T) SetActive(active bool) { o.enabled = active }

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
	// update transform
	if o.parent != nil {
		o.T.Update(&o.parent.T.World)
	} else {
		o.T.Update(nil)
	}

	// update components
	for _, component := range o.components {
		if !component.Active() {
			continue
		}
		component.Update(dt)
	}
}
