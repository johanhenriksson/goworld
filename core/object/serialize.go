package object

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Decoder interface {
	Decode(e any) error
}

type Encoder interface {
	Encode(data any) error
}

type Serializable interface {
	Serialize(Encoder) error
}

type MemorySerializer struct {
	stream []any
	index  int
}

func (m *MemorySerializer) Encode(data any) error {
	m.stream = append(m.stream, data)
	return nil
}

func (m *MemorySerializer) Decode(target any) error {
	if m.index >= len(m.stream) {
		return io.EOF
	}
	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(m.stream[m.index]))
	m.index++
	return nil
}

func Copy[T Component](pool Pool, obj T) T {
	buffer := &MemorySerializer{}

	err := Serialize(buffer, obj)
	if err != nil {
		panic(err)
	}

	kopy, err := Deserialize[T](pool, buffer)
	if err != nil {
		panic(err)
	}

	return kopy
}

func Save(key string, obj Component) error {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := Serialize(enc, obj); err != nil {
		return err
	}
	return assets.Write(key, buf.Bytes())
}

func Load[T Component](pool Pool, key string) (T, error) {
	data, err := assets.Read(key)
	if err != nil {
		var empty T
		return empty, err
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return Deserialize[T](pool, dec)
}

type ComponentState struct {
	ID      Handle
	Name    string
	Enabled bool
}

func NewComponentState(c Component) ComponentState {
	return ComponentState{
		ID:      c.ID(),
		Name:    c.Name(),
		Enabled: c.Enabled(),
	}
}

func (c ComponentState) New() Component {
	return &component{
		id:      c.ID,
		name:    c.Name,
		enabled: c.Enabled,
	}
}

type ObjectState struct {
	ComponentState
	Position vec3.T
	Rotation quat.T
	Scale    vec3.T
	Children int
}

var ErrSerialize = errors.New("serialization error")

func Serialize(enc Encoder, obj Component) error {
	kind := typeName(obj)
	serializable, ok := obj.(Serializable)
	if !ok {
		return fmt.Errorf("%w: %s is not serializable", ErrSerialize, kind)
	}
	if _, registered := types[kind]; !registered {
		return fmt.Errorf("%w: no deserializer for %s", ErrSerialize, kind)
	}
	if err := enc.Encode(kind); err != nil {
		return err
	}
	return serializable.Serialize(enc)
}

var kind string

func Deserialize[T Component](pool Pool, decoder Decoder) (T, error) {
	pool = newMappingPool(pool)
	var empty T
	if err := decoder.Decode(&kind); err != nil {
		return empty, err
	}
	typ, exists := types[kind]
	if !exists {
		return empty, fmt.Errorf("%w: no deserializer for %s", ErrSerialize, kind)
	}

	tctx := newMappingPool(pool)
	cmp, err := typ.Deserialize(tctx, decoder)
	if err != nil {
		return empty, err
	}

	cast, ok := cmp.(T)
	if !ok {
		return empty, fmt.Errorf("%w: %s is not of type %T", ErrSerialize, kind, empty)
	}
	return cast, nil
}

func (o *object) Serialize(enc Encoder) error {
	children := 0
	for _, child := range o.children {
		if _, ok := child.(Serializable); ok {
			kind := typeName(child)
			if _, registered := types[kind]; !registered {
				continue
			}
			children++
		}
	}

	if err := enc.Encode(ObjectState{
		ComponentState: NewComponentState(o),
		Position:       o.transform.Position(),
		Rotation:       o.transform.Rotation(),
		Scale:          o.transform.Scale(),
		Children:       children,
	}); err != nil {
		return err
	}

	// serialize children
	for _, child := range o.children {
		if err := Serialize(enc, child); err != nil {
			if errors.Is(err, ErrSerialize) {
				continue
			}
			return err
		}
	}
	return nil
}

func DeserializeObject(ctx Pool, dec Decoder) (Component, error) {
	var data ObjectState
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	obj := emptyObject(ctx, data.Name)
	obj.component = *data.ComponentState.New().(*component)
	obj.setEnabled(data.Enabled)
	obj.Transform().SetPosition(data.Position)
	obj.Transform().SetRotation(data.Rotation)
	obj.Transform().SetScale(data.Scale)
	ctx.assign(obj)

	// deserialize children
	for i := 0; i < data.Children; i++ {
		child, err := Deserialize[Component](ctx, dec)
		if err != nil {
			return nil, err
		}
		Attach(obj, child)
	}
	return obj, nil
}

func EncodeComponent(enc Encoder, cmp Component) error {
	state := NewComponentState(cmp)
	return enc.Encode(state)
}

func DecodeComponent(pool Pool, dec Decoder) (Component, error) {
	var data ComponentState
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return data.New(), nil
}
