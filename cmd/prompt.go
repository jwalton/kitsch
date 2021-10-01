package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/go-supportscolor"
	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/modules"
	"github.com/jwalton/kitsch-prompt/internal/shellprompt"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Show the prompt",
	Run: func(cmd *cobra.Command, args []string) {
		jobs, _ := cmd.Flags().GetInt("jobs")
		status, _ := cmd.Flags().GetInt("status")
		keymap, _ := cmd.Flags().GetString("keymap")
		shell, _ := cmd.Flags().GetString("shell")

		cmdDurationStr, _ := cmd.Flags().GetString("cmd-duration")
		cmdDuration := 0
		if cmdDurationStr != "" {
			cmdDuration, _ = strconv.Atoi(cmdDurationStr)
		}

		// Because the prompt is shown from the shell, when it is run, it
		// will not be in a TTY.  Disable TTY detection in gchalk.
		stdoutFd := os.Stdout.Fd()
		level := supportscolor.SupportsColor(stdoutFd, supportscolor.IsTTYOption(true))
		gchalk.SetLevel(level.Level)
		gchalk.Stderr.SetLevel(level.Level)

		runtimeEnv := env.New(jobs, cmdDuration, status, keymap)

		configuration, err := readConfig()
		var module modules.Module
		if err == nil {
			module, err = configuration.GetPromptModule()
		}

		if err != nil {
			fmt.Println(err)
			fmt.Print("$ ")
		} else {
			prompt := module.Execute(runtimeEnv)
			withEscapes := shellprompt.AddZeroWidthCharacterEscapes(shell, prompt.Text)
			fmt.Println(withEscapes)
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().String("shell", "", "The type of shell")
	promptCmd.Flags().StringP("cmd-duration", "d", "", "The execution duration of the last command, in milliseconds")
	promptCmd.Flags().StringP("keymap", "k", "", "The keymap of fish/zsh")
	promptCmd.Flags().IntP("jobs", "j", 0, "The number of currently running jobs")
	promptCmd.Flags().IntP("status", "s", 0, "The status code of the previously run command")
}
