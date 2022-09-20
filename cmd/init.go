package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/initscripts"
	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Returns a script which can be used to initialize " + programName + ".",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Usage()
			validShells := strings.Join(initscripts.ValidShells(), ", ")
			log.Error("init command requires the name of a shells.  Use one of: " + validShells + ".")
			os.Exit(1)
		}

		shell := args[0]

		printFullInit, err := cmd.Flags().GetBool("print-full-init")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		if !printFullInit {
			shortScript, err := initscripts.ShortInitScript(shell, cfgFile)
			if err != nil {
				cmd.PrintErrln(err.Error())
				os.Exit(1)
			}

			fmt.Println(shortScript)
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
