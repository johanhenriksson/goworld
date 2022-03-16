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
