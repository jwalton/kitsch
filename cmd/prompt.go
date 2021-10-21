package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/go-supportscolor"
	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/modules"
	"github.com/jwalton/kitsch-prompt/internal/shellprompt"
	"github.com/jwalton/kitsch-prompt/internal/styling"
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
		cmdDuration := int64(0)
		if cmdDurationStr != "" {
			cmdDuration, _ = strconv.ParseInt(cmdDurationStr, 10, 64)
		}

		// Because the prompt is shown from the shell, when it is run, it
		// will not be in a TTY.  Disable TTY detection in gchalk.
		stdoutFd := os.Stdout.Fd()
		level := supportscolor.SupportsColor(stdoutFd, supportscolor.IsTTYOption(true))
		gchalk.SetLevel(level.Level)
		gchalk.Stderr.SetLevel(level.Level)

		configuration, err := readConfig()
		if err != nil {
			println(gchalk.Red("Fatal error parsing configuration: ", err.Error()))
			os.Exit(1)
		}

		globals := modules.NewGlobals(status, cmdDuration, keymap)
		runtimeEnv := env.New(globals.CWD, jobs)

		// Load custom colors
		styles := styling.Registry{}
		for colorName, color := range configuration.Colors {
			if !strings.HasPrefix(colorName, "$") {
				runtimeEnv.Warn("Custom color \"" + colorName + "should start with $")
			}
			styles.AddCustomColor(colorName, color)
		}

		if err != nil {
			fmt.Println(err)
			fmt.Print("$ ")
		} else {
			context := modules.Context{
				Environment: runtimeEnv,
				Globals:     globals,
				Styles:      styles,
			}

			prompt := configuration.Prompt.Module.Execute(&context)
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
