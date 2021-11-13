package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdDuration(t *testing.T) {
	context := testContext("jwalton")

	forTime := func(mod Module, time int) string {
		context.Globals.PreviousCommandDuration = int64(time)
		result := mod.Execute(context)
		return result.Text
	}

	// If we're under "MinTime", should produce no output.
	mod, err := moduleFromYAML("{ type: command_duration, minTime: 2000 }")
	assert.NoError(t, err)
	assert.Equal(t, "", forTime(mod, 1000))
	assert.Equal(t, "4s", forTime(mod, 4000))
	assert.Equal(t, "1m0s", forTime(mod, 60000))
	assert.Equal(t, "1m9s", forTime(mod, 69001))

	mod, err = moduleFromYAML("{ type: command_duration, minTime: 2000, showMilliseconds: true }")
	assert.NoError(t, err)
	assert.Equal(t, "1m9s1ms", forTime(mod, 69001))
	assert.Equal(t, "2h46m40s0ms", forTime(mod, 10000000))
}
