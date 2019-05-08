package game

type VoxelId uint16

/* Voxel geometry vertex data type */
type VoxelVertex struct {
	X, Y, Z    uint8 // Vertex position relative to chunk
	Nx, Ny, Nz int8  // Normal vector
	Tx, Ty     uint8 // Tile tex coords
}

/* List of voxel verticies. Can be passed to VertexArray.Buffer */
type VoxelVertices []VoxelVertex

func (buffer VoxelVertices) Elements() int { return len(buffer) }
func (buffer VoxelVertices) Size() int     { return 8 }

/* Voxel preset data type */
type Voxel struct {
	Id     VoxelId /* Voxel type id */
	Xp, Xn TileId  /* Right, left tile ids */
	Yp, Yn TileId  /* Up, down tile ids */
	Zp, Zn TileId  /* Front, back tile ids */
}

/* Computes vertex data for this voxel.
   Parameters: X,Y,Z - chunk position
               xp, xn, yp, yn, zp, zn - draw face flags (right, left, up, down, front, back)
               data - output slice
               ts   - chunk tileset pointer */
func (voxel *Voxel) Compute(data VoxelVertices, x, y, z uint8, xp, xn, yp, yn, zp, zn bool, ts *Tileset) VoxelVertices {

	// Right (X+)
	if xp {
		right := ts.Get(voxel.Xp)
		rx, ry := uint8(right.X), uint8(right.Y) // right tile coords
		data = append(data,
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, Nx: 1, Ny: 0, Nz: 0, Tx: rx + 1, Ty: ry + 1},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, Nx: 1, Ny: 0, Nz: 0, Tx: rx + 1, Ty: ry + 0},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, Nx: 1, Ny: 0, Nz: 0, Tx: rx + 0, Ty: ry + 0},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, Nx: 1, Ny: 0, Nz: 0, Tx: rx + 1, Ty: ry + 1},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, Nx: 1, Ny: 0, Nz: 0, Tx: rx + 0, Ty: ry + 0},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, Nx: 1, Ny: 0, Nz: 0, Tx: rx + 0, Ty: ry + 1})
	}

	// Left faces (X-)
	if xn {
		left := ts.Get(voxel.Xn)
		lx, ly := uint8(left.X), uint8(left.Y) // left tile coords
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, Nx: -1, Ny: 0, Nz: 0, Tx: lx + 0, Ty: ly + 1},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, Nx: -1, Ny: 0, Nz: 0, Tx: lx + 1, Ty: ly + 0},
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, Nx: -1, Ny: 0, Nz: 0, Tx: lx + 0, Ty: ly + 0},
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, Nx: -1, Ny: 0, Nz: 0, Tx: lx + 0, Ty: ly + 1},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, Nx: -1, Ny: 0, Nz: 0, Tx: lx + 1, Ty: ly + 1},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, Nx: -1, Ny: 0, Nz: 0, Tx: lx + 1, Ty: ly + 0})
	}

	// Top faces (Y+)
	if yp {
		up := ts.Get(voxel.Yp)
		ux, uy := uint8(up.X), uint8(up.Y) // top tile coords

		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, Nx: 0, Ny: 1, Nz: 0, Tx: ux + 0, Ty: uy + 0},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, Nx: 0, Ny: 1, Nz: 0, Tx: ux + 0, Ty: uy + 1},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, Nx: 0, Ny: 1, Nz: 0, Tx: ux + 1, Ty: uy + 0},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, Nx: 0, Ny: 1, Nz: 0, Tx: ux + 1, Ty: uy + 0},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, Nx: 0, Ny: 1, Nz: 0, Tx: ux + 0, Ty: uy + 1},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, Nx: 0, Ny: 1, Nz: 0, Tx: ux + 1, Ty: uy + 1})
	}

	// Bottom faces (Y-)
	if yn {
		down := ts.Get(voxel.Yn)
		dx, dy := uint8(down.X), uint8(down.Y) // bottom tile coords
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, Nx: 0, Ny: -1, Nz: 0, Tx: dx + 0, Ty: dy + 0},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, Nx: 0, Ny: -1, Nz: 0, Tx: dx + 1, Ty: dy + 0},
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, Nx: 0, Ny: -1, Nz: 0, Tx: dx + 0, Ty: dy + 1},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, Nx: 0, Ny: -1, Nz: 0, Tx: dx + 1, Ty: dy + 0},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, Nx: 0, Ny: -1, Nz: 0, Tx: dx + 1, Ty: dy + 1},
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, Nx: 0, Ny: -1, Nz: 0, Tx: dx + 0, Ty: dy + 1})
	}

	// Front faces (Z+)
	if zp {
		front := ts.Get(voxel.Zp)
		fx, fy := uint8(front.X), uint8(front.Y) // front tile coords

		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, Nx: 0, Ny: 0, Nz: 1, Tx: fx + 1, Ty: fy + 0},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, Nx: 0, Ny: 0, Nz: 1, Tx: fx + 0, Ty: fy + 0},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, Nx: 0, Ny: 0, Nz: 1, Tx: fx + 1, Ty: fy + 1},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, Nx: 0, Ny: 0, Nz: 1, Tx: fx + 0, Ty: fy + 0},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, Nx: 0, Ny: 0, Nz: 1, Tx: fx + 0, Ty: fy + 1},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, Nx: 0, Ny: 0, Nz: 1, Tx: fx + 1, Ty: fy + 1})
	}

	// Back faces (Z-)
	if zn {
		back := ts.Get(voxel.Zn)
		bx, by := uint8(back.X), uint8(back.Y) // back tile coords
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, Nx: 0, Ny: 0, Nz: -1, Tx: bx + 0, Ty: by + 0},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, Nx: 0, Ny: 0, Nz: -1, Tx: bx + 0, Ty: by + 1},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, Nx: 0, Ny: 0, Nz: -1, Tx: bx + 1, Ty: by + 0},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, Nx: 0, Ny: 0, Nz: -1, Tx: bx + 1, Ty: by + 0},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, Nx: 0, Ny: 0, Nz: -1, Tx: bx + 0, Ty: by + 1},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, Nx: 0, Ny: 0, Nz: -1, Tx: bx + 1, Ty: by + 1})
	}

	return data
}
