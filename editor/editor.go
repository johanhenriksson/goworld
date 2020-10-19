package editor

import (
	"github.com/johanhenriksson/goworld/engine"
)

// Editor base struct
type Editor struct {
}

// NewEditor creates a new editor application
func NewEditor(app *engine.Application) *Editor {
	return &Editor{}
}

func (e *Editor) Draw(args engine.DrawArgs) {

}

func (e *Editor) Update(dt float32) {

}

// editor components:
// - arcball camera (low prio)
// - tools
//   place voxel
//     1. palette
//     2. destination box - perhaps even ghost voxel?
//     3. placement grids
//   remove voxel
// 	   1. source box
