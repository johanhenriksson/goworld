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
			Min: vec3.New(-1, -1, -1),
			Max: vec3.New(1, 1, 1),
		}
		intersect := func(origin, ray, expect vec3.T) {
			hit, point := box.Intersect(&physics.Ray{
				Origin: origin,
				Dir:    ray,
			})
			Expect(hit).To(BeTrue())
			Expect(point.ApproxEqual(expect)).To(BeTrue())
		}

		intersect(vec3.New(2, 0, 0), vec3.UnitXN, vec3.New(1, 0, 0))
		intersect(vec3.New(-2, 0, 0), vec3.UnitX, vec3.New(-1, 0, 0))

		intersect(vec3.New(0, 2, 0), vec3.UnitYN, vec3.New(0, 1, 0))
		intersect(vec3.New(0, -2, 0), vec3.UnitY, vec3.New(0, -1, 0))

		intersect(vec3.New(0, 0, 2), vec3.UnitZN, vec3.New(0, 0, 1))
		intersect(vec3.New(0, 0, -2), vec3.UnitZ, vec3.New(0, 0, -1))
	})
})
