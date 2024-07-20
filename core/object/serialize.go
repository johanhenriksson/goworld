package object

import (
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

func Copy(obj Component) Component {
	buffer := &MemorySerializer{}

	err := Serialize(buffer, obj)
	if err != nil {
		panic(err)
	}

	kopy, err := Deserialize(buffer)
	if err != nil {
		panic(err)
	}

	return kopy
}

func Save(key string, obj Component) error {
	fp, err := assets.Write(key)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := gob.NewEncoder(fp)
	return Serialize(enc, obj)
}

func Load(key string) (Component, error) {
	fp, err := assets.Open(key)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	dec := gob.NewDecoder(fp)
	return Deserialize(dec)
}

type ComponentState struct {
	ID      uint
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

func Deserialize(decoder Decoder) (Component, error) {
	var kind string
	if err := decoder.Decode(&kind); err != nil {
		return nil, err
	}
	typ, exists := types[kind]
	if !exists {
		return nil, fmt.Errorf("%w: no deserializer for %s", ErrSerialize, kind)
	}
	return typ.Deserialize(decoder)
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

func DeserializeObject(dec Decoder) (Component, error) {
	var data ObjectState
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	obj := Empty(data.Name)
	obj.setEnabled(data.Enabled)
	obj.Transform().SetPosition(data.Position)
	obj.Transform().SetRotation(data.Rotation)
	obj.Transform().SetScale(data.Scale)

	// deserialize children
	for i := 0; i < data.Children; i++ {
		child, err := Deserialize(dec)
		if err != nil {
			return nil, err
		}
		Attach(obj, child)
	}
	return obj, nil
}
