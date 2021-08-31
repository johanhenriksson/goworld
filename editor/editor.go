package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Editor base struct
type Editor struct {
	object.T

	Chunk   *game.Chunk
	Camera  *engine.Camera
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
}

// NewEditor creates a new editor application
func NewEditor(chunk *game.Chunk, camera *engine.Camera, gbuffer *render.GeometryBuffer) *Editor {
	e := &Editor{
		T:       object.New("Editor"),
		Chunk:   chunk,
		Camera:  camera,
		Palette: NewPaletteWindow(render.DefaultPalette),

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
		Color: render.DarkGrey,
	}).Create(e.T)

	e.Palette.SetPosition(vec2.New(300, 20))

	// X Construction Plane
	plane.Builder(&e.XPlane, plane.Args{
		Size:  float32(chunk.Sx),
		Color: render.Red.WithAlpha(0.25),
	}).
		Position(center.WithX(0)).
		Rotation(vec3.New(-90, 0, 90)).
		Active(false).
		Create(e.T)

	// Y Construction Plane
	plane.Builder(&e.YPlane, plane.Args{
		Size:  float32(chunk.Sy),
		Color: render.Green.WithAlpha(0.25),
	}).
		Position(center.WithY(0)).
		Active(false).
		Create(e.T)

	// Z Construction Plane
	plane.Builder(&e.ZPlane, plane.Args{
		Size:  float32(chunk.Sz),
		Color: render.Blue.WithAlpha(0.25),
	}).
		Position(center.WithZ(0)).
		Rotation(vec3.New(-90, 0, 0)).
		Active(false).
		Create(e.T)

	e.SelectTool(e.PlaceTool)

	// could we avoid this somehow?
	e.Adopt(e.PlaceTool, e.ReplaceTool, e.EraseTool, e.SampleTool)

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

	e.updateToolSelection()
	e.updateConstructPlanes()
	e.updateTool()

	// clear chunk
	if keys.Pressed(keys.N) && keys.Ctrl() {
		e.Chunk.Clear()
		e.Chunk.Light.Calculate()
		e.mesh.Compute()
	}
}

func (e *Editor) updateToolSelection() {
	// deselect tool
	if keys.Pressed(keys.Escape) {
		e.DeselectTool()
	}

	// place tool
	if keys.Pressed(keys.F) {
		e.SelectTool(e.PlaceTool)
	}

	// erase tool
	if keys.Pressed(keys.C) {
		e.SelectTool(e.EraseTool)
	}

	// replace tool
	if keys.Pressed(keys.R) {
		e.SelectTool(e.ReplaceTool)
	}

	// eyedropper tool
	if keys.Pressed(keys.T) {
		e.SelectTool(e.SampleTool)
	}
}

func (e *Editor) updateConstructPlanes() {
	// toggle construction planes
	if keys.Ctrl() {
		if keys.Pressed(keys.X) {
			e.XPlane.SetActive(!e.XPlane.Active())
		}
		if keys.Pressed(keys.Y) {
			e.YPlane.SetActive(!e.YPlane.Active())
		}
		if keys.Pressed(keys.Z) {
			e.ZPlane.SetActive(!e.ZPlane.Active())
		}
		return
	}

	m := 1
	if keys.Shift() {
		m = -1
	}

	if keys.Pressed(keys.X) && e.XPlane.Active() {
		e.xp = (e.xp + e.Chunk.Sx + m + 1) % (e.Chunk.Sx + 1)
		p := e.XPlane.Transform().Position().WithX(float32(e.xp))
		e.XPlane.Transform().SetPosition(p)
	}

	if keys.Pressed(keys.Y) && e.YPlane.Active() {
		e.yp = (e.yp + e.Chunk.Sy + m + 1) % (e.Chunk.Sy + 1)
		p := e.YPlane.Transform().Position().WithY(float32(e.yp))
		e.YPlane.Transform().SetPosition(p)
	}

	if keys.Pressed(keys.Z) && e.ZPlane.Active() {
		e.zp = (e.zp + e.Chunk.Sz + m + 1) % (e.Chunk.Sz + 1)
		p := e.ZPlane.Transform().Position().WithZ(float32(e.zp))
		e.ZPlane.Transform().SetPosition(p)
	}
}

func (e *Editor) updateTool() {
	exists, position, normal := e.cursorPositionNormal()
	if !exists {
		return
	}

	if e.Tool != nil {
		e.Tool.Hover(e, position, normal)

		// use active tool
		if mouse.Pressed(mouse.Button2) {
			e.Tool.Use(e, position, normal)
		}
	}
}

// sample world position at current mouse coords
func (e *Editor) cursorPositionNormal() (bool, vec3.T, vec3.T) {
	depth, depthExists := e.gbuffer.SampleDepth(mouse.Position)
	if !depthExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewNormal, normalExists := e.gbuffer.SampleNormal(mouse.Position)
	if !normalExists {
		return false, vec3.Zero, vec3.Zero
	}

	position := e.Camera.Unproject(vec3.Extend(
		mouse.Position.Div(e.gbuffer.Depth.Size()),
		depth,
	))

	viewInv := e.Camera.View.Invert()
	normal := viewInv.TransformDir(viewNormal)

	return true, position, normal
}
