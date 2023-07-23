package test

import (
	"fmt"

	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/onsi/gomega/format"
)

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

func BeApproxVec3(v vec3.T) *approxVec3 {
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
