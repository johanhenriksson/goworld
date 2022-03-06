package cache

import (
	"log"

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
	Buffer buffer.T
	Mesh   vertex.Mesh
}

func (m *VkMesh) Draw() error {
	return nil
}

func NewVkCache(backend vulkan.T) Meshes {
	worker := command.NewWorker(backend.Device(), vk.QueueFlags(vk.QueueTransferBit))
	return &meshes{
		maxAge: 1000,
		cache:  make(map[string]*entry),
		backend: &vkmeshes{
			backend: backend,
			worker:  worker,
		},
	}
}

func (m *vkmeshes) Instantiate(mesh vertex.Mesh, mat material.T) GpuMesh {
	bufsize := 100 * 1024 // 100k for now

	vtx := buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit))

	vtxstage := buffer.NewShared(m.backend.Device(), bufsize)
	vtxstage.Write(mesh.VertexData(), 0)

	m.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(vtxstage, vtx)
	})
	m.worker.Submit(command.SubmitInfo{})
	log.Println("waiting for transfers...")
	m.worker.Wait()
	log.Println("transfers completed")

	vtxstage.Destroy()

	return &VkMesh{
		Buffer: vtx,
		Mesh:   mesh,
	}
}

func (m *vkmeshes) Update(bmesh GpuMesh, mesh vertex.Mesh) {
}

func (m *vkmeshes) Delete(bmesh GpuMesh) {
	if vkmesh, ok := bmesh.(*VkMesh); ok {
		vkmesh.Buffer.Destroy()
	}
}

func (m *vkmeshes) Destroy() {
	m.worker.Destroy()
}
