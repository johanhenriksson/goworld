package builtin

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func init() {
	editor.Register(&light.Point{}, NewPointLightEditor)
}

type PointLightEditor struct {
	object.Object
	target *light.Point

	GUI gui.Fragment
}

func NewPointLightEditor(ctx *editor.Context, target *light.Point) *PointLightEditor {
	return object.New("PointLightEditor", &PointLightEditor{
		target: target,

		GUI: editor.InspectorGUI(
			target,
			propedit.FloatField("intensity", "Intensity", propedit.FloatProps{
				Value:    target.Intensity.Get(),
				OnChange: target.Intensity.Set,
			}),
			propedit.FloatField("range", "Range", propedit.FloatProps{
				Value:    target.Range.Get(),
				OnChange: target.Range.Set,
			}),
			propedit.Vec3Field("color", "Color", propedit.Vec3Props{
				Value: target.Color.Get().Vec3(),
				OnChange: func(v vec3.T) {
					target.Color.Set(color.FromVec3(v))
				},
			}),
		),
	})
}

func (e *PointLightEditor) Bounds() physics.Shape {
	return nil
}

func (e *PointLightEditor) Actions() []editor.Action {
	return nil
}

