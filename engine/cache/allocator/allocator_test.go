package allocator_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/engine/cache/allocator"
)

var _ = Describe("", func() {
	It("allocates!", func() {
		fl := allocator.New(1024)
		block, err := fl.Alloc(16)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.Size).To(Equal(256))

		err = fl.Free(block.Offset)
		Expect(err).ToNot(HaveOccurred())

		block2, err := fl.Alloc(257)
		Expect(err).ToNot(HaveOccurred())
		Expect(block2.Size).To(Equal(512))
	})

	It("allocates correct sizes", func() {
		fl := allocator.New(1024)
		block, err := fl.Alloc(257)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.Size).To(Equal(512))
	})

	It("assigns tiers correctly", func() {
		Expect(allocator.GetBucketTier(allocator.MinAlloc)).To(Equal(0))
		Expect(allocator.GetBucketTier(allocator.MinAlloc + 1)).To(Equal(1))
		Expect(allocator.GetBucketTier(2 * allocator.MinAlloc)).To(Equal(1))
		Expect(allocator.GetBucketTier(2*allocator.MinAlloc + 1)).To(Equal(2))
	})

	It("checks powers of two", func() {
		Expect(allocator.IsPowerOfTwo(2)).To(BeTrue())
		Expect(allocator.IsPowerOfTwo(4)).To(BeTrue())
		Expect(allocator.IsPowerOfTwo(8)).To(BeTrue())
		Expect(allocator.IsPowerOfTwo(16)).To(BeTrue())

		Expect(allocator.IsPowerOfTwo(0)).To(BeFalse())
		Expect(allocator.IsPowerOfTwo(3)).To(BeFalse())
		Expect(allocator.IsPowerOfTwo(121)).To(BeFalse())
	})
})
