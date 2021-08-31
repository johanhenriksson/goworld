package object

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"testing"
)

func TestObjectTransforms(t *testing.T) {
	a := New("A")
	b := New("B")

	a.Transform().SetPosition(vec3.New(1, 0, 0))

	a.Adopt(b)

	b.Transform().SetPosition(vec3.New(1, 0, 0))
	a.Transform().SetRotation(vec3.New(0, 90, 0))

	v := b.Transform().TransformPoint(vec3.Zero)
	e := vec3.New(1, 0, -1)
	if !v.ApproxEqual(e) {
		t.Errorf("child transform is wrong, expected transformed point %+v but was %+v", e, v)
	}
}
