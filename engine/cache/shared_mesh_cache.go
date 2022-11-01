package cache

import (
	"github.com/johanhenriksson/goworld/engine/cache/allocator"
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

// mesh cache backend
type shmeshes struct {
	backend vulkan.T
	worker  command.Worker
	buffer  buffer.T
	alloc   allocator.T
}

func NewSharedMeshCache(backend vulkan.T, size int) MeshCache {
	return NewConcurrent[vertex.Mesh, VkMesh](&shmeshes{
		backend: backend,
		worker:  backend.Transferer(),

		buffer: buffer.NewRemote(backend.Device(), size, vk.BufferUsageVertexBufferBit|vk.BufferUsageIndexBufferBit),
		alloc:  allocator.New(size),
	})
}

func (m *shmeshes) Name() string {
	return "SharedMesh"
}

func (m *shmeshes) Instantiate(mesh vertex.Mesh) VkMesh {
	cached := &sharedMesh{
		buffer: m.buffer,
	}
	m.upload(cached, mesh)
	return cached
}

func (m *shmeshes) upload(cached *sharedMesh, mesh vertex.Mesh) {
	// todo: rewrite in a thread safe manner
	// introduce a queue and a goroutine that periodically runs transfers

	vtxSize := mesh.VertexSize() * mesh.Vertices()
	idxSize := mesh.IndexSize() * mesh.Indices()
	size := vtxSize + idxSize
	vtxAlign := 0

	if cached.block.Size > 0 && (cached.block.Size < size || cached.block.Size > size*2) {
		// free
		m.alloc.Free(cached.block.Offset)
		cached.block = allocator.Block{}
	}

	if size > 0 {
		if cached.block.Size == 0 {
			// add space for an extra vertex so that we have room for alignment
			var err error
			cached.block, err = m.alloc.Alloc(size + mesh.VertexSize())
			if err != nil {
				panic(err)
			}
		}

		// allocate staging buffer
		stage := buffer.NewShared(m.backend.Device(), size)
		defer stage.Destroy()
		stage.Write(0, mesh.VertexData())
		stage.Write(vtxSize, mesh.IndexData())

		// we need to align the vertices in a multiple of the vertex size,
		// so that we can express the vertex offset in the buffer as an
		// index value.
		vtxAlign = -(cached.block.Offset % -mesh.VertexSize())

		m.worker.Queue(func(cmd command.Buffer) {
			cmd.CmdCopyBuffer(stage, m.buffer, vk.BufferCopy{
				SrcOffset: vk.DeviceSize(0),
				DstOffset: vk.DeviceSize(cached.block.Offset + vtxAlign),
				Size:      vk.DeviceSize(size),
			})
		})
		m.worker.Submit(command.SubmitInfo{})
		m.worker.Wait()
	}

	cached.elements = mesh.Indices()
	cached.idxType = vertex.IndexType(mesh.IndexSize())
	cached.idxOffset = (cached.block.Offset + vtxSize) / mesh.IndexSize()
	cached.vtxOffset = (cached.block.Offset + vtxAlign) / mesh.VertexSize()
}

func (m *shmeshes) Update(cached VkMesh, mesh vertex.Mesh) {
	shmesh := cached.(*sharedMesh)
	m.upload(shmesh, mesh)
}

func (m *shmeshes) Delete(vkmesh VkMesh) {
	cached := vkmesh.(*sharedMesh)
	vkmesh.Destroy()
	if err := m.alloc.Free(cached.block.Offset); err != nil {
		panic(err)
	}
}

func (m *shmeshes) Destroy() {
	m.buffer.Destroy()
}

type sharedMesh struct {
	block     allocator.Block
	elements  int
	vtxOffset int
	idxOffset int
	idxType   vk.IndexType
	buffer    buffer.T
}

func (m *sharedMesh) Draw(cmd command.Buffer, index int) {
	cmd.CmdBindVertexBuffer(m.buffer, 0)
	cmd.CmdBindIndexBuffers(m.buffer, 0, m.idxType)

	// index of the object properties in the ssbo
	cmd.CmdDrawIndexed(m.elements, 1, m.idxOffset, m.vtxOffset, index)
}

func (m *sharedMesh) Destroy() {
}
