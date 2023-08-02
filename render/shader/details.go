package shader

import (
	"encoding/json"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/types"
)

type InputDetails struct {
	Index int
	Type  string
}

type Details struct {
	Inputs   map[string]InputDetails
	Bindings map[string]int
	Textures []texture.Slot
}

func (d *Details) ParseInputs() (Inputs, error) {
	inputs := Inputs{}
	for name, input := range d.Inputs {
		kind, err := types.TypeFromString(input.Type)
		if err != nil {
			return nil, err
		}
		inputs[name] = Input{
			Index: input.Index,
			Type:  kind,
		}
	}
	return inputs, nil
}

func ReadDetails(path string) (*Details, error) {
	data, err := assets.ReadAll(path)
	if err != nil {
		return nil, err
	}

	details := &Details{}
	err = json.Unmarshal(data, details)
	if err != nil {
		return nil, err
	}

	return details, nil
}
