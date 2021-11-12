package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/go-supportscolor"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/config"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modules"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
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
		perf, _ := cmd.Flags().GetBool("perf")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			log.SetVerbose(true)
		}

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

		globals := modules.NewGlobals(shell, status, cmdDuration, keymap)
		runtimeEnv := env.New(globals.CWD, jobs)

		if err != nil {
			log.Error(err)
			fmt.Print("$ ")
		} else {
			result := renderPrompt(configuration, globals, runtimeEnv)
			if perf {
				renderPerf(result.ChildDurations, 0)
			}

			withEscapes := shellprompt.AddZeroWidthCharacterEscapes(globals.Shell, result.Text)
			fmt.Print(withEscapes)
		}
	},
}

// renderPrompt will render the prompt with the given configuration.
func renderPrompt(
	configuration *config.Config,
	globals modules.Globals,
	runtimeEnv env.Env,
) modules.ModuleResult {
	// Load custom colors
	styles := styling.Registry{}
	for colorName, color := range configuration.Colors {
		if !strings.HasPrefix(colorName, "$") {
			log.Warn("Custom color \"" + colorName + "must start with $")
		} else {
			styles.AddCustomColor(colorName, color)
		}
	}

	context := modules.Context{
		Environment:  runtimeEnv,
		Directory:    fileutils.NewDirectory(globals.CWD),
		Styles:       styles,
		Globals:      globals,
		ProjectTypes: configuration.ProjectsTypes,
	}

	return configuration.Prompt.Module.Execute(&context)
}

func renderPerf(durations []modules.ModuleDuration, indent int) {
	for _, duration := range durations {
		printDuration := duration.Duration.String()
		if duration.Duration > 1000000 {
			printDuration = gchalk.Yellow(printDuration)
		} else if duration.Duration > 250000000 {
			printDuration = gchalk.Red(printDuration)
		} else {
			printDuration = gchalk.Green(printDuration)
		}

		fmt.Printf("%s%s(%d:%d) - %s\n",
			strings.Repeat(" ", indent),
			duration.Module.Type,
			duration.Module.Line,
			duration.Module.Column,
			printDuration,
		)
		if len(duration.Children) > 0 {
			renderPerf(duration.Children, indent+2)
		}
	}
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().String("shell", "", "The type of shell")
	promptCmd.Flags().StringP("cmd-duration", "d", "", "The execution duration of the last command, in milliseconds")
	promptCmd.Flags().StringP("keymap", "k", "", "The keymap of fish/zsh")
	promptCmd.Flags().IntP("jobs", "j", 0, "The number of currently running jobs")
	promptCmd.Flags().IntP("status", "s", 0, "The status code of the previously run command")
	promptCmd.Flags().Bool("perf", false, "Print performance information about each module")
	promptCmd.Flags().Bool("verbose", false, "Print verbose output")
}
