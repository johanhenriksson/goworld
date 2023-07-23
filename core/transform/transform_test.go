package transform_test

import (
	. "github.com/johanhenriksson/goworld/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func TestLabel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transform Suite")
}

var _ = Describe("", func() {
	It("initializes properly", func() {
		t := transform.Identity()
		Expect(t.Forward()).To(Equal(vec3.UnitZ))
	})

	It("applies hierarchical transformation", func() {
		// values extracted from an identical scene set up in unity

		origin := transform.New(vec3.Zero, quat.Euler(30, 45, 0), vec3.One)
		camera := transform.New(vec3.New(0, 0, -10), quat.Ident(), vec3.One)

		camera.SetParent(origin)

		Expect(vec3.Distance(camera.WorldPosition(), vec3.New(-6.12, 5.0, -6.12))).To(BeNumerically("<", 0.1))
		Expect(vec3.Dot(camera.Forward(), vec3.New(0.61, -0.5, 0.61))).To(BeNumerically(">", 0.99))
	})

	It("maintains local transform when attaching to parent", func() {
		parent := transform.New(vec3.One, quat.Ident(), vec3.One)
		child := transform.New(vec3.One, quat.Ident(), vec3.One)
		child.SetParent(parent)
		Expect(child.WorldPosition()).To(Equal(vec3.One.Scaled(2)))
	})

	It("refreshes when parent is modified", func() {
		parent := transform.New(vec3.One, quat.Ident(), vec3.One)
		child := transform.New(vec3.One, quat.Ident(), vec3.One)
		child.SetParent(parent)
		parent.SetPosition(vec3.Zero)
		Expect(child.WorldPosition()).To(Equal(vec3.One))
	})

	It("sets world position relative to parent", func() {
		parent := transform.New(vec3.One, quat.Ident(), vec3.One)
		child := transform.New(vec3.One, quat.Ident(), vec3.One)
		child.SetParent(parent)
		child.SetWorldPosition(vec3.Zero)
		Expect(child.WorldPosition()).To(Equal(vec3.Zero))
		Expect(child.Position()).To(Equal(vec3.New(-1, -1, -1)))
	})

	It("sets world rotation relative to parent", func() {
		parent := transform.New(vec3.Zero, quat.Euler(0, 90, 0), vec3.One)
		child := transform.New(vec3.Zero, quat.Euler(0, 90, 0), vec3.One)
		child.SetParent(parent)
		Expect(child.WorldRotation().Euler()).To(Equal(vec3.New(0, 180, 0)))

		child.SetWorldRotation(quat.Euler(0, 90, 0))
		Expect(child.WorldRotation().Euler()).To(BeApproxVec3(vec3.New(0, 90, 0)))
		Expect(child.Rotation().Euler()).To(BeApproxVec3(vec3.Zero))
	})

	It("sets scale relative to parent", func() {
		parent := transform.New(vec3.Zero, quat.Euler(0, 90, 0), vec3.New(2, 2, 2))
		child := transform.New(vec3.Zero, quat.Euler(0, 90, 0), vec3.One)
		child.SetParent(parent)
		Expect(child.WorldScale()).To(BeApproxVec3(vec3.New(2, 2, 2)))
		child.SetWorldScale(vec3.One)
		Expect(child.WorldScale()).To(BeApproxVec3(vec3.One))
		Expect(child.Scale()).To(BeApproxVec3(vec3.New(0.5, 0.5, 0.5)))
	})
})
