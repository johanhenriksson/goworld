package engine

import (
    "fmt"
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VoxelVertex struct {
    X, Y, Z     uint8 // Chunk Position
    Nx, Ny, Nz  int8  // Normal id
    Tx, Ty      uint8 // Tile coords
}

type VoxelVertices []VoxelVertex

func (buffer VoxelVertices) Elements() int {
    return len(buffer)
}

func (buffer VoxelVertices) Size() int {
    return 8
}

func (buffer VoxelVertices) GLPtr() unsafe.Pointer {
    return gl.Ptr(buffer)
}

func GenerateVoxel(x, y, z uint8, voxel *Voxel, tileset *Tileset) VoxelVertices {
    top, bottom, front, back, left, right := tileset.Get(voxel.Yp), tileset.Get(voxel.Yn),
                                             tileset.Get(voxel.Zp), tileset.Get(voxel.Zn),
                                             tileset.Get(voxel.Xn), tileset.Get(voxel.Xp)
    v := &VoxelVertex { }
    fmt.Println("Vertex actual size:", unsafe.Sizeof(v))
    return VoxelVertices {
        // Top 
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx:  0, Ny:  1, Nz:  0, Tx: top.X+0, Ty: top.Y+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  1, Nz:  0, Tx: top.X+0, Ty: top.Y+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  0, Ny:  1, Nz:  0, Tx: top.X+1, Ty: top.Y+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  0, Ny:  1, Nz:  0, Tx: top.X+1, Ty: top.Y+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  1, Nz:  0, Tx: top.X+0, Ty: top.Y+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+1, Nx:  0, Ny:  1, Nz:  0, Tx: top.X+1, Ty: top.Y+1, },
        // Bottom 
        VoxelVertex { X: x+0, Y: y+0, Z: z+0, Nx:  0, Ny: -1, Nz:  0, Tx: bottom.X+0, Ty: bottom.Y+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny: -1, Nz:  0, Tx: bottom.X+1, Ty: bottom.Y+0, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx:  0, Ny: -1, Nz:  0, Tx: bottom.X+0, Ty: bottom.Y+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny: -1, Nz:  0, Tx: bottom.X+1, Ty: bottom.Y+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  0, Ny: -1, Nz:  0, Tx: bottom.X+1, Ty: bottom.Y+1, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx:  0, Ny: -1, Nz:  0, Tx: bottom.X+0, Ty: bottom.Y+1, },
        // Front
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: front.X+1, Ty: front.Y+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: front.X+0, Ty: front.Y+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: front.X+1, Ty: front.Y+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: front.X+0, Ty: front.Y+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: front.X+0, Ty: front.Y+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx:  0, Ny:  0, Nz:  1, Tx: front.X+1, Ty: front.Y+1, },
        // Back
        VoxelVertex { X: x+0, Y: y+0, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: back.X+0, Ty: back.Y+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: back.X+0, Ty: back.Y+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: back.X+1, Ty: back.Y+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: back.X+1, Ty: back.Y+0, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: back.X+0, Ty: back.Y+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  0, Ny:  0, Nz: -1, Tx: back.X+1, Ty: back.Y+1, },
        // Left
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx: -1, Ny:  0, Nz:  0, Tx: left.X+0, Ty: left.Y+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx: -1, Ny:  0, Nz:  0, Tx: left.X+1, Ty: left.Y+0, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+0, Nx: -1, Ny:  0, Nz:  0, Tx: left.X+0, Ty: left.Y+0, },
        VoxelVertex { X: x+0, Y: y+0, Z: z+1, Nx: -1, Ny:  0, Nz:  0, Tx: left.X+0, Ty: left.Y+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+1, Nx: -1, Ny:  0, Nz:  0, Tx: left.X+1, Ty: left.Y+1, },
        VoxelVertex { X: x+0, Y: y+1, Z: z+0, Nx: -1, Ny:  0, Nz:  0, Tx: left.X+1, Ty: left.Y+0, },
        // Right
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  1, Ny:  0, Nz:  0, Tx: right.X+1, Ty: right.Y+1, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+0, Nx:  1, Ny:  0, Nz:  0, Tx: right.X+1, Ty: right.Y+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  1, Ny:  0, Nz:  0, Tx: right.X+0, Ty: right.Y+0, },
        VoxelVertex { X: x+1, Y: y+0, Z: z+1, Nx:  1, Ny:  0, Nz:  0, Tx: right.X+1, Ty: right.Y+1, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+0, Nx:  1, Ny:  0, Nz:  0, Tx: right.X+0, Ty: right.Y+0, },
        VoxelVertex { X: x+1, Y: y+1, Z: z+1, Nx:  1, Ny:  0, Nz:  0, Tx: right.X+0, Ty: right.Y+1, },
    }
}
