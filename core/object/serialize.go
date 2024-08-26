package object

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

var ErrSerialize = errors.New("serialization error")

type Serializable interface {
	Serialize(Encoder) error
	Deserialize(Pool, Decoder) error
}

type Decoder interface {
	Decode(e any) error
}

type Encoder interface {
	Encode(data any) error
}

var serializableType = reflect.TypeOf((*Serializable)(nil)).Elem()

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

type componentState struct {
	ID      Handle
	Name    string
	Enabled bool
}

func newComponentState(c Component) componentState {
	return componentState{
		ID:      c.ID(),
		Name:    c.Name(),
		Enabled: c.Enabled(),
	}
}

type objectState struct {
	componentState
	Position vec3.T
	Rotation quat.T
	Scale    vec3.T
	Children int
}

type serializationHeader struct {
	Type   string
	Object bool
	Depth  int
}

//
// serialize
//

func Serialize(enc Encoder, obj Component) error {
	return serializeItem(enc, obj, 0)
}

func serializeItem(enc Encoder, item Component, depth int) error {
	kind := typeName(item)
	if _, exists := types[kind]; !exists {
		return fmt.Errorf("%w: type %s is not serializable", ErrSerialize, kind)
	}
	if obj, isObject := item.(Object); isObject {
		return serializeObject(enc, obj, depth)
	} else {
		return serializeComponent(enc, item, depth)
	}
}

func serializeObject(enc Encoder, obj Object, depth int) error {
	val := reflect.ValueOf(obj).Elem()
	vtype := reflect.TypeOf(obj).Elem()

	// object header
	enc.Encode(serializationHeader{
		Type:   typeName(obj),
		Object: true,
		Depth:  depth,
	})

	if vtype == baseObjectType {
		// object base
		if err := enc.Encode(objectState{
			componentState: newComponentState(obj),
			Position:       obj.Transform().Position(),
			Rotation:       obj.Transform().Rotation(),
			Scale:          obj.Transform().Scale(),
			Children:       len(obj.Children()),
		}); err != nil {
			return err
		}

		// children
		for _, child := range obj.Children() {
			if err := serializeItem(enc, child, depth+1); err != nil {
				return err
			}
		}
	} else {
		// embedded object
		var base Object
		for i := 0; i < vtype.NumField(); i++ {
			if vtype.Field(i).Anonymous {
				base = val.Field(i).Interface().(Object)
				if err := serializeObject(enc, base, depth+1); err != nil {
					return err
				}
				break
			}
		}

		if err := encodePointers(enc, val, base.Children()); err != nil {
			return err
		}
		if err := encodeFields(enc, val); err != nil {
			return err
		}
	}

	return nil
}

func serializeComponent(enc Encoder, cmp Component, depth int) error {
	val := reflect.ValueOf(cmp).Elem()
	vtype := reflect.TypeOf(cmp).Elem()

	// object header
	enc.Encode(serializationHeader{
		Type:   typeName(cmp),
		Object: false,
		Depth:  depth,
	})

	if vtype == baseComponentType {
		// base component
		return enc.Encode(newComponentState(cmp))
	} else {
		// embedded component
		for i := 0; i < vtype.NumField(); i++ {
			if vtype.Field(i).Anonymous {
				base := val.Field(i).Interface().(Component)
				if err := serializeComponent(enc, base, depth+1); err != nil {
					return err
				}
				break
			}
		}

		if err := encodeFields(enc, val); err != nil {
			return err
		}
	}

	return nil
}

// deserialization
func Deserialize[T Component](pool Pool, decoder Decoder) (T, error) {
	pool = newMappingPool(pool)

	var empty T
	result, err := deserializeItem(pool, decoder, 0)
	if err != nil {
		return empty, err
	}

	cast, ok := result.(T)
	if !ok {
		return empty, fmt.Errorf("%w: object is not of type %T", ErrSerialize, empty)
	}
	return cast, nil
}

func deserializeItem(pool Pool, dec Decoder, depth int) (result Component, err error) {
	var header serializationHeader
	if err := dec.Decode(&header); err != nil {
		return nil, err
	}

	typeInfo, exists := types[header.Type]
	if !exists {
		return nil, fmt.Errorf("%w: unknown type %s", ErrSerialize, header.Type)
	}

	if header.Object {
		result, err = deserializeObject(pool, dec, typeInfo, depth)
	} else {
		result, err = deserializeComponent(pool, dec, typeInfo, depth)
	}
	if err != nil {
		return nil, err
	}

	pool.assign(result)
	return result, nil
}

