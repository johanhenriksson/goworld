package buffer

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"
)

// Strongly typed array buffer
type Array[K any] interface {
	T

	// Set the value of element i
	Set(index int, data K)

	// Sets a range of elements, starting at i
	SetRange(index int, data []K)

	// Count returns the number of items in the array
	Count() int

	// Element returns the aligned byte size of a single element
	Element() int
}

type array[K any] struct {
	T
	element int
	count   int
}

// NewArray creates a new typed array buffer.
// When allocating arrays, the Size argument is the number of elements
func NewArray[K any](device device.T, args Args) Array[K] {
	align, maxSize := GetBufferLimits(device, args.Usage)

	var empty K
	kind := reflect.TypeOf(empty)

	element := util.Align(int(kind.Size()), align)

	count := args.Size
	size := count * element
	if size > maxSize {
		panic(fmt.Sprintf("buffer is too large for the specified usage. size: %d, max: %d", size, maxSize))
	}

	args.Size = size
	buffer := New(device, args)

	return &array[K]{
		T:       buffer,
		element: element,
		count:   count,
	}
}

func (a *array[K]) Set(index int, data K) {
	ptr := &data
	offset := index * a.element
	a.T.Write(offset, ptr)
}

func (a *array[K]) SetRange(offset int, data []K) {
	for i, el := range data {
		a.Set(offset+i, el)
	}
}

func (a *array[K]) Count() int   { return a.count }
func (a *array[K]) Element() int { return a.element }
