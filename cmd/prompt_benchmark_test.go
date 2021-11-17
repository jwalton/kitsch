package cmd

import (
	"testing"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/kitsch-prompt/internal/cache"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/config"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modules"
)

func BenchmarkPrompt(b *testing.B) {
	b.ReportAllocs()

	configuration, err := config.LoadDefaultConfig()
	if err != nil {
		b.Fatal(err)
		return
	}

	globals := modules.Globals{
		CWD:                     "/Users/jwalton",
		Home:                    "/Users/jwalton",
		Username:                "jwalton",
		UserFullName:            "Jason Walton",
		Hostname:                "lucid",
		Status:                  0,
		PreviousCommandDuration: 0,
		Shell:                   "bash",
	}

	dummyEnv := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
	}

	gchalk.SetLevel(gchalk.LevelAnsi16m)
	gchalk.Stderr.SetLevel(gchalk.LevelAnsi16m)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valueCache := cache.NewMemoryCache()
		renderPrompt(configuration, globals, dummyEnv, valueCache)
	}
}
