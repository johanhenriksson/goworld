package terrain

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/sprite"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/texture"
)

type World struct {
	object.Object
	Terrain *Map

	distance float32
	lock     sync.Mutex
	active   map[string]object.Component
	ready    chan tileSpawn
}

type tileSpawn struct {
	Object   object.Object
	Key      string
	Position vec3.T
}

// Builds a world of tiles around the active camera as it moves around
func NewWorld(terrain *Map, distance float32) *World {
	return object.New("World", &World{
		Terrain:  terrain,
		distance: distance,
		active:   make(map[string]object.Component, 100),
		ready:    make(chan tileSpawn, 100),
	})
}

func (w *World) EditorUpdate(scene object.Component, dt float32) {
	w.Update(scene, dt)
}

func (c *World) Update(scene object.Component, dt float32) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// update tiles
	c.Object.Update(scene, dt)

	// insert any new tiles
	select {
	case spawn := <-c.ready:
		c.active[spawn.Key] = spawn.Object
		object.Attach(c, spawn.Object)
	default:
	}

	// find the active camera
	root := object.Root(scene)
	cam := object.GetInChildren[*camera.Camera](root)
	if cam == nil {
		log.Println("terrain world: no active camera")
		return
	}

	size := c.Terrain.TileSize
	// half := vec3.NewI(size/2, 0, size/2)
	pos := cam.Transform().WorldPosition()
	pos.Y = 0

	// destroy tiles that are too far away
	for key, tile := range c.active {
		if tile == nil {
			// being loaded
			continue
		}
		dist := vec3.Distance(pos, tile.Transform().Position())
		if dist > c.distance*1.1 {
			tile.Destroy()
			delete(c.active, key)
		}
	}

	// create tiles close to us
	tilePos := pos.Scaled(1 / float32(size)).Floor()
	cx, cz := int(tilePos.X), int(tilePos.Z)

	steps := int(c.distance / float32(size))
	minDist := math.InfPos
	var spawn func()
	var spawnKey string
	for x := cx - steps; x < cx+steps; x++ {
		for z := cz - steps; z < cz+steps; z++ {
			// check if the center of tile would have been in range
			p := vec3.NewI(x*size, 0, z*size)
			dist := vec3.Distance(pos, p)
			if dist > c.distance {
				continue
			}
			if dist > minDist {
				continue
			}

			// check if its already active
			key := fmt.Sprintf("Tile:%d,%d", x, z)
			if v, active := c.active[key]; active && v != nil {
				continue
			}

			// spawn it
			minDist = dist
			spawnKey = key

			ix, iz := x, z
			spawn = func() {
				assetKey := fmt.Sprintf("maps/%s/tile_%d_%d", c.Terrain.Name, ix, iz)

				var tile object.Object
				if tileCmp, err := object.Load(assetKey); err == nil {
					tile = tileCmp.(object.Object)

					mesh := object.GetInChildren[*Mesh](tile)
					if mesh == nil {
						log.Println("failed to load tile:", tile.Name())
						tile = nil
					} else {
						c.Terrain.AddTile(mesh.Tile)
					}
				} else if errors.Is(err, assets.ErrNotFound) {
					tileData := c.Terrain.Tile(ix, iz, true)
					tile = DefaultWorldTile(key, c.Terrain, tileData)
					tile.Transform().SetPosition(p)
				} else {
					panic(err)
				}

				c.ready <- tileSpawn{
					Key:      key,
					Position: p,

					Object: tile,
				}
			}
		}
	}
	if spawn != nil {
		// mark key as active before we release the lock
		c.active[spawnKey] = nil
		go spawn()
	}
}

type TileBuilderFn func(terrain *Map, tileData *Tile) object.Object

func DefaultWorldTile(key string, terrain *Map, tileData *Tile) object.Object {
	tile := object.Builder(object.Empty(key)).
		Attach(NewMesh(tileData)).
		Attach(physics.NewRigidBody(0)).
		Attach(physics.NewMesh()).
		Create()

	// bushes
	type RandomSprite struct {
		Texture   string
		SizeMin   float32
		SizeMax   float32
		CountMin  int
		CountMax  int
		Collision bool
	}
	sprites := []RandomSprite{
		{
			Texture:  "sprites/objects/flower3.png",
			SizeMin:  0.5,
			SizeMax:  1.5,
			CountMax: 10,
		},
		{
			Texture:   "sprites/objects/tree1.png",
			SizeMin:   7,
			SizeMax:   12,
			CountMin:  3,
			CountMax:  9,
			Collision: true,
		},
		{
			Texture:  "sprites/objects/flower2.png",
			SizeMin:  0.5,
			SizeMax:  1.5,
			CountMin: 2,
			CountMax: 10,
		},
		{
			Texture:  "sprites/objects/bush2.png",
			SizeMin:  1,
			SizeMax:  3,
			CountMax: 10,
		},
	}
	for _, s := range sprites {
		for i := 0; i < random.Int(s.CountMin, s.CountMax); i++ {
			size := random.Range(s.SizeMin, s.SizeMax)
			prop := object.Builder(object.Empty("Prop")).
				Attach(sprite.New(sprite.Args{
					Size: vec2.New(size, size),
					Texture: texture.PathArgsRef(s.Texture, texture.Args{
						Filter:  texture.FilterNearest,
						Mipmaps: true,
					}),
				})).
				Position(vec3.New(
					random.Range(0, float32(terrain.TileSize)),
					size/2,
					random.Range(0, float32(terrain.TileSize)),
				)).
				Create()

			if s.Collision {
				object.Attach(prop, physics.NewSphere(0.8*size/2))
				object.Attach(prop, physics.NewRigidBody(0))
			}

			object.Attach(tile, prop)
		}
	}
	return tile
}
