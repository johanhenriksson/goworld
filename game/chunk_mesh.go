package game

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
)

type ChunkMesh struct {
	*engine.Mesh
	*Chunk
}

func NewChunkMesh(object *engine.Object, chunk *Chunk) *ChunkMesh {
	mesh := engine.NewMesh(assets.GetMaterialCached("color_voxels"))
	chk := &ChunkMesh{
		Mesh:  mesh,
		Chunk: chunk,
	}
	chk.ComponentBase = engine.NewComponent(object, chk)
	return chk
}

// Computes the chunk mesh and returns a pointer to it.
func (cm *ChunkMesh) Compute() {
	data := make(VoxelVertices, 0, 64)

	light := NewLightVolume(cm.Sx, cm.Sy, cm.Sz)
	for z := 0; z < light.Sz; z++ {
		for y := 0; y < light.Sy; y++ {
			for x := 0; x < light.Sx; x++ {
				if !cm.Free(x, y, z) {
					light.Block(x, y, z)
				}
			}
		}
	}
	light.Calculate()

	/* geometry pass */
	// for z := 0; z < cm.Sz; z++ {
	// 	for y := 0; y < cm.Sy; y++ {
	// 		for x := 0; x < cm.Sx; x++ {
	// 			v := cm.At(x, y, z)
	// 			if v == EmptyVoxel {
	// 				continue
	// 			}

	// 			/* Simple optimization - dont draw hidden faces */
	// 			xp := cm.At(x+1, y, z) == EmptyVoxel
	// 			xn := cm.At(x-1, y, z) == EmptyVoxel
	// 			yp := cm.At(x, y+1, z) == EmptyVoxel
	// 			yn := cm.At(x, y-1, z) == EmptyVoxel
	// 			zp := cm.At(x, y, z+1) == EmptyVoxel
	// 			zn := cm.At(x, y, z-1) == EmptyVoxel

	// 			/* Compute & append vertex data */
	// 			vertices := v.Compute(light, byte(x), byte(y), byte(z), xp, xn, yp, yn, zp, zn)
	// 			data = append(data, vertices...)
	// 		}
	// 	}
	// }

	for z := 0; z < cm.Sz; z++ {
		for y := 0; y < cm.Sy; y++ {
			for x := 0; x < cm.Sx; x++ {
				v := cm.At(x, y, z)
				if v != EmptyVoxel {
					// consider ONLY empty voxels
					continue
				}

				l := light.Brightness(x, y, z)
				lxp := light.Brightness(x+1, y, z)
				lxn := light.Brightness(x-1, y, z)
				lyp := light.Brightness(x, y+1, z)
				lyn := light.Brightness(x, y-1, z)
				lzp := light.Brightness(x, y, z+1)
				lzn := light.Brightness(x, y, z-1)

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

				lypzn := light.Brightness(x, y+1, z-1)
				lypzp := light.Brightness(x, y+1, z+1)
				lynzn := light.Brightness(x, y-1, z-1)
				lynzp := light.Brightness(x, y-1, z+1)

				if xpf {
					// xp is empty - tesselate square with X- normal
					n := byte(2)
					v1 := VoxelVertex{
						X: byte(x + 1), Y: byte(y + 1), Z: byte(z + 1), N: n,
						R: xp.R, G: xp.G, B: xp.B,
						O: byte(255 * (1 - (lzp+lyp+lypzp+l)/4)),
					}
					if ypf && zpf {
						v1.O = byte(255 * (1 - (lzp+lyp+l)/3))
					}
					v2 := VoxelVertex{
						X: byte(x + 1), Y: byte(y + 1), Z: byte(z), N: n,
						R: xp.R, G: xp.G, B: xp.B,
						O: byte(255 * (1 - (lzn+lyp+lypzn+l)/4)),
					}
					if ypf && znf {
						v2.O = byte(255 * (1 - (lzn+lyp+l)/3))
					}
					v3 := VoxelVertex{
						X: byte(x + 1), Y: byte(y), Z: byte(z + 1), N: n,
						R: xp.R, G: xp.G, B: xp.B,
						O: byte(255 * (1 - (lzp+lyn+lynzp+l)/4)),
					}
					if ynf && zpf {
						v3.O = byte(255 * (1 - (lzp+lyn+l)/4))
					}
					v4 := VoxelVertex{
						X: byte(x + 1), Y: byte(y), Z: byte(z), N: n,
						R: xp.R, G: xp.G, B: xp.B,
						O: byte(255 * (1 - (lzn+lyn+lynzn+l)/4)),
					}
					if ynf && znf {
						v4.O = byte(255 * (1 - (lzn+lyn+l)/3))
					}
					data = append(data, v2, v3, v1, v4, v3, v2)
				}

				if xnf {
					// xn is empty - tesselate square with x+ normal
					n := byte(1)
					v1 := VoxelVertex{
						X: byte(x), Y: byte(y + 1), Z: byte(z), N: n,
						R: xn.R, G: xn.G, B: xn.B,
						O: byte(255 * (1 - (lyp+lzn+lypzn+l)/4)),
					}
					if ypf && znf {
						v1.O = byte(255 * (1 - (lyp+lzn+l)/3))
					}
					v2 := VoxelVertex{
						X: byte(x), Y: byte(y + 1), Z: byte(z + 1), N: n,
						R: xn.R, G: xn.G, B: xn.B,
						O: byte(255 * (1 - (lyp+lzp+lypzp+l)/4)),
					}
					if ypf && zpf {
						v2.O = byte(255 * (1 - (lyp+lzp+l)/3))
					}
					v3 := VoxelVertex{
						X: byte(x), Y: byte(y), Z: byte(z), N: n,
						R: xn.R, G: xn.G, B: xn.B,
						O: byte(255 * (1 - (lyn+lzn+lynzn+l)/4)),
					}
					if ynf && znf {
						v3.O = byte(255 * (1 - (lyn+lzn+l)/3))
					}
					v4 := VoxelVertex{
						X: byte(x), Y: byte(y), Z: byte(z + 1), N: n,
						R: xn.R, G: xn.G, B: xn.B,
						O: byte(255 * (1 - (lyn+lzp+lynzp+l)/4)),
					}
					if ynf && zpf {
						v4.O = byte(255 * (1 - (lyn+lzp+l)/3))
					}
					data = append(data, v2, v3, v1, v4, v3, v2)
				}

				lxpzp := light.Brightness(x+1, y, z+1)
				lxpzn := light.Brightness(x+1, y, z-1)
				lxnzp := light.Brightness(x-1, y, z+1)
				lxnzn := light.Brightness(x-1, y, z-1)

				if ypf {
					n := byte(3) // YN
					v1 := VoxelVertex{
						X: byte(x + 1), Y: byte(y + 1), Z: byte(z + 1), N: n,
						R: yp.R, G: yp.G, B: yp.B,
						O: byte(255 * (1 - (lxp+lzp+lxpzp+l)/4)),
					}
					if xpf && zpf {
						v1.O = byte(255 * (1 - (lxp+lzp+l)/3))
					}
					v2 := VoxelVertex{
						X: byte(x + 1), Y: byte(y + 1), Z: byte(z), N: n,
						R: yp.R, G: yp.G, B: yp.B,
						O: byte(255 * (1 - (lxp+lzn+lxpzn+l)/4)),
					}
					if xpf && znf {
						v2.O = byte(255 * (1 - (lxp+lzn+l)/3))
					}
					v3 := VoxelVertex{
						X: byte(x), Y: byte(y + 1), Z: byte(z + 1), N: n,
						R: yp.R, G: yp.G, B: yp.B,
						O: byte(255 * (1 - (lxn+lzp+lxnzp+l)/4)),
					}
					if xnf && zpf {
						v3.O = byte(255 * (1 - (lxn+lzp+l)/3))
					}
					v4 := VoxelVertex{
						X: byte(x), Y: byte(y + 1), Z: byte(z), N: n,
						R: yp.R, G: yp.G, B: yp.B,
						O: byte(255 * (1 - (lxn+lzn+lxnzn+l)/4)),
					}
					if xnf && znf {
						v4.O = byte(255 * (1 - (lxn+lzn+l)/3))
					}
					data = append(data, v1, v3, v2, v2, v3, v4)
				}

				if ynf {
					// Y-1 is filled, add quad with Y+ normal
					n := byte(3) // YP
					v1 := VoxelVertex{
						X: byte(x + 1), Y: byte(y), Z: byte(z + 1), N: n,
						R: yn.R, G: yn.G, B: yn.B,
						O: byte(255 * (1 - (lxp+lzp+lxpzp+l)/4)),
					}
					if xpf && zpf {
						v1.O = byte(255 * (1 - (lxp+lzp+l)/3))
					}
					v2 := VoxelVertex{
						X: byte(x + 1), Y: byte(y), Z: byte(z), N: n,
						R: yn.R, G: yn.G, B: yn.B,
						O: byte(255 * (1 - (lxp+lzn+lxpzn+l)/4)),
					}
					if xpf && znf {
						v2.O = byte(255 * (1 - (lxp+lzn+l)/3))
					}
					v3 := VoxelVertex{
						X: byte(x), Y: byte(y), Z: byte(z + 1), N: n,
						R: yn.R, G: yn.G, B: yn.B,
						O: byte(255 * (1 - (lxn+lzp+lxnzp+l)/4)),
					}
					if xnf && zpf {
						v3.O = byte(255 * (1 - (lxn+lzp+l)/3))
					}
					v4 := VoxelVertex{
						X: byte(x), Y: byte(y), Z: byte(z), N: n,
						R: yn.R, G: yn.G, B: yn.B,
						O: byte(255 * (1 - (lxn+lzn+lxnzn+l)/4)),
					}
					if xnf && znf {
						v4.O = byte(255 * (1 - (lxn+lzn+l)/3))
					}
					data = append(data, v2, v3, v1, v4, v3, v2)
				}

				lxnyp := light.Brightness(x-1, y+1, z)
				lxpyp := light.Brightness(x+1, y+1, z)
				lxnyn := light.Brightness(x-1, y-1, z)
				lxpyn := light.Brightness(x+1, y-1, z)

				if zpf {
					// zp is empty - tesselate square with ZN normal
					n := byte(3)
					v1 := VoxelVertex{
						X: byte(x), Y: byte(y + 1), Z: byte(z + 1), N: n,
						R: zp.R, G: zp.G, B: zp.B,
						O: byte(255 * (1 - (lxn+lyp+lxnyp+l)/4)),
					}
					if xnf && ypf {
						v1.O = byte(255 * (1 - (lxn+lyp+l)/3))
					}
					v2 := VoxelVertex{
						X: byte(x + 1), Y: byte(y + 1), Z: byte(z + 1), N: n,
						R: zp.R, G: zp.G, B: zp.B,
						O: byte(255 * (1 - (lxp+lyp+lxpyp+l)/4)),
					}
					if xpf && ypf {
						v2.O = byte(255 * (1 - (lxp+lyp+l)/3))
					}
					v3 := VoxelVertex{
						X: byte(x), Y: byte(y), Z: byte(z + 1), N: n,
						R: zp.R, G: zp.G, B: zp.B,
						O: byte(255 * (1 - (lxn+lyn+lxnyn+l)/4)),
					}
					if xnf && ynf {
						v3.O = byte(255 * (1 - (lxn+lyn+l)/3))
					}
					v4 := VoxelVertex{
						X: byte(x + 1), Y: byte(y), Z: byte(z + 1), N: n,
						R: zp.R, G: zp.G, B: zp.B,
						O: byte(255 * (1 - (lxp+lyn+lxpyn+l)/4)),
					}
					if xpf && ynf {
						v4.O = byte(255 * (1 - (lxp+lyn+l)/3))
					}
					data = append(data, v2, v3, v1, v4, v3, v2)
				}

				if znf {
					// zn is empty - tesselate square with ZP normal
					n := byte(5)
					v1 := VoxelVertex{
						X: byte(x), Y: byte(y + 1), Z: byte(z), N: n,
						R: zn.R, G: zn.G, B: zn.B,
						O: byte(255 * (1 - (lxn+lyp+lxnyp+l)/4)),
					}
					if xnf && ypf {
						v1.O = byte(255 * (1 - (lxn+lyp+l)/3))
					}
					v2 := VoxelVertex{
						X: byte(x + 1), Y: byte(y + 1), Z: byte(z), N: n,
						R: zn.R, G: zn.G, B: zn.B,
						O: byte(255 * (1 - (lxp+lyp+lxpyp+l)/4)),
					}
					if xpf && ypf {
						v2.O = byte(255 * (1 - (lxp+lyp+l)/3))
					}
					v3 := VoxelVertex{
						X: byte(x), Y: byte(y), Z: byte(z), N: n,
						R: zn.R, G: zn.G, B: zn.B,
						O: byte(255 * (1 - (lxn+lyn+lxnyn+l)/4)),
					}
					if xnf && ynf {
						v3.O = byte(255 * (1 - (lxn+lyn+l)/3))
					}
					v4 := VoxelVertex{
						X: byte(x + 1), Y: byte(y), Z: byte(z), N: n,
						R: zn.R, G: zn.G, B: zn.B,
						O: byte(255 * (1 - (lxp+lyn+lxpyn+l)/4)),
					}
					if xpf && ynf {
						v4.O = byte(255 * (1 - (lxp+lyn+l)/3))
					}
					data = append(data, v1, v3, v2, v2, v3, v4)
				}
			}
		}
	}

	// buffer vertex data to GPU memory
	cm.Buffer("geometry", data)
}
