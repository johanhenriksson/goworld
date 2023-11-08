package terrain_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/game/terrain"
	"github.com/johanhenriksson/goworld/math/ivec2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTerrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terrain Suite")
}

var _ = Describe("Map", func() {
	It("correctly calculates tile coords", func() {
		s := 16
		Expect(terrain.Floor(-17, s)).To(Equal(-2))
		Expect(terrain.Floor(-1, s)).To(Equal(-1))
		Expect(terrain.Floor(0, s)).To(Equal(0))
		Expect(terrain.Floor(1, s)).To(Equal(0))
		Expect(terrain.Floor(16, s)).To(Equal(1))
	})

	It("applies patches across multiple tiles", func() {
		point := func(v int) terrain.Point { return terrain.Point{Height: float32(v)} }
		m := terrain.NewMap("test", 2)
		patch := &terrain.Patch{
			Offset: ivec2.New(-1, -1),
			Size:   ivec2.New(2, 2),
			Points: [][]terrain.Point{
				{point(1), point(2), point(3)},
				{point(4), point(5), point(6)},
				{point(7), point(8), point(9)},
			},
			Source: m,
		}
		m.Set(patch)

		// t00 := m.GetTile(0, 0, false)
		// Expect(t00.Point(1, 1).Height).To(Equal(float32(1)))
		// Expect(t00.Point(1, 2).Height).To(Equal(float32(4)))
		// Expect(t00.Point(2, 1).Height).To(Equal(float32(2)))
		// Expect(t00.Point(2, 2).Height).To(Equal(float32(5)))
		//
		// t11 := m.GetTile(1, 1, false)
		// Expect(t11.Point(0, 0).Height).To(Equal(float32(5)))
		// Expect(t11.Point(0, 1).Height).To(Equal(float32(8)))
		// Expect(t11.Point(1, 0).Height).To(Equal(float32(6)))
		// Expect(t11.Point(1, 1).Height).To(Equal(float32(9)))
		//
		// t10 := m.GetTile(1, 0, false)
		// Expect(t10.Point(0, 1).Height).To(Equal(float32(2)))
		// Expect(t10.Point(0, 2).Height).To(Equal(float32(5)))
		// Expect(t10.Point(1, 1).Height).To(Equal(float32(3)))
		// Expect(t10.Point(1, 2).Height).To(Equal(float32(6)))
		//
		// t01 := m.GetTile(0, 1, false)
		// Expect(t01.Point(1, 0).Height).To(Equal(float32(4)))
		// Expect(t01.Point(1, 1).Height).To(Equal(float32(7)))
		// Expect(t01.Point(2, 0).Height).To(Equal(float32(5)))
		// Expect(t01.Point(2, 1).Height).To(Equal(float32(8)))

		// verify that the expected data is returned from Patch() calls
		patch2 := m.Get(patch.Offset, patch.Size)
		Expect(patch2.Size).To(Equal(patch.Size))
		Expect(patch2.Offset).To(Equal(patch.Offset))
		Expect(patch2.Points).To(Equal(patch.Points))
	})
})
