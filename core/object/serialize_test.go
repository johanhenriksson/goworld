package object_test

import (
	. "github.com/johanhenriksson/goworld/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/object"
)

var _ = Describe("serialization", func() {
	var ctx object.Pool

	It("serializes basic objects", func() {
		a0 := object.Empty(ctx, "A")
		a1 := object.Copy(ctx, a0)

		Expect(a1.Name()).To(Equal(a0.Name()))
		Expect(a1.Enabled()).To(Equal(a0.Enabled()))
		Expect(a1.Transform().Position()).To(ApproxVec3(a0.Transform().Position()))
		Expect(a1.Transform().Rotation()).To(ApproxQuat(a0.Transform().Rotation()))
		Expect(a1.Transform().Scale()).To(ApproxVec3(a0.Transform().Scale()))
	})

	It("serializes nested objects", func() {
		a0 := object.Builder(object.Empty(ctx, "Parent")).
			Attach(object.Empty(ctx, "Child 1")).
			Attach(object.Empty(ctx, "Child 2")).
			Create()

		a1 := object.Copy(ctx, a0).(object.Object)
		children := a1.Children()
		Expect(len(children)).To(Equal(2))
		Expect(children[0].Name()).To(Equal("Child 1"))
		Expect(children[1].Name()).To(Equal("Child 2"))
	})

	It("serializes scenes", func() {
		a0 := object.Scene(ctx)
		a1 := object.Copy(ctx, a0)
		Expect(a1.Name()).To(Equal("Scene"))
	})
})
