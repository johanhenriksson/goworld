package chunk_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/voxel"
)

func TestChunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chunk Suite")
}

var _ = Describe("operations", func() {
	Context("crop", func() {
		It("correctly removes surrounding empty space", func() {
			crop := chunk.New("t", 3, 3, 3)
			crop.Set(1, 1, 1, voxel.Red)
			crop.Set(1, 0, 1, voxel.Red)
			chunk.Crop(crop)
			Expect(crop.Sx).To(Equal(1))
			Expect(crop.Sy).To(Equal(2))
			Expect(crop.Sz).To(Equal(1))
		})
	})
})
