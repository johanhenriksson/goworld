package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type MeshCache T[vertex.Mesh, VkMesh]

type VkMesh interface {
	Draw(command.Buffer, int)
	Destroy()
}

// mesh cache backend
type meshes struct {
	device device.T
	worker command.Worker
}

func NewMeshCache(dev device.T, transferer command.Worker) MeshCache {
	return New[vertex.Mesh, VkMesh](&meshes{
		device: dev,
		worker: transferer,
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
	util.Timer("Mesh:Upload")

	vtxSize := mesh.VertexSize() * mesh.Vertices()
	vtxStage := buffer.NewShared(m.device, vtxSize)

	idxSize := mesh.IndexSize() * mesh.Indices()
	idxStage := buffer.NewShared(m.device, idxSize)

	vtxStage.Write(0, mesh.VertexData())
	idxStage.Write(0, mesh.IndexData())

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
		cmd.CmdCopyBuffer(vtxStage, cached.vertices, vk.BufferCopy{
			Size: vk.DeviceSize(vtxSize),
		})
		cmd.CmdCopyBuffer(idxStage, cached.indices, vk.BufferCopy{
			Size: vk.DeviceSize(idxSize),
		})
	})
	m.worker.Submit(command.SubmitInfo{
		Marker: "MeshCache",
		Then: func() {
			vtxStage.Destroy()
			idxStage.Destroy()
		},
	})
	m.worker.Wait()

	util.Elapsed("Mesh:Upload")
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