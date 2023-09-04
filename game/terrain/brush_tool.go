package terrain

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

type BrushTool struct {
	object.Object
	Sphere   *lines.Sphere
	Radius   object.Property[float32]
	Strength object.Property[float32]

	Brush   Brush
	pressed bool
	editor  *Editor
	center  vec3.T
}

var _ editor.Tool = &BrushTool{}

func NewBrushTool(brush Brush, color color.T) *BrushTool {
	e := object.New("Brush Tool", &BrushTool{
		Sphere: lines.NewSphere(lines.SphereArgs{
			Radius: 1,
			Color:  color,
		}),
		Radius:   object.NewProperty(float32(1)),
		Strength: object.NewProperty(float32(3)),
		Brush:    brush,
	})
	e.Radius.OnChange.Subscribe(func(float32) { e.Sphere.Radius.Set(e.Radius.Get()) })
	e.Radius.Set(4)
	return e
}

func (pt *BrushTool) Use(editor *Editor, position vec3.T, dt float32) {
	if pt.Brush == nil {
		return
	}

	// calculate affected terrain patch
	pos := position.Floor()
	r := int(math.Ceil(pt.Radius.Get()))
	mx, mz := math.Max(int(pos.X)-r, 0), math.Max(int(pos.Z)-r, 0)
	Mx, Mz := math.Min(int(pos.X)+r, editor.Tile.Size-1), math.Min(int(pos.Z)+r, editor.Tile.Size-1)
	wx, wz := Mx-mx, Mz-mz

	// empty patch
	if wx <= 0 || wz <= 0 {
		return
	}

	// cut a patch of terrain
	patch := editor.Tile.Patch(ivec2.New(mx, mz), ivec2.New(wx, wz))

	// apply brush to patch
	pt.Brush.Paint(patch, position, pt.Radius.Get(), pt.Strength.Get()*dt)

	// apply patch to tile
	editor.Tile.ApplyPatch(patch)

	// recompute mesh
	editor.Recalculate()
}

func (pt *BrushTool) Hover(editor *Editor, position, normal vec3.T) {
	pt.Transform().SetPosition(position)
}

func (pt *BrushTool) CanDeselect() bool {
	return false
}

func (pt *BrushTool) ToolMouseEvent(ev mouse.Event, hover physics.RaycastHit) {
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
		ev.Consume()
	}
}

func (pt *BrushTool) Update(scene object.Component, dt float32) {
	pt.Object.Update(scene, dt)

	if pt.pressed {
		pt.Use(pt.editor, pt.center, dt)
	}
}
