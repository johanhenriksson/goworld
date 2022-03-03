package cache

// GpuMesh represents a mesh that exists in GPU memory
type GpuMesh interface {
	Draw() error
}

// ensure nilmesh implements GpuMesh
var _ GpuMesh = &nilmesh{}

// nilmesh is a mesh placeholder with a no-op draw call
type nilmesh struct{}

func (n nilmesh) Draw() error { return nil }
