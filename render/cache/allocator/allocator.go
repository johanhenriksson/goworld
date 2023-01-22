package allocator

import (
	"errors"
	"math"
)

// Minimum allocation block size
const MinAlloc = 256

var minBucketTier = int(math.Log2(MinAlloc))

var ErrOutOfMemory = errors.New("out of memory")
var ErrInvalidFree = errors.New("illegal free() call")

type Block struct {
	Offset int
	Size   int
}

type T interface {
	Alloc(int) (Block, error)
	Free(int) error
}

// Buddy allocator implementation
type buddy struct {
	free  [][]Block
	alloc map[int]int
	top   int
}

func New(size int) T {
	if !IsPowerOfTwo(size) {
		panic("allocator size must be a power of 2")
	}

	top := GetBucketTier(size)
	free := make([][]Block, top+1)
	free[top] = []Block{{Offset: 0, Size: size}}

	return &buddy{
		top:   top,
		free:  free,
		alloc: map[int]int{},
	}
}

func (f *buddy) Alloc(size int) (Block, error) {
	tier := GetBucketTier(size)
	block, err := f.getBlock(tier)
	if err != nil {
		return Block{}, err
	}
	f.alloc[block.Offset] = block.Size
	return block, nil
}

func (f *buddy) getBlock(tier int) (Block, error) {
	if tier > f.top {
		return Block{}, ErrOutOfMemory
	}

	if bucket := f.free[tier]; len(bucket) > 0 {
		lastIdx := len(bucket) - 1
		block := bucket[lastIdx]
		f.free[tier] = bucket[:lastIdx]
		return block, nil
	}

	split, err := f.getBlock(tier + 1)
	if err != nil {
		return Block{}, err
	}

	size := split.Size / 2
	f.free[tier] = append(f.free[tier], Block{
		Offset: split.Offset + size,
		Size:   size,
	})
	return Block{
		Offset: split.Offset,
		Size:   size,
	}, nil
}

func (f *buddy) Free(offset int) error {
	size, exists := f.alloc[offset]
	if !exists {
		return ErrInvalidFree
	}

	freed := Block{
		Offset: offset,
		Size:   size,
	}

	tier := GetBucketTier(size)
	f.free[tier] = append(f.free[tier], freed)

	// mark as free
	delete(f.alloc, offset)

	// merge buddies
	f.merge(tier, freed, len(f.free[tier])-1)

	return nil
}

func (f *buddy) merge(tier int, block Block, blockIdx int) {
	// nothing to merge at the top tier
	if tier >= f.top {
		return
	}
	level := f.free[tier]

	// figure out the offset of our buddy block, and the resulting offset of a merge
	buddyOffset := 0
	mergedOffset := 0
	if block.Offset%(2*block.Size) == 0 {
		// we are an even block, buddy is after
		buddyOffset = block.Offset + block.Size
		mergedOffset = block.Offset
	} else {
		// we are an odd block, buddy is before
		buddyOffset = block.Offset - block.Size
		mergedOffset = buddyOffset
	}

	// check if buddy block is allocated
	if _, allocated := f.alloc[buddyOffset]; allocated {
		// yes - then we can't merge
		return
	}

	// find the free list index of the buddy
	var buddyIdx int
	for candidateIdx, candidate := range level {
		if candidate.Offset == buddyOffset {
			buddyIdx = candidateIdx
			break
		}
	}

	// remove both blocks from free list
	// todo: implement using a linked list
	if buddyIdx > blockIdx {
		f.free[tier] = append(append(level[:blockIdx], level[blockIdx+1:buddyIdx]...), level[buddyIdx+1:]...)
	} else {
		f.free[tier] = append(append(level[:buddyIdx], level[buddyIdx+1:blockIdx]...), level[blockIdx+1:]...)
	}

	// add the merged block to free list, on the next tier
	merged := Block{
		Offset: mergedOffset,
		Size:   2 * block.Size,
	}
	f.free[tier+1] = append(f.free[tier+1], merged)

	// attempt to merge next level
	f.merge(tier+1, merged, len(f.free[tier+1])-1)
}

func GetBucketTier(size int) int {
	tier := int(math.Log2(float64(size-1))) + 1

	tier -= minBucketTier
	if tier < 0 {
		return 0
	}

	return tier
}

func IsPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}
