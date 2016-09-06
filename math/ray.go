package math

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Ray struct {
    Origin  mgl.Vec3
    Dir     mgl.Vec3
}

func (ray Ray) IntersectBox(box *AABB) (bool, mgl.Vec3) {
    // Fast Ray-Box Intersection by Andrew Woo
    // from "Graphics Gems", Academic Press, 1990

    const (
        RIGHT  = 0
        LEFT   = 1
        MIDDLE = 2
        DIM    = 3
    )

    hit         := mgl.Vec3 { }
    minB        := box.Min()
    maxB        := box.Max()
    inside      := true
    maxT        := [DIM]float32 { }
    candidate   := [DIM]float32 { }
    quadrant    := [DIM]uint8 { }
    whichPlane  := 0

    // Find candidate planes
    for i := 0; i < DIM; i++ {
        if ray.Origin[i] < minB[i] {
            quadrant[i] = LEFT
            candidate[i] = minB[i]
            inside = false
        } else if ray.Origin[i] > maxB[i] {
            quadrant[i] = RIGHT
            candidate[i] = maxB[i]
            inside = false
        } else {
            quadrant[i] = MIDDLE
        }
    }

    // ray origin is inside the bounding box
    if inside {
        return true, ray.Origin
    }

    // calculate T distance to candidate planes
    for i := 0; i < DIM; i++ {
        if quadrant[i] != MIDDLE && ray.Dir[i] != 0 {
            maxT[i] = (candidate[i] - ray.Origin[i]) / ray.Dir[i]
        } else {
            maxT[i] = -1
        }
    }

    // choose largest maxT
    for i := 0; i < DIM; i++ {
        if maxT[whichPlane] < maxT[i] {
            whichPlane = i
        }
    }

    // make sure final candidate is actually inside the bounding box
    if maxT[whichPlane] < 0 {
        return false, mgl.Vec3 { }
    }
    for i := 0; i < DIM; i++ {
        if whichPlane != i {
            hit[i] = ray.Origin[i] - maxT[whichPlane] * ray.Dir[i]
            if hit[i] < minB[i] || hit[i] > maxB[i] {
                return false, mgl.Vec3 { }
            }
        } else {
            hit[i] = candidate[i]
        }
    }

    return true, hit
}
