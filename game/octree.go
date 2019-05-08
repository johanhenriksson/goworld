package game

/*

import (
	"math"
)

type VoxelType int
type OctreeSlot int
type Vec3 [3]int

const (
	VxEmpty VoxelType = iota
	VxSplit
	VxGrass
	VxRock
)

const (
	LBF OctreeSlot = 0
	LBB = 1
	LTF = 2
	LTB = 3
	RBF = 4
	RBB = 5
	RTF = 6
	RTB = 7
)


type Octree struct {
	Type		VoxelType
	Height		int
	Size		int
	Parent		*Octree
	Child		[8]*Octree
	Position	Vec3
}

func NewOctree() *Octree {
	var children [8]*Octree
	if vtype == VxSplit {
		children = [8]*Octree { }
	}

	h := 4
	return &Octree {
		Type:   VxRock,
		Height: h,
		Parent: nil,
		Child:  children,
		Size:   int(math.Pow(4, float64(h))),
	}
}


func (ot *Octree) Split() {


}

func (ot *Octree) Set(x, y, z int, t VoxelType) bool {
	if ot.Height > 1 {
		// recurse
		if ot.Type != VxSplit {
			// split first
		}
		return ot.getChildLoc(x, y, z).Set(x, y, z, t)
	} else {
		ot.Type = t
		return true
	}
}

func (ot Octree) getChildLoc(x, y, z) *Octree {
}

func (ot Octree) getOffset(slot OctreeSlot) Vec3 {
	flags := int(slot)
	return Vec3 {
		(flags >> 2) & 1,
		(flags >> 1) & 1,
		(flags >> 0) & 1,
	}
}

func (ot Octree) getSlot(Vec3 pos) OctreeSlot {
	flags := 0
	if pos[0] > 0 {
		flags |= 1 << 2
	}
	if pos[1] > 0 {
		flags |= 1 << 1
	}
	if pos[2] > 0 {
		flags |= 1 << 0
	}
	return OctreeSlot(flags)
}

func (ot Octree) Get(x, y, z int) VoxelType {
	if ot.Type == VxSplit {


	}
	return ot.Type
}
*/
