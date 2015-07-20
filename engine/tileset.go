package engine

const TileSize = 16
const TilesetTexWidth = 4096
const TilesetTexHeight = 4096

type Tileset struct {
    Width   int
    Height  int
    Size    int
    Tiles   []Tile
}

type Tile struct {
    Id           uint16
    X            uint8
    Y            uint8
}

func CreateTileset() *Tileset {
    ts := &Tileset {
        Size:   TileSize,
        Width:  TilesetTexWidth / TileSize,
        Height: TilesetTexHeight / TileSize,
    }
    ts.Generate()
    return ts
}

func (ts *Tileset) Generate() {
    ts.Tiles = make([]Tile, ts.Width * ts.Height)
    for y := 0; y < ts.Height; y++ {
        for x := 0; x < ts.Width; x++ {
            id := y * ts.Width + x
            ts.Tiles[id] = Tile {
                Id: uint16(id),
                X:  uint8(x),
                Y:  uint8(y),
            }
        }
    }
}

func (ts *Tileset) GetId(x, y int) uint16 {
    return uint16(y * ts.Width + x)
}

func (ts *Tileset) Get(id uint16) *Tile {
    return &ts.Tiles[int(id)]
}
