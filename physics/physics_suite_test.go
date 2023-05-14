package physics_test

import (
	"log"
	"testing"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPhysics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Physics Suite")
}

var _ = Describe("physics tests", func() {
	It("creates a new dynamics world", func() {
		world := physics.NewWorld()
		world.SetGravity(vec3.New(0, -10, 0))

		boxShape := physics.NewBoxShape(vec3.One)
		box := physics.NewRigidBody(10, boxShape)
		// box.SetPosition(vec3.New(0, 2, 0))
		world.AddRigidBody(box)

		groundShape := physics.NewBoxShape(vec3.New(100, 1, 100))
		ground := physics.NewRigidBody(0, groundShape)
		ground.SetPosition(vec3.New(0, -5, 0))
		world.AddRigidBody(ground)

		steps := 100
		for i := 0; i < steps; i++ {
			world.Step(float32(1) / 60)
			log.Println(box.Position())
		}
	})
})
