package cmd

import (
	"fmt"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/config"
	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Print the JSON schema for the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.JSONSchema())
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
