package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type TextureCache T[texture.Ref, *texture.Texture]

func NewTextureCache(device *device.Device, worker command.Worker) TextureCache {
	return New[texture.Ref, *texture.Texture](&textures{
		device: device,
		worker: worker,
	})
}

type textures struct {
	device *device.Device
	worker command.Worker
}

func (t *textures) Instantiate(ref texture.Ref, callback func(*texture.Texture)) {
	var stage *buffer.Buffer
	var tex *texture.Texture

	// transfer data to texture buffer
	cmds := command.NewRecorder()
	cmds.Record(func(cmd *command.Buffer) {
		// load image data
		img := ref.ImageData()

		// args & defaults
		args := ref.TextureArgs()

		// allocate texture
		var err error
		tex, err = texture.New(t.device, ref.Key(), img.Width, img.Height, img.Format, args)
		if err != nil {
			panic(err)
		}

		// allocate staging buffer
		stage = buffer.NewCpuLocal(t.device, "staging:texture", len(img.Buffer))

		// write to staging buffer
		stage.Write(0, img.Buffer)
		stage.Flush()

		mipLevels := 1
		if args.Mipmaps {
			// verify that the format supports linear filtering
			format := t.device.GetFormatProperties(img.Format)
			if format.OptimalTilingFeatures&core1_0.FormatFeatureSampledImageFilterLinear == 0 {
				panic("cant generate mipmaps: texture format does not support linear filtering")
			}

			mipLevels = image.MipLevels(img.Width, img.Height)
		}

		cmd.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			tex.Image(),
			core1_0.ImageLayoutUndefined,
			core1_0.ImageLayoutTransferDstOptimal,
			core1_0.ImageAspectColor,
			0, mipLevels)
		cmd.CmdCopyBufferToImage(stage, tex.Image(), core1_0.ImageLayoutTransferDstOptimal)

		// generate mipmaps
		mipWidth := img.Width
		mipHeight := img.Height
		for dstLevel := 1; dstLevel < mipLevels; dstLevel++ {
			srcLevel := dstLevel - 1
			dstWidth := max(1, mipWidth/2)
			dstHeight := max(1, mipHeight/2)

			// transition source image to transfer source optimal
			cmd.CmdImageBarrier(
				core1_0.PipelineStageTransfer,
				core1_0.PipelineStageTransfer,
				tex.Image(),
				core1_0.ImageLayoutTransferDstOptimal,
				core1_0.ImageLayoutTransferSrcOptimal,
				core1_0.ImageAspectColor,
				srcLevel, 1)

			// blit to next mip level
			cmd.Ptr().CmdBlitImage(
				tex.Image().Ptr(), core1_0.ImageLayoutTransferSrcOptimal,
				tex.Image().Ptr(), core1_0.ImageLayoutTransferDstOptimal,
				[]core1_0.ImageBlit{
					{
						SrcSubresource: core1_0.ImageSubresourceLayers{
							AspectMask:     core1_0.ImageAspectColor,
							MipLevel:       srcLevel,
							BaseArrayLayer: 0,
							LayerCount:     1,
						},
						SrcOffsets: [2]core1_0.Offset3D{
							{X: 0, Y: 0, Z: 0},
							{X: mipWidth, Y: mipHeight, Z: 1},
						},
						DstSubresource: core1_0.ImageSubresourceLayers{
							AspectMask:     core1_0.ImageAspectColor,
							MipLevel:       dstLevel,
							BaseArrayLayer: 0,
							LayerCount:     1,
						},
						DstOffsets: [2]core1_0.Offset3D{
							{X: 0, Y: 0, Z: 0},
							{X: dstWidth, Y: dstHeight, Z: 1},
						},
					},
				},
				core1_0.FilterLinear) // todo: based on texture filtering settings

			// transition to shader read optimal
			cmd.CmdImageBarrier(
				core1_0.PipelineStageTransfer,
				core1_0.PipelineStageFragmentShader,
				tex.Image(),
				core1_0.ImageLayoutTransferSrcOptimal,
				core1_0.ImageLayoutShaderReadOnlyOptimal,
				core1_0.ImageAspectColor,
				srcLevel, 1)

			mipWidth = dstWidth
			mipHeight = dstHeight
		}

		// transition the top mip level to shader read optimal
		cmd.CmdImageBarrier(
			core1_0.PipelineStageTransfer,
			core1_0.PipelineStageFragmentShader,
			tex.Image(),
			core1_0.ImageLayoutTransferDstOptimal,
			core1_0.ImageLayoutShaderReadOnlyOptimal,
			core1_0.ImageAspectColor,
			mipLevels-1, 1)
	})
	t.worker.Submit(command.SubmitInfo{
		Marker:   "TextureUpload",
		Commands: cmds,
		Callback: func() {
			stage.Destroy()
			callback(tex)
		},
	})
}

func (t *textures) Delete(tex *texture.Texture) {
	tex.Destroy()
}

func (t *textures) Destroy() {

}

func (t *textures) Name() string   { return "TextureCache" }
func (t *textures) String() string { return "TextureCache" }
