package buffer

import (
	"sort"
	"sync"
)

type blockAlloc struct {
	buffer   T
	freelist []Block
	mutex    sync.Mutex
}

func NewBlockAllocator(buf T) *blockAlloc {
	if !isPowerOfTwo(uint(buf.Size())) {
		panic("buffer size must be a power of two")
	}
	m := &blockAlloc{
		buffer:   buf,
		freelist: make([]Block, 0, 128),
	}
	m.freelist = append(m.freelist, Block{
		buffer: buf,
		size:   buf.Size(),
		offset: 0,
	})
	return m
}

func (m *blockAlloc) Alloc(size int) (Block, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// closest power of two
	power := 64
	for power < size {
		power *= 2
	}
	size = power

	// find the smallest block that fits
	smallest := -1
	for i, block := range m.freelist {
		if block.size < size {
			continue
		}
		if smallest == -1 || block.size < m.freelist[smallest].size {
			smallest = i
		}
	}

	if smallest < 0 {
		return Block{}, ErrOutOfMemory
	}

	// take the smallest block
	block := m.freelist[smallest]
	m.freelist = append(m.freelist[:smallest], m.freelist[smallest+1:]...)

	// split the block until we have the smallest possible allocation
	for block.size >= size*2 {
		block.size /= 2
		m.freelist = append(m.freelist, Block{
			buffer: m.buffer,
			offset: block.offset + block.size,
			size:   block.size,
		})
	}

	return block, nil
}

func (m *blockAlloc) Free(b Block) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if b.buffer != m.buffer {
		return ErrInvalidFree
	}

	m.freelist = append(m.freelist, b)

	// sort by offset
	sort.Slice(m.freelist, func(i, j int) bool {
		return m.freelist[i].offset < m.freelist[j].offset
	})

	// defragmentation
	for {
		merged := false
		for i := 0; i < len(m.freelist)-1; i++ {
			if m.freelist[i].size != m.freelist[i+1].size {
				// only equal size blocks can be merged
				continue
			}
			if m.freelist[i].offset+m.freelist[i].size == m.freelist[i+1].offset {
				m.freelist[i].size += m.freelist[i+1].size
				m.freelist = append(m.freelist[:i+1], m.freelist[i+2:]...)
				merged = true
				break
			}
		}
		if !merged {
			break
		}
	}

	return nil
}
