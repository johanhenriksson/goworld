package upload

import (
	"fmt"
	osimage "image"
	"image/png"
	"os"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

func NewTextureSync(dev device.T, worker command.Worker, img *osimage.RGBA) (texture.T, error) {
	// allocate texture
	tex, err := texture.New(dev, texture.Args{
		Width:  img.Rect.Size().X,
		Height: img.Rect.Size().Y,
		Format: core1_0.FormatR8G8B8A8UnsignedNormalized,
		Filter: core1_0.FilterLinear,
		Wrap:   core1_0.SamplerAddressModeRepeat,
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
	worker.Submit(command.SubmitInfo{
		Marker: "TextureUpload",
		Then:   stage.Destroy,
	})
	worker.Wait()

	return tex, nil
}

func DownloadImageAsync(dev device.T, worker command.Worker, src image.T) (<-chan *osimage.RGBA, error) {
	swizzle := false
	switch src.Format() {
	case core1_0.FormatB8G8R8A8UnsignedNormalized:
		swizzle = true
	case core1_0.FormatR8G8B8A8UnsignedNormalized:
		break
	default:
		return nil, fmt.Errorf("unsupported source format")
	}

	dst, err := image.New(dev, image.Args{
		Type:    core1_0.ImageType2D,
		Width:   src.Width(),
		Height:  src.Height(),
		Depth:   1,
		Layers:  1,
		Levels:  1,
		Format:  core1_0.FormatR8G8B8A8UnsignedNormalized,
		Memory:  core1_0.MemoryPropertyHostVisible | core1_0.MemoryPropertyHostCoherent,
		Tiling:  core1_0.ImageTilingLinear,
		Usage:   core1_0.ImageUsageTransferDst,
		Sharing: core1_0.SharingModeExclusive,
		Layout:  core1_0.ImageLayoutUndefined,
	})
	if err != nil {
		return nil, err
	}

	// transfer data from texture buffer
	worker.Queue(func(cmd command.Buffer) {
		cmd.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			src,
			core1_0.ImageLayoutUndefined,
			core1_0.ImageLayoutTransferSrcOptimal,
			core1_0.ImageAspectColor)
		cmd.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			dst,
			core1_0.ImageLayoutUndefined,
			core1_0.ImageLayoutTransferDstOptimal,
			core1_0.ImageAspectColor)
		cmd.CmdCopyImage(src, core1_0.ImageLayoutTransferSrcOptimal, dst, core1_0.ImageLayoutTransferDstOptimal, core1_0.ImageAspectColor)
		cmd.CmdImageBarrier(
			core1_0.PipelineStageTransfer,
			core1_0.PipelineStageFragmentShader,
			src,
			core1_0.ImageLayoutTransferSrcOptimal,
			core1_0.ImageLayoutColorAttachmentOptimal,
			core1_0.ImageAspectColor)
		cmd.CmdImageBarrier(
			core1_0.PipelineStageTopOfPipe,
			core1_0.PipelineStageTransfer,
			dst,
			core1_0.ImageLayoutTransferDstOptimal,
			core1_0.ImageLayoutGeneral,
			core1_0.ImageAspectColor)
	})

	done := make(chan *osimage.RGBA)
	worker.Submit(command.SubmitInfo{
		Marker: "TextureDownload",
		Then: func() {
			defer dst.Destroy()
			defer close(done)

			out := osimage.NewRGBA(osimage.Rect(0, 0, dst.Width(), dst.Height()))
			dst.Memory().Read(0, out.Pix)

			// swizzle colors if required BGR -> RGB
			if swizzle {
				for i := 0; i < len(out.Pix); i += 4 {
					b := out.Pix[i]
					r := out.Pix[i+2]
					out.Pix[i] = r
					out.Pix[i+2] = b
				}
			}

			done <- out
		},
	})

	return done, nil
}

func DownloadImage(dev device.T, worker command.Worker, src image.T) (*osimage.RGBA, error) {
	img, err := DownloadImageAsync(dev, worker, src)
	if err != nil {
		return nil, err
	}
	return <-img, nil
}

func SavePng(img *osimage.RGBA, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer out.Close()
	if err := png.Encode(out, img); err != nil {
		return err
	}
	return nil
}
