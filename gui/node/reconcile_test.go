package node_test

import (
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"

	"testing"
)

type props struct {
	Text string
}

type testWidget struct {
	widget.T
	props *props
}

func (w *testWidget) Update(p any) { w.props = p.(*props) }
func (w *testWidget) Props() any   { return w.props }

func TestReconcile(t *testing.T) {
	element := &testWidget{T: widget.New("a")}
	hydrateCalls := 0
	hydrate := func(w widget.T, p *props) widget.T {
		hydrateCalls++
		return element
	}

	a := node.Builtin("a", &props{"hello"}, nil, hydrate)
	a.Hydrate("test")
	b := node.Builtin("a", &props{"world"}, nil, hydrate)

	r := node.Reconcile(a, b)
	r.Hydrate("test")

	if a.Hooks() != r.Hooks() {
		t.Error("expected hook state to be unchanged")
	}
	rp := a.Props().(*props)
	if rp.Text != "world" {
		t.Errorf("expected props to be updated")
	}

	if hydrateCalls != 1 {
		t.Errorf("expected a single hydrate call")
	}

	ep := element.Props().(*props)
	if ep.Text != "world" {
		t.Errorf("expected hydrated widget props to be updated")
	}
}

func TestReconcileChildren(t *testing.T) {
	hydrate := func(w widget.T, p *props) widget.T { return nil }

	a := node.Builtin("a", &props{}, []node.T{
		node.Builtin("b", &props{"child"}, nil, hydrate),
	}, hydrate)
	b := node.Builtin("a", &props{}, []node.T{
		node.Builtin("b", &props{"updated"}, nil, hydrate),
	}, hydrate)

	r := node.Reconcile(a, b)
	child := r.Children()[0]
	cp := child.Props().(*props)
	if cp.Text != "updated" {
		t.Error("expected child props to be updated")
	}
}
