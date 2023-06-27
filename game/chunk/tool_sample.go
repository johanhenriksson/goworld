package chunk

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type SampleTool struct {
	object.G
	Box *box.T

	Reselect func()
}

var _ editor.Tool = &SampleTool{}

func NewSampleTool() *SampleTool {
	return object.Group("Sample Tool", &SampleTool{
		Box: box.New(box.Args{
			Size:  vec3.One,
			Color: color.Purple,
		}),
	})
}

func (pt *SampleTool) Use(editor Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	voxel := editor.GetVoxel(int(target.X), int(target.Y), int(target.Z))
	editor.SelectColor(color.RGB8(voxel.R, voxel.G, voxel.B))

	if pt.Reselect != nil {
		pt.Reselect()
		pt.Reselect = nil
	}
}

func (pt *SampleTool) Hover(editor Editor, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Box.SetActive(true)
		pt.Transform().SetPosition(p.Floor())
	} else {
		pt.Box.SetActive(false)
	}
}

func (pt *SampleTool) CanDeselect() bool {
	return false
}

func (pt *SampleTool) MouseEvent(ev mouse.Event) {
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
