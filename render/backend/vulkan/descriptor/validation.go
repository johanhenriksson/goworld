package descriptor

import (
	"fmt"
	"reflect"
)

func ValidateShaderStruct(value any) error {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("value must be a struct, was %s", t.Kind())
	}

	// log.Println("uniform array of type", t.Name())
	// log.Println("  size:", d.Size)
	expectedOffset := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// log.Println("  field", field.Name, "offset:", field.Offset, "size:", field.Type.Size())
		if field.Offset != uintptr(expectedOffset) {
			return fmt.Errorf("layout causes alignment issues. expected field %s to have offset %d, was %d",
				field.Name, expectedOffset, field.Offset)
		}
		expectedOffset = int(field.Offset + field.Type.Size())
	}

	return nil
}

func Align(offset, alignment int) int {
	count := offset / alignment
	diff := offset % alignment
	if diff > 0 {
		count++
	}
	return count * alignment
}
