package batching

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType struct {
	Prop1 int
	Prop2 string
}

func TestSplitToBatches(t *testing.T) {
	data := []testType{}
	for i := 0; i < 100; i++ {
		data = append(data, testType{Prop1: i, Prop2: fmt.Sprintf("test %d", i)})
	}

	resultsBatchesOf25 := SplitToBatches(data, 25)
	resultsBatchesOf35 := SplitToBatches(data, 35)

	assert.Len(t, data, 100, "wrong number of initial data")

	assert.Len(t, resultsBatchesOf25, 4)
	for _, r := range resultsBatchesOf25 {
		assert.Len(t, r, 25)
	}

	assert.Len(t, resultsBatchesOf35, 3)
	assert.Len(t, resultsBatchesOf35[0], 35)
	assert.Len(t, resultsBatchesOf35[1], 35)
	assert.Len(t, resultsBatchesOf35[2], 30, "leftover data not correct")
}
