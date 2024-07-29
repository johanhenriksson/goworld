package swapchain

import (
	"github.com/johanhenriksson/goworld/render/sync"
)

type Context struct {
	Index          int
	ImageAvailable sync.Semaphore
	RenderComplete sync.Semaphore
}

func DummyContext() *Context {
	return &Context{}
}