func deserializeObject(pool Pool, dec Decoder, typ *Type, depth int) (Object, error) {
	if typ.rtype == baseObjectType {
		return decodeBaseObject(pool, dec, depth)
	}

	// decode base object
	dbase, err := deserializeItem(pool, dec, depth+1)
	if err != nil {
		return nil, err
	}
	base, isObject := dbase.(Object)
	if !isObject {
		return nil, fmt.Errorf("%w: %s is not an object", ErrSerialize, typeName(dbase))
	}

	// create object of the desired type
	obj := reflect.New(typ.rtype).Elem()
	setBase(obj, base)

	if err := decodePointers(pool, dec, obj, base.Children()); err != nil {
		return nil, err
	}
	if err := decodeFields(pool, dec, obj); err != nil {
		return nil, err
	}

	result := obj.Addr().Interface().(Object)
	return result, nil
}

func decodeBaseObject(pool Pool, dec Decoder, depth int) (Object, error) {
	var data objectState
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	base := &object{
		component: component{
			id:      data.ID,
			name:    data.Name,
			enabled: data.Enabled,
		},
		transform: transform.New(data.Position, data.Rotation, data.Scale),
		children:  make([]Component, 0, data.Children),
	}

	// children
	for i := 0; i < data.Children; i++ {
		child, err := deserializeItem(pool, dec, depth+1)
		if err != nil {
			return nil, err
		}
		base.children = append(base.children, child)
	}

	return base, nil
}

func deserializeComponent(pool Pool, dec Decoder, typ *Type, depth int) (Component, error) {
	if typ.rtype == baseComponentType {
		var state componentState
		if err := dec.Decode(&state); err != nil {
			return nil, err
		}
		return &component{
			id:      state.ID,
			name:    state.Name,
			enabled: state.Enabled,
		}, nil
	}

	// decode component base
	base, err := deserializeItem(pool, dec, depth+1)
	if err != nil {
		return nil, err
	}

	// create object of the desired type
	obj := reflect.New(typ.rtype).Elem()
	setBase(obj, base)

	if err := decodeFields(pool, dec, obj); err != nil {
		return nil, err
	}

	output := obj.Addr().Interface().(Component)
	pool.assign(output)
	return output, nil
}

func setBase(obj reflect.Value, base any) {
	ot := obj.Type()
	for i := 0; i < ot.NumField(); i++ {
		if ot.Field(i).Anonymous {
			// assign base object
			obj.Field(i).Set(reflect.ValueOf(base))
			return
		}
	}
	panic("object has no base")
}

//
// serializable fields
//

func decodeFields(pool Pool, dec Decoder, obj reflect.Value) error {
	for i := 0; i < obj.Type().NumField(); i++ {
		field := obj.Type().Field(i)
		if field.Anonymous {
			continue
		}
		ftype := obj.Field(i).Addr().Type()
		if !ftype.Implements(serializableType) {
			continue
		}

		// instantiate & deserialize field
		fieldval := reflect.New(field.Type)
		serializer := fieldval.Interface().(Serializable)
		if err := serializer.Deserialize(pool, dec); err != nil {
			return err
		}
		obj.Field(i).Set(fieldval.Elem())
	}
	return nil
}

func encodeFields(enc Encoder, val reflect.Value) error {
	for i := 0; i < val.Type().NumField(); i++ {
		if val.Type().Field(i).Anonymous {
			continue
		}
		ftype := val.Field(i).Addr().Type()

		if ftype.Implements(serializableType) {
			serializable := val.Field(i).Addr().Interface().(Serializable)
			if err := serializable.Serialize(enc); err != nil {
				return err
			}
		}
	}
	return nil
}

//
// pointers
//

type childRef struct {
	Index int
}

func encodePointers(enc Encoder, obj reflect.Value, children []Component) error {
	for i := 0; i < obj.Type().NumField(); i++ {
		if obj.Type().Field(i).Anonymous {
			continue
		}

		// direct reference to child
		// not legal for components
		// perhaps a bad idea even for objects due dangling pointer issues
		// what happens if the child is deleted?
		if obj.Field(i).Type().Implements(componentType) {
			index := -1
			if cmp, ok := obj.Field(i).Interface().(Component); ok {
				index = slices.Index(children, cmp)
			}
			if err := enc.Encode(childRef{
				Index: index,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func decodePointers(pool Pool, dec Decoder, obj reflect.Value, children []Component) error {
	for i := 0; i < obj.NumField(); i++ {
		if obj.Type().Field(i).Anonymous {
			continue
		}
		if !obj.Field(i).Type().Implements(componentType) {
			continue
		}

		var ref childRef
		if err := dec.Decode(&ref); err != nil {
			return err
		}

		// set child reference
		if ref.Index < 0 {
			continue
		}
		obj.Field(i).Set(reflect.ValueOf(children[ref.Index]))
	}
	return nil
}

type MemorySerializer struct {
	Stream []any
	index  int
}

func (m *MemorySerializer) Encode(data any) error {
	m.Stream = append(m.Stream, data)
	return nil
}

func (m *MemorySerializer) Decode(target any) error {
	if m.index >= len(m.Stream) {
		return io.EOF
	}
	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(m.Stream[m.index]))
	m.index++
	return nil
}
