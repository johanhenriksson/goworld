package renderer

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/graph"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render/upload"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type T interface {
	Draw(scene object.Object, time, delta float32)
	Recreate()
	Screenshot()
	Destroy()
}

type rgraph struct {
	graph.T
	app vulkan.App
}

func NewGraph(app vulkan.App) T {
	r := &rgraph{
		app: app,
	}
	r.T = graph.New(app, func(g graph.T, target pass.RenderTarget, gbuffer pass.GeometryBuffer) {
		shadows := pass.NewShadowPass(app)
		shadowNode := g.Node(shadows)

		deferred := g.Node(pass.NewDeferredPass(app, target, gbuffer, shadows))
		deferred.After(shadowNode, core1_0.PipelineStageTopOfPipe)

		forward := g.Node(pass.NewForwardPass(app, target, gbuffer))
		forward.After(deferred, core1_0.PipelineStageTopOfPipe)

		ssaoPass := pass.NewAmbientOcclusionPass(app, gbuffer)
		ssao := g.Node(ssaoPass)
		ssao.After(forward, core1_0.PipelineStageTopOfPipe)

		blurPass := pass.NewBlurPass(app, ssaoPass.Target)
		blur := g.Node(blurPass)
		blur.After(ssao, core1_0.PipelineStageTopOfPipe)

		postPass := pass.NewPostProcessPass(app, target, blurPass.Target())
		post := g.Node(postPass)
		post.After(blur, core1_0.PipelineStageTopOfPipe)

		lines := g.Node(pass.NewLinePass(app, postPass.Target(), target))
		lines.After(post, core1_0.PipelineStageTopOfPipe)

		gui := g.Node(pass.NewGuiPass(app, postPass.Target()))
		gui.After(lines, core1_0.PipelineStageTopOfPipe)

		output := g.Node(pass.NewOutputPass(app, postPass.Target()))
		output.After(gui, core1_0.PipelineStageTopOfPipe)

		// editor forward
		// editor lines
	})
	return r
}

func (r *rgraph) Screenshot() {
	idx := 0
	r.app.Device().WaitIdle()
	source := r.app.Surfaces()[idx]
	ss, err := upload.DownloadImage(r.app.Device(), r.app.Worker(idx), source)
	if err != nil {
		panic(err)
	}
	filename := fmt.Sprintf("Screenshot-%s.png", time.Now().Format("2006-01-02_15-04-05"))
	if err := upload.SavePng(ss, filename); err != nil {
		panic(err)
	}
}

func (r *rgraph) Destroy() {
	if r.T != nil {
		r.T.Destroy()
	}
}
