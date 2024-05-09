package srv

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")
var ErrType = errors.New("type mismatch")

type Pool[T any] struct {
	lock   sync.RWMutex
	data   []T
	gen    []uint32
	typeId int
	nextId int
	size   int
}

func NewPool[T any](typeId int, sizeBits int) *Pool[T] {
	size := 1 << sizeBits

	return &Pool[T]{
		typeId: typeId,
		size:   size,
		data:   make([]T, size),
		gen:    make([]uint32, size),
		nextId: 1,
	}
}

func (p *Pool[T]) Add(e T) Identity {
	// p.lock.Lock()
	// defer p.lock.Unlock()

	i := p.nextId
	if i >= p.size {
		// todo: implement a free list
		panic("pool is full")
	}
	p.nextId++

	gen := (p.gen[i]&0xFFFFFF + 1) | 0xFF000000
	p.gen[i] = gen
	p.data[i] = e

	return Identity(uint64(p.typeId)<<56 | uint64(gen)<<32 | uint64(i))
}

func (p *Pool[T]) Remove(id Identity) error {
	if id.TypeID() != p.typeId {
		return ErrType
	}

	index := id.Index()
	if index >= p.size {
		panic("index out of bounds")
	}

	// p.lock.Lock()
	// defer p.lock.Unlock()

	if id.Generation() != int(p.gen[index]) {
		return ErrNotFound
	}

	// remove alive flag
	p.gen[index] = p.gen[index] & 0x00FFFFFF
	return nil
}

func (p *Pool[T]) Get(id Identity) (T, error) {
	var empty T
	if id == None {
		return empty, ErrNotFound
	}

	if id.TypeID() != p.typeId {
		// wrong type
		return empty, ErrType
	}

	index := id.Index()
	if index >= p.size {
		panic("index out of bounds")
	}

	// p.lock.RLock()
	// defer p.lock.RUnlock()

	g := p.gen[index]
	alive := g&0xFF000000 > 0
	gen := int(g & 0x00FFFFFF)
	if !alive || id.Generation() != gen {
		// stale
		return empty, ErrNotFound
	}

	return p.data[index], nil
}

func (p *Pool[T]) Each(f func(T) bool) {
	// p.lock.Lock()
	// defer p.lock.Unlock()

	end := min(p.nextId, p.size)
	for i := 0; i < end; i++ {
		g := p.gen[i]
		dead := g&0xFF000000 == 0
		if dead {
			continue
		}

		if !f(p.data[i]) {
			break
		}
	}
}

func (p *Pool[T]) Filter(f func(T) bool) []T {
	out := make([]T, 0, 16)
	p.Each(func(e T) bool {
		if f(e) {
			out = append(out, e)
		}
		return true
	})
	return out
}
