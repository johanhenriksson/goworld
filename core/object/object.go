package object

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/transform"
)

// object is the basic building block of the scene graph
type object struct {
	transform  transform.T
	name       string
	enabled    bool
	parent     T
	components []Component
	children   []T
}

// New instantiates a new game object
func New(name string, components ...Component) T {
	obj := &object{
		transform: transform.Identity(),
		enabled:   true,
		name:      name,
		parent:    nil,
	}
	obj.Attach(components...)
	return obj
}

func (o *object) Name() string { return o.name }

// Parent returns a pointer to the parent object (or nil)
func (o *object) Parent() T { return o.parent }

// SetParent sets the parent object pointer
func (o *object) SetParent(p T) {
	o.parent = p
	o.transform.Recalculate(p.Transform())
}

// SetActive sets the objects active state
func (o *object) SetActive(active bool) {
	o.enabled = active
}

// Active indicates whether the object is currently enabled
func (o *object) Active() bool { return o.enabled }

// Collect performs a query against this objects child components
func (o *object) Collect(query *Query) {
	for _, component := range o.components {
		if !component.Active() {
			continue
		}
		if query.Match(component) {
			query.Append(component)
		}
	}
	for _, child := range o.children {
		if !child.Active() {
			continue
		}
		child.Collect(query)
	}
}

// Attach a component to this object
func (o *object) Attach(components ...Component) {
	for _, component := range components {
		// attach it
		o.components = append(o.components, component)
		component.SetObject(o)
	}
}

// Update this object and its child components
func (o *object) Update(dt float32) {
	// update components
	for _, component := range o.components {
		if !component.Active() {
			continue
		}
		component.Update(dt)
	}
	for _, child := range o.children {
		if !child.Active() {
			continue
		}
		child.Update(dt)
	}
}

func (o *object) Transform() transform.T {
	// this might need to be optimized
	// we are currently recalculating the entire transformation hierarchy every time the transform is accessed
	// pros - always up to date
	// cons - its dumb

	var pt transform.T = nil
	if o.Parent() != nil {
		pt = o.parent.Transform()
	}
	o.transform.Recalculate(pt)

	return o.transform
}

func (o *object) Adopt(children ...T) {
	for _, child := range children {
		// attach it
		o.children = append(o.children, child)
		child.SetParent(o)
	}
}

func (o *object) Children() []T {
	return o.children
}

func (o *object) KeyEvent(e keys.Event) {
	for _, component := range o.components {
		// propagate the event to components if they implement keys.Handler
		if handler, ok := component.(keys.Handler); ok {
			handler.KeyEvent(e)
			if e.Handled() {
				return
			}
		}
	}

	for _, child := range o.children {
		child.KeyEvent(e)
		if e.Handled() {
			return
		}
	}
}

func (o *object) MouseEvent(e mouse.Event) {
	for _, component := range o.components {
		// propagate the event to components if they implement mouse.Handler
		if handler, ok := component.(mouse.Handler); ok {
			handler.MouseEvent(e)
			if e.Handled() {
				return
			}
		}
	}

	for _, child := range o.children {
		child.MouseEvent(e)
		if e.Handled() {
			return
		}
	}
}