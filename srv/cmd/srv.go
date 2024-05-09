package main

import (
	"time"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/srv"
)

func main() {
	area := srv.NewSimpleArea()

	unit := srv.NewUnit("testy")
	uid := area.Join(unit)

	srv.NewAIController(unit, srv.Behavior{
		"idle": srv.NewTaskLoop(
			&srv.MoveToTask{
				Target: vec3.T{X: 3, Y: 0, Z: 0},
				Speed:  1.4,
			},
			&srv.MoveToTask{
				Target: vec3.T{X: 0, Y: 0, Z: 0},
				Speed:  1.4,
			},
		),
	})

	player := srv.NewUnit("player")
	area.Join(player)

	inputs := make(chan srv.Action)
	srv.NewPlayerController(area, player, inputs)

	time.Sleep(2 * time.Second)
	inputs <- srv.SetTargetAction{Target: uid}
	inputs <- srv.CastSpellAction{Target: uid, Spell: "fireball"}

	time.Sleep(8 * time.Second)
}
