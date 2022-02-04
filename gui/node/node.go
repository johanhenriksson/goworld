package node

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/gui/component"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/widget"
)

type T interface {
	Key() string
	Type() reflect.Type
	Props() any
	Children() []T
	SetChildren([]T)
	Hooks() *hooks.State

	Update(any)
	Render()
	Destroy()
	Hydrate() widget.T
	Hydrated() bool
}

type Args[P any] struct {
	Key      string
	Props    P
	Children []T
}

type node[K widget.T, P any] struct {
	key      string
	props    P
	kind     reflect.Type
	render   func(P) T
	hydrate  func(string, P) K
	widget   widget.T
	children []T
	hooks    hooks.State
}

func Builtin[K widget.T, P any](key string, props P, children []T, hydrate func(string, P) K) T {
	var empty K
	return &node[K, P]{
		key:      key,
		props:    props,
		kind:     reflect.TypeOf(empty),
		hydrate:  hydrate,
		children: children,
		render:   nil,
	}
}

func Component[P any](key string, props P, children []T, render func(P) T) T {
	return &node[component.T, P]{
		key:      key,
		props:    props,
		kind:     reflect.TypeOf(props),
		children: children,
		hydrate: func(key string, props P) component.T {
			return component.New(key, props)
		},
		render: render,
	}
}

func (n node[K, P]) Key() string {
	return n.key
}

func (n node[K, P]) Type() reflect.Type {
	return n.kind
}

func (n node[K, P]) Props() any {
	return n.props
}

func (n node[K, P]) Children() []T {
	return n.children
}

func (n *node[K, P]) SetChildren(children []T) {
	n.children = children
}

func (n node[K, P]) Hydrated() bool {
	return n.widget != nil
}

func (n *node[K, P]) Render() {
	if n.render == nil {
		return
	}

	hooks.Enable(&n.hooks)
	defer hooks.Disable()

	n.children = []T{
		n.render(n.props),
	}
}

func (n *node[K, P]) Update(props any) {
	n.props = props.(P)

	n.Render()
	if n.Hydrated() {
		// we are a basic element, update my props
		n.widget.Update(n.props)
	}
}

func (n *node[K, P]) Destroy() {
	if !n.Hydrated() {
		return
	}
	n.widget.Destroy()
	n.widget = nil
}

func (n *node[K, P]) Hooks() *hooks.State {
	return &n.hooks
}

func (n *node[K, P]) Hydrate() widget.T {
	if n.widget == nil {
		n.widget = n.hydrate(n.key, n.props)
	}

	// render children
	children := make([]widget.T, len(n.children))
	for i, child := range n.children {
		children[i] = child.Hydrate()
	}
	n.widget.SetChildren(children)

	return n.widget
}

func (n *node[K, P]) String() string {
	return fmt.Sprintf("Node[%s] %s", n.kind, n.key)
}
