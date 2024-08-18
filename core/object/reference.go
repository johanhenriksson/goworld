package object

import (
	"encoding/gob"
)

type Handle uint

func init() {
	gob.Register(Handle(0))
}

type Reference[T Component] struct {
	Property[Handle]
	pool Pool
}

var _ GenericProp = &Reference[Component]{}

func NewReference[T Component](cmp T) Reference[T] {
	return Reference[T]{
		Property: NewProperty(cmp.ID()),
		pool:     cmp.context(),
	}
}

func EmptyReference[T Component]() Reference[T] {
	return Reference[T]{
		Property: NewProperty[Handle](0),
		pool:     nil,
	}
}

func (r *Reference[T]) Get() (T, bool) {
	var empty T
	if r.pool == nil {
		return empty, false
	}

	cmp, ok := r.pool.Resolve(r.Property.value)
	if !ok {
		return empty, false
	}

	cast, ok := cmp.(T)
	if !ok {
		return empty, false
	}

	return cast, true
}

func (r *Reference[T]) Set(cmp T) {
	r.Property.Set(cmp.ID())
	r.pool = cmp.context()
}

func (r *Reference[T]) Serialize(enc Encoder) error {
	return enc.Encode(r.Property.value)
}

func DeserializeReference[T Component](pool Pool, dec Decoder) (Reference[T], error) {
	handle := Handle(0)
	if err := dec.Decode(&handle); err != nil {
		return Reference[T]{}, err
	}

	newHandle := pool.remap(handle)

	return Reference[T]{
		Property: NewProperty(newHandle),
		pool:     pool,
	}, nil
}
