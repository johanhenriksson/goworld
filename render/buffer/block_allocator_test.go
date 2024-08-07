package buffer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/render/buffer"
)

var _ = Describe("block allocator", func() {
	var alloc buffer.Allocator
	var buf buffer.T

	BeforeEach(func() {
		buf = &mockBuffer{size: 1024}
		alloc = buffer.NewBlockAllocator(buf)
	})

	It("allocates correctly sized blocks", func() {
		block, err := alloc.Alloc(100)
		Expect(err).To(BeNil())
		Expect(block.Offset()).To(Equal(0))
		Expect(block.Size()).To(Equal(128))

		block2, err := alloc.Alloc(100)
		Expect(err).To(BeNil())
		Expect(block2.Offset()).To(Equal(128))
		Expect(block2.Size()).To(Equal(128))

		block3, err := alloc.Alloc(500)
		Expect(err).To(BeNil())
		Expect(block3.Offset()).To(Equal(512))
		Expect(block3.Size()).To(Equal(512))
	})

	It("frees blocks", func() {
		block1, err := alloc.Alloc(1000)
		Expect(err).To(BeNil())
		alloc.Free(block1)

		// this wont work unless the block was returned
		_, err = alloc.Alloc(1000)
		Expect(err).To(BeNil())
	})

	It("merges free blocks", func() {
		block1, err := alloc.Alloc(50)
		Expect(err).To(BeNil())
		Expect(block1.Offset()).To(Equal(0))
		alloc.Free(block1)

		// this wont work unless blocks are merged
		block2, err := alloc.Alloc(1000)
		Expect(err).To(BeNil())
		Expect(block2.Offset()).To(Equal(0))
	})

	It("throws an error if the buffer size is not a power of two", func() {
		_, err := alloc.Alloc(1000)
		Expect(err).To(BeNil())

		_, err = alloc.Alloc(1000)
		Expect(err).To(MatchError(buffer.ErrOutOfMemory))
	})
})
