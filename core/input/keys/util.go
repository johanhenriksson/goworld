package keys

func Pressed(ev Event, code Code) bool {
	if ev.Action() != Press {
		return false
	}
	if ev.Code() != code {
		return false
	}
	return true
}

func PressedMods(ev Event, code Code, mods Modifier) bool {
	if !Pressed(ev, code) {
		return false
	}
	return ev.Modifier() == mods
}

func Released(ev Event, code Code) bool {
	if ev.Action() != Release {
		return false
	}
	if ev.Code() != code {
		return false
	}
	return true
}

func ReleasedMods(ev Event, code Code, mods Modifier) bool {
	if !Released(ev, code) {
		return false
	}
	return ev.Modifier() == mods
}
