package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Editor base struct
type Editor struct {
	*engine.Group

	Chunk   *game.Chunk
	Camera  *engine.Camera
	Palette *PaletteWindow
	Tool    Tool

	PlaceTool   *PlaceTool
	EraseTool   *EraseTool
	SampleTool  *SampleTool
	ReplaceTool *ReplaceTool

	XPlane *geometry.Plane
	YPlane *geometry.Plane
	ZPlane *geometry.Plane

	xp, yp, zp int

	bounds  *geometry.Box
	mesh    *game.ChunkMesh
	gbuffer *render.GeometryBuffer
}

// NewEditor creates a new editor application
func NewEditor(chunk *game.Chunk, camera *engine.Camera, gbuffer *render.GeometryBuffer) *Editor {
	e := &Editor{
		Group:   engine.NewGroup("Editor", vec3.Zero, vec3.Zero),
		Chunk:   chunk,
		Camera:  camera,
		Palette: NewPaletteWindow(render.DefaultPalette),

		PlaceTool:   NewPlaceTool(),
		EraseTool:   NewEraseTool(),
		SampleTool:  NewSampleTool(),
		ReplaceTool: NewReplaceTool(),

		XPlane: geometry.NewPlane(float32(chunk.Sx), render.Red.WithAlpha(0.25)),
		YPlane: geometry.NewPlane(float32(chunk.Sy), render.Green.WithAlpha(0.25)),
		ZPlane: geometry.NewPlane(float32(chunk.Sz), render.Blue.WithAlpha(0.25)),

		mesh:    game.NewChunkMesh(chunk),
		bounds:  geometry.NewBox(vec3.NewI(chunk.Sx, chunk.Sy, chunk.Sz), render.DarkGrey),
		gbuffer: gbuffer,
	}

	e.xp = chunk.Sx
	e.XPlane.Rotation.X = -90
	e.XPlane.Rotation.Z = 90
	e.XPlane.Position = vec3.New(float32(e.xp), float32(chunk.Sy)/2, float32(chunk.Sz)/2)

	e.YPlane.Position = vec3.New(float32(chunk.Sx)/2, float32(e.yp), float32(chunk.Sz)/2)

	e.zp = chunk.Sz
	e.ZPlane.Rotation.X = -90
	e.ZPlane.Position = vec3.New(float32(chunk.Sx)/2, float32(chunk.Sy)/2, float32(e.zp))

	e.Tool = e.PlaceTool

	e.Attach(e.mesh, e.bounds, e.XPlane, e.YPlane, e.ZPlane)

	return e
}

func (e *Editor) Collect(pass engine.DrawPass, args engine.DrawArgs) {
	e.Group.Collect(pass, args)
	engine.Collect(pass, args, e.Tool)
}

func (e *Editor) Update(dt float32) {
	e.Group.Update(dt)
	engine.Update(dt, e.Tool)

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
		e.Tool = nil
	}

	// place tool
	if keys.Pressed(keys.F) {
		e.Tool = e.PlaceTool
	}

	// erase tool
	if keys.Pressed(keys.C) {
		e.Tool = e.EraseTool
	}

	// replace tool
	if keys.Pressed(keys.R) {
		e.Tool = e.ReplaceTool
	}

	// eyedropper tool
	if keys.Pressed(keys.T) {
		e.Tool = e.SampleTool
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
		e.XPlane.Position.X = float32(e.xp)
	}

	if keys.Pressed(keys.Key2) {
		e.yp = (e.yp + e.Chunk.Sy + m + 1) % (e.Chunk.Sy + 1)
		e.YPlane.Position.Y = float32(e.yp)
	}

	if keys.Pressed(keys.Key3) {
		e.zp = (e.zp + e.Chunk.Sz + m + 1) % (e.Chunk.Sz + 1)
		e.ZPlane.Position.Z = float32(e.zp)
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