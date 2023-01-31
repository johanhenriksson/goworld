package cache

import (
	"log"
)

type T[K Key, V any] interface {
	// Fetches the given mesh & material combination from the cache
	// If it does not exist, it will be inserted.
	Fetch(K) V

	// Tick increments the age of all items in the cache
	Tick()

	Destroy()
}

type Key interface {
	Id() string
	Version() int
}

type Backend[K Key, V any] interface {
	Name() string
	Instantiate(K) V
	Update(V, K) V
	Delete(V)
	Destroy()
}

type cache[K Key, V any] struct {
	maxAge  int
	cache   map[string]*cacheLine[V]
	backend Backend[K, V]
}

type cacheLine[V any] struct {
	age     int
	version int
	value   V
}

func New[K Key, V any](backend Backend[K, V]) T[K, V] {
	return &cache[K, V]{
		maxAge:  1000,
		cache:   make(map[string]*cacheLine[V]),
		backend: backend,
	}
}

func (m *cache[K, V]) Fetch(key K) V {
	line, hit := m.cache[key.Id()]

	// not in cache - instantiate a buffered mesh
	if !hit {
		log.Println("instantiate new", m.backend.Name(), key.Id(), "version", key.Version())
		value := m.backend.Instantiate(key)
		line = &cacheLine[V]{
			value: value,
		}
		m.cache[key.Id()] = line
	}

	// version has changed, update the mesh
	if line.version != key.Version() {
		// we might want to queue this operation and run it at a more appropriate time
		line.value = m.backend.Update(line.value, key)
		line.version = key.Version()
	}

	// reset age
	line.age = 0

	return line.value
}

func (m *cache[K, V]) Tick() {
	// increment the age of every item in the cache
	for _, entry := range m.cache {
		entry.age++
	}
	// evict items
	// todo: this causes problems, since not every object is requested on each frame
	// m.evict()
}

func (m *cache[K, V]) evict() bool {
	for id, entry := range m.cache {
		// skip any meshes that have been recently used
		if entry.age < m.maxAge {
			continue
		}

		// deallocate gpu memory
		log.Println("deallocating", id, "from gpu")
		m.backend.Delete(m.cache[id].value)

		// remove cache line
		delete(m.cache, id)
		return true
	}
	return false
}

func (m *cache[K, V]) Destroy() {
	for _, line := range m.cache {
		m.backend.Delete(line.value)
	}
	m.backend.Destroy()
	m.cache = nil
}
