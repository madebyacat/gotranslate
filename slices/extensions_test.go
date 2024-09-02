package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsWithSimpleTypes(t *testing.T) {
	expectedString, expectedInt, notExpectedInt, expectedFloat := "test 2", 555, -555, -5.8
	dataStrings := []string{"test 1", expectedString, "test 3"}
	dataInts := []int{1, 2, 3, -222, 4, expectedInt, 5, 6}
	dataFloats := []float64{2.55, 0.2, expectedFloat, 22.2}

	stringsResult := Contains(dataStrings, func(x string) bool { return x == expectedString })
	intsResult := Contains(dataInts, func(x int) bool { return x == expectedInt })
	intsResultNotExist := !Contains(dataInts, func(x int) bool { return x == notExpectedInt })
	floatsResult := Contains(dataFloats, func(x float64) bool { return x == expectedFloat })

	assert.Truef(t, stringsResult, "strings should contain expected value: %v", expectedString)
	assert.Truef(t, intsResult, "ints should contain expected value: %v", expectedInt)
	assert.Truef(t, intsResultNotExist, "ints should NOT contain value: %v", notExpectedInt)
	assert.Truef(t, floatsResult, "floats should contain expected value: %v", expectedFloat)
}

func TestAllWithSimpleTypes(t *testing.T) {
	dataStrings := []string{"test 1", "test 2", "test 3"}
	dataInts := []int{1, 2, 3, 4, 555, 5, 6}
	dataFloats := []float64{2.55, 0.2, -0.2, 22.2}

	stringsResult := All(dataStrings, func(x string) bool { return x != "" })
	intsResult := All(dataInts, func(x int) bool { return x > 0 })
	intsResultNotExist := !All(dataInts, func(x int) bool { return x < 10 })
	floatsResult := All(dataFloats, func(x float64) bool { return x > -0.3 })

	assert.Truef(t, stringsResult, "expected true for All strings")
	assert.Truef(t, intsResult, "expected true for All ints")
	assert.Truef(t, intsResultNotExist, "expected false for All because at least 1 mismatch")
	assert.Truef(t, floatsResult, "expected true for All floats")
}

type testType struct {
	Prop1 int
	Prop2 string
}

func TestContainsWithComplexTypes(t *testing.T) {
	expectedProp1, expectedProp2 := 555, "test 4"
	data := []testType{
		{1, "test 1"},
		{55, "test 2"},
		{789, "test 3"},
		{555, "test 4"},
	}

	result := Contains(data, func(x testType) bool { return x.Prop1 == expectedProp1 && x.Prop2 == expectedProp2 })
	resultNotExist := Contains(data, func(x testType) bool { return x.Prop1 == -555 && x.Prop2 == expectedProp2 })

	assert.Truef(t, result, "expected Contains to be true for existing data")
	assert.Falsef(t, resultNotExist, "expected Contains to be false for non existing data")
}

func TestAllWithComplexTypes(t *testing.T) {
	data := []testType{
		{1, "test 1"},
		{55, "test 2"},
		{789, "test 3"},
		{555, "test 4"},
	}

	result := All(data, func(x testType) bool { return x.Prop1 > 0 && x.Prop2 != "teeeest" })
	resultNotExist := All(data, func(x testType) bool { return x.Prop1 > 1 })

	assert.Truef(t, result, "expected All to match predicate")
	assert.Falsef(t, resultNotExist, "expected All to mismatch with some data")
}
