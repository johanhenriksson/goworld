package vkrender

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type MeshCache cache.T[vertex.Mesh, *VkMesh]

// mesh cache backend
type vkmeshes struct {
	backend vulkan.T
	worker  command.Worker
}

type VkMesh struct {
	Vertices buffer.T
	Indices  buffer.T
	Mesh     vertex.Mesh
}

func NewMeshCache(backend vulkan.T) MeshCache {
	return cache.New[vertex.Mesh, *VkMesh](&vkmeshes{
		backend: backend,
		worker:  backend.Transferer(),
	})
}

func (m *vkmeshes) Instantiate(mesh vertex.Mesh) *VkMesh {
	bufsize := 100 * 1024 // 100k for now

	vtx := buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageVertexBufferBit)
	idx := buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageIndexBufferBit)

	vtxstage := buffer.NewShared(m.backend.Device(), bufsize)
	vtxstage.Write(0, mesh.VertexData())

	idxstage := buffer.NewShared(m.backend.Device(), bufsize)
	idxstage.Write(0, mesh.IndexData())

	m.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(vtxstage, vtx)
		cmd.CmdCopyBuffer(idxstage, idx)
	})
	m.worker.Submit(command.SubmitInfo{})
	m.worker.Wait()

	vtxstage.Destroy()
	idxstage.Destroy()

	return &VkMesh{
		Vertices: vtx,
		Indices:  idx,
		Mesh:     mesh,
	}
}

func (m *vkmeshes) Update(bmesh *VkMesh, mesh vertex.Mesh) {
}

func (m *vkmeshes) Delete(vkmesh *VkMesh) {
	vkmesh.Vertices.Destroy()
	vkmesh.Indices.Destroy()
}

func (m *vkmeshes) Destroy() {
}
