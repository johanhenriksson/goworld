package render

import (
	"image"
	"io/ioutil"
	"math"
    "fmt"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
    File string
    Size float64
    DPI float64
    Spacing float64

    src *image.Uniform
    fnt *truetype.Font
    drawer *font.Drawer
}

func (f *Font) setup() {
    f.drawer = &font.Drawer {
        Src: f.src,
        Face: truetype.NewFace(f.fnt, &truetype.Options {
            Size:    f.Size,
            DPI:     f.DPI,
            Hinting: font.HintingNone,
        }),
    }
}

func (f *Font) Render(text string, width, height float32) *Texture {
	line := math.Ceil(f.Size * f.DPI / 72)
    //height := int(f.Spacing * line)
    //width := int(float64(f.drawer.MeasureString(text)) / f.Size)

    /* Create and attach destination image */
	rgba := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	f.drawer.Dst = rgba

    /* Draw text */
    f.drawer.Dot = fixed.P(0, int(line))
    f.drawer.DrawString(text)

    fmt.Println("Font texture size W:", width, "H:", height)
    tx := CreateTexture(int32(width), int32(height))
    tx.Buffer(rgba)
    return tx
}

/** Load a truetype font */
func LoadFont(filename string, size, dpi, spacing float64) *Font {
	fontBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

    fnt := &Font {
        Size: size,
        DPI: dpi,
        Spacing: spacing,

        fnt: f,
        src: image.Black,
    }
    fnt.setup()
    return fnt
}
