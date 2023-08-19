package chunk

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

type SampleTool struct {
	object.Object
	Box object.Object

	Reselect func()
}

var _ editor.Tool = &SampleTool{}

func NewSampleTool() *SampleTool {
	padding := float32(0.05)
	return object.New("Sample Tool", &SampleTool{
		Box: object.Builder(object.Empty("Box")).
			Attach(lines.NewBox(lines.BoxArgs{
				Extents: vec3.One.Scaled(1 + padding),
				Color:   color.Purple,
			})).
			Position(vec3.New(-padding/2+0.5, -padding/2+0.5, -padding/2+0.5)).
			Create(),
	})
}

func (pt *SampleTool) Use(editor *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	voxel := editor.GetVoxel(int(target.X), int(target.Y), int(target.Z))
	editor.SelectColor(color.RGB8(voxel.R, voxel.G, voxel.B))

	if pt.Reselect != nil {
		pt.Reselect()
		pt.Reselect = nil
	}
}

func (pt *SampleTool) Hover(editor *Editor, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		object.Enable(pt.Box)
		pt.Transform().SetPosition(p.Floor())
	} else {
		object.Disable(pt.Box)
	}
}

func (pt *SampleTool) CanDeselect() bool {
	return false
}

func (pt *SampleTool) ToolMouseEvent(ev mouse.Event, hover physics.RaycastHit) {
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
