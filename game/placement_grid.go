package game

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/render"
)

type PlacementGrid struct {
	*engine.ComponentBase
	ChunkMesh *ChunkMesh
	Color     render.Color

	mesh *geometry.Lines

	/* Current height */
	Y int
}

func NewPlacementGrid(parentObject *engine.Object) *PlacementGrid {
	pg := &PlacementGrid{
		mesh:      geometry.CreateLines(),
		ChunkMesh: nil,
		Color:     render.Black,
	}
	pg.ComponentBase = engine.NewComponent(parentObject, pg)

	// find chunk
	obj, exists := parentObject.GetComponent(pg.ChunkMesh)
	if !exists {
		panic("no chunk component")
	}
	chunk := obj.(*ChunkMesh)

	// compute grid mesh
	pg.ChunkMesh = chunk
	pg.Y = 15
	pg.Compute()

	return pg
}

func (grid *PlacementGrid) Up() {
	if grid.Y < (grid.ChunkMesh.Sy - 1) {
		fmt.Println("grid up")
		grid.Y++
		grid.Compute()
	}
}

func (grid *PlacementGrid) Down() {
	if grid.Y > 0 {
		fmt.Println("grid down")
		grid.Y--
		grid.Compute()
	}
}

func (grid *PlacementGrid) Update(dt float32) {
	if keys.Pressed(keys.J) {
		grid.Down()
	}
	if keys.Pressed(keys.K) {
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

	for x := 0; x < grid.ChunkMesh.Sx; x++ {
		for z := 0; z < grid.ChunkMesh.Sz; z++ {
			if true || grid.ChunkMesh.At(x, grid.Y, z) == EmptyVoxel {
				// place box
				grid.mesh.Box(float32(x), float32(grid.Y)+0.001, float32(z), // position
					1, 0, 1, // size
					grid.Color) // color (RGBA)
			}
		}
	}

	grid.mesh.Compute()
}
