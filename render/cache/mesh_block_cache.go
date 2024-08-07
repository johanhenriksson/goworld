package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type MeshBlockCache T[vertex.Mesh, Mesh]

type meshBlockCache struct {
	device   device.T
	worker   command.Worker
	vtxData  buffer.T
	idxData  buffer.T
	vtxAlloc buffer.Allocator
	idxAlloc buffer.Allocator
}

func NewMeshBlockCache(device device.T, worker command.Worker, vtxSize, idxSize int) MeshCache {
	vtxBuf := buffer.NewRemote(device, "MeshVertexBlocks", vtxSize, core1_0.BufferUsageVertexBuffer)
	idxBuf := buffer.NewRemote(device, "MeshIndexBlocks", idxSize, core1_0.BufferUsageIndexBuffer)
	return New[vertex.Mesh, Mesh](&meshBlockCache{
		device:   device,
		worker:   worker,
		vtxData:  vtxBuf,
		idxData:  idxBuf,
		vtxAlloc: buffer.NewBlockAllocator(vtxBuf),
		idxAlloc: buffer.NewBlockAllocator(idxBuf),
	})
}

func (m *meshBlockCache) Instantiate(mesh vertex.Mesh, callback func(Mesh)) {
	var cached *meshBlock
	var vtxStage, idxStage buffer.T

	var idxType core1_0.IndexType
	switch mesh.IndexSize() {
	case 2:
		idxType = core1_0.IndexTypeUInt16
	case 4:
		idxType = core1_0.IndexTypeUInt32
	default:
		panic("illegal index type")
	}

	if mesh.IndexCount() == 0 {
		// special case for empty mesh
		callback(&meshBlock{
			indexCount: mesh.IndexCount(),
			idxType:    idxType,
		})
		return
	}

	vtxSize := mesh.VertexSize() * mesh.VertexCount()
	vtxStage = buffer.NewShared(m.device, "staging:vertex", vtxSize)

	idxSize := mesh.IndexSize() * mesh.IndexCount()
	idxStage = buffer.NewShared(m.device, "staging:index", idxSize)

	log.Printf("MeshBufferCache: Instantiating mesh %s", mesh.Key())
	log.Printf("vtxSize: %d, idxSize: %d", mesh.VertexSize(), mesh.IndexSize())

	// allocate buffers
	vertexBlock, err := m.vtxAlloc.Alloc(vtxSize + mesh.VertexSize())
	if err != nil {
		panic("failed to allocate vertex buffer")
	}
	indexBlock, err := m.idxAlloc.Alloc(idxSize + mesh.IndexSize())
	if err != nil {
		panic("failed to allocate index buffer")
	}

	vtxOffset := buffer.Align(vertexBlock.Offset(), mesh.VertexSize())
	idxOffset := buffer.Align(indexBlock.Offset(), mesh.IndexSize())

	cached = &meshBlock{
		idxType:  idxType,
		vertices: vertexBlock,
		indices:  indexBlock,

		indexCount:   mesh.IndexCount(),
		vertexOffset: vtxOffset / mesh.VertexSize(),
		firstIndex:   idxOffset / mesh.IndexSize(),
	}

	cmds := command.NewRecorder()
	cmds.Record(func(cmd command.Buffer) {
		vtxStage.Write(0, mesh.VertexData())
		vtxStage.Flush()
		idxStage.Write(0, mesh.IndexData())
		idxStage.Flush()

		cmd.CmdCopyBuffer(vtxStage, cached.vertices.Buffer(), core1_0.BufferCopy{
			Size:      vtxSize,
			DstOffset: vtxOffset,
		})
		cmd.CmdCopyBuffer(idxStage, cached.indices.Buffer(), core1_0.BufferCopy{
			Size:      idxSize,
			DstOffset: idxOffset,
		})
	})
	m.worker.Submit(command.SubmitInfo{
		Marker:   "MeshBufferCache",
		Commands: cmds,
		Callback: func() {
			vtxStage.Destroy()
			idxStage.Destroy()
			callback(cached)
		},
	})
}

func (m *meshBlockCache) Delete(mesh Mesh) {
	msh := mesh.(*meshBlock)
	m.vtxAlloc.Free(msh.vertices)
	m.idxAlloc.Free(msh.indices)
}

func (m *meshBlockCache) Destroy() {
	m.vtxData.Destroy()
	m.idxData.Destroy()
}

func (m *meshBlockCache) Name() string   { return "MeshCache" }
func (m *meshBlockCache) String() string { return "MeshCache" }
