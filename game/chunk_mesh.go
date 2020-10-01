package game

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
)

type ChunkMesh struct {
	*engine.Mesh
	*Chunk
	meshComputed chan VoxelVertices
}

func NewChunkMesh(parent *engine.Object, chunk *Chunk) *ChunkMesh {
	mesh := engine.NewMesh(parent, assets.GetMaterialCached("color_voxels"))
	chk := &ChunkMesh{
		Mesh:         mesh,
		Chunk:        chunk,
		meshComputed: make(chan VoxelVertices),
	}
	parent.Attach(chk)
	return chk
}

func (cm *ChunkMesh) Update(dt float32) {
	select {
	case newMesh := <-cm.meshComputed:
		cm.Buffer("geometry", newMesh)
	default:
	}
}

// Queues recomputation of the mesh
func (cm *ChunkMesh) Compute() {
	go func() {
		data := cm.computeVertexData()
		cm.meshComputed <- data
	}()
}

func (cm *ChunkMesh) computeVertexData() VoxelVertices {
	data := make(VoxelVertices, 0, 64)
	light := cm.Light.Brightness
	Omax := float32(220)

	for z := -1; z <= cm.Sz; z++ {
		for x := -1; x <= cm.Sx; x++ {
			for y := -1; y <= cm.Sy; y++ {
				v := cm.At(x, y, z)
				if v != EmptyVoxel {
					// consider ONLY empty voxels
					continue
				}

				xp := cm.At(x+1, y, z)
				xn := cm.At(x-1, y, z)
				yp := cm.At(x, y+1, z)
				yn := cm.At(x, y-1, z)
				zp := cm.At(x, y, z+1)
				zn := cm.At(x, y, z-1)
				xpf := xp != EmptyVoxel
				xnf := xn != EmptyVoxel
				ypf := yp != EmptyVoxel
				ynf := yn != EmptyVoxel
				zpf := zp != EmptyVoxel
				znf := zn != EmptyVoxel

				l := light(x, y, z)

				if xpf || xnf {
					lyp := light(x, y+1, z)
					lyn := light(x, y-1, z)
					lzp := light(x, y, z+1)
					lzn := light(x, y, z-1)
					lypzn := light(x, y+1, z-1)
					lypzp := light(x, y+1, z+1)
					lynzn := light(x, y-1, z-1)
					lynzp := light(x, y-1, z+1)

					if xpf {
						// xp is empty - tesselate square with X- normal
						n := byte(2)
						v1 := VoxelVertex{
							X: byte(x + 1), Y: byte(y + 1), Z: byte(z + 1), N: n,
							R: xp.R, G: xp.G, B: xp.B,
							O: byte(Omax * (1 - (lzp+lyp+lypzp+l)/4)),
						}
						if ypf && zpf {
							v1.O = byte(Omax * (1 - (lzp+lyp+l)/3))
						}
						v2 := VoxelVertex{
							X: byte(x + 1), Y: byte(y + 1), Z: byte(z), N: n,
							R: xp.R, G: xp.G, B: xp.B,
							O: byte(Omax * (1 - (lzn+lyp+lypzn+l)/4)),
						}
						if ypf && znf {
							v2.O = byte(Omax * (1 - (lzn+lyp+l)/3))
						}
						v3 := VoxelVertex{
							X: byte(x + 1), Y: byte(y), Z: byte(z + 1), N: n,
							R: xp.R, G: xp.G, B: xp.B,
							O: byte(Omax * (1 - (lzp+lyn+lynzp+l)/4)),
						}
						if ynf && zpf {
							v3.O = byte(Omax * (1 - (lzp+lyn+l)/3))
						}
						v4 := VoxelVertex{
							X: byte(x + 1), Y: byte(y), Z: byte(z), N: n,
							R: xp.R, G: xp.G, B: xp.B,
							O: byte(Omax * (1 - (lzn+lyn+lynzn+l)/4)),
						}
						if ynf && znf {
							v4.O = byte(Omax * (1 - (lzn+lyn+l)/3))
						}
						data = append(data, v2, v3, v1, v4, v3, v2)
					}

					if xnf {
						// xn is empty - tesselate square with x+ normal
						n := byte(1)
						v1 := VoxelVertex{
							X: byte(x), Y: byte(y + 1), Z: byte(z), N: n,
							R: xn.R, G: xn.G, B: xn.B,
							O: byte(Omax * (1 - (lyp+lzn+lypzn+l)/4)),
						}
						if ypf && znf {
							v1.O = byte(Omax * (1 - (lyp+lzn+l)/3))
						}
						v2 := VoxelVertex{
							X: byte(x), Y: byte(y + 1), Z: byte(z + 1), N: n,
							R: xn.R, G: xn.G, B: xn.B,
							O: byte(Omax * (1 - (lyp+lzp+lypzp+l)/4)),
						}
						if ypf && zpf {
							v2.O = byte(Omax * (1 - (lyp+lzp+l)/3))
						}
						v3 := VoxelVertex{
							X: byte(x), Y: byte(y), Z: byte(z), N: n,
							R: xn.R, G: xn.G, B: xn.B,
							O: byte(Omax * (1 - (lyn+lzn+lynzn+l)/4)),
						}
						if ynf && znf {
							v3.O = byte(Omax * (1 - (lyn+lzn+l)/3))
						}
						v4 := VoxelVertex{
							X: byte(x), Y: byte(y), Z: byte(z + 1), N: n,
							R: xn.R, G: xn.G, B: xn.B,
							O: byte(Omax * (1 - (lyn+lzp+lynzp+l)/4)),
						}
						if ynf && zpf {
							v4.O = byte(Omax * (1 - (lyn+lzp+l)/3))
						}
						data = append(data, v2, v3, v1, v4, v3, v2)
					}
				}

				if ypf || ynf {
					lxp := light(x+1, y, z)
					lxn := light(x-1, y, z)
					lzp := light(x, y, z+1)
					lzn := light(x, y, z-1)
					lxpzp := light(x+1, y, z+1)
					lxpzn := light(x+1, y, z-1)
					lxnzp := light(x-1, y, z+1)
					lxnzn := light(x-1, y, z-1)

					if ypf {
						n := byte(3) // YN
						v1 := VoxelVertex{
							X: byte(x + 1), Y: byte(y + 1), Z: byte(z + 1), N: n,
							R: yp.R, G: yp.G, B: yp.B,
							O: byte(Omax * (1 - (lxp+lzp+lxpzp+l)/4)),
						}
						if xpf && zpf {
							v1.O = byte(Omax * (1 - (lxp+lzp+l)/3))
						}
						v2 := VoxelVertex{
							X: byte(x + 1), Y: byte(y + 1), Z: byte(z), N: n,
							R: yp.R, G: yp.G, B: yp.B,
							O: byte(Omax * (1 - (lxp+lzn+lxpzn+l)/4)),
						}
						if xpf && znf {
							v2.O = byte(Omax * (1 - (lxp+lzn+l)/3))
						}
						v3 := VoxelVertex{
							X: byte(x), Y: byte(y + 1), Z: byte(z + 1), N: n,
							R: yp.R, G: yp.G, B: yp.B,
							O: byte(Omax * (1 - (lxn+lzp+lxnzp+l)/4)),
						}
						if xnf && zpf {
							v3.O = byte(Omax * (1 - (lxn+lzp+l)/3))
						}
						v4 := VoxelVertex{
							X: byte(x), Y: byte(y + 1), Z: byte(z), N: n,
							R: yp.R, G: yp.G, B: yp.B,
							O: byte(Omax * (1 - (lxn+lzn+lxnzn+l)/4)),
						}
						if xnf && znf {
							v4.O = byte(Omax * (1 - (lxn+lzn+l)/3))
						}
						data = append(data, v1, v3, v2, v2, v3, v4)
					}

					if ynf {
						// Y-1 is filled, add quad with Y+ normal
						n := byte(3) // YP
						v1 := VoxelVertex{
							X: byte(x + 1), Y: byte(y), Z: byte(z + 1), N: n,
							R: yn.R, G: yn.G, B: yn.B,
							O: byte(Omax * (1 - (lxp+lzp+lxpzp+l)/4)),
						}
						if xpf && zpf {
							v1.O = byte(Omax * (1 - l/3))
						}
						v2 := VoxelVertex{
							X: byte(x + 1), Y: byte(y), Z: byte(z), N: n,
							R: yn.R, G: yn.G, B: yn.B,
							O: byte(Omax * (1 - (lxp+lzn+lxpzn+l)/4)),
						}
						if xpf && znf {
							v2.O = byte(Omax * (1 - l/3))
						}
						v3 := VoxelVertex{
							X: byte(x), Y: byte(y), Z: byte(z + 1), N: n,
							R: yn.R, G: yn.G, B: yn.B,
							O: byte(Omax * (1 - (lxn+lzp+lxnzp+l)/4)),
						}
						if xnf && zpf {
							v3.O = byte(Omax * (1 - l/3))
						}
						v4 := VoxelVertex{
							X: byte(x), Y: byte(y), Z: byte(z), N: n,
							R: yn.R, G: yn.G, B: yn.B,
							O: byte(Omax * (1 - (lxn+lzn+lxnzn+l)/4)),
						}
						if xnf && znf {
							v4.O = byte(Omax * (1 - l/3))
						}
						data = append(data, v2, v3, v1, v4, v3, v2)
					}
				}

				if zpf || znf {
					lxp := light(x+1, y, z)
					lxn := light(x-1, y, z)
					lyp := light(x, y+1, z)
					lyn := light(x, y-1, z)
					lxnyp := light(x-1, y+1, z)
					lxpyp := light(x+1, y+1, z)
					lxnyn := light(x-1, y-1, z)
					lxpyn := light(x+1, y-1, z)

					if zpf {
						// zp is empty - tesselate square with ZN normal
						n := byte(6)
						v1 := VoxelVertex{
							X: byte(x), Y: byte(y + 1), Z: byte(z + 1), N: n,
							R: zp.R, G: zp.G, B: zp.B,
							O: byte(Omax * (1 - (lxn+lyp+lxnyp+l)/4)),
						}
						if xnf && ypf {
							v1.O = byte(Omax * (1 - (lxn+lyp+l)/3))
						}
						v2 := VoxelVertex{
							X: byte(x + 1), Y: byte(y + 1), Z: byte(z + 1), N: n,
							R: zp.R, G: zp.G, B: zp.B,
							O: byte(Omax * (1 - (lxp+lyp+lxpyp+l)/4)),
						}
						if xpf && ypf {
							v2.O = byte(Omax * (1 - (lxp+lyp+l)/3))
						}
						v3 := VoxelVertex{
							X: byte(x), Y: byte(y), Z: byte(z + 1), N: n,
							R: zp.R, G: zp.G, B: zp.B,
							O: byte(Omax * (1 - (lxn+lyn+lxnyn+l)/4)),
						}
						if xnf && ynf {
							v3.O = byte(Omax * (1 - (lxn+lyn+l)/3))
						}
						v4 := VoxelVertex{
							X: byte(x + 1), Y: byte(y), Z: byte(z + 1), N: n,
							R: zp.R, G: zp.G, B: zp.B,
							O: byte(Omax * (1 - (lxp+lyn+lxpyn+l)/4)),
						}
						if xpf && ynf {
							v4.O = byte(Omax * (1 - (lxp+lyn+l)/3))
						}
						data = append(data, v2, v3, v1, v4, v3, v2)
					}

					if znf {
						// zn is empty - tesselate square with ZP normal
						n := byte(5)
						v1 := VoxelVertex{
							X: byte(x), Y: byte(y + 1), Z: byte(z), N: n,
							R: zn.R, G: zn.G, B: zn.B,
							O: byte(Omax * (1 - (lxn+lyp+lxnyp+l)/4)),
						}
						if xnf && ypf {
							v1.O = byte(Omax * (1 - (lxn+lyp+l)/3))
						}
						v2 := VoxelVertex{
							X: byte(x + 1), Y: byte(y + 1), Z: byte(z), N: n,
							R: zn.R, G: zn.G, B: zn.B,
							O: byte(Omax * (1 - (lxp+lyp+lxpyp+l)/4)),
						}
						if xpf && ypf {
							v2.O = byte(Omax * (1 - (lxp+lyp+l)/3))
						}
						v3 := VoxelVertex{
							X: byte(x), Y: byte(y), Z: byte(z), N: n,
							R: zn.R, G: zn.G, B: zn.B,
							O: byte(Omax * (1 - (lxn+lyn+lxnyn+l)/4)),
						}
						if xnf && ynf {
							v3.O = byte(Omax * (1 - (lxn+lyn+l)/3))
						}
						v4 := VoxelVertex{
							X: byte(x + 1), Y: byte(y), Z: byte(z), N: n,
							R: zn.R, G: zn.G, B: zn.B,
							O: byte(Omax * (1 - (lxp+lyn+lxpyn+l)/4)),
						}
						if xpf && ynf {
							v4.O = byte(Omax * (1 - (lxp+lyn+l)/3))
						}
						data = append(data, v1, v3, v2, v2, v3, v4)
					}
				}
			}
		}
	}
	return data
}
