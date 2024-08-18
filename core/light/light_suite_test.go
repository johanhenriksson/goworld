package light_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
)

func TestLight(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "core/light")
}

type TestShadowStore struct{}

var _ light.ShadowmapStore = (*TestShadowStore)(nil)

func (t *TestShadowStore) Lookup(lit light.T, index int) (int, bool) {
	return index, true
}

var _ = Describe("serialization", func() {
	var pool object.Pool
	BeforeEach(func() {
		pool = object.NewPool()
	})

	It("serializes point lights", func() {
		a0 := light.NewPoint(pool, light.PointArgs{
			Color:     color.Red,
			Intensity: 1,
			Range:     2,
		})
		a1 := object.Copy(pool, a0)
		Expect(a0.ID()).ToNot(Equal(a1.ID()))

		Expect(a1.Color.Get()).To(Equal(a0.Color.Get()))
		Expect(a1.Intensity.Get()).To(Equal(a0.Intensity.Get()))
		Expect(a1.Range.Get()).To(Equal(a0.Range.Get()))

		ss := &TestShadowStore{}
		Expect(a1.LightData(ss)).To(Equal(a0.LightData(ss)))
	})

	It("serializes directional lights", func() {
		a0 := light.NewDirectional(pool, light.DirectionalArgs{
			Cascades:  4,
			Color:     color.Red,
			Intensity: 1,
			Shadows:   true,
		})
		a1 := object.Copy(pool, a0)
		Expect(a0.ID()).ToNot(Equal(a1.ID()))

		Expect(a1.Color.Get()).To(Equal(a0.Color.Get()))
		Expect(a1.Intensity.Get()).To(Equal(a0.Intensity.Get()))
		Expect(a1.Shadows.Get()).To(Equal(a0.Shadows.Get()))
		Expect(a1.Cascades.Get()).To(Equal(a0.Cascades.Get()))

		ss := &TestShadowStore{}
		Expect(a1.LightData(ss)).To(Equal(a0.LightData(ss)))
	})
})
