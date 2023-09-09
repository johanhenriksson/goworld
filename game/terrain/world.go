package terrain

import (
	"fmt"
	"log"
	"sync"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

type World struct {
	object.Object
	terrain  *Map
	distance float32

	lock   sync.Mutex
	active map[string]object.Component
	ready  chan tileSpawn
}

type tileSpawn struct {
	Object   object.Object
	Key      string
	Position vec3.T
}

// Builds a world of tiles around the active camera as it moves around
func NewWorld(terrain *Map, distance float32) *World {
	return object.New("World", &World{
		terrain:  terrain,
		distance: distance,
		active:   make(map[string]object.Component, 100),
		ready:    make(chan tileSpawn, 100),
	})
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
			log.Println("destroy tile", key)
			tile.Destroy()
			delete(c.active, key)
		}
	}

	// create tiles close to us
	size := c.terrain.TileSize
	tilePos := pos.Scaled(1 / float32(size)).Floor()
	cx, cz := int(tilePos.X), int(tilePos.Z)

	steps := int(c.distance / float32(size))
	minDist := math.InfPos
	var spawn func()
	for x := cx - steps; x < cx+steps; x++ {
		for z := cz - steps; z < cz+steps; z++ {
			// check if the tile would have been in range
			p := vec3.NewI(x*size, 0, z*size)
			dist := vec3.Distance(pos, p)
			if dist > c.distance {
				continue
			}

			// check if its already active
			key := fmt.Sprintf("Tile:%d,%d", x, z)
			_, active := c.active[key]
			if active {
				continue
			}

			// spawn it
			if dist < minDist {
				minDist = dist
				ix, iz := x, z
				spawn = func() {
					log.Println("spawn tile", key)
					c.lock.Lock()
					c.active[key] = nil
					c.lock.Unlock()

					tile := c.terrain.GetTile(ix, iz, true)
					c.ready <- tileSpawn{
						Key:      key,
						Position: p,

						Object: object.Builder(object.Empty(key)).
							Attach(NewMesh(tile)).
							Attach(physics.NewRigidBody(0)).
							Attach(physics.NewMesh()).
							Position(p).
							Create(),
					}
				}
			}
		}
	}
	if spawn != nil {
		go spawn()
	}
}
