// Package cmd contains code for the `pixdl` CLI tool.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/kitsch/config"
	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/projects"
	"github.com/spf13/cobra"
)

var userConfigDir string
var cfgFile string
var defaultConfigFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   programName,
	Short: "Ridiculously customizable shell prompt",
	Long: heredoc.Doc(`
		` + programName + ` is a program for displaying a shell prompt.  If you're seeing this,
		it's because you're trying to run "` + programName + `" without a command.  If you
		want to install kitsch:

		  TODO: installation example goes here.

		Examples:

		  # Check your configuration for errors
		  ` + programName + ` check

		  # Display what your prompt would look like using a certain config
		  ` + programName + ` show --config ./config.yaml --dry-run
	`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	var err error
	userConfigDir, err = getConfigFolder()
	if err != nil {
		userConfigDir = "~"
	}
	defaultConfigFile = filepath.Join(userConfigDir, "kitsch.yaml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+defaultConfigFile+")")
	rootCmd.PersistentFlags().Bool("verbose", false, "Use verbose output")
}

func readConfig() (*config.Config, error) {
	var configuration *config.Config
	var err error

	if cfgFile != "" {
		configuration, err = config.LoadConfigFromFile(cfgFile, false)
		if err != nil {
			log.Error("Error loading config file "+cfgFile+": ", err)
		}
	}

	if configuration == nil && cfgFile != defaultConfigFile {
		configuration, err = config.LoadConfigFromFile(defaultConfigFile, false)
		if err != nil && !os.IsNotExist(err) {
			log.Error("Error loading config file "+defaultConfigFile+": ", err)
		}
	}

	if configuration == nil {
		configuration, err = config.LoadDefaultConfig()
		if err != nil {
			log.Error("Error loading default config: ", err)
		}
	}

	// Merge in default project types.
	configuration.ProjectsTypes = projects.MergeProjectTypes(
		configuration.ProjectsTypes,
		projects.DefaultProjectTypes,
		true,
	)

	return configuration, err
}

// getConfigFolder returns the folder that contains configuration
// information (e.g. "~/.config/kitsch" on Mac or Linux,
// "C:\Users\<User>\AppData\Roaming\kitsch\kitsch" on PC).
func getConfigFolder() (string, error) {
	configDir := config.GetConfigFolder(programName, programName)

	// Create the folder if it doesn't exist.
	err := os.MkdirAll(configDir, 0750)
	if err != nil {
		return "", fmt.Errorf("error creating config folder %v: %w", configDir, err)
	}

	return configDir, nil
}
