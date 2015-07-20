package engine

import (
    "github.com/johanhenriksson/goworld/geometry"
)

type Voxel struct {
    Id      uint16
    /* Texture Coords */
    Xp, Xn  uint16
    Yp, Yn  uint16
    Zp, Zn  uint16
}

type Chunk struct {
    Size    int
    Tileset *Tileset
    Data    []*Voxel
}

func (chk *Chunk) At(x, y, z int) *Voxel {
    s := chk.Size
    s2 := s * s
    pos := z * s2 + y * s + x
    if pos < 0 || pos > s2 * s {
        return nil
    }
    return chk.Data[pos]
}

func (chk *Chunk) Store(vbo *geometry.VertexBuffer) {
    s := chk.Size
    //cache := make(map[uint16]*Voxel)
    for z := 0; z < s; z++ {
        for y := 0; y < s; y++ {
            for x := 0; x < s; x++ {
                v := chk.At(x,y,z)
                if v == nil {
                    continue
                }

                /* Compute vertex data */
            }
        }
    }
}
