package vertex

import (
	"fmt"
	"strconv"
	"strings"
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
