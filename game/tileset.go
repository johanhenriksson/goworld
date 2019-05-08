package game

const TileSize = 16
const TilesetTexWidth = 4096
const TilesetTexHeight = 4096

type TileId uint16
type TileCoord uint8

type Tileset struct {
	Width  int
	Height int
	Size   int
	Tiles  []Tile
}

type Tile struct {
	Id TileId
	X  TileCoord
	Y  TileCoord
}

func CreateTileset() *Tileset {
	ts := &Tileset{
		Size:   TileSize,
		Width:  TilesetTexWidth / TileSize,
		Height: TilesetTexHeight / TileSize,
	}
	ts.Generate()
	return ts
}

func (ts *Tileset) Generate() {
	ts.Tiles = make([]Tile, ts.Width*ts.Height)
	for y := 0; y < ts.Height; y++ {
		for x := 0; x < ts.Width; x++ {
			id := y*ts.Width + x
			ts.Tiles[id] = Tile{
				Id: TileId(id),
				X:  TileCoord(x),
				Y:  TileCoord(y),
			}
		}
	}
}

func (ts *Tileset) GetId(x, y int) TileId {
	return TileId(y*ts.Width + x)
}

func (ts *Tileset) Get(id TileId) *Tile {
	return &ts.Tiles[int(id)]
}
