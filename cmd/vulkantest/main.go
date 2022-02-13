package main

import (
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

func main() {
	backend := vulkan.New("goworld: vulkan", 0)
	defer backend.Destroy()

	wnd, err := window.New(backend, window.Args{
		Title:  "goworld: vulkan",
		Width:  500,
		Height: 500,
	})
	if err != nil {
		panic(err)
	}

	for !wnd.ShouldClose() {
		// aquire backbuffer image
		backend.Aquire()

		// draw
		backend.Submit([]vk.CommandBuffer{})

		backend.Present()

		wnd.Poll()
	}
}
