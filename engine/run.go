package engine

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/window"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type SceneFunc func(*Renderer, object.T)

type Args struct {
	Title     string
	Width     int
	Height    int
	SceneFunc SceneFunc
}

func Run(args Args) {
	fmt.Println("goworld")

	// cpu profiling
	flag.Parse()
	if *cpuprofile != "" {
		os.MkdirAll("profiling", 0755)
		ppath := fmt.Sprintf("profiling/%s", *cpuprofile)
		f, err := os.Create(ppath)
		if err != nil {
			panic(err)
		}
		fmt.Println("writing cpu profiling output to", ppath)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// create a window
	wnd, err := window.New(window.Args{
		Title:  args.Title,
		Width:  args.Width,
		Height: args.Height,
	})
	if err != nil {
		panic(err)
	}

	// initialize graphics pipeline
	renderer := NewRenderer(wnd)

	// create scene
	scene := object.New("Scene")
	wnd.SetInputHandler(scene)
	args.SceneFunc(renderer, scene)

	// run the render loop
	fmt.Println("Ready")
	for !wnd.ShouldClose() {
		scene.Update(0.030)
		renderer.Draw(scene)
		wnd.SwapBuffers()
	}
}
