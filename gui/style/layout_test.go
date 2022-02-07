package style_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
)

func assertSize(t *testing.T, wgt widget.T, w, h float32) {
	t.Helper()
	expected := vec2.New(w, h)
	if wgt.Size() != expected {
		t.Errorf("expected %s to have size %v, was %v", wgt.Key(), expected, wgt.Size())
	}
}

func assertPosition(t *testing.T, wgt widget.T, x, y float32) {
	t.Helper()
	expected := vec2.New(x, y)
	if wgt.Position() != expected {
		t.Errorf("expected %s to have position %v, was %v", wgt.Key(), expected, wgt.Position())
	}
}
