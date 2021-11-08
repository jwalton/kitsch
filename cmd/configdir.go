package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configDirCmd = &cobra.Command{
	Use:   "configdir",
	Short: "Print the location of the configuration directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(userConfigDir)
	},
}

func init() {
	rootCmd.AddCommand(configDirCmd)
}
