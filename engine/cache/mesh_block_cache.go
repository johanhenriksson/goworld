package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type meshBlockCache struct {
	device   *device.Device
	worker   command.Worker
	vtxData  buffer.T
	idxData  buffer.T
	vtxAlloc buffer.Allocator
	idxAlloc buffer.Allocator
}

func NewMeshBlockCache(device *device.Device, worker command.Worker, vtxSize, idxSize int) MeshCache {
	vtxBuf := buffer.NewGpuLocal(device, "MeshVertexBlocks", vtxSize, core1_0.BufferUsageVertexBuffer)
	idxBuf := buffer.NewGpuLocal(device, "MeshIndexBlocks", idxSize, core1_0.BufferUsageIndexBuffer)
	return New[vertex.Mesh, *GpuMesh](&meshBlockCache{
		device:   device,
		worker:   worker,
		vtxData:  vtxBuf,
		idxData:  idxBuf,
		vtxAlloc: buffer.NewBlockAllocator(vtxBuf),
		idxAlloc: buffer.NewBlockAllocator(idxBuf),
	})
}

func (m *meshBlockCache) Instantiate(mesh vertex.Mesh, callback func(*GpuMesh)) {
	var cached *GpuMesh
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
		callback(&GpuMesh{
			Key:        mesh.Key(),
			IndexCount: mesh.IndexCount(),
			IndexType:  idxType,
		})
		return
	}

	vtxSize := mesh.VertexSize() * mesh.VertexCount()
	vtxStage = buffer.NewCpuLocal(m.device, "staging:vertex", vtxSize)

	idxSize := mesh.IndexSize() * mesh.IndexCount()
	idxStage = buffer.NewCpuLocal(m.device, "staging:index", idxSize)

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

	cached = &GpuMesh{
		Key:       mesh.Key(),
		IndexType: idxType,
		Vertices:  vertexBlock,
		Indices:   indexBlock,

		IndexCount:   mesh.IndexCount(),
		VertexOffset: vtxOffset / mesh.VertexSize(),
		IndexOffset:  idxOffset / mesh.IndexSize(),
	}

	cmds := command.NewRecorder()
	cmds.Record(func(cmd *command.Buffer) {
		vtxStage.Write(0, mesh.VertexData())
		vtxStage.Flush()
		idxStage.Write(0, mesh.IndexData())
		idxStage.Flush()

		cmd.CmdCopyBuffer(vtxStage, cached.Vertices.Buffer(), core1_0.BufferCopy{
			Size:      vtxSize,
			DstOffset: vtxOffset,
		})
		cmd.CmdCopyBuffer(idxStage, cached.Indices.Buffer(), core1_0.BufferCopy{
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

func (m *meshBlockCache) Delete(mesh *GpuMesh) {
	if mesh.IndexCount == 0 {
		return
	}
	if err := m.vtxAlloc.Free(mesh.Vertices); err != nil {
		panic(err)
	}
	if err := m.idxAlloc.Free(mesh.Indices); err != nil {
		panic(err)
	}
}

func (m *meshBlockCache) Destroy() {
	m.vtxData.Destroy()
	m.idxData.Destroy()
}

func (m *meshBlockCache) Name() string   { return "MeshCache" }
func (m *meshBlockCache) String() string { return "MeshCache" }
