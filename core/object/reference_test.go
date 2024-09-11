package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ObjectWithReference struct {
	Object
	Reference Ref[Object]
}

var _ = Describe("", func() {
	var pool Pool
	var a *ObjectWithReference
	var b Object

	BeforeEach(func() {
		Register[*ObjectWithReference](Type{})

		pool = NewPool()
		a = NewObject(pool, "a", &ObjectWithReference{})
		b = Empty(pool, "b")
		Attach(a, b)
		a.Reference.Set(b)
	})

	It("serializes empty references", func() {
		obj := NewObject(pool, "obj", &ObjectWithReference{})
		out := Copy(pool, obj)
		Expect(out).ToNot(BeNil())
		Expect(out.Len()).To(Equal(0))
		ref, ok := out.Reference.Get()
		Expect(ok).To(BeFalse())
		Expect(ref).To(BeNil())
	})

	It("serializes correctly", func() {
		sa := Copy(pool, a)

		Expect(sa.ID()).ToNot(Equal(a.ID()))
		Expect(sa.Len()).To(Equal(1))

		sbref, ok := sa.Reference.Get()
		Expect(ok).To(BeTrue(), "handle should be valid")

		sb := sa.Child(0).(Object)
		Expect(sb.ID()).To(Equal(sbref.ID()))
	})
})
