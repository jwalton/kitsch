package cmd

import (
	"fmt"
	"os"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/config"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [file]",
	Short: "Check configuration file for errors",
	Long: `Checks a configuration file for errors.  If no filename is given,
it will check the default configuration file.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.SetVerbose(true)

		configFile := ""
		if len(args) > 0 {
			configFile = args[0]
		} else if fileutils.FileExists(defaultConfigFile) {
			configFile = defaultConfigFile
		} else {
			log.Error("No configuration file found.")
			os.Exit(1)
		}

		fmt.Println("Checking config file: " + configFile)

		contents, err := os.ReadFile(configFile)
		if err != nil {
			log.Error("Could not read configuration file " + configFile + ": " + err.Error())
			os.Exit(1)
		}

		err = config.ValidateConfiguration(contents)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		fmt.Println(gchalk.BrightGreen("OK"))
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
