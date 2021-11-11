package projects

import (
	"github.com/jwalton/kitsch-prompt/internal/kitsch/condition"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
)

// DefaultProjectTypes is a default list of project types, in priority order.
var DefaultProjectTypes = []ProjectType{
	{
		Name: "java",
		Condition: condition.Condition{
			IfExtensions: []string{"java"},
		},
		ToolSymbol: "java",
		ToolVersion: getters.CustomGetter{
			Type: "custom",
			From: "java -Xinternalversion",
			// Based on https://stackoverflow.com/questions/66601929/how-can-i-determine-whether-the-installed-java-supports-modules-or-not
			Regex: `\(([\d\.]+)[^\d\.]?[^\s]*\)(:?, built|from)`,
		},
	},
	{
		Name: "go",
		Condition: condition.Condition{
			IfFiles:      []string{"go.mod"},
			IfExtensions: []string{"go"},
		},
		ToolSymbol: "go",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "go version",
			Regex: `go version go(\d+\.\d+\.\d+)`,
		},
	},
	{
		Name: "rust",
		Condition: condition.Condition{
			IfFiles:      []string{"Cargo.toml"},
			IfExtensions: []string{"rs"},
		},
		ToolSymbol: "rustc",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "rustc version",
			Regex: `rustc (\d+\.\d+\.\d+)`,
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
		Condition: condition.Condition{
			IfFiles: []string{"yarn.lock"},
		},
		ToolSymbol: "node",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "node --version",
			Regex: `v(.*)`,
		},
		PackageManagerSymbol: "yarn",
		PackageManagerVersion: getters.CustomGetter{
			Type: "custom",
			// TODO: This is pretty slow - should use the same trick we use for NPM.
			// With yarn v1.x installed via Brew, we want to follow that yarn symlink,
			// then look in ../libexec/package.json.  On the alpine docker container,
			// the file is ../package.json.  For some reason in my `v14.18.1` node installed
			// via nvm, `yarn` is in v14.18.1/lib/node_modules/corepack/dist???
			From: "yarn --version",
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
		Condition: condition.Condition{
			IfFiles: []string{"package.json"},
		},
		ToolSymbol: "node",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "node --version",
			Regex: `v(.*)`,
		},
		PackageManagerSymbol:  "npm",
		PackageManagerVersion: npmVersionGetter{},
		PackageVersion: getters.CustomGetter{
			Type:          "file",
			From:          "package.json",
			As:            "json",
			ValueTemplate: "{{ .version }}",
		},
	},
	{
		Name: "deno",
		Condition: condition.Condition{
			IfFiles: []string{"mod.ts"},
		},
		ToolSymbol: "deno",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "deno --version",
			Regex: `deno (\d+\.\d+\.\d+)`,
		},
	},
	{
		Name: "helm",
		Condition: condition.Condition{
			IfFiles: []string{"Chart.yaml"},
		},
		ToolSymbol: "helm",
		ToolVersion: getters.CustomGetter{
			Type:  "custom",
			From:  "helm version",
			Regex: `^version.BuildInfo{Version:"v(\d+\.\d+\.\d+)"`,
		},
	},
}
