package physics_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/vertex"
)

var _ = Describe("Optimize", func() {
	It("correctly reduces the mesh", func() {
		vertices := []physics.Vertex{
			{vec3.Zero},
			{vec3.Zero},
			{vec3.New(1, 1, 1)},
			{vec3.Zero},
			{vec3.One},
		}
		indices := []uint32{
			4, 1, 2, 3, 0,
		}

		A := vertex.NewTriangles("test", vertices, indices)
		C := physics.CollisionMesh(A)

		m := C.(vertex.MutableMesh[physics.Vertex, uint32])
		Expect(m.Vertices()).To(HaveLen(2))
		Expect(m.Vertices()).To(Equal([]physics.Vertex{{vec3.One}, {vec3.Zero}}))
		Expect(m.Indices()).To(Equal([]uint32{0, 1, 0, 1, 1}))
	})
})
