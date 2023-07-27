package object_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/object"
)

type WithProp struct {
	object.Component
	Prop object.Property[int]
}

var _ = Describe("Properties", func() {
	Context("basics", func() {
		It("stores data", func() {
			cmp := object.NewComponent(&WithProp{
				Prop: object.NewProperty(1337),
			})
			Expect(cmp.Prop.Get()).To(Equal(1337), "wrong default value")

			cmp.Prop.Set(42)
			Expect(cmp.Prop.Get()).To(Equal(42), "setter should update stored value")
		})
	})

	Context("generic property functions", func() {
		It("collects and modifies generic properties", func() {
			cmp := object.NewComponent(&WithProp{
				Prop: object.NewProperty(1337),
			})

			props := object.Properties(cmp)
			Expect(props).To(HaveLen(1))

			props[0].SetAny(12)

			Expect(cmp.Prop.Get()).To(Equal(12))
		})
	})
})
