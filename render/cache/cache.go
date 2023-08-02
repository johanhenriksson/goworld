package cache

import (
	"sync"

	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/util"
)

type T[K Key, V Value] interface {
	// TryFetch returns a value if it exists and is ready to use.
	// Resets the age of the cache line
	// Returns a bool indicating whether the value exists.
	TryFetch(K) (V, bool)

	// Fetch returns a value, waiting until its becomes available if it does not yet exist
	// Resets the age of the cache line
	Fetch(K) V

	// MaxAge returns the number of ticks until unused lines are evicted
	MaxAge() int

	// Tick increments the age of all cache lines, and evicts those
	// that have not been accessed in maxAge ticks or more.
	Tick()

	// Destroy the cache and all data held in it.
	Destroy()
}

type Key interface {
	Key() string
	Version() int
}

type Value interface{}

type Backend[K Key, V Value] interface {
	// Instantiate the resource referred to by Key.
	// Must execute on a background goroutine
	Instantiate(K, func(V))

	Delete(V)
	Destroy()
	Name() string
}

type line[V Value] struct {
	value     V
	age       int
	version   int
	available bool
	wait      chan struct{}
}

type cache[K Key, V Value] struct {
	backend Backend[K, V]
	data    map[string]*line[V]
	worker  *command.ThreadWorker
	lock    *sync.RWMutex
	maxAge  int
	async   bool
}

func New[K Key, V Value](backend Backend[K, V]) T[K, V] {
	c := &cache[K, V]{
		backend: backend,
		data:    map[string]*line[V]{},
		worker:  command.NewThreadWorker(backend.Name(), 100, false),
		lock:    &sync.RWMutex{},
		maxAge:  100,
		async:   false,
	}
	return c
}

func (c cache[K, V]) MaxAge() int { return c.maxAge }

func (c *cache[K, V]) get(key K) (*line[V], bool) {
	c.lock.RLock()
	ln, hit := c.data[key.Key()]
	c.lock.RUnlock()
	return ln, hit
}

func (c *cache[K, V]) init(key K) *line[V] {
	ln := &line[V]{
		available: false,
		wait:      make(chan struct{}),
	}
	c.lock.Lock()
	c.data[key.Key()] = ln
	c.lock.Unlock()
	return ln
}

func (c *cache[K, V]) fetch(key K) *line[V] {
	ln, hit := c.get(key)
	if !hit {
		ln = c.init(key)
	}

	// check if a newer version has been requested
	// since the initial line has version 0, this always happens on the first request.
	if ln.version != key.Version() {
		// update version immediately, so that duplicate instantiantions wont happen.
		// note that the previous version will be returned until the new one is available
		ln.version = key.Version()

		// instantiate new version
		c.backend.Instantiate(key, func(value V) {
			if ln.available {
				// available implies that we have a previous value
				// however, it is most likely in use rendering the in-flight frame!
				// deleting it here may cause a segfault
				c.deleteLater(ln.value)
			}
			ln.value = value

			// if its the very first time this item is requested, signal any waiting
			// synchronous fetch that it is ready.
			if !ln.available {
				ln.available = true
				close(ln.wait)
			}
		})
	}

	// reset age
	ln.age = 0

	return ln
}

func (c *cache[K, V]) TryFetch(key K) (V, bool) {
	if !c.async {
		return c.Fetch(key), true
	}

	ln := c.fetch(key)

	// not available yet - return nothing
	if !ln.available {
		var empty V
		return empty, false
	}

	return ln.value, true
}

func (c *cache[K, V]) Fetch(key K) V {
	ln := c.fetch(key)

	// not available yet - wait for it.
	if !ln.available {
		<-ln.wait
	}

	return ln.value
}

func (c *cache[K, V]) deleteLater(value V) {
	// we can reuse the eviction mechanic to delete values later
	// simply attach it to a cache line with a random key that will never be accessed
	c.lock.Lock()
	defer c.lock.Unlock()
	randomKey := "trash-" + util.NewUUID(8)
	c.data[randomKey] = &line[V]{
		value:     value,
		available: true,
		age:       c.maxAge - 10, // delete in 10 frames
	}
}

func (c *cache[K, V]) Tick() {
	// eviction
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, line := range c.data {
		line.age++
		if line.age > c.maxAge {
			delete(c.data, key)

			// delete any instantiated object
			if line.available {
				c.backend.Delete(line.value)
			}
		}
	}
}

func (c *cache[K, V]) Destroy() {
	c.lock.Lock()
	defer c.lock.Unlock()

	// flush any pending work
	c.worker.Flush()

	// destroy all the data in the cache
	for _, line := range c.data {
		// if the cache line is pending creation, we must wait for it to complete
		// before destroying it. failing to do so may cause a segfault
		if !line.available {
			<-line.wait
		}
		c.backend.Delete(line.value)
	}

	c.backend.Destroy()
}
