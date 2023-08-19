package chunk

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

type EraseTool struct {
	object.Object
	Box object.Object
}

var _ editor.Tool = &EraseTool{}

func NewEraseTool() *EraseTool {
	padding := float32(0.05)
	return object.New("Erase Tool", &EraseTool{
		Box: object.Builder(object.Empty("Box")).
			Attach(lines.NewBox(lines.BoxArgs{
				Extents: vec3.One.Scaled(1 + padding),
				Color:   color.Red,
			})).
			Position(vec3.New(-padding/2+0.5, -padding/2+0.5, -padding/2+0.5)).
			Create(),
	})
}

func (pt *EraseTool) Use(editor *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	editor.SetVoxel(int(target.X), int(target.Y), int(target.Z), voxel.Empty)

	// recompute mesh
	editor.Recalculate()
}

func (pt *EraseTool) Hover(editor *Editor, position, normal vec3.T) {
	// parent actually refers to the editor right now
	// tools should be attached to their own object
	// they could potentially share positioning logic
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}

func (pt *EraseTool) CanDeselect() bool {
	return false
}

func (pt *EraseTool) ToolMouseEvent(ev mouse.Event, hover physics.RaycastHit) {
	if hover.Shape == nil {
		return
	}

	editor := object.GetInParents[*Editor](pt)
	if editor == nil {
		// hm?
		return
	}

	pos := editor.Transform().Unproject(hover.Point)
	norm := editor.Transform().UnprojectDir(hover.Normal)

	if ev.Action() == mouse.Move {
		pt.Hover(editor, pos, norm)
	}

	if ev.Action() == mouse.Press && ev.Button() == mouse.Button1 {
		pt.Use(editor, pos, norm)
		ev.Consume()
	}
}
