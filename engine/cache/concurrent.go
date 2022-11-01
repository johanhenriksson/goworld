package cache

type concurrent[K Key, V any] struct {
	backend Backend[K, V]
	cache   map[string]*cacheLine[V]
	work    []cacheWork[K, V]
	maxAge  int
}

type cacheOp int

const createOp cacheOp = 1
const updateOp cacheOp = 2

type cacheWork[K Key, V any] struct {
	Op    cacheOp
	Key   K
	Value V
}

func NewConcurrent[K Key, V any](backend Backend[K, V]) T[K, V] {
	return &concurrent[K, V]{
		backend: backend,
		cache:   make(map[string]*cacheLine[V], 64),
		work:    make([]cacheWork[K, V], 0, 32),
		maxAge:  1000,
	}
}

func (c *concurrent[K, V]) Fetch(key K) V {
	line, hit := c.cache[key.Id()]

	if !hit {
		// its not in the cache! queue instantiation work
		// todo: handle multiple fetches of the same item
		c.work = append(c.work, cacheWork[K, V]{
			Op:  createOp,
			Key: key,
		})

		// return an empty object until
		var empty V
		return empty
	}

	if key.Version() > line.version {
		// the cached version is outdated - queue update work
		// todo: handle updates fetches of the same item
		c.work = append(c.work, cacheWork[K, V]{
			Op:  updateOp,
			Key: key,
		})
	}

	// reset access counter
	line.age = 0

	return line.value
}

func (c *concurrent[K, V]) Tick() {
	// increment age of all cache objects
	for _, line := range c.cache {
		line.age++
	}

	// run evictions
	// todo: this causes problems, since not every object is requested on each frame
	// c.evict()

	// perform work
	c.process()
}

func (c *concurrent[K, V]) process() {
	for _, work := range c.work {
		switch work.Op {
		case createOp:
			c.process_create(work.Key)
		case updateOp:
			c.process_update(work.Key)
		default:
			panic("invalid cache work operation")
		}
	}
	c.work = c.work[:0]
}

func (c *concurrent[K, V]) process_create(key K) {
	value := c.backend.Instantiate(key)
	c.cache[key.Id()] = &cacheLine[V]{
		version: key.Version(),
		value:   value,
	}
}

func (c *concurrent[K, V]) process_update(key K) {
	line := c.cache[key.Id()]
	c.backend.Update(line.value, key)
	line.version = key.Version()
}

func (c *concurrent[K, V]) evict() bool {
	for id, entry := range c.cache {
		// skip any meshes that have been recently used
		if entry.age < c.maxAge {
			continue
		}

		// destroy object
		c.backend.Delete(c.cache[id].value)

		// remove cache line
		delete(c.cache, id)
		return true
	}
	return false
}

func (c *concurrent[K, V]) Destroy() {
	for _, line := range c.cache {
		c.backend.Delete(line.value)
	}
	c.backend.Destroy()
}
