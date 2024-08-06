package buffer

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/util"
)

// Strongly typed array buffer
type Array[K any] interface {
	T

	// Set the value of element i and flushes the buffer.
	Set(index int, data K)

	// Sets a range of elements, starting at i and flushes the buffer.
	SetRange(index int, data []K)

	// Count returns the number of items in the array
	Count() int

	// Stride returns the aligned byte size of a single element
	Stride() int
}

type array[K any] struct {
	T
	stride int
	count  int
}

// NewArray creates a new typed array buffer.
// When allocating arrays, the Size argument is the number of elements
func NewArray[K any](device device.T, args Args) Array[K] {
	align, maxSize := GetBufferLimits(device, args.Usage)

	var empty K
	kind := reflect.TypeOf(empty)
	sizeof := int(kind.Size())

	stride := util.Align(sizeof, align)
	count := args.Size
	size := count * stride
	if size > maxSize {
		panic(fmt.Sprintf("buffer is too large for the specified usage. size: %d, max: %d", size, maxSize))
	}

	args.Size = size
	buffer := New(device, args)

	return &array[K]{
		T:      buffer,
		stride: stride,
		count:  count,
	}
}

func (a *array[K]) Set(index int, data K) {
	a.Write(index*a.stride, &data)
	a.Flush()
}

func (a *array[K]) SetRange(offset int, data []K) {
	for i, el := range data {
		a.Write((i+offset)*a.stride, &el)
	}
	a.Flush()
}

func (a *array[K]) Count() int  { return a.count }
func (a *array[K]) Stride() int { return a.stride }
