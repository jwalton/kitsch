// Package sampleconfig contains sample configuration files.
package sampleconfig

import (
	// embed required for sample configs below.
	_ "embed"
)

// DefaultConfig is the default configuration, as YAML data.
//go:embed default.yaml
var DefaultConfig []byte
