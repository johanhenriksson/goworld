package gui_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func TestReconcileProps(t *testing.T) {
	a1 := rect.New("A", &rect.Props{})
	a2 := rect.New("A", &rect.Props{})

	if !Reconcile(a1, a2) {
		t.Error("expected reconciliation to succeed")
	}
	if a1.Props() == a2.Props() {
		t.Error("props should not have been replaced")
	}

	a3 := rect.New("A", &rect.Props{
		Color: color.Red,
	})
	if !Reconcile(a1, a3) {
		t.Error("expected reconciliation to succeed")
	}
	if a1.Props() != a3.Props() {
		t.Error("props should have been replaced")
	}
}

func TestReconcileChildren(t *testing.T) {
	b := rect.New("B1", &rect.Props{})
	a1 := rect.New("A", &rect.Props{}, b)
	a2 := rect.New("A", &rect.Props{}, rect.New("B2", &rect.Props{}), rect.New("B1", &rect.Props{}))
	if !Reconcile(a1, a2) {
		t.Error("expected reconciliation to succeed")
	}

	ch := a1.Children()
	if ch[1] != b {
		t.Error("expected B1 to be reused at index 1")
	}
}
