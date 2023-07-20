package object_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/object"
)

type WithProp struct {
	object.Component
	Prop *object.Property[int]
}

var _ = Describe("Properties", func() {
	It("works", func() {
		cmp := object.NewComponent(&WithProp{
			Prop: object.NewProperty(1337),
		})

		props := object.Properties(cmp)
		Expect(props).To(HaveLen(1))

		props[0].SetAny(12)

		Expect(cmp.Prop.Get()).To(Equal(12))
	})
})
