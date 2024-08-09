package command

type Draw struct {
	VertexCount    uint32
	InstanceCount  uint32
	VertexOffset   int32
	InstanceOffset uint32
}

type DrawBuffer interface {
	CmdDraw(Draw)
}

type DrawIndexed struct {
	IndexCount     uint32
	InstanceCount  uint32
	IndexOffset    uint32
	VertexOffset   int32
	InstanceOffset uint32
}

type DrawIndexedBuffer interface {
	CmdDrawIndexed(DrawIndexed)
}
