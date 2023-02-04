package node

import (
	"fmt"
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/component"
)

type T interface {
	Key() string
	Type() reflect.Type
	Props() any
	Children() []T
	SetChildren([]T)
	Hooks() *hooks.State

	Inject(T)

	Update(any)
	Render(*hooks.State)
	Destroy()
	Hydrate(parentKey string) widget.T
	Hydrated() bool
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
	injected []T
}

func Builtin[K widget.T, P any](key string, props P, children []T, hydrate func(widget.T, P) K) T {
	return &node[K, P]{
		key:      key,
		props:    props,
		kind:     reflect.TypeOf(props),
		hydrate:  hydrate,
		children: children,
		render:   nil,
	}
}

func Component[P any](key string, props P, render func(P) T) T {
	return &node[component.T, P]{
		key:   key,
		props: props,
		kind:  reflect.TypeOf(props),
		hydrate: func(w widget.T, props P) component.T {
			return component.Create(w, props)
		},
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

func (n *node[K, P]) Inject(node T) {
	n.injected = append(n.injected, node)
}

func (n *node[K, P]) Hydrated() bool {
	return n.widget != nil
}

func (n *node[K, P]) Render(hook *hooks.State) {
	if n.render == nil {
		return
	}
	if hook == nil {
		hook = &n.hooks
	}

	hooks.Enable(hook)
	defer hooks.Disable()

	n.children = append(n.injected, n.render(n.props))
}

func (n *node[K, P]) Update(props any) {
	n.props = props.(P)

	if n.Hydrated() {
		// we are a basic element, update my props
		n.widget.Update(props)
	}
}

func (n *node[K, P]) Destroy() {
	if !n.Hydrated() {
		return
	}
	log.Println("destroy node", n.key)
	for _, child := range n.injected {
		child.Destroy()
	}
	for _, child := range n.children {
		child.Destroy()
	}
	n.widget.Destroy()
	n.widget = nil
}

func (n *node[K, P]) Hooks() *hooks.State {
	return &n.hooks
}

func (n *node[K, P]) Hydrate(parentKey string) widget.T {
	key := joinKeys(n.key, parentKey)
	if n.widget == nil {
		n.widget = n.hydrate(widget.New(key), n.props)
	}

	// render children
	children := make([]widget.T, 0, len(n.children)+len(n.injected))
	for _, child := range n.injected {
		children = append(children, child.Hydrate(key))
	}
	for _, child := range n.children {
		children = append(children, child.Hydrate(key))
	}
	n.widget.SetChildren(children)

	return n.widget
}

func (n *node[K, P]) String() string {
	return fmt.Sprintf("Node[%s] %s", n.kind, n.key)
}

func joinKeys(parent, child string) string {
	p := len(parent)
	buffer := make([]byte, p+len(child)+1)
	copy(buffer, []byte(parent))
	buffer[p] = '/'
	copy(buffer[p+1:], child)
	return string(buffer)
}
