package cache

import (
	"fmt"

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
	stage   buffer.T
	buffer  buffer.T
	alloc   allocator.T
}

func NewSharedMeshCache(backend vulkan.T, size int) MeshCache {
	stagesize := 1024 * 1024 // 1MB for now

	alloc := allocator.New(size)
	alloc.Alloc(4096)

	return New[vertex.Mesh, VkMesh](&shmeshes{
		backend: backend,
		worker:  backend.Transferer(),

		stage:  buffer.NewShared(backend.Device(), stagesize),
		buffer: buffer.NewRemote(backend.Device(), size, vk.BufferUsageVertexBufferBit|vk.BufferUsageIndexBufferBit),
		alloc:  alloc,
	})
}

func (m *shmeshes) ItemName() string {
	return "Mesh"
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

	if size > m.buffer.Size() {
		panic("mesh is too large for stage buffer")
	}

	m.stage.Write(0, mesh.VertexData())
	m.stage.Write(vtxSize, mesh.IndexData())

	fmt.Println("shmesh: upload", size, "bytes. current block:", cached.block.Size)

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

		// we need to align the vertices in a multiple of the vertex size,
		// so that we can express the vertex offset in the buffer as an
		// index value.
		vtxAlign = -(cached.block.Offset % -mesh.VertexSize())

		m.worker.Queue(func(cmd command.Buffer) {
			cmd.CmdCopyBuffer(m.stage, m.buffer, vk.BufferCopy{
				SrcOffset: vk.DeviceSize(0),
				DstOffset: vk.DeviceSize(cached.block.Offset + vtxAlign),
				Size:      vk.DeviceSize(size),
			})
		})
		m.worker.Submit(command.SubmitInfo{})
		m.worker.Wait()
	}

	cached.elements = mesh.Indices()
	cached.idxType = vk.IndexTypeUint16
	cached.idxOffset = (cached.block.Offset + vtxSize) / 2
	cached.vtxOffset = (cached.block.Offset + vtxAlign) / mesh.VertexSize()
}

func (m *shmeshes) Update(cached VkMesh, mesh vertex.Mesh) {
	vkmesh := cached.(*sharedMesh)
	m.upload(vkmesh, mesh)
}

func (m *shmeshes) Delete(vkmesh VkMesh) {
	cached := vkmesh.(*sharedMesh)
	vkmesh.Destroy()
	if err := m.alloc.Free(cached.block.Offset); err != nil {
		panic(err)
	}
}

func (m *shmeshes) Destroy() {
	m.stage.Destroy()
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
