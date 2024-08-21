package object

import (
	"encoding/gob"
	"log"
	"reflect"
)

type Handle uint

func init() {
	gob.Register(Handle(0))
}

type ReferenceProp interface {
	GenericProp
	handle() Handle
	setHandle(Pool, Handle)
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

// marks the type as a reference property
func (r *Reference[T]) handle() Handle {
	return r.Property.value
}

func (r *Reference[T]) setHandle(pool Pool, h Handle) {
	r.Property.Set(h)
	r.pool = pool
}

func (r *Reference[T]) setValue() {} // ensure Reference does not implement ValueProp

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

//
// serialization
//

type referenceProp struct {
	Handle Handle
}

func encodeReferences(enc Encoder, val reflect.Value) error {
	// reference to component
	// .Addr() since property methods have pointer receivers
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).Addr().Type().Implements(refPropType) {
			continue
		}

		log.Println("encoding reference", val.Type().Field(i).Name)
		var handle Handle
		if ref, ok := val.Field(i).Addr().Interface().(ReferenceProp); ok {
			handle = ref.handle()
		}

		if err := enc.Encode(referenceProp{
			Handle: handle,
		}); err != nil {
			return err
		}
	}
	return nil
}

func decodeReferences(pool Pool, dec Decoder, val reflect.Value) error {
	// reference to component
	// .Addr() since property methods have pointer receivers
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).Addr().Type().Implements(refPropType) {
			continue
		}

		log.Println("decoding reference", val.Type().Field(i).Name)
		var ref referenceProp
		if err := dec.Decode(&ref); err != nil {
			return err
		}

		newHandle := pool.remap(ref.Handle)
		log.Println("remapped handle", ref.Handle, "->", newHandle)
		val.Field(i).Addr().Interface().(ReferenceProp).setHandle(pool, newHandle)
	}
	return nil
}
