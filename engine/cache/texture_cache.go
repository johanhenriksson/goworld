package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type TextureCache T[texture.Ref, texture.T]

// mesh cache backend
type textures struct {
	backend vulkan.T
	worker  command.Worker
}

func NewTextureCache(backend vulkan.T) TextureCache {
	return New[texture.Ref, texture.T](&textures{
		backend: backend,
		worker:  backend.Transferer(),
	})
}

func (t *textures) ItemName() string {
	return "Texture"
}

func (t *textures) Instantiate(ref texture.Ref) texture.T {
	img := ref.Load()

	stage := buffer.NewShared(t.backend.Device(), len(img.Pix))
	stage.Write(0, img.Pix)

	tex, err := texture.New(t.backend.Device(), texture.Args{
		Width:  img.Rect.Size().X,
		Height: img.Rect.Size().Y,
		Format: vk.FormatR8g8b8a8Unorm,
		Filter: vk.FilterLinear,
		Wrap:   vk.SamplerAddressModeRepeat,
	})
	if err != nil {
		panic(err)
	}

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

	return tex
}

func (m *textures) Update(tex texture.T, ref texture.Ref) {
}

func (m *textures) Delete(tex texture.T) {
	tex.Destroy()
}

func (m *textures) Destroy() {
}
