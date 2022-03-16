package util_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/util"
)

type AlignCase struct {
	Offset    int
	Alignment int
	Expected  int
}

func TestAlign(t *testing.T) {
	cases := []AlignCase{
		{23, 64, 64},
		{64, 64, 64},
		{72, 64, 128},
	}
	for i, testcase := range cases {
		actual := Align(testcase.Offset, testcase.Alignment)
		if actual != testcase.Expected {
			t.Errorf("case %d: expected %d, was %d", i, testcase.Expected, actual)
		}
	}
}

func TestValidateAlign(t *testing.T) {
	type FailingStruct struct {
		A bool
		B int
	}
	err := ValidateAlignment(FailingStruct{})
	if err == nil {
		t.Error("expected struct to fail alignment test")
	}

	type PassingStruct struct {
		A int
		B float32
	}
	err = ValidateAlignment(PassingStruct{})
	if err != nil {
		t.Error("expected struct to pass alignment test")
	}
}
