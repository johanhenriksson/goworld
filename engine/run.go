package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
)

type SceneFunc func(renderer.T, object.T)
type RendererFunc func() renderer.T

type Args struct {
	Title    string
	Width    int
	Height   int
	Backend  window.GlfwBackend
	Renderer RendererFunc
}

func Run(args Args, scenefuncs ...SceneFunc) {
	log.Println("goworld")

	go RunProfilingServer(6060)
	interrupt := NewInterrupter()

	backend := args.Backend
	if backend == nil {
		panic("no backend provided")
	}
	defer backend.Destroy()

	if args.Renderer == nil {
		panic("no renderer given")
	}

	// create a window
	wnd, err := window.New(backend, window.Args{
		Title:  args.Title,
		Width:  args.Width,
		Height: args.Height,
	})
	if err != nil {
		panic(err)
	}

	var renderer renderer.T
	recreateRenderer := func() {
		if renderer != nil {
			renderer.Destroy()
		}
		renderer = args.Renderer()
	}
	recreateRenderer()
	defer func() {
		renderer.Destroy()
	}()

	// create scene
	scene := object.New("Scene")
	wnd.SetInputHandler(scene)
	for _, scenefunc := range scenefuncs {
		scenefunc(renderer, scene)
	}

	// run the render loop
	log.Println("ready")

	lastFrameTime := time.Now()
	for interrupt.Running() && !wnd.ShouldClose() {
		wnd.Poll()

		w, h := wnd.Size()
		screen := render.Screen{
			Width:  w,
			Height: h,
			Scale:  wnd.Scale(),
		}

		// find the first active camera
		camera := query.New[camera.T]().First(scene)
		if camera == nil {
			fmt.Println("no active camera in the scene")
			continue
		}

		// draw
		context, err := backend.Aquire()
		if err != nil {
			log.Println("swapchain recreated?? recreating renderer")
			recreateRenderer()
			continue
		}

		args := createRenderArgs(screen, camera)
		args.Context = context

		renderer.Draw(args, scene)
		backend.Present()

		// update scene
		endFrameTime := time.Now()
		elapsed := endFrameTime.Sub(lastFrameTime)
		lastFrameTime = endFrameTime
		scene.Update(float32(elapsed.Seconds()))
	}
}

func createRenderArgs(screen render.Screen, cam camera.T) render.Args {
	return render.Args{
		Projection: cam.Projection(),
		View:       cam.View(),
		VP:         cam.ViewProj(),
		MVP:        cam.ViewProj(),
		Transform:  mat4.Ident(),
		Position:   cam.Transform().WorldPosition(),
		Clear:      cam.ClearColor(),
		Viewport:   screen,
	}
}
