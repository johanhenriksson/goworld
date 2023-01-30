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

	vk "github.com/vulkan-go/vulkan"
)

func NewTextureSync(dev device.T, worker command.Worker, img *osimage.RGBA) (texture.T, error) {
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

func DownloadImageAsync(dev device.T, worker command.Worker, src image.T) (<-chan *osimage.RGBA, error) {
	swizzle := false
	switch src.Format() {
	case vk.FormatB8g8r8a8Unorm:
		swizzle = true
	case vk.FormatR8g8b8a8Unorm:
		break
	default:
		return nil, fmt.Errorf("unsupported source format")
	}

	dst, err := image.New(dev, image.Args{
		Type:    vk.ImageType2d,
		Width:   src.Width(),
		Height:  src.Height(),
		Depth:   1,
		Layers:  1,
		Levels:  1,
		Format:  vk.FormatR8g8b8a8Unorm,
		Memory:  vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit,
		Tiling:  vk.ImageTilingLinear,
		Usage:   vk.ImageUsageTransferDstBit,
		Sharing: vk.SharingModeExclusive,
		Layout:  vk.ImageLayoutUndefined,
	})
	if err != nil {
		return nil, err
	}

	// transfer data from texture buffer
	worker.Queue(func(cmd command.Buffer) {
		cmd.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			src,
			vk.ImageLayoutUndefined,
			vk.ImageLayoutTransferSrcOptimal,
			vk.ImageAspectColorBit)
		cmd.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			dst,
			vk.ImageLayoutUndefined,
			vk.ImageLayoutTransferDstOptimal,
			vk.ImageAspectColorBit)
		cmd.CmdCopyImage(src, vk.ImageLayoutTransferSrcOptimal, dst, vk.ImageLayoutTransferDstOptimal, vk.ImageAspectColorBit)
		cmd.CmdImageBarrier(
			vk.PipelineStageTransferBit,
			vk.PipelineStageFragmentShaderBit,
			src,
			vk.ImageLayoutTransferSrcOptimal,
			vk.ImageLayoutColorAttachmentOptimal,
			vk.ImageAspectColorBit)
		cmd.CmdImageBarrier(
			vk.PipelineStageTopOfPipeBit,
			vk.PipelineStageTransferBit,
			dst,
			vk.ImageLayoutTransferDstOptimal,
			vk.ImageLayoutGeneral,
			vk.ImageAspectColorBit)
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
