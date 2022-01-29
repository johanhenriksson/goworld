package texture

import (
	"image"
	"image/png" // png support
	"os"

	ogl "github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/backend/gl"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/texture"
)

type gltexture struct {
	id       uint32
	width    int
	height   int
	format   texture.Format
	internal texture.Format
	filter   texture.Filter
	wrap     texture.WrapMode
	datatype gl.Type
	miplevel int32
	border   int
}

// New creates a new 2D texture and sets some sane defaults
func New(width, height int) texture.T {
	var id uint32
	ogl.GenTextures(1, &id)

	tx := &gltexture{
		id:       id,
		width:    width,
		height:   height,
		format:   ogl.RGBA,
		internal: ogl.RGBA,
		datatype: gl.UInt8,
		border:   0,
	}
	tx.Bind()

	/* Texture parameters - pass as parameters? */
	tx.SetFilter(texture.LinearFilter)
	tx.SetWrapMode(texture.ClampWrap)

	return tx
}

func (tx *gltexture) Filter() texture.Filter         { return tx.filter }
func (tx *gltexture) WrapMode() texture.WrapMode     { return tx.wrap }
func (tx *gltexture) Format() texture.Format         { return tx.format }
func (tx *gltexture) InternalFormat() texture.Format { return tx.internal }
func (tx *gltexture) Width() int                     { return tx.width }
func (tx *gltexture) Height() int                    { return tx.height }
func (tx *gltexture) DataType() types.Type           { return tx.datatype.Cast() }

func (tx *gltexture) Size() vec2.T {
	return vec2.New(float32(tx.width), float32(tx.height))
}

func (tx *gltexture) SetFilter(filter texture.Filter) {
	tx.filter = filter
	gl.SetTexture2DFilter(filter, filter)
}

func (tx *gltexture) SetWrapMode(mode texture.WrapMode) {
	tx.wrap = mode
	gl.SetTexture2DWrapMode(mode, mode)
}

func (tx *gltexture) SetFormat(fmt texture.Format) {
	tx.format = fmt
}

func (tx *gltexture) SetInternalFormat(fmt texture.Format) {
	tx.internal = fmt
}

func (tx *gltexture) SetDataType(t types.Type) {
	tx.datatype = gl.TypeCast(t)
}

// Use binds this texture to the given texture slot
func (tx *gltexture) Use(slot int) {
	if err := gl.ActiveTexture(gl.TextureSlot(slot)); err != nil {
		panic(err)
	}
	tx.Bind()
}

// Bind texture to the currently active texture slot
func (tx *gltexture) Bind() {
	ogl.BindTexture(ogl.TEXTURE_2D, tx.id)
	switch ogl.GetError() {
	case ogl.INVALID_ENUM:
		panic("texture target is not one of the allowable values")
	case ogl.INVALID_VALUE:
		panic("texture is not a name returned from a previous call to glGenTextures")
	case ogl.INVALID_OPERATION:
		panic("texture was previously created with a target that doesn't match that of target.")
	}
}

// FrameBufferTarget attaches this texture to the current frame buffer object
func (tx *gltexture) FrameBufferTarget(attachment uint32) {
	ogl.FramebufferTexture(ogl.FRAMEBUFFER, attachment, tx.id, tx.miplevel)
}

// Clear the texture
func (tx *gltexture) Clear() {
	tx.Bind()
	ogl.TexImage2D(
		ogl.TEXTURE_2D,
		0,
		int32(tx.internal), // gl.RGBA,
		int32(tx.width), int32(tx.height),
		0,
		uint32(tx.format),   //gl.RGBA,
		uint32(tx.datatype), // gl.UNSIGNED_BYTE,
		nil)                 // null ptr
}

func (tx *gltexture) Resize(width, height int) {
	// ensure that the size has actually changed first
	if tx.width == width && tx.height == height {
		return
	}

	tx.width = width
	tx.height = height
	tx.Clear()
}

// BufferImage buffers texture data from an image object
func (tx *gltexture) BufferImage(img *image.RGBA) {
	tx.width = img.Rect.Size().X
	tx.height = img.Rect.Size().Y
	tx.datatype = ogl.UNSIGNED_BYTE
	ogl.TexImage2D(
		ogl.TEXTURE_2D,
		0,
		int32(tx.internal),
		int32(tx.width), int32(tx.height),
		0,
		uint32(tx.format),
		uint32(tx.datatype),
		ogl.Ptr(img.Pix))
}

// BufferFloats buffers texture data from a float array
func (tx *gltexture) BufferFloats(img []float32) {
	tx.datatype = gl.Float
	ogl.TexImage2D(
		ogl.TEXTURE_2D,
		0,
		int32(tx.internal),
		int32(tx.width), int32(tx.height),
		0,
		uint32(tx.format),
		uint32(tx.datatype),
		ogl.Ptr(&img[0]))
}

func (tx *gltexture) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(tx.width), int(tx.height))
}

func (tx *gltexture) ToImage() *image.RGBA {
	tx.Bind()
	img := image.NewRGBA(tx.Bounds())
	ogl.GetTexImage(ogl.TEXTURE_2D, 0, ogl.RGBA, ogl.UNSIGNED_BYTE, ogl.Ptr(img.Pix))
	return img
}

func (tx *gltexture) Save(filename string) error {
	img := tx.ToImage()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
