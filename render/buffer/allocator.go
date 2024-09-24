package buffer

import (
	"errors"

	"github.com/johanhenriksson/goworld/render/device"
)

var ErrOutOfMemory = errors.New("out of memory")
var ErrInvalidFree = errors.New("illegal free() call")

type Allocator interface {
	Alloc(size int) (Block, error)
	Free(b Block) error
}

type Block struct {
	buffer T
	size   int
	offset int
}

func EntireBuffer(buffer T) Block {
	return Block{buffer, buffer.Size(), 0}
}

func (b *Block) Buffer() T   { return b.buffer }
func (b *Block) Size() int   { return b.size }
func (b *Block) Offset() int { return b.offset }

func (b *Block) Address() device.Address {
	return b.buffer.Address() + device.Address(b.offset)
}

func isPowerOfTwo(x uint) bool {
	return (x & (x - 1)) == 0
}

func Align(offset, align int) int {
	remainder := offset % align
	if remainder == 0 {
		return offset
	}
	return offset + (align - remainder)
}
