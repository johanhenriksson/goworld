package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type MeshCache T[vertex.Mesh, Mesh]

type meshes struct {
	device device.T
	worker command.Worker
}

func NewMeshCache(device device.T, worker command.Worker) MeshCache {
	return New[vertex.Mesh, Mesh](&meshes{
		device: device,
		worker: worker,
	})
}

func (m *meshes) Instantiate(mesh vertex.Mesh, callback func(Mesh)) {
	cached := &vkMesh{
		key:      mesh.Key(),
		elements: mesh.Indices(),
		idxType:  core1_0.IndexTypeUInt16,
	}
	if cached.elements == 0 {
		// special case for empty meshes
		callback(cached)
		return
	}

	vtxSize := mesh.VertexSize() * mesh.Vertices()
	vtxStage := buffer.NewShared(m.device, vtxSize)

	idxSize := mesh.IndexSize() * mesh.Indices()
	idxStage := buffer.NewShared(m.device, idxSize)

	vtxStage.Write(0, mesh.VertexData())
	idxStage.Write(0, mesh.IndexData())

	// allocate buffers
	cached.vertices = buffer.NewRemote(m.device, vtxSize, core1_0.BufferUsageVertexBuffer)
	cached.indices = buffer.NewRemote(m.device, idxSize, core1_0.BufferUsageIndexBuffer)

	m.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(vtxStage, cached.vertices, core1_0.BufferCopy{
			Size: vtxSize,
		})
		cmd.CmdCopyBuffer(idxStage, cached.indices, core1_0.BufferCopy{
			Size: idxSize,
		})
	})
	m.worker.OnComplete(func() {
		vtxStage.Destroy()
		idxStage.Destroy()
		callback(cached)
	})
}

func (m *meshes) Submit() {
	m.worker.Submit(command.SubmitInfo{Marker: "MeshCache"})
}

func (m *meshes) Delete(mesh Mesh) {
	mesh.Destroy()
}

func (m *meshes) Destroy() {}

func (m *meshes) Name() string   { return "MeshCache" }
func (m *meshes) String() string { return "MeshCache" }
