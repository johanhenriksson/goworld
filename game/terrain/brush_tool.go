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
	Brush    Brush

	terrain  *Map
	pressed  bool
	center   vec3.T
	heldtime float32
}

var _ editor.Tool = &BrushTool{}

func NewBrushTool(terrain *Map, brush Brush, color color.T) *BrushTool {
	e := object.New("Brush Tool", &BrushTool{
		Sphere: lines.NewSphere(lines.SphereArgs{
			Radius: 1,
			Color:  color,
		}),
		Radius:   object.NewProperty(float32(1)),
		Strength: object.NewProperty(float32(1)),
		Brush:    brush,

		terrain: terrain,
	})
	e.Radius.OnChange.Subscribe(func(float32) { e.Sphere.Radius.Set(e.Radius.Get()) })
	e.Radius.Set(4)
	return e
}

func (pt *BrushTool) Use(position vec3.T, dt float32) {
	if pt.Brush == nil {
		return
	}

	// calculate affected terrain patch
	pos := position.Floor()
	r := int(math.Ceil(pt.Radius.Get()))
	mx, mz := int(pos.X)-r, int(pos.Z)-r
	Mx, Mz := int(pos.X)+r, int(pos.Z)+r
	wx, wz := Mx-mx, Mz-mz

	// cut a patch of terrain
	patch := pt.terrain.Get(ivec2.New(mx, mz), ivec2.New(wx, wz))

	// apply brush to patch
	pt.Brush.Paint(patch, position, pt.Radius.Get(), pt.Strength.Get()*dt)

	// apply patch to tile
	pt.terrain.Set(patch)
}

func (pt *BrushTool) Hover(position vec3.T) {
	pt.Transform().SetWorldPosition(position)
}

func (pt *BrushTool) CanDeselect() bool {
	return false
}

func (pt *BrushTool) ToolMouseEvent(ev mouse.Event, hover physics.RaycastHit) {
	if ev.Action() == mouse.Scroll {
		radiusSensitivity := float32(0.33)
		radius := math.Clamp(pt.Radius.Get()-ev.Scroll().Y*radiusSensitivity, 0.5, 64)
		pt.Radius.Set(radius)
	}

	if hover.Shape == nil {
		return
	}

	if ev.Action() == mouse.Move {
		pt.Hover(hover.Point)
	}

	if ev.Button() == mouse.Button1 {
		if ev.Action() == mouse.Press {
			pt.pressed = true
		} else if ev.Action() == mouse.Release {
			pt.pressed = false
		}
		pt.center = hover.Point
		ev.Consume()
	}
}

func (pt *BrushTool) Update(scene object.Component, dt float32) {
	pt.Object.Update(scene, dt)

	if pt.pressed {
		pt.heldtime += dt

		const updatesPerSecond = 15
		const updateHoldTime float32 = 1.0 / updatesPerSecond
		if pt.heldtime > updateHoldTime {
			pt.Use(pt.center, updateHoldTime)
			pt.heldtime -= updateHoldTime
		}
	}
}
