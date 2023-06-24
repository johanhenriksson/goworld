package transform_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

		camera.Recalculate(origin)

		Expect(vec3.Distance(camera.WorldPosition(), vec3.New(-6.12, 5.0, -6.12))).To(BeNumerically("<", 0.1))
		Expect(vec3.Dot(camera.Forward(), vec3.New(0.61, -0.5, 0.61))).To(BeNumerically(">", 0.99))
	})
})
