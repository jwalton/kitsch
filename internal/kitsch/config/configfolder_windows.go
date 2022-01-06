//go:build windows
// +build windows

package config

import (
	"os"
	"path/filepath"
)

// GetConfigFolder returns the path to the config folder for the application.
func GetConfigFolder(vendor string, application string) string {
	rootDir := os.Getenv("APPDATA")

	if rootDir == "" {
		var err error
		rootDir, err = os.UserHomeDir()
		if err != nil {
			rootDir = os.Getenv("HOME")
		}
	}

	return filepath.Join(rootDir, vendor, application)
}
