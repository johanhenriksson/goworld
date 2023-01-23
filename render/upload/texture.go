package upload

import (
	"image"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"

	vk "github.com/vulkan-go/vulkan"
)

func NewTextureSync(dev device.T, worker command.Worker, img *image.RGBA) (texture.T, error) {
	// allocate texture
	tex, err := texture.New(dev, texture.Args{
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
	stage := buffer.NewShared(dev, len(img.Pix))

	// write to staging buffer
	stage.Write(0, img.Pix)

	// transfer data to texture buffer
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
	worker.Submit(command.SubmitInfo{
		Marker: "TextureUpload",
		Then:   stage.Destroy,
	})
	worker.Wait()

	return tex, nil
}
