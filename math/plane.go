package math

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Plane struct {
    Normal  mgl.Vec3
    D       float32
}
