package cache

import (
	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type TextureCache T[texture.Ref, texture.T]

func NewTextureCache(device device.T, worker command.Worker) TextureCache {
	return New[texture.Ref, texture.T](&textures{
		device: device,
		worker: worker,
	})
}

type textures struct {
	device device.T
	worker command.Worker
}

func (t *textures) Instantiate(ref texture.Ref, callback func(texture.T)) {
	var stage buffer.T
	var tex texture.T

	// transfer data to texture buffer
	t.worker.Queue(func(cmd command.Buffer) {
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
		stage = buffer.NewShared(t.device, "staging:texture", len(img.Buffer))

		// write to staging buffer
		stage.Write(0, img.Buffer)
		stage.Flush()

		cmd.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			tex.Image(),
			core1_0.ImageLayoutUndefined,
			core1_0.ImageLayoutTransferDstOptimal,
			core1_0.ImageAspectColor)
		cmd.CmdCopyBufferToImage(stage, tex.Image(), core1_0.ImageLayoutTransferDstOptimal)
		cmd.CmdImageBarrier(
			core1_0.PipelineStageTransfer,
			core1_0.PipelineStageFragmentShader,
			tex.Image(),
			core1_0.ImageLayoutTransferDstOptimal,
			core1_0.ImageLayoutShaderReadOnlyOptimal,
			core1_0.ImageAspectColor)
	})
	t.worker.Submit(command.SubmitInfo{
		Marker: "TextureUpload",
		Callback: func() {
			stage.Destroy()
			callback(tex)
		},
	})
}

func (t *textures) Delete(tex texture.T) {
	tex.Destroy()
}

func (t *textures) Destroy() {

}

func (t *textures) Name() string   { return "TextureCache" }
func (t *textures) String() string { return "TextureCache" }
