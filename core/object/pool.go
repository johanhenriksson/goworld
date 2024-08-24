package object

import "log"

var GlobalPool = NewPool()

type Pool interface {
	Resolve(Handle) (Component, bool)

	assign(Component)
	remap(Handle) Handle
	release(Handle)
	unwrap() Pool
}

type context struct {
	handles map[Handle]Component
	next    Handle
}

func NewPool() Pool {
	return &context{
		handles: map[Handle]Component{},
		next:    Handle(1),
	}
}

func (c *context) Resolve(h Handle) (Component, bool) {
	obj, ok := c.handles[h]
	return obj, ok
}

func (c *context) assign(obj Component) {
	handle := c.remap(obj.ID())
	obj.setHandle(c, handle)
	c.handles[handle] = obj
}

func (c *context) release(h Handle) {
	delete(c.handles, h)
}

func (c *context) remap(h Handle) Handle {
	if h != 0 {
		return h
	}
	return c.nextHandle()
}

func (c *context) unwrap() Pool {
	return c
}

func (c *context) nextHandle() Handle {
	handle := c.next
	c.next++
	return handle
}

// mappingPool remaps existing handles to new handles in the
// underlying object context. This is useful when deserializing objects
// to avoid conflicts with existing objects.
type mappingPool struct {
	Pool
	mapping map[Handle]Handle
}

func newMappingPool(pool Pool) Pool {
	// prevent nesting
	if mpool, ok := pool.(*mappingPool); ok {
		return mpool
	}
	return &mappingPool{
		Pool:    pool,
		mapping: map[Handle]Handle{},
	}
}

func (c *mappingPool) assign(obj Component) {
	handle := obj.ID()
	if handle == 0 {
		// having handle == 0 means the object is new.
		// this should not happen when using a translating context.
		panic("translating context should only be used for deserialization")
	}
	newHandle := c.remap(handle)
	obj.setHandle(c, newHandle)
	c.Pool.assign(obj)
}

func (c *mappingPool) remap(h Handle) Handle {
	if newHandle, ok := c.mapping[h]; ok {
		return newHandle
	}
	newHandle := c.Pool.remap(0) // remap(0) returns a new handle
	c.mapping[h] = newHandle
	return newHandle
}

func (c *mappingPool) unwrap() Pool {
	return c.Pool
}
