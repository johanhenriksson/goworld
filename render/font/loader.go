package font

import (
	"fmt"
	"log"
	"sync"

	"github.com/golang/freetype/truetype"
	fontlib "golang.org/x/image/font"

	"github.com/johanhenriksson/goworld/assets"
)

var parseCache map[string]*truetype.Font = make(map[string]*truetype.Font, 32)
var faceCache map[string]T = make(map[string]T, 128)

func loadTruetypeFont(name string) (*truetype.Font, error) {
	// check parsed font cache
	if fnt, exists := parseCache[name]; exists {
		return fnt, nil
	}

	fontBytes, err := assets.ReadAll(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	// add to cache
	parseCache[name] = fnt
	return fnt, nil
}

func Load(name string, size int, scale float32) T {
	key := fmt.Sprintf("%s:%dx%.2f", name, size, scale)
	if font, exists := faceCache[key]; exists {
		return font
	}

	log.Printf("+ font %s %dpt x%.2f\n", name, size, scale)

	ttf, err := loadTruetypeFont(name)
	if err != nil {
		panic(err)
	}

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
		fnt:    ttf,
		face:   face,
		glyphs: make(map[rune]*Glyph, 128),
		drawer: &fontlib.Drawer{Face: face},
		mutex:  &sync.Mutex{},
	}

	faceCache[key] = fnt
	return fnt
}
