package cache

import (
	"sync"
	"time"

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
	MaxAge() time.Duration

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

type line struct {
	value     any
	age       time.Duration
	version   int
	available bool
	wait      chan struct{}
	permanent bool
}

type cache[K Key, V Value] struct {
	backend Backend[K, V]
	data    map[string]*line
	lock    *sync.RWMutex
	async   bool

	maxAge   time.Duration
	lastTick time.Time
}

func New[K Key, V Value](backend Backend[K, V]) T[K, V] {
	c := &cache[K, V]{
		backend: backend,
		data:    map[string]*line{},
		lock:    &sync.RWMutex{},

		maxAge:   10 * time.Second,
		lastTick: time.Now(),

		// async caches should be a setting that can be disabled for testing purposes.
		// with async enabled, its difficult to render a single frame deterministically
		async: false,
	}
	return c
}

func (c cache[K, V]) MaxAge() time.Duration { return c.maxAge }

func (c *cache[K, V]) get(key K) (*line, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ln, hit := c.data[key.Key()]
	return ln, hit
}

func (c *cache[K, V]) init(key K) *line {
	ln := &line{
		available: false,
		wait:      make(chan struct{}),
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key.Key()] = ln
	return ln
}

func (c *cache[K, V]) fetch(key K) *line {
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
				c.deleteLater(ln.value.(V))
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
	var empty V
	if !ln.available {
		return empty, false
	}

	return ln.value.(V), true
}

func (c *cache[K, V]) Fetch(key K) V {
	ln := c.fetch(key)

	// not available yet - wait for it.
	if !ln.available {
		<-ln.wait
	}

	return ln.value.(V)
}

func (c *cache[K, V]) deleteLater(value V) {
	// we can reuse the eviction mechanic to delete values later
	// simply attach it to a cache line with a random key that will never be accessed
	c.lock.Lock()
	defer c.lock.Unlock()
	randomKey := "trash-" + util.NewUUID(8)
	c.data[randomKey] = &line{
		value:     value,
		available: true,
		age:       c.maxAge - time.Second, // delete in 1 second
	}
}

func (c *cache[K, V]) Tick() {
	// eviction
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	delta := now.Sub(c.lastTick)
	c.lastTick = now

	for key, line := range c.data {
		if line.permanent {
			continue
		}
		line.age += delta
		if line.age > c.maxAge {
			delete(c.data, key)

			// delete any instantiated object
			if line.available {
				c.backend.Delete(line.value.(V))
			}
		}
	}
}

func (c *cache[K, V]) Destroy() {
	c.lock.Lock()
	defer c.lock.Unlock()

	// destroy all the data in the cache
	for _, line := range c.data {
		// if the cache line is pending creation, we must wait for it to complete
		// before destroying it. failing to do so may cause a segfault
		if !line.available {
			<-line.wait
		}
		c.backend.Delete(line.value.(V))
	}

	c.backend.Destroy()
}
