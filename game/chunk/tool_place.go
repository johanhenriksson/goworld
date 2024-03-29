package chunk

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

type PlaceTool struct {
	object.Object
	Box object.Object
}

var _ editor.Tool = &PlaceTool{}

func NewPlaceTool() *PlaceTool {
	padding := float32(0.05)
	return object.New("Place Tool", &PlaceTool{
		Box: object.Builder(object.Empty("Box")).
			Attach(lines.NewBox(lines.BoxArgs{
				Extents: vec3.One.Scaled(1 + padding),
				Color:   color.Blue,
			})).
			Position(vec3.New(-padding/2+0.5, -padding/2+0.5, -padding/2+0.5)).
			Create(),
	})
}

func (pt *PlaceTool) Use(editor *Editor, position, normal vec3.T) {
	target := position.Add(normal.Scaled(0.5))
	x, y, z := int(target.X), int(target.Y), int(target.Z)

	if editor.GetVoxel(x, y, z) != voxel.Empty {
		return
	}

	clr := editor.SelectedColor()
	log.Println("place", clr, "at", x, y, z)
	editor.SetVoxel(x, y, z, voxel.New(clr))

	// recompute mesh
	editor.Recalculate()
}

func (pt *PlaceTool) Hover(editor *Editor, position, normal vec3.T) {
	p := position.Add(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}

func (pt *PlaceTool) CanDeselect() bool {
	return false
}

func (pt *PlaceTool) ToolMouseEvent(ev mouse.Event, hover physics.RaycastHit) {
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
