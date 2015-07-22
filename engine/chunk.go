package engine

import (
    "github.com/johanhenriksson/goworld/geometry"
)

type Chunk struct {
    Size    int
    Tileset *Tileset
    Data    []*Voxel

    vao     *geometry.VertexArray
    vbo     *geometry.VertexBuffer
}

func CreateChunk(size int, ts *Tileset) *Chunk {
    chk := &Chunk {
        Size: size,
        Tileset: ts,
        Data: make([]*Voxel, size * size * size),

        vao: geometry.CreateVertexArray(),
        vbo: geometry.CreateVertexBuffer(),
    }
    return chk
}

func (chk *Chunk) Clear() {
    for i := 0; i < len(chk.Data); i++ {
        chk.Data[i] = nil
    }
}

func (chk *Chunk) offset(x, y, z int) int {
    s := chk.Size
    s2 := s * s
    pos := z * s2 + y * s + x
    if pos < 0 || pos > s2 * s {
        panic("Voxel out of bounds")
    }
    return pos
}

func (chk *Chunk) At(x, y, z int) *Voxel {
    return chk.Data[chk.offset(x,y,z)]
}

func (chk *Chunk) Set(x, y, z int, voxel *Voxel) {
    chk.Data[chk.offset(x,y,z)] = voxel
}

func (chk *Chunk) Compute() *Mesh {
    s := chk.Size
    data := make(VoxelVertices, 0, 1)

    for z := 0; z < s; z++ {
        for y := 0; y < s; y++ {
            for x := 0; x < s; x++ {
                v := chk.At(x,y,z)
                if v == nil {
                    continue
                }

                data = v.Compute(uint8(x), uint8(y), uint8(z), data, chk.Tileset)
            }
        }
    }

    chk.vao.Length = int32(len(data))
    chk.vao.Bind()
    chk.vbo.Buffer(data)

    return CreateMesh(chk.vao, chk.Tileset.Material)
}
