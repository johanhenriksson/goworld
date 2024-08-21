package object

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ObjectWithReference struct {
	Object
	Reference Reference[Object]
}

var _ = Describe("", func() {
	var pool Pool
	var a *ObjectWithReference
	var b Object

	BeforeEach(func() {
		Register[*ObjectWithReference](TypeInfo{})

		pool = NewPool()
		a = New(pool, "a", &ObjectWithReference{})
		b = Empty(pool, "b")
		Attach(a, b)
		a.Reference.Set(b)
	})

	It("encodes references", func() {
		s := &MemorySerializer{}
		val := reflect.ValueOf(a).Elem()
		err := encodeReferences(s, val)
		Expect(err).To(BeNil())

		Expect(s.Stream).To(HaveLen(1))

		err = decodeReferences(pool, s, val)
		Expect(err).To(BeNil())

		ref, ok := a.Reference.Get()
		Expect(ok).To(BeTrue())
		Expect(ref.ID()).To(Equal(b.ID()))
	})

	It("serializes empty references", func() {
		obj := New(pool, "obj", &ObjectWithReference{})
		out := Copy(pool, obj)
		Expect(out).ToNot(BeNil())
		Expect(out.Children()).To(HaveLen(0))
		ref, ok := out.Reference.Get()
		Expect(ok).To(BeFalse())
		Expect(ref).To(BeNil())
	})

	It("serializes correctly", func() {
		sa := Copy(pool, a)

		Expect(sa.ID()).ToNot(Equal(a.ID()))
		Expect(sa.Children()).To(HaveLen(1))

		sbref, ok := sa.Reference.Get()
		Expect(ok).To(BeTrue(), "handle should be valid")

		sb := sa.Children()[0].(Object)
		Expect(sb.ID()).To(Equal(sbref.ID()))
	})
})
