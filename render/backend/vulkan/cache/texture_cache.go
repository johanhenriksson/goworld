package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/texture"

	vk "github.com/vulkan-go/vulkan"
)

type TextureCache cache.T[texture.Ref, texture.T]

// mesh cache backend
type vktextures struct {
	backend vulkan.T
	worker  command.Worker
}

func NewTextureCache(backend vulkan.T) TextureCache {
	return cache.New[texture.Ref, texture.T](&vktextures{
		backend: backend,
		worker:  backend.Transferer(),
	})
}

func (t *vktextures) Instantiate(ref texture.Ref) texture.T {
	img := ref.Load()

	stage := buffer.NewShared(t.backend.Device(), len(img.Pix))
	stage.Write(0, img.Pix)

	tex := texture.New(t.backend.Device(), texture.Args{
		Width:  img.Rect.Size().X,
		Height: img.Rect.Size().Y,
		Format: vk.FormatR8g8b8a8Unorm,
		Filter: vk.FilterLinear,
		Wrap:   vk.SamplerAddressModeRepeat,
	})

	t.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			tex.Image(),
			vk.ImageLayoutUndefined,
			vk.ImageLayoutTransferDstOptimal,
			vk.ImageAspectColorBit)
		cmd.CmdCopyBufferToImage(stage, tex.Image(), vk.ImageLayoutTransferDstOptimal)
		cmd.CmdImageBarrier(
			vk.PipelineStageTransferBit,
			vk.PipelineStageFragmentShaderBit,
			tex.Image(),
			vk.ImageLayoutTransferDstOptimal,
			vk.ImageLayoutShaderReadOnlyOptimal,
			vk.ImageAspectColorBit)
	})
	t.worker.Submit(command.SubmitInfo{})
	t.worker.Wait()

	stage.Destroy()

	log.Println("buffered texture", ref.Id())

	return tex
}

func (m *vktextures) Update(tex texture.T, ref texture.Ref) {
}

func (m *vktextures) Delete(tex texture.T) {
	tex.Destroy()
}

func (m *vktextures) Destroy() {
}
