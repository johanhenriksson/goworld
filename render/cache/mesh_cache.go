package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
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
		vertices: buffer.NewRemote(m.device, vtxSize, core1_0.BufferUsageVertexBuffer),
		indices:  buffer.NewRemote(m.device, idxSize, core1_0.BufferUsageIndexBuffer),
	}
	m.upload(cached, mesh)

	return cached
}

func (m *meshes) upload(cached *vkMesh, mesh vertex.Mesh) {
	// todo: rewrite in a thread safe manner
	// introduce a queue and a goroutine that periodically runs transfers

	vtxSize := mesh.VertexSize() * mesh.Vertices()
	vtxStage := buffer.NewShared(m.device, vtxSize)

	idxSize := mesh.IndexSize() * mesh.Indices()
	idxStage := buffer.NewShared(m.device, idxSize)

	vtxStage.Write(0, mesh.VertexData())
	idxStage.Write(0, mesh.IndexData())

	// reallocate buffers if required
	if cached.vertices.Size() < vtxSize || cached.vertices.Size() > 2*vtxSize {
		cached.vertices.Destroy()
		cached.vertices = buffer.NewRemote(m.device, vtxSize, core1_0.BufferUsageVertexBuffer)
	}
	if cached.indices.Size() < idxSize || cached.indices.Size() > 2*idxSize {
		cached.indices.Destroy()
		cached.indices = buffer.NewRemote(m.device, idxSize, core1_0.BufferUsageIndexBuffer)
	}

	m.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(vtxStage, cached.vertices, core1_0.BufferCopy{
			Size: vtxSize,
		})
		cmd.CmdCopyBuffer(idxStage, cached.indices, core1_0.BufferCopy{
			Size: idxSize,
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

	cached.elements = mesh.Indices()
	cached.idxType = core1_0.IndexTypeUInt16
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
	idxType  core1_0.IndexType
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
