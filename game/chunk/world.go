package chunk

import (
	"fmt"
	"log"
	"sync"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type World struct {
	object.Object
	size      int
	distance  float32
	generator Generator

	lock   *sync.Mutex
	active map[string]object.Component
	ready  chan *T
}

// Builds a world of chunks around the active camera as it moves around
func NewWorld(size int, generator Generator, distance float32) *World {
	return object.New("World", &World{
		size:      size,
		generator: generator,
		distance:  distance,
		active:    make(map[string]object.Component, 100),
		ready:     make(chan *T, 100),
		lock:      &sync.Mutex{},
	})
}

func (c *World) Update(scene object.Component, dt float32) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// update chunks
	c.Object.Update(scene, dt)

	// find the active camera
	root := object.Root(scene)
	cam := object.GetInChildren[*camera.T](root)
	if cam == nil {
		return
	}

	pos := cam.Transform().WorldPosition()
	pos.Y = 0

	// insert any new chunks
	select {
	case chk := <-c.ready:
		key := fmt.Sprintf("Chunk:%d,%d", chk.Cx, chk.Cz)
		chonk := object.Builder(object.Empty(key)).
			Attach(NewMesh(chk)).
			Attach(box.New(box.Args{
				Size:  vec3.NewI(c.size, c.size, c.size),
				Color: color.Purple,
			})).
			Position(vec3.NewI(chk.Cx*c.size, 0, chk.Cz*c.size)).
			Parent(c).
			Create()
		c.active[key] = chonk
	default:
	}

	// destroy chunks that are too far away
	for key, chunk := range c.active {
		if chunk == nil {
			// being loaded
			continue
		}
		dist := vec3.Distance(pos, chunk.Transform().Position())
		if dist > c.distance*1.1 {
			log.Println("destroy chunk", key)
			chunk.Destroy()
			delete(c.active, key)
		}
	}

	// create chunks close to us
	chunkPos := pos.Scaled(1 / float32(c.size)).Floor()
	cx, cz := int(chunkPos.X), int(chunkPos.Z)

	steps := int(c.distance / float32(c.size))
	minDist := math.InfPos
	var spawn func()
	for x := cx - steps; x < cx+steps; x++ {
		for z := cz - steps; z < cz+steps; z++ {
			// check if the chunk would have been in range
			p := vec3.NewI(x*c.size, 0, z*c.size)
			dist := vec3.Distance(pos, p)
			if dist > c.distance {
				continue
			}

			// check if its already active
			key := fmt.Sprintf("Chunk:%d,%d", x, z)
			_, active := c.active[key]
			if active {
				continue
			}

			// spawn it
			if dist < minDist {
				minDist = dist
				ix, iz := x, z
				spawn = func() {
					log.Println("spawn chunk", key)
					c.lock.Lock()
					c.active[key] = nil
					c.lock.Unlock()

					chunkData := Generate(c.generator, c.size, ix, iz)
					c.ready <- chunkData
				}
			}
		}
	}
	if spawn != nil {
		go spawn()
	}
}
