package object_test

import (
	"fmt"
	"testing"

	"github.com/johanhenriksson/goworld/core/object"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

func TestObject2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Object2 Suite")
}

var _ = Describe("Object", func() {
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
})
