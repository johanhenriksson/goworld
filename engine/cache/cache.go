package cache

import (
	"log"
)

type T[I Item, O any] interface {
	// Fetches the given mesh & material combination from the cache
	// If it does not exist, it will be inserted.
	Fetch(I) O

	// Evict the first cached mesh that is older than max age
	// Returns true if a mesh is evicted
	Evict() bool

	// Tick increments the age of all items in the cache
	Tick()

	Destroy()
}

type Item interface {
	Id() string
	Version() int
}

type Cachable interface {
}

type Backend[I Item, O any] interface {
	Instantiate(I) O
	Update(O, I)
	Delete(O)
	Destroy()
}

type cache[I Item, O any] struct {
	maxAge  int
	cache   map[string]*cache_line[O]
	backend Backend[I, O]
}

type cache_line[O any] struct {
	age     int
	version int
	item    O
}

func New[I Item, O any](backend Backend[I, O]) T[I, O] {
	return &cache[I, O]{
		maxAge:  1000,
		cache:   make(map[string]*cache_line[O]),
		backend: backend,
	}
}

func (m *cache[I, O]) Fetch(item I) O {
	line, hit := m.cache[item.Id()]

	// not in cache - instantiate a buffered mesh
	if !hit {
		log.Println("buffering new mesh", item.Id(), "version", item.Version())
		vao := m.backend.Instantiate(item)
		line = &cache_line[O]{
			item: vao,
		}
		m.cache[item.Id()] = line
	}

	// version has changed, update the mesh
	if line.version != item.Version() {
		// we might want to queue this operation and run it at a more appropriate time
		log.Println("updating existing mesh", item.Id(), "to version", item.Version())
		m.backend.Update(line.item, item)
		line.version = item.Version()
	}

	// reset age
	line.age = 0

	return line.item
}

func (m *cache[I, O]) Tick() {
	// increment the age of every item in the cache
	for _, entry := range m.cache {
		entry.age++
	}
}

func (m *cache[I, O]) Evict() bool {
	for id, entry := range m.cache {
		// skip any meshes that have been recently used
		if entry.age < m.maxAge {
			continue
		}

		// deallocate gpu memory
		log.Println("deallocating", id, "from gpu")
		m.backend.Delete(m.cache[id].item)

		// remove cache line
		delete(m.cache, id)
		return true
	}
	return false
}

func (m *cache[I, O]) Destroy() {
	for _, line := range m.cache {
		m.backend.Delete(line.item)
	}
	m.backend.Destroy()
	m.cache = make(map[string]*cache_line[O])
}
