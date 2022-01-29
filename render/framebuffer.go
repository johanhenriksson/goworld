package render

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render/color"
)

// DrawBuffer holds a target texture of a frame buffer object
type DrawBuffer struct {
	Target  uint32 // GL attachment enum (DEPTH_ATTACHMENT, COLOR_ATTACHMENT etc)
	Texture *Texture
}

// FrameBuffer holds information about an OpenGL frame buffer object
type FrameBuffer struct {
	Buffers []DrawBuffer
	Width   int
	Height  int
	id      uint32
	mipLvl  int32
	targets []uint32
}

// ScreenBuffer is the frame buffer of the screen
var ScreenBuffer = FrameBuffer{
	Buffers: []DrawBuffer{},
	Width:   0,
	Height:  0,
	id:      0,
	mipLvl:  0,
	targets: []uint32{gl.COLOR_ATTACHMENT0},
}

// NewBuffer creates a new frame buffer texture and attaches it to the given target.
// Returns a pointer to the created texture object. FBO must be bound first.
func (f *FrameBuffer) NewBuffer(target, internalFormat, format, datatype uint32) *Texture {
	// Create texture object
	texture := CreateTexture(f.Width, f.Height)
	texture.Format = format
	texture.InternalFormat = internalFormat
	texture.DataType = datatype
	texture.Clear()

	// attach texture
	f.AttachBuffer(target, texture)

	return texture
}

// AttachBuffer attaches a texture to the given frame buffer target
func (f *FrameBuffer) AttachBuffer(target uint32, texture *Texture) {
	// Set texture as frame buffer target
	texture.FrameBufferTarget(target)

	// Attach to frame buffer
	f.Buffers = append(f.Buffers, DrawBuffer{
		Target:  target,
		Texture: texture,
	})

	if target != gl.DEPTH_ATTACHMENT {
		// add the target to the list of enabled draw buffers
		f.targets = append(f.targets, target)
	}
}

// CreateFrameBuffer creates a new frame buffer object with a given size
func CreateFrameBuffer(width, height int) *FrameBuffer {
	f := &FrameBuffer{
		Width:   width,
		Height:  height,
		Buffers: []DrawBuffer{},
		targets: make([]uint32, 0, 8),
	}
	gl.GenFramebuffers(1, &f.id)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	return f
}

// Bind the frame buffer for drawing
func (f *FrameBuffer) Bind() {
	// set viewport size equal to buffer size
	SetViewport(0, 0, f.Width, f.Height)

	// bind this frame buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
}

// Unbind the frame buffer
func (f *FrameBuffer) Unbind() {
	// finish drawing
	gl.Flush()

	// unbind
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	// gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
}

// Delete the frame buffer object
func (f *FrameBuffer) Delete() {
	if f.id == 0 {
		panic("Cant delete framebuffer 0")
	}
	gl.DeleteFramebuffers(1, &f.id)
	f.id = 0
}

// Sample the color buffer at a given coordinate.
func (f *FrameBuffer) Sample(target uint32, x, y int) color.T {
	pixel := make([]float32, 4)
	gl.ReadBuffer(target)
	gl.ReadPixels(int32(x), int32(y), 1, 1, gl.RGBA, gl.FLOAT, unsafe.Pointer(&pixel[0]))
	return color.RGBA(pixel[0], pixel[1], pixel[2], pixel[3])
}

// SampleDepth samples the depth buffer at a given coordinate.
func (f *FrameBuffer) SampleDepth(x, y int) float32 {
	float := float32(0)
	gl.ReadPixels(int32(x), int32(y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, unsafe.Pointer(&float))
	return float
}

// DrawBuffers sets up all the attached buffers for drawing
func (f *FrameBuffer) DrawBuffers() {
	gl.DrawBuffers(int32(len(f.targets)), &f.targets[0])
}

func (f *FrameBuffer) Resize(width, height int) {
	// ensure that the size has actually changed first
	if f.Width == width && f.Height == height {
		return
	}

	f.Width = width
	f.Height = height
	for _, buffer := range f.Buffers {
		buffer.Texture.Resize(width, height)
	}
}
