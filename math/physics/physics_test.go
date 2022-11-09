package physics_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPhysics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Physics Suite")
}

var _ = Describe("Sphere", func() {
	It("correctly intersects a sphere", func() {
		sphere := physics.Sphere{
			Center: vec3.Zero,
			Radius: 8,
		}
		ray := physics.Ray{
			Origin: vec3.New(-16, 1.7, -0.17),
			Dir:    vec3.New(1, 0, 0).Normalized(),
		}
		hit, _ := sphere.Intersect(&ray)
		Expect(hit).To(BeTrue())
		// Expect(point.X).To(Equal(float32(-5)))
	})
})

var _ = Describe("Box", func() {
	It("intersects a box", func() {
		box := physics.Box{
			Min: vec3.New(-8, -8, -8),
			Max: vec3.New(8, 8, 8),
		}
		ray := physics.Ray{
			Origin: vec3.New(-2.7, 21.6, -14.3),
			Dir:    vec3.New(0.124, -0.66, 0.74).Normalized(),
		}
		hit, _ := box.Intersect(&ray)
		Expect(hit).To(BeTrue())
	})
})
