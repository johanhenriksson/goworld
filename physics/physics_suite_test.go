package physics

import (
	. "github.com/johanhenriksson/goworld/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func TestPhysics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "physics")
}

var _ = Describe("general physics tests", func() {
	var (
		scene object.Object
		world *World
	)

	BeforeEach(func() {
		scene = object.Scene()
		world = NewWorld()
		object.Attach(scene, world)
	})

	Context("rigidbody dynamics", func() {
		var (
			body *RigidBody
			box  *Box
			obj  object.Object
		)

		BeforeEach(func() {
			body = NewRigidBody(1)
			box = NewBox(vec3.One)
			obj = object.Builder(object.Empty("physics object")).
				Attach(body).
				Attach(box).
				Parent(scene).
				Create()
		})

		It("connects and objects to the physics world", func() {
			Expect(body.world).To(Equal(world))
			Expect(body.shape).To(Equal(box))
		})

		It("simulates rigidbody movement", func() {
			// ... run simulation ...
			world.Update(scene, 0.1)

			Expect(obj.Transform().Position().Y).To(BeNumerically("<", 0), "the box should have fallen")
		})

		It("returns correct raycast results", func() {
			// raycast towards the origin should hit the box
			hit, ok := world.Raycast(vec3.New(0, 5, 0), vec3.Zero, All)
			Expect(ok).To(BeTrue())
			Expect(hit.Point).To(ApproxVec3(vec3.New(0, 0.5, 0)))
		})
	})

	Context("kinematic rigidbodies", func() {
		var (
			body *RigidBody
			box  *Box
			obj  object.Object
		)

		BeforeEach(func() {
			body = NewRigidBody(0)
			box = NewBox(vec3.One)
			obj = object.Builder(object.Empty("physics object")).
				Attach(body).
				Attach(box).
				Parent(scene).
				Create()
		})

		It("connects and objects to the physics world", func() {
			Expect(body.world).To(Equal(world))
			Expect(body.shape).To(Equal(box))
		})

		It("kinematic rigidbodies do not move", func() {
			// ... run simulation ...
			world.Update(scene, 0.1)

			Expect(obj.Transform().Position().Y).To(BeNumerically("~", 0), "the box should not have fallen")
		})

		It("returns correct raycast results", func() {
			// raycast towards the origin should hit the box
			hit, ok := world.Raycast(vec3.New(0, 5, 0), vec3.Zero, All)
			Expect(ok).To(BeTrue())
			Expect(hit.Point).To(ApproxVec3(vec3.New(0, 0.5, 0)))
		})
	})
})
