package object

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"testing"
)

func TestObjectTransforms(t *testing.T) {
	a := New("A")
	b := New("B")

	a.SetPosition(vec3.New(1, 0, 0))
	b.SetPosition(vec3.New(1, 0, 0))

	a.Attach(b)
	// a.Update(0)

	v := b.TransformPoint(vec3.Zero)
	if v.X != 2 {
		t.Errorf("child transform is wrong, was %f", v.X)
	}
}
