package object_test

import (
	"log"

	. "github.com/johanhenriksson/goworld/core/object"
	. "github.com/johanhenriksson/goworld/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
	"testing"

	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type A struct {
	Object
	B *B
}

type B struct {
	Component
}

func NewB(pool Pool) *B {
	return NewComponent(pool, &B{})
}

func NewA(pool Pool) *A {
	return NewObject(pool, "a", &A{
		B: NewB(pool),
	})
}

func TestObject(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "core/object")
}

var _ = Describe("Object", func() {
	pool := NewPool()

	It("generates proper string keys", func() {
		a := Empty(pool, "a")
		key := Key("hello", a)
		Expect(key[:5]).To(Equal("hello"))
		Expect(key[5]).To(Equal(byte('-')))
		Expect(key[6:]).To(Equal(fmt.Sprintf("%d", a.ID())))
	})

	It("attaches & detaches children", func() {
		a := Empty(pool, "A")
		b := Empty(pool, "B")

		Attach(a, b)
		Expect(a.Children()).To(HaveLen(1))
		Expect(a.Parent()).To(BeNil())
		Expect(b.Children()).To(HaveLen(0))
		Expect(b.Parent()).To(Equal(a))

		Detach(b)
		Expect(a.Children()).To(HaveLen(0))
		Expect(a.Parent()).To(BeNil())
		Expect(b.Children()).To(HaveLen(0))
		Expect(b.Parent()).To(BeNil())
	})

	It("instantiates component structs", func() {
		b := NewComponent(pool, &B{})
		Expect(b.Component).ToNot(BeNil())

		a := NewObject(pool, "A", &A{
			B: b,
		})
		Expect(a.Children()).To(HaveLen(1))
		Expect(b.Parent()).To(Equal(a))
	})

	It("correctly creates a key string", func() {
		b := NewComponent(pool, &B{})
		key := Key("hello", b)
		Expect(key).To(Equal(fmt.Sprintf("hello-%x", b.ID())))
	})

	Context("ghost object", func() {
		It("follows the target object", func() {
			target := Empty(pool, "target")
			ghost := Ghost(pool, "ghost", target.Transform())

			pos := vec3.New(10, 20, 30)
			target.Transform().SetPosition(pos)
			Expect(ghost.Transform().WorldPosition()).To(ApproxVec3(pos))
		})

		It("propagates events", func() {
			target := Empty(pool, "target")
			ghost := Ghost(pool, "ghost", target.Transform())

			triggered := false
			ghost.Transform().OnChange().Subscribe(func(t transform.T) {
				triggered = true
			})

			target.Transform().SetScale(vec3.New(2, 2, 2))
			Expect(triggered).To(BeTrue())
		})

		It("maintains heirarchical properties", func() {
			target1 := Builder(Empty(pool, "target1")).Position(vec3.One).Create()
			target2 := Builder(Empty(pool, "target2")).Position(vec3.One).Create()
			Attach(target1, target2)

			ghost1 := Ghost(pool, "ghost1", target1.Transform())
			ghost2 := Ghost(pool, "ghost2", target2.Transform())
			Attach(ghost1, ghost2)

			// simulate a component attached to target1/ghost1
			ghostc := Ghost(pool, "ghostc", target1.Transform())
			Attach(ghost1, ghostc)

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
