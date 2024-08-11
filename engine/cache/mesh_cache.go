package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type MeshCache T[vertex.Mesh, *GpuMesh]

type meshCache struct {
	device *device.Device
	worker command.Worker
}

func NewMeshCache(device *device.Device, worker command.Worker) MeshCache {
	return New[vertex.Mesh, *GpuMesh](&meshCache{
		device: device,
		worker: worker,
	})
}

func (m *meshCache) Instantiate(mesh vertex.Mesh, callback func(*GpuMesh)) {
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

	cached = &GpuMesh{
		Key:          mesh.Key(),
		IndexCount:   mesh.IndexCount(),
		IndexType:    idxType,
		IndexOffset:  0,
		VertexOffset: 0,
	}
	if cached.IndexCount == 0 {
		// special case for empty meshes
		callback(cached)
		return
	}

	cmds := command.NewRecorder()
	cmds.Record(func(cmd *command.Buffer) {
		vtxSize := mesh.VertexSize() * mesh.VertexCount()
		vtxStage = buffer.NewShared(m.device, "staging:vertex", vtxSize)

		idxSize := mesh.IndexSize() * mesh.IndexCount()
		idxStage = buffer.NewShared(m.device, "staging:index", idxSize)

		vtxStage.Write(0, mesh.VertexData())
		vtxStage.Flush()
		idxStage.Write(0, mesh.IndexData())
		idxStage.Flush()

		// allocate buffers
		vtxBuffer := buffer.NewRemote(m.device, mesh.Key()+":vertex", vtxSize, core1_0.BufferUsageVertexBuffer)
		idxBuffer := buffer.NewRemote(m.device, mesh.Key()+":index", idxSize, core1_0.BufferUsageIndexBuffer)
		cached.Vertices = buffer.EntireBuffer(vtxBuffer)
		cached.Indices = buffer.EntireBuffer(idxBuffer)

		cmd.CmdCopyBuffer(vtxStage, cached.Vertices.Buffer(), core1_0.BufferCopy{
			Size: vtxSize,
		})
		cmd.CmdCopyBuffer(idxStage, cached.Indices.Buffer(), core1_0.BufferCopy{
			Size: idxSize,
		})
	})
	m.worker.Submit(command.SubmitInfo{
		Marker:   "MeshCache",
		Commands: cmds,
		Callback: func() {
			vtxStage.Destroy()
			idxStage.Destroy()
			callback(cached)
		},
	})
}

func (m *meshCache) Delete(mesh *GpuMesh) {
	if mesh.IndexCount == 0 {
		return
	}
	mesh.Vertices.Buffer().Destroy()
	mesh.Indices.Buffer().Destroy()
}

func (m *meshCache) Destroy() {}

func (m *meshCache) Name() string   { return "MeshCache" }
func (m *meshCache) String() string { return "MeshCache" }
