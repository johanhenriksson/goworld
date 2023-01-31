package object_test

import (
	"github.com/johanhenriksson/goworld/core/object"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type A struct {
	object.T
	B *B
}

type B struct {
	object.T
}

func NewB() *B {
	return object.New(&B{})
}

func NewA() *A {
	return object.New(&A{
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

	It("instantiates object structs", func() {
		b := object.New(&B{})
		Expect(b.T).ToNot(BeNil())

		a := object.New(&A{
			B: b,
		})
		Expect(a.Children()).To(HaveLen(1))
		Expect(b.Parent()).To(Equal(a))
	})

	It("supports object interfaces")
})
