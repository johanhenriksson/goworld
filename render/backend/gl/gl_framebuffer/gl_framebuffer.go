package gl_framebuffer

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/math/vec2"
	gltex "github.com/johanhenriksson/goworld/render/backend/gl/gl_texture"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type glframebuf struct {
	id      uint32
	width   int
	height  int
	targets []framebuffer.Target
	buffers []framebuffer.Buffer
}

// CreateFrameBuffer creates a new frame buffer object with a given size
func New(width, height int) framebuffer.T {
	f := &glframebuf{
		width:   width,
		height:  height,
		buffers: []framebuffer.Buffer{},
		targets: make([]framebuffer.Target, 0, 8),
	}
	gl.GenFramebuffers(1, &f.id)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	return f
}

func (f *glframebuf) Width() int   { return f.width }
func (f *glframebuf) Height() int  { return f.height }
func (f *glframebuf) Size() vec2.T { return vec2.NewI(f.width, f.height) }

// NewBuffer creates a new frame buffer texture and attaches it to the given target.
// Returns a pointer to the created texture object. FBO must be bound first.
func (f *glframebuf) NewBuffer(target framebuffer.Target, internalFormat, format texture.Format, datatype types.Type) texture.T {
	// Create texture object
	texture := gltex.New(f.width, f.height)
	texture.SetFormat(format)
	texture.SetInternalFormat(internalFormat)
	texture.SetDataType(datatype)
	texture.Clear()

	// attach texture
	f.AttachBuffer(target, texture)

	return texture
}

// AttachBuffer attaches a texture to the given frame buffer target
func (f *glframebuf) AttachBuffer(target framebuffer.Target, tex texture.T) {
	gl.FramebufferTexture(
		gl.FRAMEBUFFER,
		uint32(target),
		uint32(tex.ID()),
		int32(tex.MipLevel()))

	// Attach to frame buffer
	f.buffers = append(f.buffers, framebuffer.Buffer{
		Target:  target,
		Texture: tex,
	})

	// todo: we want to avoid multiple depth attachments
	if target != gl.DEPTH_ATTACHMENT {
		// add the target to the list of enabled draw buffers
		f.targets = append(f.targets, target)
	}
}

// Bind the frame buffer for drawing
func (f *glframebuf) Bind() {
	// set viewport size equal to buffer size
	gl.Viewport(0, 0, int32(f.width), int32(f.height))

	// bind this frame buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
}

// Unbind the frame buffer
func (f *glframebuf) Unbind() {
	// finish drawing
	gl.Flush()

	// unbind
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	// gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
}

// Delete the frame buffer object
func (f *glframebuf) Delete() {
	if f.id == 0 {
		panic("Cant delete framebuffer 0")
	}
	gl.DeleteFramebuffers(1, &f.id)
	f.id = 0
}

// Sample the color buffer at a given coordinate.
func (f *glframebuf) Sample(target framebuffer.Target, pos vec2.T) (color.T, bool) {
	pixel := make([]float32, 4)
	gl.ReadBuffer(uint32(target))
	gl.ReadPixels(int32(pos.X), int32(pos.Y), 1, 1, gl.RGBA, gl.FLOAT, unsafe.Pointer(&pixel[0]))
	return color.RGBA(pixel[0], pixel[1], pixel[2], pixel[3]), true
}

// SampleDepth samples the depth buffer at a given coordinate.
func (f *glframebuf) SampleDepth(pos vec2.T) (float32, bool) {
	float := float32(0)
	gl.ReadPixels(int32(pos.X), int32(pos.Y), 1, 1, gl.DEPTH_COMPONENT, gl.FLOAT, unsafe.Pointer(&float))
	return float, true
}

// DrawBuffers sets up all the attached buffers for drawing
func (f *glframebuf) DrawBuffers() {
	gl.DrawBuffers(int32(len(f.targets)), (*uint32)(&f.targets[0]))
}

func (f *glframebuf) Resize(width, height int) {
	// ensure that the size has actually changed first
	if f.width == width && f.height == height {
		return
	}

	f.width = width
	f.height = height
	for _, buffer := range f.buffers {
		buffer.Texture.Resize(width, height)
	}
}
