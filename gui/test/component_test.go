package test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGUI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gui")
}

type ComponentProps struct{}

func Component(key string, props ComponentProps) node.T {
	return node.Component(key, props, func(ComponentProps) node.T {
		return rect.New(key, rect.Props{
			Children: []node.T{
				rect.New("a", rect.Props{}),
				rect.New("b", rect.Props{}),
			},
		})
	})
}

type NestedProps struct {
	StartClosed bool
}

func NestedComponent(key string, props NestedProps) node.T {
	return node.Component(key, props, func(props NestedProps) node.T {
		closed, _ := hooks.UseState(props.StartClosed)

		var children []node.T
		if !closed {
			children = []node.T{
				Component("c1", ComponentProps{}),
				Component("c2", ComponentProps{}),
			}
		}
		return rect.New(key, rect.Props{
			Children: children,
		})
	})
}

var _ = Describe("", func() {
	Context("components", func() {
		It("hydrates components", func() {
			root := Component("root", ComponentProps{})
			node.Reconcile(nil, root)
			w := root.Hydrate("test")
			Expect(w.Children()).To(HaveLen(2))
		})

		It("hydrates nested components", func() {
			root := NestedComponent("root", NestedProps{})
			node.Reconcile(nil, root)
			w := root.Hydrate("test")
			Expect(w.Children()).To(HaveLen(2))
		})

		It("handles hook state changes", func() {
			root := NestedComponent("root", NestedProps{})
			root.Hooks().Write(0, true) // open component
			tree := node.Reconcile(nil, root)
			w := tree.Hydrate("test")
			Expect(w.Children()).To(HaveLen(0))

			tree.Hooks().Write(0, false) // close component again
			tree = node.Reconcile(tree, root)
			w = tree.Hydrate("test")
			Expect(w.Children()).To(HaveLen(2))
		})
	})
})
