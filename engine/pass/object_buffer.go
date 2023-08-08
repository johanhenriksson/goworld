package pass

import (
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/descriptor"
)

type ObjectBuffer struct {
	buffer []uniform.Object
}

func NewObjectBuffer(capacity int) *ObjectBuffer {
	return &ObjectBuffer{
		buffer: make([]uniform.Object, 0, capacity),
	}
}

func (b *ObjectBuffer) Flush(desc *descriptor.Storage[uniform.Object]) {
	desc.SetRange(0, b.buffer)
}

func (b *ObjectBuffer) Reset() {
	b.buffer = b.buffer[:0]
}

func (b *ObjectBuffer) Store(light uniform.Object) int {
	index := len(b.buffer)
	b.buffer = append(b.buffer, light)
	return index
}
