package game

import (
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

type PlacementGrid struct {
    *engine.ComponentBase
    Chunk     *Chunk

    mesh      *geometry.Lines

    /* Current height */
    Y     int
}

func NewPlacementGrid(parentObject *engine.Object) *PlacementGrid {
    pg := &PlacementGrid {
        mesh: geometry.CreateLines(),
        Chunk: nil,
    }
    pg.ComponentBase = engine.NewComponent(parentObject, pg)

    // find chunk
    obj, exists := parentObject.GetComponent(pg.Chunk)
    if !exists {
        panic("no chunk component")
    }
    chunk := obj.(*Chunk)

    // compute grid mesh
    pg.Chunk = chunk
    pg.Y = chunk.Size - 1
    pg.Compute()

    return pg
}

func (grid *PlacementGrid) Up() {
    if grid.Y < (grid.Chunk.Size - 1) {
        grid.Y += 1
        grid.Compute()
    }
}

func (grid *PlacementGrid) Down() {
    if grid.Y < (grid.Chunk.Size - 1) {
        grid.Y += 1
        grid.Compute()
    }
}

func (grid *PlacementGrid) Update(dt float32) {
}

func (grid *PlacementGrid) Draw(args render.DrawArgs) {
    grid.mesh.Draw(args)
}

/* Compute grid mesh - draw an empty box for every empty
 * voxel in the current layer */
func (grid *PlacementGrid) Compute() {
    for x := 0; x < grid.Chunk.Size; x++ {
        for z := 0; z < grid.Chunk.Size; z++ {
            if grid.Chunk.At(x, grid.Y, z) == nil {
                // place box
                grid.mesh.Box(float32(x), float32(grid.Y), float32(z), // position
                    1, 1, 1, // size
                    1, 1, 1, 1) // color (RGBA)
            }
        }
    }

    grid.mesh.Compute()
}
