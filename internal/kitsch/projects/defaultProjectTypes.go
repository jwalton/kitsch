package projects

import (
	"github.com/jwalton/kitsch-prompt/internal/kitsch/condition"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
)

// DefaultProjectTypes is a default list of project types, in priority order.
var DefaultProjectTypes = []ProjectType{
	{
		Name: "java",
		Conditions: condition.Conditions{
			IfExtensions: []string{"java"},
		},
		ToolSymbol: "java",
		ToolVersion: getters.CustomGetter{
			Type: "custom",
			From: "java -Xinternalversion",
			// Based on https://stackoverflow.com/questions/66601929/how-can-i-determine-whether-the-installed-java-supports-modules-or-not
			Regex: `\(([\d\.]+)[^\d\.]?[^\s]*\)(:?, built|from)`,
			Cache: getters.CacheSettings{Enabled: true},
		},
	},
	{
		Name: "go",
		Conditions: condition.Conditions{
			IfFiles:      []string{"go.mod"},
			IfExtensions: []string{"go"},
		},
		ToolSymbol: "go",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "go version",
			Regex: `go version go(\d+\.\d+\.\d+)`,
			Cache: getters.CacheSettings{Enabled: true},
		},
	},
	{
		Name: "rust",
		Conditions: condition.Conditions{
			IfFiles:      []string{"Cargo.toml"},
			IfExtensions: []string{"rs"},
		},
		ToolSymbol: "rustc",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "rustc version",
			Regex: `rustc (\d+\.\d+\.\d+)`,
			Cache: getters.CacheSettings{Enabled: true},
		},
		PackageVersion: getters.CustomGetter{
			Type:          "file",
			From:          "Cargo.toml",
			As:            "toml",
			ValueTemplate: "{{ .package.version }}",
		},
	},
	{
		Name: "node-yarn",
		Conditions: condition.Conditions{
			IfFiles: []string{"yarn.lock"},
		},
		ToolSymbol: "node",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "node --version",
			Regex: `v(.*)`,
			Cache: getters.CacheSettings{Enabled: true},
		},
		PackageManagerSymbol: "yarn",
		PackageManagerVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "yarn --version",
			Cache: getters.CacheSettings{Enabled: true},
		},
		PackageVersion: getters.CustomGetter{
			Type:          "file",
			From:          "package.json",
			As:            "json",
			ValueTemplate: "{{ .version }}",
		},
	},
	{
		Name: "node",
		Conditions: condition.Conditions{
			IfFiles: []string{"package.json"},
		},
		ToolSymbol: "node",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "node --version",
			Regex: `v(.*)`,
			Cache: getters.CacheSettings{Enabled: true},
		},
		PackageManagerSymbol: "npm",
		PackageManagerVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "npm --version",
			Cache: getters.CacheSettings{Enabled: true},
		},
		PackageVersion: getters.CustomGetter{
			Type:          "file",
			From:          "package.json",
			As:            "json",
			ValueTemplate: "{{ .version }}",
		},
	},
	{
		Name: "deno",
		Conditions: condition.Conditions{
			IfFiles: []string{"mod.ts"},
		},
		ToolSymbol: "deno",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "deno --version",
			Regex: `deno (\d+\.\d+\.\d+)`,
			Cache: getters.CacheSettings{Enabled: true},
		},
	},
	{
		Name: "helm",
		Conditions: condition.Conditions{
			IfFiles: []string{"Chart.yaml"},
		},
		ToolSymbol: "helm",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "helm version",
			Regex: `^version.BuildInfo{Version:"v(\d+\.\d+\.\d+)"`,
			Cache: getters.CacheSettings{Enabled: true},
		},
	},
}
