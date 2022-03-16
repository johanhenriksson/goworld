package buffer

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"
)

type Item[K any] interface {
	T

	// Set the data in the buffer
	Set(data K)
}

type item[K any] struct {
	T
}

// NewItem creates a new typed single-item buffer.
// When allocating items, the Size argument is ignored
func NewItem[K any](device device.T, args Args) Item[K] {
	align, maxSize := GetBufferLimits(device, args.Usage)

	var empty K
	kind := reflect.TypeOf(empty)

	element := util.Align(int(kind.Size()), align)
	if element > maxSize {
		panic(fmt.Sprintf("buffer is too large for the specified usage. size: %d, max: %d", element, maxSize))
	}

	args.Size = element
	buffer := New(device, args)

	return &item[K]{
		T: buffer,
	}
}

func (i *item[K]) Set(data K) {
	ptr := &data
	i.Write(0, ptr)
}
