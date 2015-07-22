package engine

import (
    "math"
    mgl "github.com/go-gl/mathgl/mgl32"
    //"github.com/johanhenriksson/goworld/util"
)

type Transform struct {
    Matrix      mgl.Mat4
    Position    mgl.Vec3
    Rotation    mgl.Vec3
    Scale       mgl.Vec3
    Forward     mgl.Vec3
    Right       mgl.Vec3
    Up          mgl.Vec3
}

func CreateTransform(x, y, z float32) *Transform {
    t := &Transform {
        Matrix:   mgl.Ident4(),
        Position: mgl.Vec3 { x,y,z },
        Rotation: mgl.Vec3 { 0,0,0 },
        Scale:    mgl.Vec3 { 1,1,1 },
    }
    t.Update(0)
    return t
}

func (t *Transform) Update(dt float32) {
    /* Update transform */
    rad         := t.Rotation.Mul(math.Pi / 180.0)
    rotation    := mgl.AnglesToQuat(rad[0], rad[1], rad[2], mgl.XYZ).Mat4()
    scaling     := mgl.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
    translation := mgl.Translate3D(t.Position[0], t.Position[1], t.Position[2])

    m := scaling.Mul4(rotation.Mul4(translation))

    /* Grab axis vectors */
    t.Right[0]   =  m[4*0+0]
    t.Right[1]   =  m[4*1+0]
    t.Right[2]   =  m[4*2+0]
    t.Up[0]      =  m[4*0+1]
    t.Up[1]      =  m[4*1+1]
    t.Up[2]      =  m[4*2+1]
    t.Forward[0] = -m[4*0+2]
    t.Forward[1] = -m[4*1+2]
    t.Forward[2] = -m[4*2+2]

    t.Matrix = m
}

func (t *Transform) Translate(offset mgl.Vec3) {
    t.Position = t.Position.Add(offset)
}
