package cmd

import (
	"os"
	"runtime"
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/initscripts"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Prints instructions about how to setup " + programName + ".",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var shell string

		if len(args) > 0 {
			shell = args[0]
		} else if runtime.GOOS == "windows" {
			shell = "powershell"
		} else {
			shellType := os.Getenv("SHELL")
			if strings.HasSuffix(shellType, "/zsh") {
				shell = "zsh"
			} else if strings.HasSuffix(shellType, "/bash") {
				shell = "bash"
			} else {
				shell = "unknown"
			}
		}

		shellDetected := len(args) == 0
		initscripts.ShowSetupInstructions(programName, githubRepo, website, shell, shellDetected)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
