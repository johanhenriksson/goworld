package engine

import (
	"log"
	"runtime"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type SceneFunc func(object.Object)
type RendererFunc func(vulkan.App, vulkan.Target) renderer.T

type Args struct {
	Title    string
	Width    int
	Height   int
	Renderer RendererFunc
}

func Run(args Args, scenefuncs ...SceneFunc) {
	log.Println("goworld")
	runtime.LockOSThread()

	go RunProfilingServer(6060)
	interrupt := NewInterrupter()

	backend := vulkan.New("goworld", 0)
	defer backend.Destroy()

	if args.Renderer == nil {
		args.Renderer = renderer.NewGraph
	}

	// create a window
	wnd, err := vulkan.NewWindow(backend, vulkan.WindowArgs{
		Title:  args.Title,
		Width:  args.Width,
		Height: args.Height,
		Frames: 3,
	})
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	// create renderer
	renderer := args.Renderer(backend, wnd)
	defer func() {
		renderer.Destroy()
	}()

	// create scene
	scene := object.Empty("Scene")
	wnd.SetInputHandler(scene)
	for _, scenefunc := range scenefuncs {
		scenefunc(scene)
	}

	object.Attach(scene, NewStatsGUI())

	// run the render loop
	log.Println("ready")

	counter := NewFrameCounter(60)
	for interrupt.Running() && !wnd.ShouldClose() {
		// update scene
		wnd.Poll()
		counter.Update()
		scene.Update(scene, counter.Delta())

		// draw
		renderer.Draw(scene, counter.Elapsed(), counter.Delta())
	}
}

func RunGC() {
	start := time.Now()
	runtime.GC()
	elapsed := time.Since(start)
	if elapsed.Milliseconds() > 1 {
		log.Printf("slow GC cycle: %.2fms", elapsed.Seconds()*1000)
	}
}
