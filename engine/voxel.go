package engine

import (
    "github.com/johanhenriksson/goworld/geometry"
)

type VoxelId uint16

type Voxel struct {
    Id      VoxelId
    /* Texture Coords */
    Xp, Xn  TileId
    Yp, Yn  TileId
    Zp, Zn  TileId
}

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

func (chk *Chunk) At(x, y, z int) *Voxel {
    s := chk.Size
    s2 := s * s
    pos := z * s2 + y * s + x
    if pos < 0 || pos > s2 * s {
        return nil
    }
    return chk.Data[pos]
}

func (chk *Chunk) Set(x, y, z int, voxel *Voxel) {
    s := chk.Size
    s2 := s * s
    pos := z * s2 + y * s + x
    if pos < 0 || pos > s2 * s {
        return
    }
    chk.Data[pos] = voxel
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

    mesh := CreateMesh(chk.vao, chk.Tileset.Material)

    return mesh
}

func (voxel *Voxel) Compute(x, y, z uint8, data VoxelVertices, ts *Tileset) VoxelVertices {
    up, down, front, back, left, right := ts.Get(voxel.Yp), ts.Get(voxel.Yn),
                                          ts.Get(voxel.Zp), ts.Get(voxel.Zn),
                                          ts.Get(voxel.Xn), ts.Get(voxel.Xp)

    ux, uy := uint8(up.X), uint8(up.Y)
    dx, dy := uint8(down.X), uint8(down.Y)
    fx, fy := uint8(front.X), uint8(front.Y)
    bx, by := uint8(back.X), uint8(back.Y)
    lx, ly := uint8(left.X), uint8(left.Y)
    rx, ry := uint8(right.X), uint8(right.Y)

    return append(data,
        // Top 
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+0,
            Nx:  0,
            Ny:  1,
            Nz:  0,
            Tx: ux+0,
            Ty: uy+0, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+1,
            Nx:  0,
            Ny:  1,
            Nz:  0,
            Tx: ux+0,
            Ty: uy+1, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+0,
            Nx:  0,
            Ny:  1,
            Nz:  0,
            Tx: ux+1,
            Ty: uy+0, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+0,
            Nx:  0,
            Ny:  1,
            Nz:  0,
            Tx: ux+1,
            Ty: uy+0, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+1,
            Nx:  0,
            Ny:  1,
            Nz:  0,
            Tx: ux+0,
            Ty: uy+1, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+1,
            Nx:  0,
            Ny:  1,
            Nz:  0,
            Tx: ux+1,
            Ty: uy+1, },
        // Bottom 
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+0,
            Nx:  0,
            Ny: -1,
            Nz:  0,
            Tx: dx+0,
            Ty: dy+0, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+0,
            Nx:  0,
            Ny: -1,
            Nz:  0,
            Tx: dx+1,
            Ty: dy+0, },
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+1,
            Nx:  0,
            Ny: -1,
            Nz:  0,
            Tx: dx+0,
            Ty: dy+1, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+0,
            Nx:  0,
            Ny: -1,
            Nz:  0,
            Tx: dx+1,
            Ty: dy+0, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+1,
            Nx:  0,
            Ny: -1,
            Nz:  0,
            Tx: dx+1,
            Ty: dy+1, },
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+1,
            Nx:  0,
            Ny: -1,
            Nz:  0,
            Tx: dx+0,
            Ty: dy+1, },
        // Front
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+1,
            Nx:  0,
            Ny:  0,
            Nz:  1,
            Tx: fx+1,
            Ty: fy+0, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+1,
            Nx:  0,
            Ny:  0,
            Nz:  1,
            Tx: fx+0,
            Ty: fy+0, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+1,
            Nx:  0,
            Ny:  0,
            Nz:  1,
            Tx: fx+1,
            Ty: fy+1, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+1,
            Nx:  0,
            Ny:  0,
            Nz:  1,
            Tx: fx+0,
            Ty: fy+0, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+1,
            Nx:  0,
            Ny:  0,
            Nz:  1,
            Tx: fx+0,
            Ty: fy+1, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+1,
            Nx:  0,
            Ny:  0,
            Nz:  1,
            Tx: fx+1,
            Ty: fy+1, },
        // Back
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+0,
            Nx:  0,
            Ny:  0,
            Nz: -1,
            Tx: bx+0,
            Ty: by+0, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+0,
            Nx:  0,
            Ny:  0,
            Nz: -1,
            Tx: bx+0,
            Ty: by+1, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+0,
            Nx:  0,
            Ny:  0,
            Nz: -1,
            Tx: bx+1,
            Ty: by+0, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+0,
            Nx:  0,
            Ny:  0,
            Nz: -1,
            Tx: bx+1,
            Ty: by+0, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+0,
            Nx:  0,
            Ny:  0,
            Nz: -1,
            Tx: bx+0,
            Ty: by+1, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+0,
            Nx:  0,
            Ny:  0,
            Nz: -1,
            Tx: bx+1,
            Ty: by+1, },
        // Left
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+1,
            Nx: -1,
            Ny:  0,
            Nz:  0,
            Tx: lx+0,
            Ty: ly+1, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+0,
            Nx: -1,
            Ny:  0,
            Nz:  0,
            Tx: lx+1,
            Ty: ly+0, },
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+0,
            Nx: -1,
            Ny:  0,
            Nz:  0,
            Tx: lx+0,
            Ty: ly+0, },
        VoxelVertex {
            X: x+0,
            Y: y+0,
            Z: z+1,
            Nx: -1,
            Ny:  0,
            Nz:  0,
            Tx: lx+0,
            Ty: ly+1, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+1,
            Nx: -1,
            Ny:  0,
            Nz:  0,
            Tx: lx+1,
            Ty: ly+1, },
        VoxelVertex {
            X: x+0,
            Y: y+1,
            Z: z+0,
            Nx: -1,
            Ny:  0,
            Nz:  0,
            Tx: lx+1,
            Ty: ly+0, },
        // Right
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+1,
            Nx:  1,
            Ny:  0,
            Nz:  0,
            Tx: rx+1,
            Ty: ry+1, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+0,
            Nx:  1,
            Ny:  0,
            Nz:  0,
            Tx: rx+1,
            Ty: ry+0, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+0,
            Nx:  1,
            Ny:  0,
            Nz:  0,
            Tx: rx+0,
            Ty: ry+0, },
        VoxelVertex {
            X: x+1,
            Y: y+0,
            Z: z+1,
            Nx:  1,
            Ny:  0,
            Nz:  0,
            Tx: rx+1,
            Ty: ry+1, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+0,
            Nx:  1,
            Ny:  0,
            Nz:  0,
            Tx: rx+0,
            Ty: ry+0, },
        VoxelVertex {
            X: x+1,
            Y: y+1,
            Z: z+1,
            Nx:  1,
            Ny:  0,
            Nz:  0,
            Tx: rx+0,
            Ty: ry+1, })
}
