package modules

import (
	"fmt"

	"github.com/jwalton/kitsch/internal/kitsch/condition"
	"github.com/jwalton/kitsch/internal/kitsch/log"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas CommonConfig

// CommonConfig is common configuration for all modules.
type CommonConfig struct {
	// Type is the type of this module.
	Type string `yaml:"type"`
	// ID is a unique identifier for this module.  IDs are unique only within the
	// parent block.
	ID string `yaml:"id"`
	// Style is the style to apply to this module.
	Style string `yaml:"style"`
	// Template is a golang template to use to render the output of this module.
	Template string `yaml:"template"`
	// Conditions are conditions that must be met for this module to execute.
	Conditions condition.Conditions `yaml:"conditions"`
}

// Validate checks for common configuration errors in the CommonConfig, and prints
// errors to the log if any are found.
func (config *CommonConfig) Validate(context *Context, prefix string) {
	_, err := context.Styles.Get(config.Style)
	if err != nil {
		log.Warn(fmt.Sprintf("%s: Error parsing style: %v", prefix, err))
	}
}

func getCommonConfig(node *yaml.Node) (CommonConfig, error) {
	var config CommonConfig
	if node == nil {
		return config, fmt.Errorf("cannot get type of empty node")
	}

	err := node.Decode(&config)
	if err != nil {
		return config, err
	}

	if config.Type == "" {
		return config, fmt.Errorf("object is missing type (%d:%d)", node.Line, node.Column)
	}

	if config.ID == "" {
		config.ID = config.Type
	}

	return config, nil
}
