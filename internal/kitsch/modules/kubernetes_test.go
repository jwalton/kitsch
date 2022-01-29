package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestKubernetes(t *testing.T) {
	mod := moduleFromYAML(heredoc.Doc(`
		type: kubernetes
	`)).(KubernetesModule)

	mod.configFileContents = []byte(heredoc.Doc(`
		apiVersion: v1
		kind: Config
		contexts:
		  - name: prod
		    context:
		      cluster: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		      user: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		current-context: prod
	`))

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	expectedData := kubernetesModuleData{
		OriginalContext: "prod",
		Context:         "prod",
		Namespace:       "",
	}

	assert.Equal(t, expectedData, result.Data)
	assert.Equal(t, "☸ prod", result.Text)
}

func TestKubernetesWithAlias(t *testing.T) {
	mod := moduleFromYAML(heredoc.Doc(`
		type: kubernetes
		contextAliases:
		  prod: production
	`)).(KubernetesModule)

	mod.configFileContents = []byte(heredoc.Doc(`
		apiVersion: v1
		kind: Config
		contexts:
		  - name: prod
		    context:
		      cluster: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		      user: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		current-context: prod
	`))

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	expectedData := kubernetesModuleData{
		OriginalContext: "prod",
		Context:         "production",
		Namespace:       "",
	}

	assert.Equal(t, expectedData, result.Data)
	assert.Equal(t, "☸ production", result.Text)
}

func TestKubernetesWithNamespace(t *testing.T) {
	mod := moduleFromYAML(heredoc.Doc(`
		type: kubernetes
	`)).(KubernetesModule)

	mod.configFileContents = []byte(heredoc.Doc(`
		apiVersion: v1
		kind: Config
		contexts:
		  - name: prod
		    context:
		      cluster: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		      user: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		      namespace: kube-system
		current-context: prod
	`))

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	expectedData := kubernetesModuleData{
		OriginalContext: "prod",
		Context:         "prod",
		Namespace:       "kube-system",
	}

	assert.Equal(t, expectedData, result.Data)
	assert.Equal(t, "☸ prod", result.Text)
}

func TestKubernetesWithDefaultNamespace(t *testing.T) {
	mod := moduleFromYAML(heredoc.Doc(`
		type: kubernetes
	`)).(KubernetesModule)

	mod.configFileContents = []byte(heredoc.Doc(`
		apiVersion: v1
		kind: Config
		contexts:
		  - name: prod
		    context:
		      cluster: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		      user: arn:aws:eks:us-east-1:00000000:cluster/my-prod-cluster
		      namespace: default
		current-context: prod
	`))

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	expectedData := kubernetesModuleData{
		OriginalContext: "prod",
		Context:         "prod",
		Namespace:       "",
	}

	assert.Equal(t, expectedData, result.Data)
	assert.Equal(t, "☸ prod", result.Text)
}

func TestKubernetesWithMissingContext(t *testing.T) {
	mod := moduleFromYAML(heredoc.Doc(`
		type: kubernetes
	`)).(KubernetesModule)

	mod.configFileContents = []byte(heredoc.Doc(`
		apiVersion: v1
		kind: Config
		current-context: prod
	`))

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	expectedData := kubernetesModuleData{
		OriginalContext: "prod",
		Context:         "prod",
		Namespace:       "",
	}

	assert.Equal(t, expectedData, result.Data)
	assert.Equal(t, "☸ prod", result.Text)
}

func TestKubernetesWithCorruptConfigFile(t *testing.T) {
	mod := moduleFromYAML(heredoc.Doc(`
		type: kubernetes
	`)).(KubernetesModule)

	mod.configFileContents = []byte(`unexpected`)

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	expectedData := kubernetesModuleData{
		OriginalContext: "",
		Context:         "",
		Namespace:       "",
	}

	assert.Equal(t, expectedData, result.Data)
	assert.Equal(t, "", result.Text)
}
