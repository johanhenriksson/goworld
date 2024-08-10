package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/geometry/sprite"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

func init() {
	editor.Register(&light.Point{}, NewPointLightEditor)
}

type PointLightEditor struct {
	object.Object
	target *light.Point

	Shape  *physics.Sphere
	Body   *physics.RigidBody
	Sprite *sprite.Mesh
	Bounds *lines.Sphere
	GUI    gui.Fragment
}

func NewPointLightEditor(ctx *editor.Context, lit *light.Point) *PointLightEditor {
	editor := object.New("PointLightEditor", &PointLightEditor{
		Object: object.Ghost(lit.Name(), lit.Transform()),
		target: lit,

		Bounds: lines.NewSphere(lines.SphereArgs{
			Radius: lit.Range.Get(),
			Color:  color.Yellow,
		}),

		Shape: physics.NewSphere(1),
		Body:  physics.NewRigidBody(0),
		Sprite: sprite.New(sprite.Args{
			Size: vec2.New(1, 1),
			Texture: texture.PathArgsRef("editor/sprites/light.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
		}),

		GUI: editor.PropertyEditorFragment(gui.FragmentLast, func() node.T {
			return editor.Inspector(
				lit,
				propedit.ColorField("color", "Color", propedit.ColorProps{
					Value:    lit.Color.Get(),
					OnChange: lit.Color.Set,
				}),
				propedit.FloatField("intensity", "Intensity", propedit.FloatProps{
					Value:    lit.Intensity.Get(),
					OnChange: lit.Intensity.Set,
				}),
				propedit.FloatField("range", "Radius", propedit.FloatProps{
					Value:    lit.Range.Get(),
					OnChange: lit.Range.Set,
				}),
				propedit.FloatField("falloff", "Falloff", propedit.FloatProps{
					Value:    lit.Falloff.Get(),
					OnChange: lit.Falloff.Set,
				}),
			)
		}),
	})

	// todo: unsubscribe at some point
	lit.Range.OnChange.Subscribe(editor.Bounds.Radius.Set)

	return editor
}

func (e *PointLightEditor) Target() object.Component { return e.target }

func (e *PointLightEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
	object.Enable(e.Bounds)
}

func (e *PointLightEditor) Deselect(ev mouse.Event) bool {
	object.Disable(e.GUI)
	object.Disable(e.Bounds)
	return true
}

func (e *PointLightEditor) Actions() []editor.Action {
	return nil
}
