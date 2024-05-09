package main

import (
	"time"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/srv"
)

func main() {

	unit := srv.NewUnit("testy")

	srv.NewAIController(unit, []srv.Task{
		&srv.MoveToTask{
			Target: vec3.T{X: 3, Y: 0, Z: 0},
			Speed:  1.4,
		},
		&srv.MoveToTask{
			Target: vec3.T{X: 0, Y: 0, Z: 0},
			Speed:  1.4,
		},
	})

	time.Sleep(8 * time.Second)
}
