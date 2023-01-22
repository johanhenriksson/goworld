package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type MeshCache T[vertex.Mesh, VkMesh]

type VkMesh interface {
	Draw(command.Buffer, int)
	Destroy()
}

// mesh cache backend
type meshes struct {
	device   device.T
	worker   command.Worker
	idxstage buffer.T
	vtxstage buffer.T
}

func NewMeshCache(dev device.T, transferer command.Worker) MeshCache {
	stagesize := 100 * 1024 // 100k for now

	return New[vertex.Mesh, VkMesh](&meshes{
		device: dev,
		worker: transferer,

		vtxstage: buffer.NewShared(dev, stagesize),
		idxstage: buffer.NewShared(dev, stagesize),
	})
}

func (m *meshes) Name() string {
	return "Mesh"
}

func (m *meshes) Instantiate(mesh vertex.Mesh) VkMesh {
	vtxSize := mesh.VertexSize() * mesh.Vertices()
	idxSize := mesh.IndexSize() * mesh.Indices()

	cached := &vkMesh{
		vertices: buffer.NewRemote(m.device, vtxSize, vk.BufferUsageVertexBufferBit),
		indices:  buffer.NewRemote(m.device, idxSize, vk.BufferUsageIndexBufferBit),
	}
	m.upload(cached, mesh)

	return cached
}

func (m *meshes) upload(cached *vkMesh, mesh vertex.Mesh) {
	// todo: rewrite in a thread safe manner
	// introduce a queue and a goroutine that periodically runs transfers

	vtxSize := mesh.VertexSize() * mesh.Vertices()
	if vtxSize > m.vtxstage.Size() {
		panic("mesh is too large for vertex stage buffer")
	}
	idxSize := mesh.IndexSize() * mesh.Indices()
	if idxSize > m.idxstage.Size() {
		panic("mesh is too large for index stage buffer")
	}

	m.vtxstage.Write(0, mesh.VertexData())
	m.idxstage.Write(0, mesh.IndexData())

	// reallocate buffers if required
	if cached.vertices.Size() < vtxSize || cached.vertices.Size() > 2*vtxSize {
		cached.vertices.Destroy()
		cached.vertices = buffer.NewRemote(m.device, vtxSize, vk.BufferUsageVertexBufferBit)
	}
	if cached.indices.Size() < idxSize || cached.indices.Size() > 2*idxSize {
		cached.indices.Destroy()
		cached.indices = buffer.NewRemote(m.device, idxSize, vk.BufferUsageIndexBufferBit)
	}

	m.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(m.vtxstage, cached.vertices, vk.BufferCopy{
			Size: vk.DeviceSize(vtxSize),
		})
		cmd.CmdCopyBuffer(m.idxstage, cached.indices, vk.BufferCopy{
			Size: vk.DeviceSize(idxSize),
		})
	})
	m.worker.Submit(command.SubmitInfo{})
	m.worker.Wait()

	cached.elements = mesh.Indices()
	cached.idxType = vk.IndexTypeUint16
}

func (m *meshes) Update(cached VkMesh, mesh vertex.Mesh) VkMesh {
	vkmesh := cached.(*vkMesh)
	m.upload(vkmesh, mesh)
	return vkmesh
}

func (m *meshes) Delete(vkmesh VkMesh) {
	vkmesh.Destroy()
}

func (m *meshes) Destroy() {
	m.vtxstage.Destroy()
	m.idxstage.Destroy()
}

type vkMesh struct {
	elements int
	idxType  vk.IndexType
	vertices buffer.T
	indices  buffer.T
}

func (m *vkMesh) Draw(cmd command.Buffer, index int) {
	cmd.CmdBindVertexBuffer(m.vertices, 0)
	cmd.CmdBindIndexBuffers(m.indices, 0, m.idxType)

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(m.elements, 1, 0, 0, index)
}

func (m *vkMesh) Destroy() {
	m.vertices.Destroy()
	m.indices.Destroy()
}
