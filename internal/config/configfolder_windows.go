// +build windows

package config

import "os"

// GetConfigFolder returns the path to the config folder for the application.
func GetConfigFolder(vendor string, application string) string {
	return filepath.Join(os.Getenv("APPDATA"), vendor, application)
}
