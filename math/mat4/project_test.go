package mat4_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type TransformTest struct {
	Input  vec3.T
	Output vec3.T
}

func AssertTransforms(t *testing.T, transform T, cases []TransformTest) {
	t.Helper()
	for _, c := range cases {
		point := transform.TransformPoint(c.Input)
		if !point.ApproxEqual(c.Output) {
			t.Errorf("expected %v was %v", c.Output, point)
		}
	}
}

func TestOrthographicVK(t *testing.T) {
	proj := OrthographicVK(0, 10, 0, 10, -1, 1)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(5, 5, 0), vec3.New(0, 0, 0.5)},
		{vec3.New(5, 5, 1), vec3.New(0, 0, 0)},
		{vec3.New(5, 5, -1), vec3.New(0, 0, 1)},
		{vec3.New(0, 0, -1), vec3.New(-1, -1, 1)},
	})
}

func TestOrthographicLH(t *testing.T) {
	proj := OrthographicLH(0, 10, 0, 10, -1, 1)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(5, 5, 0), vec3.New(0, 0, 0)},
		{vec3.New(5, 5, 1), vec3.New(0, 0, 1)},
		{vec3.New(5, 5, -1), vec3.New(0, 0, -1)},
		{vec3.New(0, 0, -1), vec3.New(-1, -1, -1)},
		{vec3.New(10, 10, -1), vec3.New(1, 1, -1)},
	})
}

func TestOrthographic(t *testing.T) {
	proj := Orthographic(0, 10, 0, 10, -1, 1)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(5, 5, 0), vec3.New(0, 0, 0)},
		{vec3.New(5, 5, 1), vec3.New(0, 0, -1)},
		{vec3.New(5, 5, -1), vec3.New(0, 0, 1)},
		{vec3.New(0, 0, -1), vec3.New(-1, -1, 1)},
		{vec3.New(10, 10, -1), vec3.New(1, 1, 1)},
	})
}

func TestPerspectiveVK(t *testing.T) {
	proj := PerspectiveVK(45, 1, 1, 100)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(0, 0, 1), vec3.New(0, 0, 0)},
		{vec3.New(0, 0, 100), vec3.New(0, 0, 1)},
	})
}

func TestPerspectiveLH(t *testing.T) {
	proj := PerspectiveLH(45, 1, 1, 100)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(0, 0, 1), vec3.New(0, 0, 1)},
		{vec3.New(0, 0, 100), vec3.New(0, 0, -1)},
	})
}

func TestPerspective(t *testing.T) {
	proj := Perspective(45, 1, 1, 100)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(0, 0, 1.00), vec3.New(0, 0, -1)},
		{vec3.New(0, 0, 100), vec3.New(0, 0, 1)},
	})
}
