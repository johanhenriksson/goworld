package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/light"
	. "github.com/johanhenriksson/goworld/core/object"
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
	editor.RegisterEditor(&light.Directional{}, NewDirectionalLightEditor)
}

type DirectionalLightEditor struct {
	Object
	target *light.Directional

	Shape  *physics.Sphere
	Body   *physics.RigidBody
	Sprite *sprite.Mesh
	Bounds *lines.Box
	GUI    gui.Fragment
}

func NewDirectionalLightEditor(ctx *editor.Context, lit *light.Directional) *DirectionalLightEditor {
	editor := NewObject(ctx.Objects, "DirectionalLightEditor", &DirectionalLightEditor{
		Object: Ghost(ctx.Objects, lit.Name(), lit.Transform()),
		target: lit,

		Bounds: lines.NewBox(ctx.Objects, lines.BoxArgs{
			Extents: vec3.New(10, 10, 1),
			Color:   color.Yellow,
		}),

		Shape: physics.NewSphere(ctx.Objects, 1),
		Body:  physics.NewRigidBody(ctx.Objects, 0),
		Sprite: sprite.New(ctx.Objects, sprite.Args{
			Size: vec2.New(1, 1),
			Texture: texture.PathArgsRef("editor/sprites/light.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
		}),

		GUI: editor.PropertyEditorFragment(ctx.Objects, gui.FragmentLast, func() node.T {
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

func (e *DirectionalLightEditor) Target() Component { return e.target }

func (e *DirectionalLightEditor) Select(ev mouse.Event) {
	Enable(e.GUI)
	Enable(e.Bounds)
}

func (e *DirectionalLightEditor) Deselect(ev mouse.Event) bool {
	Disable(e.GUI)
	Disable(e.Bounds)
	return true
}

func (e *DirectionalLightEditor) Actions() []editor.Action {
	return nil
}
