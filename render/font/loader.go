package font

import (
	"fmt"
	"log"
	"sync"

	"github.com/golang/freetype/truetype"
	fontlib "golang.org/x/image/font"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/util"
)

var parseCache map[string]*truetype.Font = make(map[string]*truetype.Font, 32)
var faceCache map[string]T = make(map[string]T, 128)

func loadTruetypeFont(filename string) (*truetype.Font, error) {
	// check parsed font cache
	if fnt, exists := parseCache[filename]; exists {
		return fnt, nil
	}

	fontBytes, err := assets.Read(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	// add to cache
	parseCache[filename] = fnt
	return fnt, nil
}

func Load(filename string, size int, scale float32) T {
	key := fmt.Sprintf("%s:%dx%.2f", filename, size, scale)
	if font, exists := faceCache[key]; exists {
		return font
	}

	ttf, err := loadTruetypeFont(filename)
	if err != nil {
		panic(err)
	}

	name := ttf.Name(truetype.NameIDFontFullName)
	log.Printf("+ font %s %dpt x%.2f\n", name, size, scale)

	dpi := 72.0 * scale
	face := truetype.NewFace(ttf, &truetype.Options{
		Size:       float64(size),
		DPI:        float64(dpi),
		Hinting:    fontlib.HintingFull,
		SubPixelsX: 8,
		SubPixelsY: 8,
	})

	fnt := &font{
		size:   float32(size),
		scale:  scale,
		name:   name,
		face:   face,
		glyphs: util.NewSyncMap[rune, *Glyph](),
		kern:   util.NewSyncMap[runepair, float32](),
		drawer: &fontlib.Drawer{Face: face},
		mutex:  &sync.Mutex{},
	}

	faceCache[key] = fnt
	return fnt
}
