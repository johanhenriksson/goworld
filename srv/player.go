package srv

import "log"

// Player controller receives network inputs from the player and controls the character.
// its perhaps also responsible for sending world events to the player client?
type PlayerController struct {
	area      Area
	character Actor
	inputs    <-chan Action
}

func NewPlayerController(area Area, character Actor, inputs <-chan Action) *PlayerController {
	pc := &PlayerController{
		area:      area,
		character: character,
		inputs:    inputs,
	}
	go pc.loop()
	return pc
}

func (c *PlayerController) loop() {
	events := EventBuffer()
	unsub := c.character.Subscribe(c, func(ev Event) {
		events <- ev
	})
	defer unsub()

	for {
		select {
		// input event
		case input := <-c.inputs:
			log.Println("input", input)
			switch input := input.(type) {
			case SetTargetAction:
				input.Unit = c.character.ID()
				c.area.Action(input)
			case CastSpellAction:
				input.Unit = c.character.ID()
				c.area.Action(input)
			}

		// character event
		case ev := <-events:
			switch ev := ev.(type) {
			case PositionUpdateEvent:
				// forward to client
				log.Println("position update", ev.Position)
			}
		}
	}
}

type Action interface {
	Apply(Area) error
}

type CastSpellAction struct {
	Unit   Identity
	Target Identity
	Spell  string
}

func (action CastSpellAction) Apply(area Area) error {
	caster, err := area.Actor(action.Unit)
	if err != nil {
		return err
	}
	target, err := area.Actor(action.Target)
	if err != nil {
		return err
	}
	log.Println(caster.Name(), "casting spell", action.Spell, "on", target.Name())
	return nil
}

type SetTargetAction struct {
	Unit   Identity
	Target Identity
}

func (action SetTargetAction) Apply(area Area) error {
	actor, err := area.Actor(action.Unit)
	if err != nil {
		return err
	}
	actor.SetTarget(action.Target)
	log.Println(actor.Name(), "set target to", action.Target)
	return nil
}
