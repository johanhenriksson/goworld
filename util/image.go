package util

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"
)

/* Loads an image from file. Returns an RGBA image object */
func LoadImage(file string) (*image.RGBA, error) {
	imgFile, err := os.Open(ExePath + file)
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
