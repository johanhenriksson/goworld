package cache

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// the key must be a reference, not the mesh itself

type MeshCache T[assets.Mesh, *GpuMesh]

type meshCache struct {
	device *device.Device
	worker command.Worker
}

func NewMeshCache(device *device.Device, worker command.Worker) MeshCache {
	return New[assets.Mesh, *GpuMesh](&meshCache{
		device: device,
		worker: worker,
	})
}

func (m *meshCache) Instantiate(ref assets.Mesh, callback func(*GpuMesh)) {
	var cached *GpuMesh
	var vtxStage, idxStage buffer.T

	mesh := ref.LoadMesh(assets.FS)

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
		key:          mesh.Key(),
		indexCount:   mesh.IndexCount(),
		indexType:    idxType,
		indexOffset:  0,
		vertexOffset: 0,
	}
	if cached.indexCount == 0 {
		// special case for empty meshes
		callback(cached)
		return
	}

	cmds := command.NewRecorder()
	cmds.Record(func(cmd *command.Buffer) {
		vtxSize := mesh.VertexSize() * mesh.VertexCount()
		vtxStage = buffer.NewCpuLocal(m.device, "staging:vertex", vtxSize)

		idxSize := mesh.IndexSize() * mesh.IndexCount()
		idxStage = buffer.NewCpuLocal(m.device, "staging:index", idxSize)

		vtxStage.Write(0, mesh.VertexData())
		vtxStage.Flush()
		idxStage.Write(0, mesh.IndexData())
		idxStage.Flush()

		// allocate buffers
		vtxBuffer := buffer.NewGpuLocal(m.device, mesh.Key()+":vertex", vtxSize, core1_0.BufferUsageVertexBuffer)
		idxBuffer := buffer.NewGpuLocal(m.device, mesh.Key()+":index", idxSize, core1_0.BufferUsageIndexBuffer)
		cached.vertices = buffer.EntireBuffer(vtxBuffer)
		cached.indices = buffer.EntireBuffer(idxBuffer)

		cmd.CmdCopyBuffer(vtxStage, cached.vertices.Buffer(), core1_0.BufferCopy{
			Size: vtxSize,
		})
		cmd.CmdCopyBuffer(idxStage, cached.indices.Buffer(), core1_0.BufferCopy{
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
	if mesh.indexCount == 0 {
		return
	}
	mesh.vertices.Buffer().Destroy()
	mesh.indices.Buffer().Destroy()
}

func (m *meshCache) Destroy() {}

func (m *meshCache) Name() string   { return "MeshCache" }
func (m *meshCache) String() string { return "MeshCache" }
