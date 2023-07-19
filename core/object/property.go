package object

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/core/events"
)

type PropType interface{}

type Property[T PropType] struct {
	value   T
	def     T
	kind    reflect.Type
	changed *events.Event[T]
}

func NewProperty[T PropType](def T) *Property[T] {
	var empty T
	return &Property[T]{
		value:   def,
		def:     def,
		kind:    reflect.TypeOf(empty),
		changed: events.New[T](),
	}
}

func (p *Property[T]) Get() T {
	return p.value
}

func (p *Property[T]) Set(value T) {
	p.value = value
	p.changed.Emit(value)
}

func (p *Property[T]) String() string {
	return fmt.Sprintf("%v", p.value)
}

func (p *Property[T]) Type() reflect.Type {
	return p.kind
}

func (p *Property[T]) OnChange() *events.Event[T] {
	return p.changed
}
