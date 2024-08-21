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

type ValueProp interface {
	GenericProp
	setValue(any)
}

type Property[T PropValue] struct {
	value T
	def   T
	kind  reflect.Type

	OnChange events.Event[T]
}

var _ GenericProp = &Property[int]{}
var _ ValueProp = &Property[int]{}

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

func (p *Property[T]) setValue(value any) {
	p.value = value.(T)
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

type valueProp struct {
	Value any
}

func encodeProperties(enc Encoder, val reflect.Value) error {
	// property
	// .Addr() since property methods have pointer receivers
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Addr().Type().Implements(valuePropType) {
			var value any
			if prop, ok := val.Field(i).Addr().Interface().(GenericProp); ok {
				value = prop.GetAny()
			}
			if err := enc.Encode(valueProp{
				Value: value,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func decodeProperties(pool Pool, dec Decoder, val reflect.Value) error {
	// property
	// .Addr() since property methods have pointer receivers
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Addr().Type().Implements(valuePropType) {
			var prop valueProp
			if err := dec.Decode(&prop); err != nil {
				return err
			}
			val.Field(i).Addr().Interface().(ValueProp).setValue(prop.Value)
		}
	}
	return nil
}
