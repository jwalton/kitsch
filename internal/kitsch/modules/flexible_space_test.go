package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFlexibleSpaces(t *testing.T) {
	result := processFlexibleSpaces(10, "a"+flexibleSpaceMarker+"b", "")
	assert.Equal(t, "a        b", result)

	result = processFlexibleSpaces(10, "a"+flexibleSpaceMarker+"b\n$ ", "")
	assert.Equal(t, "a        b\n$ ", result)

	result = processFlexibleSpaces(10, "ab", "")
	assert.Equal(t, "ab", result)

	result = processFlexibleSpaces(10, "a"+flexibleSpaceMarker+"b"+flexibleSpaceMarker+"c", "")
	assert.Equal(t, "a   b    c", result)

	result = processFlexibleSpaces(10, "a"+flexibleSpaceMarker+"b"+flexibleSpaceMarker+"c", "foo")
	assert.Equal(t, "afoobfooc", result)
}
