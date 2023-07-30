package object_test

import (
	. "github.com/johanhenriksson/goworld/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
	"testing"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type A struct {
	object.Object
	B *B
}

type B struct {
	object.Component
}

func NewB() *B {
	return object.NewComponent(&B{})
}

func NewA() *A {
	return object.New("a", &A{
		B: NewB(),
	})
}

func TestObject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "core/object")
}

var _ = Describe("Object", func() {
	It("generates proper string keys", func() {
		a := object.Empty("a")
		key := object.Key("hello", a)
		Expect(key[:5]).To(Equal("hello"))
		Expect(key[5]).To(Equal(byte('-')))
		Expect(key[6:]).To(Equal(fmt.Sprintf("%02x", a.ID())))
	})

	It("attaches & detaches children", func() {
		a := object.Empty("A")
		b := object.Empty("B")

		object.Attach(a, b)
		Expect(a.Children()).To(HaveLen(1))
		Expect(a.Parent()).To(BeNil())
		Expect(b.Children()).To(HaveLen(0))
		Expect(b.Parent()).To(Equal(a))

		object.Detach(b)
		Expect(a.Children()).To(HaveLen(0))
		Expect(a.Parent()).To(BeNil())
		Expect(b.Children()).To(HaveLen(0))
		Expect(b.Parent()).To(BeNil())
	})

	It("instantiates component structs", func() {
		b := object.NewComponent(&B{})
		Expect(b.Component).ToNot(BeNil())

		a := object.New("A", &A{
			B: b,
		})
		Expect(a.Children()).To(HaveLen(1))
		Expect(b.Parent()).To(Equal(a))
	})

	It("correctly creates a key string", func() {
		b := object.NewComponent(&B{})
		key := object.Key("hello", b)
		Expect(key).To(Equal(fmt.Sprintf("hello-%x", b.ID())))
	})

	Context("ghost object", func() {
		It("follows the target object", func() {
			target := object.Empty("target")
			ghost := object.Ghost("ghost", target.Transform())

			pos := vec3.New(10, 20, 30)
			target.Transform().SetPosition(pos)
			Expect(ghost.Transform().WorldPosition()).To(BeApproxVec3(pos))
		})

		It("propagates events", func() {
			target := object.Empty("target")
			ghost := object.Ghost("ghost", target.Transform())

			triggered := false
			ghost.Transform().OnChange().Subscribe(func(t transform.T) {
				triggered = true
			})

			target.Transform().SetScale(vec3.New(2, 2, 2))
			Expect(triggered).To(BeTrue())
		})

		It("maintains heirarchical properties", func() {
			target1 := object.Builder(object.Empty("target1")).Position(vec3.One).Create()
			target2 := object.Builder(object.Empty("target2")).Position(vec3.One).Create()
			object.Attach(target1, target2)

			ghost1 := object.Ghost("ghost1", target1.Transform())
			ghost2 := object.Ghost("ghost2", target2.Transform())
			object.Attach(ghost1, ghost2)

			// simulate a component attached to target1/ghost1
			ghostc := object.Ghost("ghostc", target1.Transform())
			object.Attach(ghost1, ghostc)

			// validate transform
			Expect(target1.Transform().WorldPosition()).To(Equal(vec3.New(1, 1, 1)))
			Expect(target2.Transform().WorldPosition()).To(Equal(vec3.New(2, 2, 2)))
			Expect(ghost1.Transform().WorldPosition()).To(Equal(target1.Transform().WorldPosition()))
			Expect(ghost2.Transform().WorldPosition()).To(Equal(target2.Transform().WorldPosition()))

			triggerTarget := false
			target2.Transform().OnChange().Subscribe(func(t transform.T) {
				triggerTarget = true
			})
			triggerGhost := false
			ghost2.Transform().OnChange().Subscribe(func(t transform.T) {
				triggerGhost = true
			})

			// modifying component transform
			ghostc.Transform().SetPosition(vec3.Zero)

			// should update everything:
			Expect(target1.Transform().WorldPosition()).To(Equal(vec3.New(0, 0, 0)))
			Expect(target2.Transform().WorldPosition()).To(Equal(vec3.New(1, 1, 1)))
			Expect(ghost1.Transform().WorldPosition()).To(Equal(target1.Transform().WorldPosition()))
			Expect(ghost2.Transform().WorldPosition()).To(Equal(target2.Transform().WorldPosition()))

			Expect(triggerTarget).To(BeTrue())
			Expect(triggerGhost).To(BeTrue())
		})
	})
})
