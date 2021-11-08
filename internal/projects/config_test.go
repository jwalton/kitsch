package projects

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/condition"
	"github.com/jwalton/kitsch-prompt/internal/getters"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadProjects(t *testing.T) {
	doc := `
  - name: "test"
    condition:
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
				Condition: condition.Condition{
					IfFiles: []string{"package.json"},
				},
				ToolSymbol: "Node",
				ToolVersion: getters.CustomGetter{
					Type: "custom",
					From: "node --version",
				},
			},
		},
		projects,
	)
}
