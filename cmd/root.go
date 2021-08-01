// Package cmd contains code for the `pixdl` CLI tool.
package cmd

import (
	"log"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kitsch",
	Short: "Ridiculously customizable shell prompt",
	Long: heredoc.Doc(`
		kitsch is a program for displaying a shell prompt.  If you're seeing this,
		it's because you're trying to run "kitsch" without a command.  If you
		want to install kitsch:

		  TODO: installation example goes here.

		Examples:

		  # Check your configuration for errors
		  kitsch check

		  # Display what your prompt would look like using a certain config
		  kitsch show --config ./config.yaml --dry-run
	`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	// TODO: Is Cobra a bit heavyweight here?
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pixdl.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "Use verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO: Read in config file.
}
