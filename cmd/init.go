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
		script, err := initscripts.InitScript(args[0])

		if err != nil {
			cmd.PrintErrln(err.Error())
			os.Exit(1)
		}

		fmt.Println(script)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
