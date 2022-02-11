package main

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"

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

	device := backend.Device()

	input := []int{1, 2, 3}
	stage := buffer.NewShared(device, 8*3)
	defer stage.Destroy()
	stage.Write(input, 0)

	remote := buffer.NewRemote(device, 8*3)
	defer remote.Destroy()

	output := buffer.NewShared(device, 8*3)
	defer output.Destroy()

	cmdbuf := backend.CmdPool().Allocate(vk.CommandBufferLevelPrimary)
	defer cmdbuf.Destroy()
	cmdbuf.Begin()
	cmdbuf.CopyBuffer(stage, remote, vk.BufferCopy{SrcOffset: 0, DstOffset: 0, Size: 24})
	cmdbuf.CopyBuffer(remote, output, vk.BufferCopy{SrcOffset: 0, DstOffset: 0, Size: 24})
	cmdbuf.End()

	cmdbuf.SubmitSync(device.GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit)))

	result := make([]int, 3)
	output.Read(result, 0)
	fmt.Println("read back", result)

	for !wnd.ShouldClose() {
		// aquire backbuffer image
		backend.Aquire()

		// draw
		backend.Submit([]vk.CommandBuffer{})

		backend.Present()

		wnd.Poll()
	}
}
