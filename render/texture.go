package render

import (
    "os"
    "image"
    "image/draw"
	_ "image/png"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
    Id      uint32
}

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

func (tx *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, tx.Id)
}

func (tx *Texture) Buffer(img *image.RGBA) {
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Rect.Size().X),
		int32(img.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix))
}

func LoadTexture(file string) (*Texture, error) {
    img, err := LoadImage(file)
    if err != nil {
        return nil, err
    }

    tx := CreateTexture()
    tx.Buffer(img)
    return tx, nil
}

func LoadImage(file string) (*image.RGBA, error) {
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
