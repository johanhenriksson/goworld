package mat4_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMat4(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "math/mat4")
}

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

func TestOrthographicRZ(t *testing.T) {
	proj := OrthographicRZ(0, 10, 0, 10, -1, 1)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(5, 5, 0), vec3.New(0, 0, 0.5)},
		{vec3.New(5, 5, 1), vec3.New(0, 0, 0)},
		{vec3.New(5, 5, -1), vec3.New(0, 0, 1)},
		{vec3.New(0, 0, -1), vec3.New(-1, -1, 1)},
	})
}

func TestPerspectiveVK(t *testing.T) {
	proj := Perspective(45, 1, 1, 100)
	AssertTransforms(t, proj, []TransformTest{
		{vec3.New(0, 0, 1), vec3.New(0, 0, 0)},
		{vec3.New(0, 0, 100), vec3.New(0, 0, 1)},
	})
}

var _ = Describe("LookAt (LH)", func() {
	It("correctly projects", func() {
		proj := LookAt(vec3.Zero, vec3.UnitZ, vec3.UnitY)
		Expect(proj.Forward().ApproxEqual(vec3.UnitZ)).To(BeTrue())
	})
})
