package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var slice = []string{"a", "b", "c"}

func TestSliceContains(t *testing.T) {
	assert.True(t, SliceContains(slice, "a"))
	assert.False(t, SliceContains(slice, "d"))
}
