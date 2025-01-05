package descriptor

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrDescriptorType = errors.New("invalid descriptor struct")

var setInterface = reflect.TypeOf((*Set)(nil)).Elem()

type Resolver interface {
	Descriptor(string) (int, bool)
}

type DescriptorBinding struct {
	Name string
	Descriptor
}

func ParseDescriptorStruct[S Set](template S) ([]Descriptor, error) {
	ptr := reflect.ValueOf(template)
	if ptr.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("%w: template must be a pointer to struct", ErrDescriptorType)
	}

	templateStruct := ptr.Elem()
	structName := templateStruct.Type().Name()
	if templateStruct.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: template %s must be a pointer to struct", ErrDescriptorType, structName)
	}

	descriptors := make([]Descriptor, 0, templateStruct.NumField())
	for i := 0; i < templateStruct.NumField(); i++ {
		field := templateStruct.Type().Field(i)
		templateField := templateStruct.Field(i)

		if field.Name == "Set" {
			// Field named Set must be an embedding of descriptor.Set
			if !templateField.IsNil() {
				return nil, fmt.Errorf("%w: %s member called Set must be nil", ErrDescriptorType, structName)
			}
			if !templateField.Type().Implements(setInterface) {
				return nil, fmt.Errorf("%w: %s member called Set must implement descriptor.Set", ErrDescriptorType, structName)
			}
			if !field.Anonymous {
				return nil, fmt.Errorf("%w: %s member called Set must be an anonymous field", ErrDescriptorType, structName)
			}
		} else {
			// template field must be a non-nil pointer
			if templateField.Kind() != reflect.Pointer {
				return nil, fmt.Errorf("%w: %s.%s is not a pointer, was %s", ErrDescriptorType, structName, field.Name, templateField.Kind())
			}
			if templateField.IsNil() {
				return nil, fmt.Errorf("%w: %s.%s is must not be nil", ErrDescriptorType, structName, field.Name)
			}

			// ensure the value is a Descriptor interface
			if !templateField.CanInterface() {
				return nil, fmt.Errorf("%w: %s.%s is not an interface", ErrDescriptorType, structName, field.Name)
			}
			descriptor, isDescriptor := templateField.Interface().(Descriptor)
			if !isDescriptor {
				return nil, fmt.Errorf("%w: %s.%s is not a Descriptor", ErrDescriptorType, structName, field.Name)
			}

			// ensure only the last descriptor element is of variable length
			_, isVariableLength := descriptor.(VariableDescriptor)
			if isVariableLength {
				isLast := i == templateStruct.NumField()-1
				if !isLast {
					return nil, fmt.Errorf("%w: %s.%s is variable length, but not the last element", ErrDescriptorType, structName, field.Name)
				}
			}

			descriptors = append(descriptors, descriptor)
		}
	}

	return descriptors, nil
}

// CopyDescriptorStruct instantiates a descriptor struct according to the given template.
// Assumes that the template has passed validation beforehand.
func CopyDescriptorStruct[S Set](template S, blank Set) (S, []Descriptor) {
	// dereference
	ptr := reflect.ValueOf(template)
	templateStruct := ptr.Elem()

	copyPtr := reflect.New(templateStruct.Type())

	descriptors := make([]Descriptor, 0, templateStruct.NumField())
	for i := 0; i < templateStruct.NumField(); i++ {
		copyField := copyPtr.Elem().Field(i)
		fieldName := templateStruct.Type().Field(i).Name

		if fieldName == "Set" {
			// store Set embedding reference
			copyField.Set(reflect.ValueOf(blank))
		} else {
			templateField := templateStruct.Field(i)

			// create a copy of the template field's value
			fieldValue := templateField.Elem()
			valueCopy := reflect.New(fieldValue.Type())
			valueCopy.Elem().Set(fieldValue)

			// write the value to the copied struct
			copyField.Set(valueCopy)

			// cast the copied value to a Descriptor interface
			descriptor := valueCopy.Interface().(Descriptor)

			descriptors = append(descriptors, descriptor)
		}
	}

	// finally, cast to correct type
	copy := copyPtr.Interface().(S)
	return copy, descriptors
}

func InitDescriptorStruct[S Set](set S, blank Set) []Descriptor {
	descriptors, err := ParseDescriptorStruct(set)
	if err != nil {
		panic(err)
	}

	v := reflect.ValueOf(set).Elem()
	v.FieldByName("Set").Set(reflect.ValueOf(blank))

	return descriptors
}
