package object

import (
	"encoding/gob"
)

type Handle uint

func init() {
	gob.Register(Handle(0))
}

type Ref[T Component] struct {
	Property[Handle]
	pool Pool
}

var _ GenericProp = &Ref[Component]{}

func NewRef[T Component](cmp T) Ref[T] {
	return Ref[T]{
		Property: NewProperty(cmp.ID()),
		pool:     cmp.Pool().unwrap(),
	}
}

func EmptyRef[T Component]() Ref[T] {
	return Ref[T]{
		Property: NewProperty[Handle](0),
		pool:     nil,
	}
}

func (r *Ref[T]) Get() (T, bool) {
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

func (r *Ref[T]) Set(cmp T) {
	r.Property.Set(cmp.ID())
	r.pool = cmp.Pool().unwrap()
}

func (r *Ref[T]) Serialize(enc Encoder) error {
	return enc.Encode(r.Property.value)
}

//
// serialization
//

func (r *Ref[T]) Deserialize(pool Pool, dec Decoder) error {
	if err := r.Property.Deserialize(pool, dec); err != nil {
		return err
	}
	r.value = pool.remap(r.value)
	r.pool = pool.unwrap()
	return nil
}
