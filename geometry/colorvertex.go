package geometry

import (
    "github.com/johanhenriksson/goworld/render"
)

/** Color vertex. Used in solid-color elements */
type ColorVertex struct {
    X, Y, Z     float32 // 12 bytes
    render.Color               // 16 bytes
} // 28 bytes

type ColorVertices []ColorVertex

func (buffer ColorVertices) Elements() int { return len(buffer) }
func (buffer ColorVertices) Size()     int { return 28 }
