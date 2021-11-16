package modules

import (
	"testing/fstest"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
	"gopkg.in/yaml.v3"
)

// createTextContext creates a Context with reasonable defaults that can
// be passed in to modules when unit testing.
func testContext(username string) *Context {
	fsys := fstest.MapFS{}

	return &Context{
		Environment: &env.DummyEnv{
			Env: map[string]string{
				"USER": username,
				"HOME": "/Users/" + username,
			},
		},
		Directory: fileutils.NewDirectoryTestFS("/Users/"+username, fsys),
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

func moduleFromYAML(data string) (Module, error) {
	var moduleSpec ModuleSpec
	err := yaml.Unmarshal([]byte(data), &moduleSpec)
	if err != nil {
		return nil, err
	}
	return moduleSpec.Module, nil
}

func moduleFromYAMLMust(data string) Module {
	module, err := moduleFromYAML(data)
	if err != nil {
		panic(err)
	}
	return module
}
