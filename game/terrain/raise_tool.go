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
	Sphere   *lines.Sphere
	Radius   object.Property[float32]
	Strength object.Property[float32]

	pressed bool
	editor  *Editor
	center  vec3.T
	normal  vec3.T
}

var _ editor.Tool = &RaiseTool{}

func NewRaiseTool() *RaiseTool {
	e := object.New("Raise Tool", &RaiseTool{
		Sphere: lines.NewSphere(lines.SphereArgs{
			Radius: 1,
			Color:  color.Yellow,
		}),
		Radius:   object.NewProperty(float32(1)),
		Strength: object.NewProperty(float32(3)),
	})
	e.Radius.OnChange.Subscribe(func(float32) { e.Sphere.Radius.Set(e.Radius.Get()) })
	e.Radius.Set(4)
	return e
}

func (pt *RaiseTool) Use(editor *Editor, position, normal vec3.T, dt float32) {
	pos := position.Floor()

	// copy potential area of effect

	log.Println("raise around", pos, pt.Radius.Get())
	r := int(math.Ceil(pt.Radius.Get()))
	mx, mz := math.Max(int(pos.X)-r, 0), math.Max(int(pos.Z)-r, 0)
	Mx, Mz := math.Min(int(pos.X)+r, editor.Tile.Size-1), math.Min(int(pos.Z)+r, editor.Tile.Size-1)
	wx, wz := Mx-mx, Mz-mz

	points := make([][]Point, wz)
	for z := mz; z < Mz; z++ {
		points[z-mz] = make([]Point, wx)
		for x := mx; x < Mx; x++ {
			points[z-mz][x-mx] = editor.Tile.Point(x, z)
		}
	}

	// apply brush operation
	pt.Brush(points, position, mx, mz, pt.Strength.Get()*dt)

	// apply copied points to tile
	for z := mz; z < Mz; z++ {
		for x := mx; x < Mx; x++ {
			editor.Tile.SetPoint(x, z, points[z-mz][x-mx])
		}
	}

	// recompute mesh
	editor.Recalculate()
}

func (pt *RaiseTool) Brush(brush [][]Point, center vec3.T, ox, oz int, dt float32) {
	// implement brush operation
	// apply operations on copied points
	for z := 0; z < len(brush); z++ {
		for x := 0; x < len(brush[z]); x++ {
			weight := 1 - vec2.NewI(ox+x, oz+z).Sub(center.XZ()).Length()/pt.Radius.Get()
			weight = math.Max(0, weight)
			weight = weight * weight

			brush[z][x].Height += dt * weight
		}
	}
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

	if ev.Button() == mouse.Button1 {
		if ev.Action() == mouse.Press {
			pt.pressed = true
		} else if ev.Action() == mouse.Release {
			pt.pressed = false
		}
		pt.editor = editor
		pt.center = pos
		pt.normal = norm
		ev.Consume()
	}
}

func (pt *RaiseTool) Update(scene object.Component, dt float32) {
	pt.Object.Update(scene, dt)

	if pt.pressed {
		pt.Use(pt.editor, pt.center, pt.normal, dt)
	}
}
