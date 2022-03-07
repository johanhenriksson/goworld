package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"

	vk "github.com/vulkan-go/vulkan"
)

type Textures interface {
	Fetch(path string) vk_texture.T
	Destroy()
}

type vktextures struct {
	backend vulkan.T
	worker  command.Worker
	cache   map[string]vk_texture.T
}

func NewVkTextures(backend vulkan.T) Textures {
	return &vktextures{
		backend: backend,
		worker:  backend.Transferer(),
		cache:   make(map[string]vk_texture.T),
	}
}

func (t *vktextures) Fetch(path string) vk_texture.T {
	if cached, hit := t.cache[path]; hit {
		return cached
	}

	img, err := render.ImageFromFile(path)
	if err != nil {
		panic(err)
	}

	stage := buffer.NewShared(t.backend.Device(), len(img.Pix))
	stage.Write(img.Pix, 0)

	tex := vk_texture.New(t.backend.Device(), vk_texture.Args{
		Width:  img.Rect.Size().X,
		Height: img.Rect.Size().Y,
		Format: vk.FormatA8b8g8r8UintPack32,
		Filter: vk.FilterLinear,
		Wrap:   vk.SamplerAddressModeRepeat,
	})

	t.worker.Queue(func(cmd command.Buffer) {
		cmd.CmdImageBarrier(
			vk.PipelineStageFlags(vk.PipelineStageTopOfPipeBit),
			vk.PipelineStageFlags(vk.PipelineStageTransferBit),
			tex.Image(),
			vk.ImageLayoutUndefined,
			vk.ImageLayoutTransferDstOptimal)
		cmd.CmdCopyBufferToImage(stage, tex.Image(), vk.ImageLayoutTransferDstOptimal)
	})
	t.worker.Submit(command.SubmitInfo{})
	t.worker.Wait()

	stage.Destroy()

	log.Println("buffered texture", path)
	t.cache[path] = tex

	return tex
}

func (t *vktextures) Destroy() {
	for _, tex := range t.cache {
		tex.Destroy()
	}
	t.cache = make(map[string]vk_texture.T)
}
