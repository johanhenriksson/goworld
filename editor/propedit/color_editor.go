package propedit

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func init() {
	Register[color.T](func(key, name string, prop object.GenericProp) node.T {
		return ColorField(key, name, ColorProps{
			Value:    prop.GetAny().(color.T),
			OnChange: func(c color.T) { prop.SetAny(c) },
		})
	})
}

type ColorProps struct {
	Value    color.T
	OnChange func(color.T)
}

func ColorField(key string, title string, props ColorProps) node.T {
	return Vec3Field(key, title, Vec3Props{
		Value:    props.Value.Vec3(),
		OnChange: func(v vec3.T) { props.OnChange(color.FromVec3(v)) },
	})
}
