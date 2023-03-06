package image

import (
	imglib "image"
	"image/draw"

	// image codecs
	_ "image/png"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Data struct {
	Width  int
	Height int
	Format core1_0.Format
	Buffer []byte
}

func LoadFile(file string) (*Data, error) {
	imgFile, err := assets.Open(file)
	if err != nil {
		return nil, err
	}
	img, _, err := imglib.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := imglib.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, imglib.Point{0, 0}, draw.Src)

	return &Data{
		Width:  rgba.Rect.Size().X,
		Height: rgba.Rect.Size().Y,
		Format: FormatRGBA8Unorm,
		Buffer: rgba.Pix,
	}, nil
}
