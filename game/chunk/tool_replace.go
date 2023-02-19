package chunk

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type ReplaceTool struct {
	object.T
	Box *box.T
}

func NewReplaceTool() *ReplaceTool {
	return object.New(&ReplaceTool{
		Box: box.New(box.Args{
			Size:  vec3.One,
			Color: color.Yellow,
		}),
	})
}

func (pt *ReplaceTool) Use(editor Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	editor.SetVoxel(int(target.X), int(target.Y), int(target.Z), voxel.New(editor.SelectedColor()))

	// recompute mesh
	editor.Recalculate()
}

func (pt *ReplaceTool) Hover(editor Editor, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}

func (pt *ReplaceTool) CanDeselect() bool {
	return false
}

func (pt *ReplaceTool) MouseEvent(ev mouse.Event) {
	editor, exists := object.FindInParents[Editor](pt)
	if !exists {
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
