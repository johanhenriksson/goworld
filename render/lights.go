package render

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Attenuation struct {
    Constant    float32
    Linear      float32
    Quadratic   float32
}

type PointLight struct {
    Range       float32
    Position    mgl.Vec3
    Color       mgl.Vec3
}
