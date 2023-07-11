package chunk

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type EraseTool struct {
	object.G
	Box *box.T
}

var _ editor.Tool = &EraseTool{}

func NewEraseTool() *EraseTool {
	return object.Group("Erase Tool", &EraseTool{
		Box: box.New(box.Args{
			Size:  vec3.One,
			Color: color.Red,
		}),
	})
}

func (pt *EraseTool) Use(editor Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	editor.SetVoxel(int(target.X), int(target.Y), int(target.Z), voxel.Empty)

	// recompute mesh
	editor.Recalculate()
}

func (pt *EraseTool) Hover(editor Editor, position, normal vec3.T) {
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

func (pt *EraseTool) MouseEvent(ev mouse.Event) {
	editor := object.GetInParents[Editor](pt)
	if editor == nil {
		// hm?
		return
	}

	if ev.Action() == mouse.Move {
		if exists, pos, normal := editor.CursorPositionNormal(ev.Position()); exists {
			pt.Hover(editor, pos, normal)
		}
	}

	if ev.Action() == mouse.Press && ev.Button() == mouse.Button1 {
		if exists, pos, normal := editor.CursorPositionNormal(ev.Position()); exists {
			pt.Use(editor, pos, normal)
			ev.Consume()
		}
	}
}
