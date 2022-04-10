package projects

import (
	"github.com/jwalton/kitsch/internal/kitsch/condition"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
)

// DefaultProjectTypes is a default list of project types, in priority order.
var DefaultProjectTypes = []ProjectType{
	{
		Name:  "java",
		Style: "brightRed",
		Conditions: &condition.Conditions{
			IfExtensions: []string{"java"},
		},
		ToolSymbol: "java",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type: getters.TypeCustom,
				From: "java -Xinternalversion",
				// Based on https://stackoverflow.com/questions/66601929/how-can-i-determine-whether-the-installed-java-supports-modules-or-not
				Regex: `\(([\d\.]+)[^\d\.]?[^\s]*\)(:?, built|from)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
	{
		Name:  "go",
		Style: "brightCyan",
		Conditions: &condition.Conditions{
			IfFiles:      []string{"go.mod"},
			IfExtensions: []string{"go"},
		},
		ToolSymbol: "go",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "go version",
				Regex: `go version go([\d+.]+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
	{
		Name:  "rust",
		Style: "brightRed",
		Conditions: &condition.Conditions{
			IfFiles:      []string{"Cargo.toml"},
			IfExtensions: []string{"rs"},
		},
		ToolSymbol: "rustc",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "rustc version",
				Regex: `rustc (\d+\.\d+\.\d+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
		PackageVersion: []getters.Getter{
			getters.CustomGetter{
				Type:          getters.TypeFile,
				From:          "Cargo.toml",
				As:            getters.AsTOML,
				ValueTemplate: "{{ .package.version }}",
			},
		},
	},
	{
		Name:  "node-yarn",
		Style: "brightBlue",
		Conditions: &condition.Conditions{
			IfFiles: []string{"yarn.lock"},
		},
		ToolSymbol: "node",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "volta which node",
				Regex: `image/node/(\d+\.\d+\.\d+)/bin`,
				Cache: getters.CacheSettings{
					Enabled: true,
					Files:   []string{"./package.json", "${VOLTA_HOME}/tools/user/platform.json"},
				},
			},
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "node --version",
				Cache: getters.CacheSettings{Enabled: true},
				Regex: `v(.*)`,
			},
		},
		PackageManagerSymbol: "yarn",
		PackageManagerVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "volta which yarn",
				Regex: `image/yarn/(\d+\.\d+\.\d+)/bin`,
				Cache: getters.CacheSettings{
					Enabled: true,
					Files:   []string{"./package.json", "${VOLTA_HOME}/tools/user/platform.json"},
				},
			},
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "yarn --version",
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
		PackageVersion: []getters.Getter{
			getters.CustomGetter{
				Type:          getters.TypeFile,
				From:          "package.json",
				As:            getters.AsTOML,
				ValueTemplate: "{{ .version }}",
			},
		},
	},
	{
		Name:  "node",
		Style: "brightYellow",
		Conditions: &condition.Conditions{
			IfFiles: []string{"package.json"},
		},
		ToolSymbol: "node",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "volta which node",
				Regex: `image/node/(\d+\.\d+\.\d+)/bin`,
				Cache: getters.CacheSettings{
					Enabled: true,
					Files:   []string{"./package.json", "${VOLTA_HOME}/tools/user/platform.json"},
				},
			},
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "node --version",
				Cache: getters.CacheSettings{Enabled: true},
				Regex: `v(.*)`,
			},
		},
		PackageManagerSymbol: "npm",
		PackageManagerVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "volta which npm",
				Regex: `image/npm/(\d+\.\d+\.\d+)/bin`,
				Cache: getters.CacheSettings{
					Enabled: true,
					Files:   []string{"./package.json", "${VOLTA_HOME}/tools/user/platform.json"},
				},
			},
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "npm --version",
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
		PackageVersion: []getters.Getter{
			getters.CustomGetter{
				Type:          getters.TypeFile,
				From:          "package.json",
				As:            getters.AsTOML,
				ValueTemplate: "{{ .version }}",
			},
		},
	},
	{
		Name:  "deno",
		Style: "brightGreen",
		Conditions: &condition.Conditions{
			IfFiles: []string{"mod.ts"},
		},
		ToolSymbol: "deno",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "deno --version",
				Regex: `deno (\d+\.\d+\.\d+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
	{
		Name:  "python",
		Style: "brightYellow",
		Conditions: &condition.Conditions{
			IfFiles:      []string{"requirements.txt", "Pipfile", "pyproject.toml"},
			IfExtensions: []string{"py"},
		},
		ToolSymbol: "python",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "python3 --version",
				Regex: `^Python (\d+\.\d+\.\d+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "python --version",
				Regex: `^Python (\d+\.\d+\.\d+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "python2 --version",
				Regex: `^Python (\d+\.\d+\.\d+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
	{
		Name:  "php",
		Style: "#8993bb",
		Conditions: &condition.Conditions{
			IfFiles:      []string{"composer.json", ".php-version"},
			IfExtensions: []string{"php"},
		},
		ToolSymbol: "php",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "php --version",
				Regex: `^PHP (\d+\.\d+\.\d+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
	{
		Name:  "ruby",
		Style: "brightRed",
		Conditions: &condition.Conditions{
			IfFiles:      []string{"Gemfile", ".ruby-version"},
			IfExtensions: []string{"rb"},
		},
		ToolSymbol: "ruby",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "ruby --version",
				Regex: `^ruby (\d+\.\d+\.[0-9a-zA-Z]+)`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
		PackageManagerSymbol: "gem",
		PackageManagerVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "gem --version",
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
	{
		Name:  "helm",
		Style: "brightWhite",
		Conditions: &condition.Conditions{
			IfFiles: []string{"Chart.yaml"},
		},
		ToolSymbol: "helm",
		ToolVersion: []getters.Getter{
			getters.CustomGetter{
				Type:  getters.TypeCustom,
				From:  "helm version",
				Regex: `^version.BuildInfo{Version:"v(\d+\.\d+\.\d+)"`,
				Cache: getters.CacheSettings{Enabled: true},
			},
		},
	},
}
