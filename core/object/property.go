package object

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/core/events"
)

type PropValue interface {
	// ~int | ~uint | ~float32 | string | bool
}

type GenericProp interface {
	Type() reflect.Type
	GetAny() any
	SetAny(any)
}

type Property[T PropValue] struct {
	value T
	def   T
	kind  reflect.Type

	// OnChange executes callbacks every time the property value is changed.
	// The callback is called with the new value.
	// Only the propertys owner object should subscribe to this event.
	OnChange events.Event[T]
}

var _ GenericProp = &Property[int]{}

func NewProperty[T PropValue](def T) Property[T] {
	var empty T
	return Property[T]{
		value: def,
		def:   def,
		kind:  reflect.TypeOf(empty),
	}
}

func (p *Property[T]) Get() T {
	return p.value
}

func (p *Property[T]) GetAny() any {
	return p.value
}

func (p *Property[T]) Set(value T) {
	p.value = value
	p.OnChange.Emit(value)
}

func (p *Property[T]) SetAny(value any) {
	if cast, ok := value.(T); ok {
		p.Set(cast)
	}
}

func (p *Property[T]) String() string {
	return fmt.Sprintf("%v", p.value)
}

func (p *Property[T]) Type() reflect.Type {
	return p.kind
}

type PropInfo struct {
	GenericProp
	Key  string
	Name string
}

func Properties(target Component) []PropInfo {
	t := reflect.TypeOf(target).Elem()
	v := reflect.ValueOf(target).Elem()

	properties := make([]PropInfo, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			// anonymous fields are not considered
			continue
		}
		if !field.IsExported() {
			// only exported fields can be properties
			continue
		}

		value := v.Field(i)
		if prop, isProp := value.Addr().Interface().(GenericProp); isProp {
			// todo: tags

			properties = append(properties, PropInfo{
				GenericProp: prop,

				Key:  field.Name,
				Name: field.Name,
			})
		}
	}

	return properties
}

//
// serialization
//

type EncodedProp interface {
	Encode() ([]byte, error)
	Decode([]byte) (PropValue, error)
}

var encodedPropType = reflect.TypeOf((EncodedProp)(nil))

func (p *Property[T]) Serialize(enc Encoder) error {
	if reflect.TypeOf(p.value).Implements(encodedPropType) {
		// use the custom serialization
		encoder := reflect.ValueOf(p.value).Interface().(EncodedProp)
		bytes, err := encoder.Encode()
		if err != nil {
			return err
		}
		return enc.Encode(bytes)
	} else {
		// use the default serialization
		return enc.Encode(p.value)
	}
}

func (p *Property[T]) Deserialize(pool Pool, dec Decoder) error {
	if reflect.TypeOf(p.value).Implements(encodedPropType) {
		// use the custom serialization
		decoder := reflect.ValueOf(p.value).Interface().(EncodedProp)
		var bytes []byte
		value, err := decoder.Decode(bytes)
		if err != nil {
			return err
		}
		p.value = value.(T)
		return nil
	} else {
		// use the default serialization
		var value T
		if err := dec.Decode(&value); err != nil {
			return err
		}
		p.value = value
		return nil
	}
}
