//go:build !windows
// +build !windows

package config

import (
	"os"
	"path/filepath"
)

// GetConfigFolder returns the path to the config folder for the application.
func GetConfigFolder(vendor string, application string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".config", application)
}
