package texture

import (
	"image"
	"image/draw"
	"os"

	_ "image/png"
)

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
