package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/styling"
)

// createTextContext creates a Context with reasonable defaults that can
// be passed in to modules when unit testing.
func testContext(username string) *Context {
	return &Context{
		Environment: &env.DummyEnv{
			Env: map[string]string{
				"USER": username,
				"HOME": "/Users/" + username,
			},
		},
		Globals: Globals{
			CWD:                     "/Users/" + username,
			Home:                    "/Users/" + username,
			Username:                username,
			UserFullName:            "Jason Walton",
			Hostname:                "lucid",
			Status:                  0,
			PreviousCommandDuration: 0,
			Shell:                   "bash",
		},
		Styles: styling.Registry{},
	}
}
