package render

import (
	"image"
	"image/draw"
	"os"

	gltex "github.com/johanhenriksson/goworld/render/backend/gl/gl_texture"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

// TextureFromImage is a helper method to create an OpenGL texture from an image object */
func TextureFromImage(img *image.RGBA) texture.T {
	width := img.Rect.Size().X
	height := img.Rect.Size().Y
	tx := gltex.New(width, height)
	tx.BufferImage(img)
	return tx
}

// TextureFromFile loads a texture from file */
func TextureFromFile(file string) (texture.T, error) {
	img, err := ImageFromFile(file)
	if err != nil {
		return nil, err
	}
	return TextureFromImage(img), nil
}

// TextureFromColor creates a 1x1 texture from a color
func TextureFromColor(color color.T) texture.T {
	tx := gltex.New(1, 1)
	tx.BufferFloats([]float32{color.R, color.G, color.B, color.A})
	return tx
}

// ImageFromFile loads an image from a file. Returns an RGBA image object */
func ImageFromFile(file string) (*image.RGBA, error) {
	// todo: http support?

	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}
