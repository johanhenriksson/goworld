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
	editor.Register(&light.Directional{}, NewDirectionalLightEditor)
}

type DirectionalLightEditor struct {
	object.Object
	target light.T

	GUI gui.Fragment
}

func NewDirectionalLightEditor(ctx *editor.Context, target *light.Directional) *DirectionalLightEditor {
	return object.New("DirectionalLightEditor", &DirectionalLightEditor{
		target: target,

		GUI: editor.InspectorGUI(
			target,
			propedit.FloatField("intensity", "Intensity", propedit.FloatProps{
				Value:    target.Intensity.Get(),
				OnChange: target.Intensity.Set,
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

func (e *DirectionalLightEditor) Bounds() physics.Shape {
	return nil
}

func (e *DirectionalLightEditor) Actions() []editor.Action {
	return nil
}
