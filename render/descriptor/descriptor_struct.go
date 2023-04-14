package descriptor

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrDescriptorType = errors.New("invalid descriptor struct")

type Resolver interface {
	Descriptor(string) (int, bool)
}

func ParseDescriptorStruct[S Set](template S) (map[string]Descriptor, error) {
	ptr := reflect.ValueOf(template)
	if ptr.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("%w: template must be a pointer to struct", ErrDescriptorType)
	}

	templateStruct := ptr.Elem()
	structName := templateStruct.Type().Name()
	if templateStruct.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: template %s must be a pointer to struct", ErrDescriptorType, structName)
	}

	descriptors := make(map[string]Descriptor)
	for i := 0; i < templateStruct.NumField(); i++ {
		fieldName := templateStruct.Type().Field(i).Name
		templateField := templateStruct.Field(i)

		if fieldName == "Set" {
			// Field named Set must be an embedding of descriptor.Set
			if !templateField.IsNil() {
				return nil, fmt.Errorf("%w: %s member called Set must be nil", ErrDescriptorType, structName)
			}
		} else {
			// template field must be a non-nil pointer
			if templateField.Kind() != reflect.Pointer {
				return nil, fmt.Errorf("%w: %s.%s is not a pointer, was %s", ErrDescriptorType, structName, fieldName, templateField.Kind())
			}
			if templateField.IsNil() {
				return nil, fmt.Errorf("%w: %s.%s is must not be nil", ErrDescriptorType, structName, fieldName)
			}

			// ensure the value is a Descriptor interface
			if !templateField.CanInterface() {
				return nil, fmt.Errorf("%w: %s.%s is not an interface", ErrDescriptorType, structName, fieldName)
			}
			descriptor, isDescriptor := templateField.Interface().(Descriptor)
			if !isDescriptor {
				return nil, fmt.Errorf("%w: %s.%s is not a Descriptor", ErrDescriptorType, structName, fieldName)
			}

			// ensure only the last descriptor element is of variable length
			_, isVariableLength := descriptor.(VariableDescriptor)
			if isVariableLength {
				isLast := i == templateStruct.NumField()-1
				if !isLast {
					return nil, fmt.Errorf("%w: %s.%s is variable length, but not the last element", ErrDescriptorType, structName, fieldName)
				}
			}

			descriptors[fieldName] = descriptor
		}
	}

	return descriptors, nil
}

// CopyDescriptorStruct instantiates a descriptor struct according to the given template.
// Assumes that the template has passed validation beforehand.
func CopyDescriptorStruct[S Set](template S, blank Set, resolver Resolver) (S, []Descriptor) {
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

			// bind it to our blank descriptor set
			binding, exists := resolver.Descriptor(fieldName)
			if !exists {
				panic(fmt.Errorf("unresolved descriptor: %s", fieldName))
			}
			descriptor.Bind(blank, binding)
			descriptors = append(descriptors, descriptor)
		}
	}

	// finally, cast to correct type
	copy := copyPtr.Interface().(S)
	return copy, descriptors
}
