package editor

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	object.Component

	SelectedColor() color.T
	GetVoxel(x, y, z int) voxel.T
	SetVoxel(x, y, z int, v voxel.T)
	SelectTool(Tool)
	SelectColor(color.T)
	Recalculate()
	InBounds(p vec3.T) bool
}

// Editor base struct
type editor struct {
	object.Component

	Chunk  *chunk.T
	Camera camera.T
	Tool   Tool

	color  color.T
	object object.T

	PlaceTool   *PlaceTool
	EraseTool   *EraseTool
	SampleTool  *SampleTool
	ReplaceTool *ReplaceTool

	XPlane *plane.T
	YPlane *plane.T
	ZPlane *plane.T

	xp, yp, zp int

	bounds  *box.T
	mesh    *chunk.Mesh
	gbuffer pass.BufferOutput

	cursorPos    vec3.T
	cursorNormal vec3.T
}

// NewEditor creates a new editor application
func NewEditor(chk *chunk.T, cam camera.T, gbuffer pass.BufferOutput) T {
	e := &editor{
		Component: object.NewComponent(),

		Chunk:  chk,
		Camera: cam,

		PlaceTool:   NewPlaceTool(),
		EraseTool:   NewEraseTool(),
		SampleTool:  NewSampleTool(),
		ReplaceTool: NewReplaceTool(),

		gbuffer: gbuffer,
		color:   color.Red,
		object:  object.New("Editor"),
	}

	dimensions := vec3.NewI(chk.Sx, chk.Sy, chk.Sz)
	center := dimensions.Scaled(0.5)

	box.Builder(&e.bounds, box.Args{
		Size:  dimensions,
		Color: color.DarkGrey,
	}).
		Parent(e.object).
		Create()

	// X Construction Plane
	plane.Builder(&e.XPlane, plane.Args{
		Size:  float32(chk.Sx),
		Color: color.Red.WithAlpha(0.25),
	}).
		Parent(e.object).
		Position(center.WithX(0)).
		Rotation(vec3.New(-90, 0, 90)).
		Active(false).
		Create()

	// Y Construction Plane
	plane.Builder(&e.YPlane, plane.Args{
		Size:  float32(chk.Sy),
		Color: color.Green.WithAlpha(0.25),
	}).
		Parent(e.object).
		Position(center.WithY(0)).
		Active(false).
		Create()

	// Z Construction Plane
	plane.Builder(&e.ZPlane, plane.Args{
		Size:  float32(chk.Sz),
		Color: color.Blue.WithAlpha(0.25),
	}).
		Parent(e.object).
		Position(center.WithZ(0)).
		Rotation(vec3.New(-90, 0, 0)).
		Active(false).
		Create()

	e.object.Adopt(e.PlaceTool, e.ReplaceTool, e.EraseTool, e.SampleTool)
	e.ReplaceTool.SetActive(false)
	e.EraseTool.SetActive(false)
	e.SampleTool.SetActive(false)

	e.SelectTool(e.PlaceTool)

	return e
}

func (e *editor) Name() string {
	return "Editor"
}

func (e *editor) AttachTo(obj object.T) {
	e.Component.AttachTo(obj)
	obj.Adopt(e.object)

	e.mesh = query.New[*chunk.Mesh]().First(obj)
	if e.mesh == nil {
		panic("parent object has no ChunkMesh attached")
	}
}

func (e *editor) GetVoxel(x, y, z int) voxel.T {
	return e.Chunk.At(x, y, z)
}

func (e *editor) SetVoxel(x, y, z int, v voxel.T) {
	e.Chunk.Set(x, y, z, v)
}

func (e *editor) SelectColor(c color.T) {
	e.color = c
}

func (e *editor) SelectedColor() color.T {
	return e.color
}

func (e *editor) DeselectTool() {
	if e.Tool != nil {
		e.Tool.SetActive(false)
		e.Tool = nil
	}
}

func (e *editor) SelectTool(tool Tool) {
	e.DeselectTool()
	e.Tool = tool
	e.Tool.SetActive(true)
}

