package cache

import (
	"log"
)

type T[K Key, V Value] interface {
	Fetch(K) V
	Delete(V)
	Tick(int)
	Destroy()
}

type Key interface {
	Id() string
	Version() int
}

type Value interface{}

type Backend[K Key, V Value] interface {
	Instantiate(K, func(V))
	Delete(V)
	Destroy()
}

type line[V Value] struct {
	value   V
	age     int
	version int
	ready   bool
}

type job[K Key, V Value] struct {
	key  K
	done func(V)
}

type cache[K Key, V Value] struct {
	backend Backend[K, V]
	data    map[string]*line[V]
	remove  map[int][]V
	work    chan job[K, V]
	frame   int
	buffer  int
}

func New[K Key, V Value](backend Backend[K, V]) T[K, V] {
	c := &cache[K, V]{
		backend: backend,
		data:    map[string]*line[V]{},
		remove:  map[int][]V{},
		buffer:  100,
	}
	go c.worker()
	return c
}

func (c *cache[K, V]) worker() {
	c.work = make(chan job[K, V], c.buffer)
	for {
		job, more := <-c.work
		if !more {
			break
		}
		c.backend.Instantiate(job.key, job.done)
	}
	log.Println(c.backend, "exited")
}

func (c *cache[K, V]) Fetch(key K) V {
	var empty V

	ln, hit := c.data[key.Id()]
	if !hit {
		ln = &line[V]{
			ready: false,
		}
		c.data[key.Id()] = ln
	}

	if ln.version != key.Version() {
		ln.version = key.Version()
		c.work <- job[K, V]{
			key: key,
			done: func(value V) {
				if ln.ready {
					// ready implies that we have a previous value - schedule deletion
					c.Delete(ln.value)
				}

				ln.value = value
				ln.ready = true
			},
		}
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

// Should be called immediately after aquiring a new frame, passing the index of the aquired frame.
// Releases any unused resources associated with that frame index.
func (c *cache[K, V]) Tick(frameIndex int) {
	c.frame = frameIndex

	// eviction
	for key, line := range c.data {
		line.age++
		if line.age > 100 {
			c.backend.Delete(line.value)
			delete(c.data, key)
		}
	}
}

func (c *cache[K, V]) Destroy() {
	close(c.work)
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
