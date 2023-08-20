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
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

func init() {
	editor.Register(&light.Directional{}, NewDirectionalLightEditor)
}

type DirectionalLightEditor struct {
	object.Object
	target *light.Directional

	Shape  *physics.Sphere
	Body   *physics.RigidBody
	Sprite *sprite.Mesh
	Bounds *lines.Box
	GUI    gui.Fragment
}

func NewDirectionalLightEditor(ctx *editor.Context, lit *light.Directional) *DirectionalLightEditor {
	editor := object.New("DirectionalLightEditor", &DirectionalLightEditor{
		Object: object.Ghost(lit.Name(), lit.Transform()),
		target: lit,

		Bounds: lines.NewBox(lines.BoxArgs{
			Extents: vec3.New(10, 10, 1),
			Color:   color.Yellow,
		}),

		Shape: physics.NewSphere(1),
		Body:  physics.NewRigidBody(0),
		Sprite: sprite.New(sprite.Args{
			Size: vec2.New(1, 1),
			Texture: texture.PathArgsRef("textures/ui/light.png", texture.Args{
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
			)
		}),
	})

	return editor
}

func (e *DirectionalLightEditor) Target() object.Component { return e.target }

func (e *DirectionalLightEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
	object.Enable(e.Bounds)
}

func (e *DirectionalLightEditor) Deselect(ev mouse.Event) bool {
	object.Disable(e.GUI)
	object.Disable(e.Bounds)
	return true
}

func (e *DirectionalLightEditor) Actions() []editor.Action {
	return nil
}
