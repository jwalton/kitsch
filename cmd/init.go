package cmd

import (
	"fmt"
	"os"

	"github.com/jwalton/kitsch-prompt/internal/initscripts"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:       "init",
	Short:     "Returns a script which can be used to initialize kitsch-prompt.",
	ValidArgs: []string{"bash", "zsh"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]

		cfgOptions := ""
		if cfgFile != "" {
			cfgOptions = fmt.Sprintf("--config \"%s\" ", cfgFile)
		}

		printFullInit, err := cmd.Flags().GetBool("print-full-init")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		if !printFullInit {
			executable, err := os.Executable()
			if err != nil {
				executable = cmd.Parent().CommandPath()
			}
			fmt.Printf("source <(\"%s\" init %s %s--print-full-init)\n", executable, shell, cfgOptions)
		} else {
			script, err := initscripts.InitScript(shell, cfgFile)
			if err != nil {
				cmd.PrintErrln(err.Error())
				os.Exit(1)
			}

			fmt.Println(script)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().Bool("print-full-init", false, "Print the main initialization script")
}
