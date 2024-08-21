package object

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"slices"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

var refPropType = reflect.TypeOf((*ReferenceProp)(nil)).Elem()
var valuePropType = reflect.TypeOf((*ValueProp)(nil)).Elem()

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
	if obj, isObject := item.(Object); isObject {
		_, err := serializeObject(enc, obj, depth)
		return err
	} else {
		return serializeComponent(enc, item, depth)
	}
}

func serializeObject(enc Encoder, obj Object, depth int) ([]Component, error) {
	val := reflect.ValueOf(obj).Elem()
	vtype := reflect.TypeOf(obj).Elem()

	log.Println(depth, "serializing object of type", typeName(obj))
	defer log.Println(depth, "end object")

	// object header
	enc.Encode(serializationHeader{
		Type:   typeName(obj),
		Object: true,
		Depth:  depth,
	})

	if vtype == baseObjectType {
		return encodeBaseObject(enc, obj, depth)
	}

	// embedded object
	var err error
	var children []Component
	for i := 0; i < vtype.NumField(); i++ {
		if vtype.Field(i).Anonymous {
			log.Println(depth, "object base of type", vtype.Field(i).Name)
			base := val.Field(i).Interface().(Object)
			if children, err = serializeObject(enc, base, depth+1); err != nil {
				return nil, err
			}
			break
		}
	}

	if err := encodePointers(enc, val, children); err != nil {
		return nil, err
	}
	if err := encodeReferences(enc, val); err != nil {
		return nil, err
	}
	if err := encodeProperties(enc, val); err != nil {
		return nil, err
	}

	return children, nil
}

func encodeBaseObject(enc Encoder, obj Object, depth int) ([]Component, error) {
	// count serializable children
	children := make([]Component, 0, len(obj.Children()))
	for _, child := range obj.Children() {
		if _, ok := child.(Serializable); ok {
			children = append(children, child)
		}
	}

	// object base
	if err := enc.Encode(ObjectState{
		ComponentState: NewComponentState(obj),
		Position:       obj.Transform().Position(),
		Rotation:       obj.Transform().Rotation(),
		Scale:          obj.Transform().Scale(),
		Children:       len(children),
	}); err != nil {
		return nil, err
	}

	// children
	for i, child := range children {
		log.Println(depth, "child", i)
		if err := serializeItem(enc, child, depth+1); err != nil {
			return nil, err
		}
	}

	return children, nil
}

func serializeComponent(enc Encoder, cmp Component, depth int) error {
	val := reflect.ValueOf(cmp).Elem()
	vtype := reflect.TypeOf(cmp).Elem()

	log.Println(depth, "serializing component of type", typeName(cmp))
	defer log.Println(depth, "end component")

	// object header
	enc.Encode(serializationHeader{
		Type:   typeName(cmp),
		Object: false,
		Depth:  depth,
	})

	if vtype == baseComponentType {
		return enc.Encode(NewComponentState(cmp))
	}

	for i := 0; i < vtype.NumField(); i++ {
		if vtype.Field(i).Anonymous {
			log.Println(depth, "component base of type", vtype.Field(i).Name)
			base := val.Field(i).Interface().(Component)
			if err := serializeComponent(enc, base, depth+1); err != nil {
				return err
			}
			break
		}
	}

	if err := encodeReferences(enc, val); err != nil {
		return err
	}
	if err := encodeProperties(enc, val); err != nil {
		return err
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

func deserializeItem(pool Pool, dec Decoder, depth int) (Component, error) {
	var header serializationHeader
	if err := dec.Decode(&header); err != nil {
		return nil, err
	}

	typeInfo, exists := types[header.Type]
	if !exists {
		return nil, fmt.Errorf("%w: unknown type %s", ErrSerialize, header.Type)
	}

	var err error
	var result Component
	if header.Object {
		log.Println(depth, "deserializing object", header.Type)
		result, err = deserializeObject(pool, dec, typeInfo, depth)
		log.Println(depth, "end object")
	} else {
		log.Println(depth, "deserializing component", header.Type)
		result, err = deserializeComponent(pool, dec, typeInfo, depth)
		log.Println(depth, "end object")
	}
	if err != nil {
		return nil, err
	}

	pool.assign(result)
	return result, nil
}

func deserializeObject(pool Pool, dec Decoder, typ TypeInfo, depth int) (Object, error) {
	var base Object
	if typ.rtype == baseObjectType {
		return decodeBaseObject(pool, dec, depth)
	} else {
		dbase, err := deserializeItem(pool, dec, depth+1)
		if err != nil {
			return nil, err
		}
		var isObject bool
		base, isObject = dbase.(Object)
		if !isObject {
			return nil, fmt.Errorf("%w: %s is not an object", ErrSerialize, typeName(dbase))
		}
	}

	// create object of the desired type
	obj := reflect.New(typ.rtype).Elem()
	setBase(obj, base)

	if err := decodePointers(pool, dec, obj, base.Children()); err != nil {
		return nil, err
	}
	if err := decodeReferences(pool, dec, obj); err != nil {
		return nil, err
	}
	if err := decodeProperties(pool, dec, obj); err != nil {
		return nil, err
	}

	result := obj.Addr().Interface().(Object)
	return result, nil
}

func decodeBaseObject(pool Pool, dec Decoder, depth int) (Object, error) {
	var data ObjectState
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
		log.Println(depth, "child", i)
		child, err := deserializeItem(pool, dec, depth+1)
		if err != nil {
			return nil, err
		}
		base.children = append(base.children, child)
	}

	return base, nil
}

func deserializeComponent(pool Pool, dec Decoder, typ TypeInfo, depth int) (Component, error) {
	var base Component
	if typ.rtype == baseComponentType {
		var state ComponentState
		if err := dec.Decode(&state); err != nil {
			return nil, err
		}
		return state.New(), nil
	} else {
		var err error
		base, err = deserializeItem(pool, dec, depth+1)
		if err != nil {
			return nil, err
		}
	}

	// create object of the desired type
	obj := reflect.New(typ.rtype).Elem()
	setBase(obj, base)

	if err := decodeReferences(pool, dec, obj); err != nil {
		return nil, err
	}
	if err := decodeProperties(pool, dec, obj); err != nil {
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
				log.Println("direct reference")
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

//
// deprecated
//

func (o *object) Serialize(enc Encoder) error {
	log.Println("old object serialize!")
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
