package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type MeshCache T[vertex.Mesh, Mesh]

type meshCache struct {
	device device.T
	worker command.Worker
}

func NewMeshCache(device device.T, worker command.Worker) MeshCache {
	return New[vertex.Mesh, Mesh](&meshCache{
		device: device,
		worker: worker,
	})
}

func (m *meshCache) Instantiate(mesh vertex.Mesh, callback func(Mesh)) {
	var cached *cachedMesh
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

	cached = &cachedMesh{
		key:        mesh.Key(),
		indexCount: mesh.IndexCount(),
		idxType:    idxType,
	}
	if cached.indexCount == 0 {
		// special case for empty meshes
		callback(cached)
		return
	}

	cmds := command.NewRecorder()
	cmds.Record(func(cmd command.Buffer) {
		vtxSize := mesh.VertexSize() * mesh.VertexCount()
		vtxStage = buffer.NewShared(m.device, "staging:vertex", vtxSize)

		idxSize := mesh.IndexSize() * mesh.IndexCount()
		idxStage = buffer.NewShared(m.device, "staging:index", idxSize)

		vtxStage.Write(0, mesh.VertexData())
		vtxStage.Flush()
		idxStage.Write(0, mesh.IndexData())
		idxStage.Flush()

		// allocate buffers
		cached.vertices = buffer.NewRemote(m.device, mesh.Key()+":vertex", vtxSize, core1_0.BufferUsageVertexBuffer)
		cached.indices = buffer.NewRemote(m.device, mesh.Key()+":index", idxSize, core1_0.BufferUsageIndexBuffer)

		cmd.CmdCopyBuffer(vtxStage, cached.vertices, core1_0.BufferCopy{
			Size: vtxSize,
		})
		cmd.CmdCopyBuffer(idxStage, cached.indices, core1_0.BufferCopy{
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

func (m *meshCache) Delete(mesh Mesh) {
	vkmesh := mesh.(*cachedMesh)
	if vkmesh.vertices != nil {
		vkmesh.vertices.Destroy()
		vkmesh.vertices = nil
	}
	if vkmesh.indices != nil {
		vkmesh.indices.Destroy()
		vkmesh.indices = nil
	}
}

func (m *meshCache) Destroy() {}

func (m *meshCache) Name() string   { return "MeshCache" }
func (m *meshCache) String() string { return "MeshCache" }
