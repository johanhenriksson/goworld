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
	var obj *WithProp
	BeforeEach(func() {
		obj = object.NewComponent(&WithProp{
			Prop: object.NewProperty(1337),
		})
	})

	Context("basics", func() {
		It("stores data", func() {
			Expect(obj.Prop.Get()).To(Equal(1337), "wrong default value")

			obj.Prop.Set(42)
			Expect(obj.Prop.Get()).To(Equal(42), "setter should update stored value")
		})

		It("raises an OnChanged event", func() {
			var event int
			obj.Prop.OnChange.Subscribe(func(i int) {
				event = i
			})

			obj.Prop.Set(42)
			Expect(event).To(Equal(42))
		})
	})

	Context("generic property functions", func() {
		It("collects and modifies generic properties", func() {
			props := object.Properties(obj)
			Expect(props).To(HaveLen(1))

			props[0].SetAny(12)

			Expect(obj.Prop.Get()).To(Equal(12))
		})
	})
})
