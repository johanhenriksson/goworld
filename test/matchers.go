package test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
	"strings"

	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/upload"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/format"
)

//
// quaternion matcher
//

func ApproxQuat(q quat.T) *approxQuat {
	return &approxQuat{q}
}

type approxQuat struct {
	Expected interface{}
}

func (matcher *approxQuat) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	expect, expectOk := matcher.Expected.(quat.T)
	if !expectOk {
		return false, fmt.Errorf("expected a quat.T value")
	}
	actualV, actualOk := actual.(quat.T)
	if !actualOk {
		return false, fmt.Errorf("expected a quat.T value")
	}
	return expect.ApproxEqual(actualV), nil
}

func (matcher *approxQuat) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}
	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *approxQuat) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}

//
// vec3 matcher
//

func ApproxVec3(v vec3.T) *approxVec3 {
	return &approxVec3{v}
}

type approxVec3 struct {
	Expected interface{}
}

func (matcher *approxVec3) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	expect, expectOk := matcher.Expected.(vec3.T)
	if !expectOk {
		return false, fmt.Errorf("expected a vec3.T value")
	}
	actualV, actualOk := actual.(vec3.T)
	if !actualOk {
		return false, fmt.Errorf("expected a vec3.T value")
	}
	return expect.ApproxEqual(actualV), nil
}

func (matcher *approxVec3) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}
	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *approxVec3) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}

//
// image matcher
//

func ApproxImage(expected any) *approxImage {
	return &approxImage{
		Expected:  expected,
		Threshold: 10,
	}
}

type approxImage struct {
	Expected  interface{}
	Threshold int
}

func (matcher *approxImage) Match(actualValue interface{}) (success bool, err error) {
	if actualValue == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	path, pathOk := matcher.Expected.(string)
	if !pathOk {
		return false, fmt.Errorf("expected an image path")
	}
	infile, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer infile.Close()

	expected, _, err := image.Decode(infile)
	if err != nil {
		return false, err
	}

	actual, actualOk := actualValue.(image.Image)
	if !actualOk {
		return false, fmt.Errorf("expected an *image.RGBA value")
	}

	if expected.Bounds() != actual.Bounds() {
		return false, fmt.Errorf("image bounds dont match, expected: %s but was %s",
			expected.Bounds(), actual.Bounds())
	}

	diff := func(a, b uint32) int {
		a, b = a&0xFF, b&0xFF
		if a > b {
			return int(a - b)
		} else {
			return int(b - a)
		}
	}

	for y := 0; y < expected.Bounds().Max.Y; y++ {
		for x := 0; x < expected.Bounds().Max.X; x++ {
			expectedColor := expected.At(x, y)
			actualColor := actual.At(x, y)

			er, eg, eb, ea := expectedColor.RGBA()
			ar, ag, ab, aa := actualColor.RGBA()
			delta := diff(er, ar) + diff(eg, ag) + diff(eb, ab) + diff(ea, aa)

			if delta > matcher.Threshold {
				writeFailureImage(expected, actual)
				return false, fmt.Errorf("colors at (%d,%d) dont match (delta: %d), expected: %v but was %v",
					x, y, delta, expectedColor, actualColor)
			}
		}
	}

	return true, nil
}

func writeFailureImage(expected, actual image.Image) {
	width, height := expected.Bounds().Dx(), expected.Bounds().Dy()
	combined := image.NewRGBA(image.Rect(0, 0, 3*width, height))

	draw.Draw(combined, expected.Bounds(), expected, image.ZP, draw.Src)
	draw.Draw(combined, image.Rect(width, 0, 2*width, height), actual, image.ZP, draw.Src)

	for y := 0; y < expected.Bounds().Max.Y; y++ {
		for x := 0; x < expected.Bounds().Max.X; x++ {
			expectedColor := expected.At(x, y)
			actualColor := actual.At(x, y)
			if expectedColor == actualColor {
				expected.(*image.RGBA).Set(x, y, color.Transparent)
			}
		}
	}
	draw.Draw(combined, image.Rect(2*width, 0, 3*width, height), expected, image.ZP, draw.Src)

	testName := strings.Join(ginkgo.CurrentSpecReport().ContainerHierarchyTexts, "_")
	testName += "_" + ginkgo.CurrentSpecReport().LeafNodeText

	if err := upload.SavePng(actual, testName+"_actual.png"); err != nil {
		panic(err)
	}
	if err := upload.SavePng(combined, testName+"_failure.png"); err != nil {
		panic(err)
	}
}

func (matcher *approxImage) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}
	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *approxImage) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}
