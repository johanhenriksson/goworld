package uniform

import (
	"github.com/johanhenriksson/goworld/render/descriptor"
)

type ObjectBuffer struct {
	buffer []Object
}

func NewObjectBuffer(capacity int) *ObjectBuffer {
	return &ObjectBuffer{
		buffer: make([]Object, 0, capacity),
	}
}

func (b *ObjectBuffer) Size() int {
	return cap(b.buffer)
}

func (b *ObjectBuffer) Flush(desc *descriptor.Storage[Object]) {
	desc.SetRange(0, b.buffer)
}

func (b *ObjectBuffer) Reset() {
	b.buffer = b.buffer[:0]
}

func (b *ObjectBuffer) Store(obj Object) int {
	index := len(b.buffer)
	b.buffer = append(b.buffer, obj)
	return index
}
