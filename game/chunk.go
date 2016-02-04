package game

import (
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
)

/* Chunks are smallest renderable units of voxel geometry */
type Chunk struct {
    Size    int
    Tileset *Tileset
    Data    []*Voxel

    vao     *render.VertexArray
    vbo     *render.VertexBuffer
    mesh    *engine.Mesh
}

func CreateChunk(size int, ts *Tileset) *Chunk {
    chk := &Chunk {
        Size: size,
        Tileset: ts,
        Data: make([]*Voxel, size * size * size),

        vao: render.CreateVertexArray(),
        vbo: render.CreateVertexBuffer(),
    }
    return chk
}

/* Clears all voxel data in this chunk */
func (chk *Chunk) Clear() {
    for i := 0; i < len(chk.Data); i++ {
        chk.Data[i] = nil
    }
}

/* Returns the slice offset for a given set of coordinates, as
   well as a bool indicating whether the position is within bounds. 
   If the point is out of bounds, zero is returned */
func (chk *Chunk) offset(x, y, z int) (int, bool) {
    s := chk.Size
    if x < 0 || x >= s || y < 0 || y >= s || z < 0 || z >= s {
        return 0, false
    }
    s2 := s * s
    pos := z * s2 + y * s + x
    return pos, true
}

/* Returns a pointer to the voxel defintion at the given position.
   If the space is empty, nil is returned */
func (chk *Chunk) At(x, y, z int) *Voxel {
    pos, ok := chk.offset(x,y,z)
    if !ok {
        return nil
    }
    return chk.Data[pos]
}

/* Sets a voxel. If it is outside bounds, nothing happens */
func (chk *Chunk) Set(x, y, z int, voxel *Voxel) {
    pos, ok := chk.offset(x,y,z)
    if !ok {
        return
    }
    chk.Data[pos] = voxel
}

func (chk *Chunk) Update(dt float32) {
}

func (chk *Chunk) Draw(args render.DrawArgs) {
    chk.vao.Draw()
}

/* Recomputes the chunk mesh and returns a pointer to it. */
func (chk *Chunk) Compute() {
    s := chk.Size
    data := make(VoxelVertices, 0, 1)

    for z := 0; z < s; z++ {
        for y := 0; y < s; y++ {
            for x := 0; x < s; x++ {
                v := chk.At(x,y,z)
                if v == nil {
                    /* Empty space */
                    continue
                }

                /* Simple optimization - dont draw hidden faces */
                xp := chk.At(x+1,y,z) == nil
                xn := chk.At(x-1,y,z) == nil
                yp := chk.At(x,y+1,z) == nil
                yn := chk.At(x,y-1,z) == nil
                zp := chk.At(x,y,z+1) == nil
                zn := chk.At(x,y,z-1) == nil

                /* Compute & append vertex data */
                data = v.Compute(data, uint8(x), uint8(y), uint8(z), xp, xn, yp, yn, zp, zn, chk.Tileset)
            }
        }
    }

    /* Buffer to GPU */
    chk.vao.Length = int32(len(data))
    chk.vao.Bind()
    chk.vbo.Buffer(data)
}
