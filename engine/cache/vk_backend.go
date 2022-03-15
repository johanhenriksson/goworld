package cache

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

// glmeshes is a mesh buffer backend for OpenGL
type vkmeshes struct {
	backend vulkan.T
	worker  command.Worker
}

type VkMesh struct {
	Vertices buffer.T
	Indices  buffer.T
	Mesh     vertex.Mesh
}

func (m *VkMesh) Draw() error {
	return nil
}

func NewVkCache(backend vulkan.T) Meshes {
	return &meshes{
		maxAge: 1000,
		cache:  make(map[string]*entry),
		backend: &vkmeshes{
			backend: backend,
			worker:  backend.Transferer(),
		},
	}
}

func (m *vkmeshes) Instantiate(mesh vertex.Mesh, mat material.T) GpuMesh {
	bufsize := 100 * 1024 // 100k for now

	vtx := buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit))
	idx := buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageFlags(vk.BufferUsageIndexBufferBit))

	vtxstage := buffer.NewShared(m.backend.Device(), bufsize)
	vtxstage.Write(mesh.VertexData(), 0)

	idxstage := buffer.NewShared(m.backend.Device(), bufsize)
	idxstage.Write(mesh.IndexData(), 0)

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

func (m *vkmeshes) Update(bmesh GpuMesh, mesh vertex.Mesh) {
}

func (m *vkmeshes) Delete(bmesh GpuMesh) {
	if vkmesh, ok := bmesh.(*VkMesh); ok {
		vkmesh.Vertices.Destroy()
		vkmesh.Indices.Destroy()
	}
}

func (m *vkmeshes) Destroy() {
}