func (e *editor) Update(dt float32) {
	// e.T.Update(dt)
	// engine.Update(dt, e.Tool)

	// clear chunk
	// if keys.Pressed(keys.N) && keys.Ctrl() {
	// 	e.Chunk.Clear()
	// 	e.Chunk.Light.Calculate()
	// 	e.mesh.Compute()
	// }
}

// sample world position at current mouse coords
func (e *editor) cursorPositionNormal(cursor vec2.T) (bool, vec3.T, vec3.T) {
	viewNormal, normalExists := e.gbuffer.SampleNormal(cursor)
	if !normalExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewPosition, positionExists := e.gbuffer.SamplePosition(cursor)
	if !positionExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewInv := e.Camera.ViewInv()
	normal := viewInv.TransformDir(viewNormal)
	position := viewInv.TransformPoint(viewPosition)

	return true, position, normal
}

func (e *editor) InBounds(p vec3.T) bool {
	p = p.Floor()
	outside := p.X < 0 || p.Y < 0 || p.Z < 0 || int(p.X) >= e.Chunk.Sx || int(p.Y) >= e.Chunk.Sy || int(p.Z) >= e.Chunk.Sz
	return !outside
}

func (e *editor) KeyEvent(ev keys.Event) {
	// clear chunk hotkey
	if keys.PressedMods(ev, keys.N, keys.Ctrl) {
		e.Chunk.Clear()
		e.Chunk.Light.Calculate()
		e.mesh.Compute()
	}

	// deselect tool
	if keys.Pressed(ev, keys.Escape) {
		e.DeselectTool()
	}

	// place tool
	if keys.Pressed(ev, keys.F) {
		e.SelectTool(e.PlaceTool)
	}

	// erase tool
	if keys.Pressed(ev, keys.C) {
		e.SelectTool(e.EraseTool)
	}

	// replace tool
	if keys.Pressed(ev, keys.R) {
		e.SelectTool(e.ReplaceTool)
	}

	// eyedropper tool
	if keys.Pressed(ev, keys.T) {
		e.SelectTool(e.SampleTool)
	}

	// toggle construction planes
	if keys.PressedMods(ev, keys.X, keys.Ctrl) {
		e.XPlane.Object().SetActive(!e.XPlane.Object().Active())
	}
	if keys.PressedMods(ev, keys.Y, keys.Ctrl) {
		e.YPlane.Object().SetActive(!e.YPlane.Object().Active())
	}
	if keys.PressedMods(ev, keys.Z, keys.Ctrl) {
		e.ZPlane.Object().SetActive(!e.ZPlane.Object().Active())
	}

	m := 1
	if ev.Modifier(keys.Shift) {
		m = -1
	}

	if keys.Pressed(ev, keys.X) && e.XPlane.Active() {
		e.xp = (e.xp + e.Chunk.Sx + m + 1) % (e.Chunk.Sx + 1)
		p := e.XPlane.Transform().Position().WithX(float32(e.xp))
		e.XPlane.Transform().SetPosition(p)
	}

	if keys.Pressed(ev, keys.Y) && e.YPlane.Active() {
		e.yp = (e.yp + e.Chunk.Sy + m + 1) % (e.Chunk.Sy + 1)
		p := e.YPlane.Transform().Position().WithY(float32(e.yp))
		e.YPlane.Transform().SetPosition(p)
	}

	if keys.Pressed(ev, keys.Z) && e.ZPlane.Active() {
		e.zp = (e.zp + e.Chunk.Sz + m + 1) % (e.Chunk.Sz + 1)
		p := e.ZPlane.Transform().Position().WithZ(float32(e.zp))
		e.ZPlane.Transform().SetPosition(p)
	}
}

func (e *editor) MouseEvent(ev mouse.Event) {
	switch ev.Action() {
	case mouse.Move:
		if exists, pos, normal := e.cursorPositionNormal(ev.Position()); exists {
			e.cursorPos = pos
			e.cursorNormal = normal
			e.Tool.Hover(e, pos, normal)
		}

	case mouse.Press:
		if ev.Button() == mouse.Button2 {
			e.Tool.Use(e, e.cursorPos, e.cursorNormal)
		}
	}
}

func (e *editor) Recalculate() {
	e.Chunk.Light.Calculate()
	if e.mesh != nil {
		e.mesh.Compute()
	}
}
