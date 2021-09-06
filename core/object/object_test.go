package object

import (
	"testing"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type TransformTest struct {
	Apos vec3.T
	Arot vec3.T
	Bpos vec3.T
	Brot vec3.T
	Epos vec3.T
	Erot vec3.T
}

func TestObjectTransforms(t *testing.T) {
	a := New("A")
	b := New("B")
	a.Adopt(b)

	cases := []TransformTest{
		{
			Apos: vec3.New(1, 0, 0),
			Arot: vec3.New(0, 0, 0),
			Bpos: vec3.New(1, 0, 0),
			Brot: vec3.New(0, 0, 0),
			Epos: vec3.New(2, 0, 0),
			Erot: vec3.New(0, 0, 1),
		},
	}

	for idx, c := range cases {
		a.Transform().SetPosition(c.Apos)
		a.Transform().SetRotation(c.Arot)
		b.Transform().SetPosition(c.Bpos)
		b.Transform().SetRotation(c.Brot)

		pos := b.Transform().WorldPosition()
		if !pos.ApproxEqual(c.Epos) {
			t.Errorf("%d position is wrong, expected transformed point %+v but was %+v", idx, c.Epos, pos)
		}
		rot := b.Transform().Forward()
		if !rot.ApproxEqual(c.Erot) {
			t.Errorf("%d forward is wrong, expected transformed direction %+v but was %+v", idx, c.Erot, rot)
		}
	}
}
