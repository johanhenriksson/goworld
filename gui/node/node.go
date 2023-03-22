package node

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/widget"
)

type T interface {
	Key() string
	Type() reflect.Type

	Update(any)
	Props() any
	Hooks() *hooks.State

	Children() []T
	SetChildren([]T)
	Prepend(T)
	Append(T)

	Expand(*hooks.State)
	Hydrate(string) widget.T
	Destroy()
}

type node[P any] struct {
	key      string
	props    P
	kind     reflect.Type
	render   func(P) T
	hydrate  func(string, P) widget.T
	widget   widget.T
	children []T
	hooks    hooks.State
}

func Builtin[P any](key string, props P, children []T, hydrate func(string, P) widget.T) T {
	kind := reflect.TypeOf(props)
	n := Alloc[P](globalPool, kind)
	n.key = key
	n.props = props
	n.hydrate = hydrate
	n.children = children
	return n
}

func Component[P any](key string, props P, render func(P) T) T {
	kind := reflect.TypeOf(props)
	n := Alloc[P](globalPool, kind)
	n.key = key
	n.props = props
	n.render = render
	return n
}

func (n *node[P]) Key() string {
	return n.key
}

func (n *node[P]) Type() reflect.Type {
	return n.kind
}

func (n *node[P]) Props() any {
	return n.props
}

func (n *node[P]) Children() []T {
	return n.children
}

func (n *node[P]) SetChildren(children []T) {
	n.children = children
}

func (n *node[P]) Append(child T) {
	n.SetChildren(append(n.children, child))
}

func (n *node[P]) Prepend(child T) {
	n.SetChildren(append([]T{child}, n.children...))
}

func (n *node[P]) hydrated() bool {
	return n.widget != nil
}

func (n *node[P]) Update(props any) {
	n.props = props.(P)

	if n.render == nil {
		// the node is a built-in element, simply update its props
		if n.hydrated() {
			n.widget.Update(props)
		}
	} else {
		// this node is a potentially stateful component - what do we do?
		// the update might have caused changes to the entire subtree
	}
}

func (n *node[P]) Destroy() {
	for _, child := range n.children {
		child.Destroy()
	}
	if n.hydrated() {
		n.widget.Destroy()
		n.widget = nil
	}
	n.hooks = hooks.State{}
	n.render = nil
	n.hydrate = nil
	n.children = nil
	Free(globalPool, n)
}

func (n *node[P]) Hooks() *hooks.State {
	return &n.hooks
}

// Expand component & child nodes using its hook state.
// If the node is a component, its render function will be called
// to create any dynamic child nodes. This does not cause hydration.
func (n *node[P]) Expand(hook *hooks.State) {
	if n.render == nil {
		return
	}
	if hook == nil {
		hook = &n.hooks
	}

	hooks.Enable(hook)
	defer hooks.Disable()

	n.children = []T{n.render(n.props)}
}

// Hydrates the widgets represented by the node and all of its children.
func (n *node[P]) Hydrate(parentKey string) widget.T {
	// components should never be hydrated directly.
	// we expect to have the root element of the component as a single child at
	// this point, which can be hydrated and returned as this nodes widget,
	// effectively collapsing their nodes in the widget tree.
	if n.render != nil {
		if len(n.children) != 1 {
			panic("expected component to have a single child. did you forget to reconcile?")
		}
		component := n.children[0]
		n.widget = component.Hydrate(parentKey)
		return n.widget
	}

	// the node is a built-in element, hydrate it if it does not exist
	// note: this is the only place where hydration performs any actual work
	if n.widget == nil {
		key := joinKeys(n.key, parentKey)
		n.widget = n.hydrate(key, n.props)
	}

	// the children array might have changed, so we iterate it and hydrate everything.
	// calling Hydrate on a node that is already fully hydrated is basically a no-op, so its fine

	// the logic for optimizing this process currently exists within
	// Rect, since its the only element thay may have child elements.
	// perhaps it would make sense to extract it, and perhaps place it here?
	children := make([]widget.T, len(n.children))
	for i, child := range n.children {
		// hydrate child. if its already hydrated, this is basically a no-op
		children[i] = child.Hydrate(n.widget.Key())
	}
	n.widget.SetChildren(children)

	return n.widget
}

func (n *node[P]) String() string {
	return fmt.Sprintf("Node[%s] %s", n.kind, n.key)
}

// efficient method to concatenate key strings
func joinKeys(parent, child string) string {
	p := len(parent)
	buffer := make([]byte, p+len(child)+1)
	copy(buffer, []byte(parent))
	buffer[p] = '/'
	copy(buffer[p+1:], child)
	return string(buffer)
}
