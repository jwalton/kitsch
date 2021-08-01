package cmd

import (
	"fmt"

	"github.com/jwalton/kitsch-prompt/internal/config"
	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/modules"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Show the prompt",
	Run: func(cmd *cobra.Command, args []string) {
		jobs, _ := cmd.Flags().GetInt("jobs")
		cmdDuration, _ := cmd.Flags().GetInt("cmd-duration")
		status, _ := cmd.Flags().GetInt("status")
		keymap, _ := cmd.Flags().GetString("keymap")

		runtimeEnv := env.NewEnv(jobs, cmdDuration, status, keymap)

		yamlConfig := `
prompt:
  type: block
  join: " "
  modules:
    - type: time
      style: brightBlack
    - type: directory
      style: brightBlue
      template: "[{{.directory}}]"
    - type: git
      style: brightYellow
    - type: prompt
      style: brightBlue
`

		var module modules.Module
		configuration, err := config.ReadConfig(yamlConfig)
		if err == nil {
			module, err = configuration.GetPromptModule()
		}

		if err != nil {
			fmt.Println(err)
			fmt.Print("$ ")
		} else {
			result := module.Execute(runtimeEnv)
			fmt.Println(result.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().IntP("cmd-duration", "d", 0, "The execution duration of the last command, in milliseconds")
	promptCmd.Flags().StringP("keymap", "k", "", "The keymap of fish/zsh")
	promptCmd.Flags().IntP("jobs", "j", 0, "The number of currently running jobs")
	promptCmd.Flags().IntP("status", "s", 0, "The status code of the previously run command")
}
