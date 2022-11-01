package upload

import (
	"image"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

func NewTextureSync(backend vulkan.T, img *image.RGBA) (texture.T, error) {
	// allocate texture
	tex, err := texture.New(backend.Device(), texture.Args{
		Width:  img.Rect.Size().X,
		Height: img.Rect.Size().Y,
		Format: vk.FormatR8g8b8a8Unorm,
		Filter: vk.FilterLinear,
		Wrap:   vk.SamplerAddressModeRepeat,
	})
	if err != nil {
		return nil, err
	}

	// allocate staging buffer
	stage := buffer.NewShared(backend.Device(), len(img.Pix))
	defer stage.Destroy()

	// write to staging buffer
	stage.Write(0, img.Pix)

	// transfer data to texture buffer
	worker := backend.Transferer()
	worker.Queue(func(cmd command.Buffer) {
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
	worker.Submit(command.SubmitInfo{})
	worker.Wait()

	return tex, nil
}
