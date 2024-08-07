package spline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/math/spline"
	"github.com/johanhenriksson/goworld/math/vec2"
)

var _ = Describe("linear splines", func() {
	When("there are no points", func() {
		It("always return 0", func() {
			linear := spline.NewLinear()
			Expect(linear.Eval(0)).To(BeZero())
			Expect(linear.Eval(0.5)).To(BeZero())
			Expect(linear.Eval(1)).To(BeZero())
		})
	})

	When("there is one point", func() {
		It("always return the point", func() {
			linear := spline.NewLinear(vec2.New(0.5, 1))
			Expect(linear.Eval(0)).To(Equal(float32(1)))
			Expect(linear.Eval(0.5)).To(Equal(float32(1)))
			Expect(linear.Eval(1)).To(Equal(float32(1)))
		})
	})

	When("there are more than one point", func() {
		It("returns the linear interpolation between the points", func() {
			linear := spline.NewLinear(vec2.New(0, 0), vec2.New(1, 1))
			Expect(linear.Eval(0)).To(Equal(float32(0)))
			Expect(linear.Eval(0.5)).To(Equal(float32(0.5)))
			Expect(linear.Eval(1)).To(Equal(float32(1)))
		})

		It("returns the first point if t < X1", func() {
			linear := spline.NewLinear(vec2.New(0, 0), vec2.New(1, 1))
			Expect(linear.Eval(-1)).To(Equal(float32(0)))
		})

		It("returns the last point if t > Xn", func() {
			linear := spline.NewLinear(vec2.New(0, 0), vec2.New(1, 1))
			Expect(linear.Eval(2)).To(Equal(float32(1)))
		})
	})
})
