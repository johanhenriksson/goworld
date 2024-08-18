package object_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/object"
)

type ObjectWithReference struct {
	object.Object
	Reference object.Reference[object.Object]
}

var _ object.Serializable = (*ObjectWithReference)(nil)

func (o *ObjectWithReference) Serialize(enc object.Encoder) error {
	if err := object.Serialize(enc, o.Object); err != nil {
		return err
	}
	if err := o.Reference.Serialize(enc); err != nil {
		return err
	}
	return nil
}

func DeserializeObjectWithReference(ctx object.Pool, dec object.Decoder) (object.Component, error) {
	obj, err := object.Deserialize[object.Object](ctx, dec)
	if err != nil {
		return nil, err
	}
	o := &ObjectWithReference{
		Object: obj,
	}
	o.Reference, err = object.DeserializeReference[object.Object](ctx, dec)
	return o, err
}

var _ = FDescribe("", func() {
	BeforeEach(func() {
		object.Register[*ObjectWithReference](object.TypeInfo{
			Name:        "ObjectWithReference",
			Deserialize: DeserializeObjectWithReference,
		})
	})

	It("serializes correctly", func() {
		ctx := object.NewPool()
		a := object.New(ctx, "a", &ObjectWithReference{})
		b := object.Empty(ctx, "b")
		object.Attach(a, b)
		a.Reference.Set(b)

		sa := object.Copy(ctx, a)
		Expect(sa.ID()).ToNot(Equal(a.ID()))
		Expect(sa.Children()).To(HaveLen(1))

		sbref, ok := sa.Reference.Get()
		Expect(ok).To(BeTrue(), "handle should be valid")

		sb := sa.Children()[0].(object.Object)
		Expect(sb.ID()).To(Equal(sbref.ID()))
	})
})
