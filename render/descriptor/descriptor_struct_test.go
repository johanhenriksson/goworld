package descriptor_test

import (
	"errors"
	"testing"

	. "github.com/johanhenriksson/goworld/render/descriptor"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type TestSet struct {
	Set
	Diffuse *Sampler
}

func TestParseDescriptors(t *testing.T) {
	set := TestSet{
		Diffuse: &Sampler{
			Stages: core1_0.StageAll,
		},
	}
	desc, err := ParseDescriptorStruct(&set)
	if err != nil {
		t.Error(err)
	}
	if len(desc) != 1 {
		t.Error("expected to find diffuse descriptor")
	}
}

func TestParseDescriptorsNil(t *testing.T) {
	set := TestSet{
		Diffuse: nil,
	}
	_, err := ParseDescriptorStruct(&set)
	if !errors.Is(err, ErrDescriptorType) {
		t.Errorf("expected nil set error, was %s", err)
	}
}

func TestParseDescriptorsNonPointer(t *testing.T) {
	type FailSet struct {
		Set
		Diffuse Sampler
	}
	set := FailSet{}
	_, err := ParseDescriptorStruct(&set)
	if !errors.Is(err, ErrDescriptorType) {
		t.Errorf("expected non pointer descriptor error, was %s", err)
	}
}
