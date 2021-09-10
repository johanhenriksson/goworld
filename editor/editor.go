package editor

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

// Editor base struct
type Editor struct {
	object.T

	Chunk   *game.Chunk
	Camera  camera.T
	Palette *PaletteWindow
	Tool    Tool

	PlaceTool   *PlaceTool
	EraseTool   *EraseTool
	SampleTool  *SampleTool
	ReplaceTool *ReplaceTool

	XPlane *plane.T
	YPlane *plane.T
	ZPlane *plane.T

	xp, yp, zp int

	bounds  *box.T
	mesh    *game.ChunkMesh
	gbuffer *render.GeometryBuffer

	cursorPos    vec3.T
	cursorNormal vec3.T
}

// NewEditor creates a new editor application
func NewEditor(chunk *game.Chunk, cam camera.T, gbuffer *render.GeometryBuffer) *Editor {
	e := &Editor{
		T:       object.New("Editor"),
		Chunk:   chunk,
		Camera:  cam,
		Palette: NewPaletteWindow(color.DefaultPalette),

		PlaceTool:   NewPlaceTool(),
		EraseTool:   NewEraseTool(),
		SampleTool:  NewSampleTool(),
		ReplaceTool: NewReplaceTool(),

		mesh:    game.NewChunkMesh(chunk),
		gbuffer: gbuffer,
	}

	dimensions := vec3.NewI(chunk.Sx, chunk.Sy, chunk.Sz)
	center := dimensions.Scaled(0.5)

	box.Builder(&e.bounds, box.Args{
		Size:  dimensions,
		Color: color.DarkGrey,
	}).
		Parent(e).
		Create()

	e.Palette.SetPosition(vec2.New(300, 20))

	// X Construction Plane
	plane.Builder(&e.XPlane, plane.Args{
		Size:  float32(chunk.Sx),
		Color: color.Red.WithAlpha(0.25),
	}).
		Parent(e).
		Position(center.WithX(0)).
		Rotation(vec3.New(-90, 0, 90)).
		Active(false).
		Create()

	// Y Construction Plane
	plane.Builder(&e.YPlane, plane.Args{
		Size:  float32(chunk.Sy),
		Color: color.Green.WithAlpha(0.25),
	}).
		Parent(e).
		Position(center.WithY(0)).
		Active(false).
		Create()

	// Z Construction Plane
	plane.Builder(&e.ZPlane, plane.Args{
		Size:  float32(chunk.Sz),
		Color: color.Blue.WithAlpha(0.25),
	}).
		Parent(e).
		Position(center.WithZ(0)).
		Rotation(vec3.New(-90, 0, 0)).
		Active(false).
		Create()

	e.Adopt(e.PlaceTool, e.ReplaceTool, e.EraseTool, e.SampleTool)

	e.SelectTool(e.PlaceTool)

	e.Attach(e.mesh)

	return e
}

func (e *Editor) DeselectTool() {
	if e.Tool != nil {
		e.Tool.SetActive(false)
		e.Tool = nil
	}
}

func (e *Editor) SelectTool(tool Tool) {
	e.DeselectTool()
	e.Tool = tool
	e.Tool.SetActive(true)
}

func (e *Editor) Update(dt float32) {
	e.T.Update(dt)
	// engine.Update(dt, e.Tool)

	// clear chunk
	// if keys.Pressed(keys.N) && keys.Ctrl() {
	// 	e.Chunk.Clear()
	// 	e.Chunk.Light.Calculate()
	// 	e.mesh.Compute()
	// }
}

// sample world position at current mouse coords
func (e *Editor) cursorPositionNormal(cursor vec2.T) (bool, vec3.T, vec3.T) {
	depth, depthExists := e.gbuffer.SampleDepth(cursor)
	if !depthExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewNormal, normalExists := e.gbuffer.SampleNormal(cursor)
	if !normalExists {
		return false, vec3.Zero, vec3.Zero
	}

	position := e.Camera.Unproject(vec3.Extend(
		cursor.Div(e.gbuffer.Depth.Size()),
		depth,
	))

	view := e.Camera.View()
	viewInv := view.Invert()
	normal := viewInv.TransformDir(viewNormal)

	return true, position, normal
}

func (e *Editor) InBounds(p vec3.T) bool {
	p = p.Floor()
	outside := p.X < 0 || p.Y < 0 || p.Z < 0 || int(p.X) >= e.Chunk.Sx || int(p.Y) >= e.Chunk.Sy || int(p.Z) >= e.Chunk.Sz
	return !outside
}

func (e *Editor) KeyEvent(ev keys.Event) {
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
	if ev.Modifier()&keys.Shift == keys.Shift {
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

func (e *Editor) MouseEvent(ev mouse.Event) {
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
