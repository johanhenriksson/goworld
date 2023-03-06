package propedit

import (
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Vec3Props struct {
	Value    vec3.T
	OnChange func(vec3.T)
}

func Vec3(key string, props Vec3Props) node.T {
	return node.Component(key, props, func(props Vec3Props) node.T {
		return rect.New(key, rect.Props{
			Style: rect.Style{
				Layout:         style.Row{},
				AlignItems:     style.AlignCenter,
				JustifyContent: style.JustifySpaceBetween,
				Grow:           style.Grow(1),
				Shrink:         style.Shrink(1),
				Basis:          style.Pct(100),
			},
			Children: []node.T{
				Float("x", FloatProps{
					Label: "X",
					Value: props.Value.X,
					OnChange: func(x float32) {
						if props.OnChange != nil {
							props.OnChange(vec3.New(x, props.Value.Y, props.Value.Z))
						}
					},
				}),
				Float("y", FloatProps{
					Label: "Y",
					Value: props.Value.Y,
					OnChange: func(y float32) {
						if props.OnChange != nil {
							props.OnChange(vec3.New(props.Value.X, y, props.Value.Z))
						}
					},
				}),
				Float("z", FloatProps{
					Label: "Z",
					Value: props.Value.Z,
					OnChange: func(z float32) {
						if props.OnChange != nil {
							props.OnChange(vec3.New(props.Value.X, props.Value.Y, z))
						}
					},
				}),
			},
		})
	})
}
