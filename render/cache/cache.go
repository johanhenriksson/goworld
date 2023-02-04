package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render/command"
)

type T[K Key, V Value] interface {
	Submit()
	Fetch(K) V
	Delete(V)
	Tick(int)
	Destroy()
}

type Key interface {
	Key() string
	Version() int
}

type Value interface{}

type Backend[K Key, V Value] interface {
	Submit()
	Instantiate(K, func(V))
	Delete(V)
	Destroy()
	Name() string
}

type line[V Value] struct {
	value   V
	age     int
	version int
	ready   bool
}

type cache[K Key, V Value] struct {
	backend Backend[K, V]
	data    map[string]*line[V]
	remove  map[int][]V
	frame   int
	worker  *command.ThreadWorker
}

func New[K Key, V Value](backend Backend[K, V]) T[K, V] {
	c := &cache[K, V]{
		backend: backend,
		data:    map[string]*line[V]{},
		remove:  map[int][]V{},
		worker:  command.NewThreadWorker(backend.Name(), 100, false),
	}
	return c
}

func (c *cache[K, V]) Fetch(key K) V {
	var empty V

	ln, hit := c.data[key.Key()]
	if !hit {
		ln = &line[V]{
			ready: false,
		}
		c.data[key.Key()] = ln
	}

	if ln.version != key.Version() {
		ln.version = key.Version()
		c.worker.Invoke(func() {
			c.backend.Instantiate(key, func(value V) {
				if ln.ready {
					// ready implies that we have a previous value - schedule deletion
					c.Delete(ln.value)
				}
				ln.value = value
				ln.ready = true
			})
		})
	}

	// reset age
	ln.age = 0

	// not available yet
	if !ln.ready {
		return empty
	}

	return ln.value
}

func (c *cache[K, V]) Delete(value V) {
	c.remove[c.frame] = append(c.remove[c.frame], value)
}

// Submit pending work
func (c *cache[K, V]) Submit() {
	c.backend.Submit()
}

// Should be called immediately after aquiring a new frame, passing the index of the aquired frame.
// Releases any unused resources associated with that frame index.
func (c *cache[K, V]) Tick(frameIndex int) {
	c.frame = frameIndex

	// eviction
	for key, line := range c.data {
		line.age++
		if line.age > 100 {
			log.Println(c.backend, "evict", line.value)
			c.Delete(line.value)
			delete(c.data, key)
		}
	}

	if len(c.remove) > 1 {
		for _, value := range c.remove[c.frame] {
			c.backend.Delete(value)
		}
		c.remove[c.frame] = nil
	}
}

func (c *cache[K, V]) Destroy() {
	for _, queue := range c.remove {
		for _, value := range queue {
			c.backend.Delete(value)
		}
	}
	for _, line := range c.data {
		if line.ready {
			c.backend.Delete(line.value)
		}
	}
	c.backend.Destroy()
}
