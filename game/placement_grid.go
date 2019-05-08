package game

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/render"
)

type PlacementGrid struct {
	*engine.ComponentBase
	Chunk *ColorChunk

	mesh *geometry.Lines

	/* Current height */
	Y int
}

func NewPlacementGrid(parentObject *engine.Object) *PlacementGrid {
	pg := &PlacementGrid{
		mesh:  geometry.CreateLines(),
		Chunk: nil,
	}
	pg.ComponentBase = engine.NewComponent(parentObject, pg)

	// find chunk
	obj, exists := parentObject.GetComponent(pg.Chunk)
	if !exists {
		panic("no chunk component")
	}
	chunk := obj.(*ColorChunk)

	// compute grid mesh
	pg.Chunk = chunk
	pg.Y = 8
	pg.Compute()

	return pg
}

func (grid *PlacementGrid) Up() {
	if grid.Y < (grid.Chunk.Size - 1) {
		fmt.Println("grid up")
		grid.Y += 1
		grid.Compute()
	}
}

func (grid *PlacementGrid) Down() {
	if grid.Y > 0 {
		fmt.Println("grid down")
		grid.Y -= 1
		grid.Compute()
	}
}

func (grid *PlacementGrid) Update(dt float32) {
	if engine.KeyReleased(engine.KeyJ) {
		grid.Down()
	}
	if engine.KeyReleased(engine.KeyK) {
		grid.Up()
	}
}

func (grid *PlacementGrid) Draw(args render.DrawArgs) {
	grid.mesh.Draw(args)
}

/* Compute grid mesh - draw an empty box for every empty
 * voxel in the current layer */
func (grid *PlacementGrid) Compute() {
	grid.mesh.Clear()

	for x := 0; x < grid.Chunk.Size; x++ {
		for z := 0; z < grid.Chunk.Size; z++ {
			if true || grid.Chunk.At(x, grid.Y, z) == nil {
				// place box
				grid.mesh.Box(float32(x), float32(grid.Y), float32(z), // position
					1, 0, 1, // size
					0, 0, 0, 0.35) // color (RGBA)
			}
		}
	}

	grid.mesh.Compute()
}
