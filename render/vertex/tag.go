package vertex

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/johanhenriksson/goworld/render/types"
)

type Tag struct {
	Name      string
	Type      string
	Count     int
	Normalize bool
}

func ParseTag(tag string) (Tag, error) {
	p := strings.Split(tag, ",")
	if len(p) < 3 || len(p) > 4 {
		return Tag{}, fmt.Errorf("invalid vertex tag")
	}
	norm := false
	name := strings.Trim(p[0], " ")
	kind := strings.Trim(p[1], " ")

	count, err := strconv.Atoi(p[2])
	if err != nil {
		return Tag{}, fmt.Errorf("expected count to be a number")
	}

	if len(p) == 4 && p[3] == "normalize" {
		norm = true
	}
	return Tag{
		Name:      name,
		Type:      kind,
		Count:     count,
		Normalize: norm,
	}, nil
}

func ParsePointers(data interface{}) Pointers {
	var el reflect.Type

	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Struct {
		el = t
	} else if t.Kind() == reflect.Slice {
		el = t.Elem()
	} else {
		panic("must be struct or slice")
	}

	size := int(el.Size())

	offset := 0
	pointers := make(Pointers, 0, el.NumField())
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		if f.Name == "_" {
			// skip struct layout fields
			continue
		}
		tagstr := f.Tag.Get("vtx")
		if tagstr == "skip" {
			continue
		}
		tag, err := ParseTag(tagstr)
		if err != nil {
			fmt.Printf("tag error on %s.%s: %s\n", el.String(), f.Name, err)
			continue
		}

		kind, err := types.TypeFromString(tag.Type)
		if err != nil {
			panic(fmt.Errorf("invalid GL type: %s", tag.Type))
		}

		ptr := Pointer{
			Binding:     -1,
			Name:        tag.Name,
			Source:      kind,
			Destination: kind,
			Elements:    tag.Count,
			Normalize:   tag.Normalize,
			Offset:      offset,
			Stride:      size,
		}

		pointers = append(pointers, ptr)

		offset += kind.Size() * tag.Count
	}

	return pointers
}
