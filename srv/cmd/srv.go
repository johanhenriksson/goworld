package main

import (
	"time"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/srv"
)

func main() {
	area := srv.NewSimpleArea()

	client := &srv.DummyClient{
		Token: srv.ClientToken{
			Character: "player",
		},
	}

	realm := srv.NewRealm(area)
	realm.Accept(client)

	npc := srv.NewUnit("testy")
	uid := area.Join(npc)

	srv.NewAIController(npc, srv.Behavior{
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

	time.Sleep(2 * time.Second)

	id := client.Actor.ID()
	client.Action(srv.SetTargetAction{Unit: id, Target: uid})
	client.Action(srv.CastSpellAction{Unit: id, Target: uid, Spell: "fireball"})

	time.Sleep(8 * time.Second)
}
