package engine

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VoxelId uint16

/* Voxel geometry vertex data type */
type VoxelVertex struct {
    X, Y, Z     uint8 // Vertex position relative to chunk
    Nx, Ny, Nz  int8  // Normal vector
    Tx, Ty      uint8 // Tile tex coords
}

/* List of voxel verticies. Can be passed to VertexArray.Buffer */
type VoxelVertices []VoxelVertex

func (buffer VoxelVertices) Elements() int { return len(buffer) }
func (buffer VoxelVertices) Size()     int { return 8 }
func (buffer VoxelVertices) GLPtr()    unsafe.Pointer { return gl.Ptr(buffer) }

/* Voxel preset data type */
type Voxel struct {
    Id      VoxelId /* Voxel type id */
    Xp, Xn  TileId  /* Right, left tile ids */
    Yp, Yn  TileId  /* Up, down tile ids */
    Zp, Zn  TileId  /* Front, back tile ids */
}

func (voxel *Voxel) Compute(x, y, z uint8, data VoxelVertices, ts *Tileset) VoxelVertices {
    right, left, up, down, front, back := ts.Get(voxel.Xp), ts.Get(voxel.Xn),
                                          ts.Get(voxel.Yp), ts.Get(voxel.Yn),
                                          ts.Get(voxel.Zp), ts.Get(voxel.Zn)

    ux, uy := uint8(up.X),    uint8(up.Y)
    dx, dy := uint8(down.X),  uint8(down.Y)
    fx, fy := uint8(front.X), uint8(front.Y)
    bx, by := uint8(back.X),  uint8(back.Y)
    lx, ly := uint8(left.X),  uint8(left.Y)
    rx, ry := uint8(right.X), uint8(right.Y)

    return append(data,
        // Top 
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx:  0, Ny:  1, Nz:  0, Tx: ux+0, Ty: uy+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  1, Nz:  0, Tx: ux+0, Ty: uy+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  0, Ny:  1, Nz:  0, Tx: ux+1, Ty: uy+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  0, Ny:  1, Nz:  0, Tx: ux+1, Ty: uy+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  1, Nz:  0, Tx: ux+0, Ty: uy+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+1, Nx:  0, Ny:  1, Nz:  0, Tx: ux+1, Ty: uy+1, },
        // Bottom 
        VoxelVertex { X: x+0, Y: y+0, Z: z+0, Nx:  0, Ny: -1, Nz:  0, Tx: dx+0, Ty: dy+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny: -1, Nz:  0, Tx: dx+1, Ty: dy+0, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx:  0, Ny: -1, Nz:  0, Tx: dx+0, Ty: dy+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny: -1, Nz:  0, Tx: dx+1, Ty: dy+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  0, Ny: -1, Nz:  0, Tx: dx+1, Ty: dy+1, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx:  0, Ny: -1, Nz:  0, Tx: dx+0, Ty: dy+1, },
        // Front
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: fx+1, Ty: fy+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: fx+0, Ty: fy+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: fx+1, Ty: fy+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: fx+0, Ty: fy+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: fx+0, Ty: fy+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: fx+1, Ty: fy+1, },
        // Back
        VoxelVertex { X: x+0, Y: y+0, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: bx+0, Ty: by+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: bx+0, Ty: by+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: bx+1, Ty: by+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: bx+1, Ty: by+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: bx+0, Ty: by+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: bx+1, Ty: by+1, },
        // Left
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx: -1, Ny:  0, Nz:  0, Tx: lx+0, Ty: ly+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx: -1, Ny:  0, Nz:  0, Tx: lx+1, Ty: ly+0, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+0, Nx: -1, Ny:  0, Nz:  0, Tx: lx+0, Ty: ly+0, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx: -1, Ny:  0, Nz:  0, Tx: lx+0, Ty: ly+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx: -1, Ny:  0, Nz:  0, Tx: lx+1, Ty: ly+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx: -1, Ny:  0, Nz:  0, Tx: lx+1, Ty: ly+0, },
        // Right
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  1, Ny:  0, Nz:  0, Tx: rx+1, Ty: ry+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  1, Ny:  0, Nz:  0, Tx: rx+1, Ty: ry+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  1, Ny:  0, Nz:  0, Tx: rx+0, Ty: ry+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  1, Ny:  0, Nz:  0, Tx: rx+1, Ty: ry+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  1, Ny:  0, Nz:  0, Tx: rx+0, Ty: ry+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+1, Nx:  1, Ny:  0, Nz:  0, Tx: rx+0, Ty: ry+1, })
}
