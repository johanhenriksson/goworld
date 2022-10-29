package cache

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type MeshCache cache.T[vertex.Mesh, VkMesh]

type VkMesh interface {
	Draw(command.Buffer, int)
	Destroy()
}

// mesh cache backend
type meshes struct {
	backend  vulkan.T
	worker   command.Worker
	idxstage buffer.T
	vtxstage buffer.T
}

func NewMeshCache(backend vulkan.T) MeshCache {
	stagesize := 100 * 1024 // 100k for now

	return cache.New[vertex.Mesh, VkMesh](&meshes{
		backend: backend,
		worker:  backend.Transferer(),

		vtxstage: buffer.NewShared(backend.Device(), stagesize),
		idxstage: buffer.NewShared(backend.Device(), stagesize),
	})
}

func (m *meshes) ItemName() string {
	return "Mesh"
}

func (m *meshes) Instantiate(mesh vertex.Mesh) VkMesh {
	bufsize := 100 * 1024 // 100k for now

	cached := &vkMesh{
		vertices: buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageVertexBufferBit),
		indices:  buffer.NewRemote(m.backend.Device(), bufsize, vk.BufferUsageIndexBufferBit),
	}
	m.upload(cached, mesh)

	return cached
}

func (m *meshes) upload(cached *vkMesh, mesh vertex.Mesh) {
	// todo: rewrite in a thread safe manner
	// introduce a queue and a goroutine that periodically runs transfers

	m.vtxstage.Write(0, mesh.VertexData())
	m.idxstage.Write(0, mesh.IndexData())

	m.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(m.vtxstage, cached.vertices)
		cmd.CmdCopyBuffer(m.idxstage, cached.indices)
	})
	m.worker.Submit(command.SubmitInfo{})
	m.worker.Wait()

	cached.elements = mesh.Elements()
	cached.idxType = vk.IndexTypeUint16
	cached.idxOffset = 0
	cached.vtxOffset = 0
}

func (m *meshes) Update(cached VkMesh, mesh vertex.Mesh) {
	vkmesh := cached.(*vkMesh)
	m.upload(vkmesh, mesh)
}

func (m *meshes) Delete(vkmesh VkMesh) {
	vkmesh.Destroy()
}

func (m *meshes) Destroy() {
	m.vtxstage.Destroy()
	m.idxstage.Destroy()
}

type vkMesh struct {
	elements  int
	vtxOffset int
	idxOffset int
	idxType   vk.IndexType
	vertices  buffer.T
	indices   buffer.T
}

func (m *vkMesh) Draw(cmd command.Buffer, index int) {
	cmd.CmdBindVertexBuffer(m.vertices, 0)
	cmd.CmdBindIndexBuffers(m.indices, 0, m.idxType)

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(m.elements, 1, m.idxOffset, m.vtxOffset, index)
}

func (m *vkMesh) Destroy() {
	m.vertices.Destroy()
	m.indices.Destroy()
}
