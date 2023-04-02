package node_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
)

type props struct {
	Text string
}

func TestNode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Node Suite")
}

var _ = Describe("reconcile nodes", func() {
	It("reconciles elements of the same type", func() {
		element := widget.Dummy("element")
		hydrateCalls := 0
		hydrate := func(key string, p *props) widget.T {
			hydrateCalls++
			return element
		}

		a := node.Builtin("a", &props{"hello"}, nil, hydrate)
		a.Hydrate("test")
		b := node.Builtin("a", &props{"world"}, nil, hydrate)

		r := node.Reconcile(a, b)
		r.Hydrate("test")

		Expect(a.Hooks()).To(Equal(r.Hooks()))

		// we expect As props to be updated
		rp := a.Props().(*props)
		Expect(rp.Text).To(Equal("world"))

		// only a single hydration should have occured
		Expect(hydrateCalls).To(Equal(1))
	})

	It("reconciles children", func() {
		hydrate := func(key string, p *props) widget.T { return nil }

		a := node.Builtin("a", &props{}, []node.T{
			node.Builtin("b", &props{"child"}, nil, hydrate),
		}, hydrate)
		b := node.Builtin("a", &props{}, []node.T{
			node.Builtin("b", &props{"updated"}, nil, hydrate),
		}, hydrate)

		r := node.Reconcile(a, b)
		child := r.Children()[0]
		cp := child.Props().(*props)
		Expect(cp.Text).To(Equal("updated"))
	})

	Context("handles nil elements", func() {
		It("nil becomes element", func() {
			hydrate := func(key string, p *props) widget.T { return nil }

			a := node.Builtin("a", &props{}, []node.T{
				nil,
			}, hydrate)
			b := node.Builtin("a", &props{}, []node.T{
				node.Builtin("b", &props{"new"}, nil, hydrate),
			}, hydrate)

			r := node.Reconcile(a, b)
			child := r.Children()[0]
			cp := child.Props().(*props)
			Expect(cp.Text).To(Equal("new"))
		})

		It("element becomes nil", func() {
			hydrate := func(key string, p *props) widget.T { return nil }

			a := node.Builtin("a", &props{}, []node.T{
				node.Builtin("b", &props{"new"}, nil, hydrate),
			}, hydrate)
			b := node.Builtin("a", &props{}, []node.T{
				nil,
			}, hydrate)

			r := node.Reconcile(a, b)
			Expect(r.Children()).To(HaveLen(0))
		})
	})
})
