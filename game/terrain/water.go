package terrain

import (
	"fmt"
	"log"
	"sync"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Water struct {
	object.Object

	size     int
	distance float32
	lock     sync.Mutex
	active   map[string]object.Component
	ready    chan tileSpawn
}

// Builds a world of tiles around the active camera as it moves around
func NewWater(size int, distance float32) *Water {
	return object.New("Water", &Water{
		size:     size,
		distance: distance,
		active:   make(map[string]object.Component, 100),
		ready:    make(chan tileSpawn, 100),
	})
}

func (w *Water) EditorUpdate(scene object.Component, dt float32) {
	w.Update(scene, dt)
}

func (c *Water) Update(scene object.Component, dt float32) {
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
			// log.Println("destroy tile", key)
			tile.Destroy()
			delete(c.active, key)
		}
	}

	// create tiles close to us
	tilePos := pos.Scaled(1 / float32(c.size)).Floor()
	cx, cz := int(tilePos.X), int(tilePos.Z)

	steps := int(c.distance / float32(c.size))
	minDist := math.InfPos
	var spawn func()
	var spawnKey string
	for x := cx - steps; x < cx+steps; x++ {
		for z := cz - steps; z < cz+steps; z++ {
			// check if the center of tile would have been in range
			p := vec3.NewI(x*c.size, 0, z*c.size)
			dist := vec3.Distance(pos, p)
			if dist > c.distance {
				continue
			}
			if dist > minDist {
				continue
			}

			// check if its already active
			key := fmt.Sprintf("Water:%d,%d", x, z)
			if v, active := c.active[key]; active && v != nil {
				continue
			}

			// spawn it
			minDist = dist
			spawnKey = key

			spawn = func() {
				c.ready <- tileSpawn{
					Key:      key,
					Position: p,

					Object: object.Builder(NewWaterTile(float32(c.size))).
						Position(p).
						Create(),
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

type WaterTile struct {
	object.Object

	Mesh *plane.Mesh

	frame    int
	animTick float32
	nextAnim float32
}

func NewWaterTile(size float32) *WaterTile {
	mesh := plane.New(plane.Args{
		Size: vec2.New(size, size),
		Mat: &material.Def{
			Pass:         material.Forward,
			Shader:       "forward/water",
			VertexFormat: vertex.T{},
			DepthTest:    true,
			DepthWrite:   false,
			DepthFunc:    core1_0.CompareOpLessOrEqual,
			Primitive:    vertex.Triangles,
			CullMode:     vertex.CullBack,
			Transparent:  true,
		},
	})
	mesh.SetTexture(texture.Diffuse, texture.PathArgsRef("textures/terrain/water1.png", texture.Args{}))

	return object.New("Water", &WaterTile{
		Mesh: mesh,

		animTick: 0.2,
		frame:    1,
	})
}

func (w *WaterTile) Update(scene object.Component, dt float32) {
	w.nextAnim -= dt
	if w.nextAnim < 0 {
		w.nextAnim = random.Range(0.6, 1.2)
		w.frame++

		w.Mesh.SetTexture(texture.Diffuse, texture.PathArgsRef(fmt.Sprintf("textures/terrain/water%d.png", w.frame%2+1), texture.Args{
			Filter: texture.FilterNearest,
		}))
	}
}
