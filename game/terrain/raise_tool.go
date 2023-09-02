package terrain

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

type RaiseTool struct {
	object.Object
	Sphere *lines.Sphere
	Radius object.Property[float32]
}

var _ editor.Tool = &RaiseTool{}

func NewRaiseTool() *RaiseTool {
	e := object.New("Raise Tool", &RaiseTool{
		Sphere: lines.NewSphere(lines.SphereArgs{
			Radius: 1,
			Color:  color.Yellow,
		}),
		Radius: object.NewProperty(float32(1)),
	})
	e.Radius.OnChange.Subscribe(func(float32) { e.Sphere.Radius.Set(e.Radius.Get()) })
	e.Radius.Set(2)
	return e
}

func (pt *RaiseTool) Use(editor *Editor, position, normal vec3.T) {
	pos := position.Floor()

	log.Println("raise around", pos, pt.Radius.Get())
	r := int(math.Ceil(pt.Radius.Get()))
	for z := -r; z <= r; z++ {
		for x := -r; x <= r; x++ {
			px, pz := int(pos.X)+x, int(pos.Z)+z
			weight := 1 - vec2.NewI(px, pz).Sub(position.XZ()).Length()/pt.Radius.Get()
			weight = math.Max(0, weight)

			p := editor.Tile.Point(px, pz)
			p.Height += weight
			editor.Tile.SetPoint(px, pz, p)
		}
	}
	editor.Recalculate()
}

func (pt *RaiseTool) Hover(editor *Editor, position, normal vec3.T) {
	pt.Transform().SetPosition(position)
}

func (pt *RaiseTool) CanDeselect() bool {
	return false
}

func (pt *RaiseTool) ToolMouseEvent(ev mouse.Event, hover physics.RaycastHit) {
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
