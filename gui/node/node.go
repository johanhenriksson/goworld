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

type node[K widget.T, P any] struct {
	key      string
	props    P
	kind     reflect.Type
	render   func(P) T
	hydrate  func(widget.T, P) K
	widget   widget.T
	children []T
	hooks    hooks.State
}

func Builtin[K widget.T, P any](key string, props P, children []T, hydrate func(widget.T, P) K) T {
	return &node[K, P]{
		key:      key,
		props:    props,
		kind:     reflect.TypeOf(props),
		hydrate:  hydrate,
		children: children,
	}
}

func Component[P any](key string, props P, render func(P) T) T {
	return &node[widget.T, P]{
		key:    key,
		props:  props,
		kind:   reflect.TypeOf(props),
		render: render,
	}
}

func (n *node[K, P]) Key() string {
	return n.key
}

func (n *node[K, P]) Type() reflect.Type {
	return n.kind
}

func (n *node[K, P]) Props() any {
	return n.props
}

func (n *node[K, P]) Children() []T {
	return n.children
}

func (n *node[K, P]) SetChildren(children []T) {
	n.children = children
}

func (n *node[K, P]) Append(child T) {
	n.SetChildren(append(n.children, child))
}

func (n *node[K, P]) Prepend(child T) {
	n.SetChildren(append([]T{child}, n.children...))
}

func (n *node[K, P]) hydrated() bool {
	return n.widget != nil
}

func (n *node[K, P]) Update(props any) {
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

func (n *node[K, P]) Destroy() {
	for _, child := range n.children {
		child.Destroy()
	}
	if !n.hydrated() {
		return
	}
	n.widget.Destroy()
	n.widget = nil
}

func (n *node[K, P]) Hooks() *hooks.State {
	return &n.hooks
}

// Expand component & child nodes using its hook state.
// If the node is a component, its render function will be called
// to create any dynamic child nodes. This does not cause hydration.
func (n *node[K, P]) Expand(hook *hooks.State) {
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
func (n *node[K, P]) Hydrate(parentKey string) widget.T {
	// check if we are a component or a built-in element
	if n.render != nil {
		// components should never be hydrated directly.
		// we expect to have the root element of the component as a single child at
		// this point, which can be hydrated and returned as this nodes widget,
		// effectively collapsing their nodes in the widget tree.
		if len(n.children) != 1 {
			panic("expected component to have a single child. did you forget to reconcile?")
		}
		component := n.children[0]
		n.widget = component.Hydrate(parentKey)
	} else {
		// this node is a built-in element, hydrate it if it does not exist
		key := joinKeys(n.key, parentKey)
		if n.widget == nil {
			n.widget = n.hydrate(widget.New(key), n.props)
		}

		// rehydrate children if required
		// the logic for optimizing this process currently exists within
		// Rect, since its the only element thay may have child elements.
		// perhaps it would make sense to extract it, and perhaps place it here?
		children := make([]widget.T, 0, len(n.children))
		for _, child := range n.children {
			hydrated := child.Hydrate(key)
			children = append(children, hydrated)
		}
		n.widget.SetChildren(children)
	}

	return n.widget
}

func (n *node[K, P]) String() string {
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
