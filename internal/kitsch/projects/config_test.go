package projects

import (
	"testing"

	"github.com/jwalton/kitsch/internal/kitsch/condition"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadProjects(t *testing.T) {
	doc := `
  - name: "test"
    conditions:
      ifFiles: ["package.json"]
    toolSymbol: Node
    toolVersion:
      type: custom
      from: "node --version"
`

	projects := []ProjectType{}
	err := yaml.Unmarshal([]byte(doc), &projects)

	assert.Nil(t, err)
	assert.Equal(t,
		[]ProjectType{
			{
				Name: "test",
				Conditions: condition.Conditions{
					IfFiles: []string{"package.json"},
				},
				ToolSymbol: "Node",
				ToolVersion: getters.CustomGetter{
					Type: getters.TypeCustom,
					From: "node --version",
				},
			},
		},
		projects,
	)
}

var from = []ProjectType{
	{
		Name: "java",
		Conditions: condition.Conditions{
			IfExtensions: []string{".java"},
		},
		ToolSymbol: "Java",
		ToolVersion: getters.CustomGetter{
			Type: getters.TypeCustom,
			From: "java --version",
		},
	},
	{
		Name: "node",
		Conditions: condition.Conditions{
			IfFiles: []string{"package.json"},
		},
		ToolSymbol: "Node",
		ToolVersion: getters.CustomGetter{
			Type: getters.TypeCustom,
			From: "node --version",
		},
	},
}

func TestMergeEmptyProjectTypes(t *testing.T) {
	to := []ProjectType{}
	to = MergeProjectTypes(to, from, true)
	assert.Equal(t, from, to)
}

func projectTypesFromYAML(doc string) []ProjectType {
	projectTypes := []ProjectType{}
	err := yaml.Unmarshal([]byte(doc), &projectTypes)
	if err != nil {
		panic(err)
	}
	return projectTypes
}

func TestMergeReorderedProjectTypes(t *testing.T) {
	to := projectTypesFromYAML(`
  - name: node
  - name: java
`)

	to = MergeProjectTypes(to, from, true)
	assert.Equal(t,
		[]ProjectType{from[1], from[0]},
		to,
	)
}

func TestMergeAlteredProjectTypes(t *testing.T) {
	to := []ProjectType{
		{Name: "node", ToolSymbol: "JS"},
	}
	to = MergeProjectTypes(to, from, false)
	assert.Equal(t,
		[]ProjectType{
			{
				Name: "node",
				Conditions: condition.Conditions{
					IfFiles: []string{"package.json"},
				},
				ToolSymbol: "JS",
				ToolVersion: getters.CustomGetter{
					Type: getters.TypeCustom,
					From: "node --version",
				},
			},
		},
		to,
	)
}
