package editor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Editor base struct
type Editor struct {
	*object.T

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

	bounds  *geometry.Box
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

		XPlane: plane.New(plane.Args{
			Size:  float32(chunk.Sx),
			Color: render.Red.WithAlpha(0.25),
		}),
		YPlane: plane.New(plane.Args{
			Size:  float32(chunk.Sy),
			Color: render.Green.WithAlpha(0.25),
		}),
		ZPlane: plane.New(plane.Args{
			Size:  float32(chunk.Sz),
			Color: render.Blue.WithAlpha(0.25),
		}),

		mesh:    game.NewChunkMesh(chunk),
		bounds:  geometry.NewBox(vec3.NewI(chunk.Sx, chunk.Sy, chunk.Sz), render.DarkGrey),
		gbuffer: gbuffer,
	}

	e.xp = chunk.Sx
	e.XPlane.SetRotation(vec3.New(-90, 0, 90))
	e.XPlane.SetPosition(vec3.New(float32(e.xp), float32(chunk.Sy)/2, float32(chunk.Sz)/2))

	e.YPlane.SetPosition(vec3.New(float32(chunk.Sx)/2, float32(e.yp), float32(chunk.Sz)/2))

	e.zp = chunk.Sz
	e.ZPlane.SetRotation(vec3.New(-90, 0, 0))
	e.ZPlane.SetPosition(vec3.New(float32(chunk.Sx)/2, float32(chunk.Sy)/2, float32(e.zp)))

	e.Tool = e.PlaceTool

	// could we avoid this somehow?
	e.Attach(e.mesh, e.bounds,
		e.XPlane, e.YPlane, e.ZPlane,
		e.PlaceTool, e.EraseTool)

	return e
}

func (e *Editor) DeselectTool() {
	if e.Tool != nil {
		fmt.Println("Disable", e.Tool)
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

	if keys.Pressed(keys.N) && keys.Ctrl() {
		e.Chunk.Clear()
		e.Chunk.Light.Calculate()
		e.mesh.Compute()
	}

	m := 1
	if keys.Shift() {
		m = -1
	}

	if keys.Pressed(keys.Key1) {
		e.xp = (e.xp + e.Chunk.Sx + m + 1) % (e.Chunk.Sx + 1)
		e.XPlane.SetPosition(vec3.New(float32(e.xp), 0, 0))
	}

	if keys.Pressed(keys.Key2) {
		e.yp = (e.yp + e.Chunk.Sy + m + 1) % (e.Chunk.Sy + 1)
		e.YPlane.SetPosition(vec3.New(0, float32(e.yp), 0))
	}

	if keys.Pressed(keys.Key3) {
		e.zp = (e.zp + e.Chunk.Sz + m + 1) % (e.Chunk.Sz + 1)
		e.ZPlane.SetPosition(vec3.New(0, 0, float32(e.zp)))
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
