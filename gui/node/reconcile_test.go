package node_test

import (
	. "github.com/johanhenriksson/goworld/gui/node"
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
	hydrate := func(key string, p *props) widget.T {
		hydrateCalls++
		return element
	}

	a := Builtin("a", &props{"hello"}, nil, hydrate)
	a.Hydrate()
	b := Builtin("a", &props{"world"}, nil, hydrate)

	r := Reconcile(a, b)
	r.Hydrate()

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
	hydrate := func(key string, p *props) widget.T { return nil }

	a := Builtin("a", &props{}, []T{
		Builtin("b", &props{"child"}, nil, hydrate),
	}, hydrate)
	b := Builtin("a", &props{}, []T{
		Builtin("b", &props{"updated"}, nil, hydrate),
	}, hydrate)

	r := Reconcile(a, b)
	child := r.Children()[0]
	cp := child.Props().(*props)
	if cp.Text != "updated" {
		t.Error("expected child props to be updated")
	}
}
