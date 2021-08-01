// Package config represents a configuration file.
package config

import (
	"fmt"

	"github.com/jwalton/kitsch-prompt/internal/modules"
	"gopkg.in/yaml.v3"
)

// Config represents a configuration file.
type Config struct {
	// Prompt is the module to use to display the prompt.
	Prompt yaml.Node
}

// GetPromptModule returns the root prompt module for the configuration.
func (c *Config) GetPromptModule() (modules.Module, error) {
	if c.Prompt.IsZero() {
		return nil, fmt.Errorf("configuration is missing prompt")
	}
	return modules.CreateModule(&c.Prompt)
}

// ReadConfig reads configuration from a YAML string.
func ReadConfig(data string) (Config, error) {
	var config Config
	err := yaml.Unmarshal([]byte(data), &config)
	return config, err
}
