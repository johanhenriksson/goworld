package engine

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Attenuation struct {
    Constant    float32
    Linear      float32
    Quadratic   float32
}

type Light struct {
    Attenuation Attenuation
    Position    mgl.Vec3
    Color       mgl.Vec3
    Range       float32
}

