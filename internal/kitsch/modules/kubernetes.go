package modules

import (
	"os"
	"path/filepath"

	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas KubernetesModule

// KubernetesModule lets us know what Kubernetes context we are currently in.
type KubernetesModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=kubernetes"`
	// Symbol is a symbol to show if a Kubernetes context is detected.  Defaults to "☸ "
	Symbol string `yaml:"symbol"`
	// ContextAliases is a map where keys are context names and values are the
	// value we want to show.  If the value is an empty string, we will not
	// show anything.
	ContextAliases map[string]string `yaml:"contextAliases"`
	// ConfigFile is the path to the kubectl config file.  Defaults to "~/.kube/config".
	ConfigFile string `yaml:"configFile"`
	// configFileContents is the contents of the kubectl config file. If this value
	// is not empty, we'll use this as the contents of the kubectl config file instead
	// of reading them from ConfigFile.  This is used for unit testing.
	configFileContents []byte
}

// kubernetesModuleData is the template variables returned by the KubernetesModule.
type kubernetesModuleData struct {
	// OriginalContext is the raw "current-context" from the config file.
	OriginalContext string
	// Context is the context to display.  If "OriginalContext" maps to a ContextAlias,
	// this will be the alias, otherwise this will be the same as "OriginalContext".
	Context string
	// Namespace is the current namespace.  If not namespace is set or is
	// "default", this will be an empty string.
	Namespace string
}

type kubectlConfig struct {
	// CurrentContext is the name of the current context.
	CurrentContext string `yaml:"current-context"`
	// Contexts is the list of all contexts.
	Contexts []struct {
		// Name is the name of the context.
		Name string `yaml:"name"`
		// Context is the context.
		Context struct {
			// Cluster is the cluster.
			Cluster string `yaml:"cluster"`
			// User is the user for this cluster.
			User string `yaml:"user"`
			// Namespace is the namespace for this context, if present.
			Namespace string `yaml:"namespace"`
		}
	} `yaml:"contexts"`
}

func (mod KubernetesModule) loadConfigFile(homedir string) *kubectlConfig {
	// Use mod.configFileContents if it's set, otherwise read the config file.
	configFileContents := mod.configFileContents
	if configFileContents == nil {
		configFile := mod.ConfigFile
		if configFile == "" {
			configFile = filepath.Join(homedir, ".kube", "config")
		}

		var err error
		configFileContents, err = os.ReadFile(configFile)
		if err != nil {
			// Config file doesn't exist, or can't be read.
			return nil
		}
	}

	// Parse the config file.
	config := kubectlConfig{}
	err := yaml.Unmarshal(configFileContents, &config)
	if err != nil {
		log.Warn("Could not parse kubectl config file:", err)
		return nil
	}

	return &config
}

// Execute the module.
func (mod KubernetesModule) Execute(context *Context) ModuleResult {
	text := ""
	data := kubernetesModuleData{}

	config := mod.loadConfigFile(context.Globals.Home)
	if config != nil && config.CurrentContext != "" {
		data.OriginalContext = config.CurrentContext

		if alias, ok := mod.ContextAliases[config.CurrentContext]; ok {
			data.Context = alias
		} else {
			data.Context = config.CurrentContext
		}

		// Find the context.
		for _, context := range config.Contexts {
			if context.Name == config.CurrentContext {
				if context.Context.Namespace != "default" {
					data.Namespace = context.Context.Namespace
				}
				break
			}
		}

		if data.Context != "" {
			text = mod.Symbol + data.Context
		}
	}

	return ModuleResult{DefaultText: text, Data: data}
}

func init() {
	registerModule(
		"kubernetes",
		registeredModule{
			jsonSchema: schemas.KubernetesModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := KubernetesModule{
					Type:       "kubernetes",
					Symbol:     "☸ ",
					ConfigFile: "",
				}
				err := node.Decode(&module)
				return module, err
			},
		},
	)
}
