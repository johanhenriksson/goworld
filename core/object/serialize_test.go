package object

import (
	"slices"

	. "github.com/johanhenriksson/goworld/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("serialization", func() {
	type ObjectWithComponent struct {
		Object

		// Reference to a component that is a direct child of the object
		Pointer Component

		// Reference to any component in the scene
		Ref Ref[Component]

		Value Property[int]
	}

	var pool Pool
	BeforeEach(func() {
		pool = NewPool()
		Register[*ObjectWithComponent](Type{})
	})

	Context("objects", func() {
		It("serializes basic objects", func() {
			a0 := Empty(pool, "A")
			a1 := Copy(pool, a0)

			Expect(a1.Name()).To(Equal(a0.Name()))
			Expect(a1.Enabled()).To(Equal(a0.Enabled()))
			Expect(a1.Transform().Position()).To(ApproxVec3(a0.Transform().Position()))
			Expect(a1.Transform().Rotation()).To(ApproxQuat(a0.Transform().Rotation()))
			Expect(a1.Transform().Scale()).To(ApproxVec3(a0.Transform().Scale()))
		})

		It("serializes nested objects", func() {
			a0 := Builder(Empty(pool, "Parent")).
				Attach(Empty(pool, "Child 1")).
				Attach(Empty(pool, "Child 2")).
				Create()

			a1 := Copy(pool, a0).(Object)
			children := slices.Collect(a1.Children())
			Expect(len(children)).To(Equal(2))
			Expect(children[0].Name()).To(Equal("Child 1"))
			Expect(children[1].Name()).To(Equal("Child 2"))
		})

		It("serializes scenes", func() {
			a0 := Scene(pool)
			a1 := Copy(pool, a0)
			Expect(a1.Name()).To(Equal("Scene"))
		})

		It("serializes child references", func() {
			obj := NewObject(pool, "ObjectWithComponent", &ObjectWithComponent{
				Pointer: Empty(pool, "Child"),
				Value:   NewProperty(123),
			})
			obj.Ref.Set(obj.Pointer)

			enc := &MemorySerializer{}
			Serialize(enc, obj)

			// object header
			// child type
			// child object header
			// child
			// object
			// reference to Thing
			// reference prop -> null
			// value prop -> 0
			Expect(enc.Stream).To(HaveLen(8))

			obj, err := Deserialize[*ObjectWithComponent](pool, enc)
			Expect(err).ToNot(HaveOccurred())

			// pointer should be set and point to the child
			Expect(obj.Pointer).ToNot(BeNil())
			Expect(obj.Len()).To(Equal(1))
			Expect(obj.Pointer.ID()).To(Equal(obj.Child(0).ID()))

			// value should be preserved
			Expect(obj.Value.Get()).To(Equal(123))

			// reference should be set and point to the child
			targetRef, ok := obj.Ref.Get()
			Expect(ok).To(BeTrue())
			Expect(targetRef.ID()).To(Equal(obj.Pointer.ID()))
		})
	})

	Context("components", func() {
		type A struct {
			Component
		}
		type B struct {
			*A
		}

		BeforeEach(func() {
			Register[*A](Type{})
			Register[*B](Type{})
		})

		It("serializes base components", func() {
			a0 := NewComponent(pool, &A{})
			a1 := Copy(pool, a0)
			Expect(a1).ToNot(BeNil())
		})

		It("serializes derived components", func() {
			b0 := NewComponent(pool, &B{
				A: NewComponent(pool, &A{}),
			})
			b1 := Copy(pool, b0)
			Expect(b1).ToNot(BeNil())
		})
	})
})
