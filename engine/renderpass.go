package engine

type RenderPass interface {
    DrawPass(*Scene)
}
