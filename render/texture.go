package render

import (
    "os"
    "image"
    "image/draw"
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
    "github.com/johanhenriksson/goworld/util"
)

type Texture struct {
    Id      uint32
    Width   int32
    Height  int32
}

/* Creates a new GL texture and sets basic options */
func CreateTexture() *Texture {
	var id uint32
	gl.GenTextures(1, &id)

    tx := &Texture {
        Id: id,
    }
    tx.Bind(gl.TEXTURE0)

    /* Texture parameters */
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

    return tx
}

/* Binds this texture to the given slot and activates it */
func (tx *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, tx.Id)
}

/* Buffers texture data to GPU memory */
func (tx *Texture) Buffer(img *image.RGBA) {
	tx.Width  = int32(img.Rect.Size().X)
    tx.Height = int32(img.Rect.Size().Y)

    /* Buffer image data */
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
        tx.Width, tx.Height,
		0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix))
}

/* Loads a texture from file */
func LoadTexture(file string) (*Texture, error) {
    img, err := LoadImage(file)
    if err != nil {
        return nil, err
    }

    tx := CreateTexture()
    tx.Buffer(img)
    return tx, nil
}

/* Loads an image from file. Returns an RGBA image object */
func LoadImage(file string) (*image.RGBA, error) {
	imgFile, err := os.Open(util.ExePath + file)
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
