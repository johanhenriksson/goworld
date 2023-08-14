package light_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLight(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "core/light")
}

var _ = Describe("serialization", func() {
	It("serializes point lights", func() {
		a0 := light.NewPoint(light.PointArgs{
			Color:     color.Red,
			Intensity: 1,
			Range:     2,
		})
		cpy := object.Copy(a0)
		a1, ok := cpy.(*light.Point)
		Expect(ok).To(BeTrue(), "result should be a point light")

		Expect(a1.Color.Get()).To(Equal(a0.Color.Get()))
		Expect(a1.Intensity.Get()).To(Equal(a0.Intensity.Get()))
		Expect(a1.Range.Get()).To(Equal(a0.Range.Get()))
	})

	It("serializes directional lights", func() {
		a0 := light.NewDirectional(light.DirectionalArgs{
			Color:     color.Red,
			Intensity: 1,
			Shadows:   true,
		})
		cpy := object.Copy(a0)
		a1, ok := cpy.(*light.Directional)
		Expect(ok).To(BeTrue(), "result should be a directional light")

		Expect(a1.Color.Get()).To(Equal(a0.Color.Get()))
		Expect(a1.Intensity.Get()).To(Equal(a0.Intensity.Get()))
		Expect(a1.Shadows.Get()).To(Equal(a0.Shadows.Get()))
	})
})
