package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdDuration(t *testing.T) {
	mod := CmdDurationModule{MinTime: 2000}
	context := testContext("jwalton")

	forTime := func(time int) string {
		context.Globals.PreviousCommandDuration = int64(time)
		result := mod.Execute(context)
		return result.Text
	}

	// If we're under "MinTime", should produce no output.
	assert.Equal(t, "", forTime(1000))
	assert.Equal(t, "4s", forTime(4000))
	assert.Equal(t, "1m0s", forTime(60000))
	assert.Equal(t, "1m9s", forTime(69001))

	mod.ShowMilliseconds = true
	assert.Equal(t, "1m9s1ms", forTime(69001))
	assert.Equal(t, "2h46m40s0ms", forTime(10000000))
}
