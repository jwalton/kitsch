package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/go-supportscolor"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modules"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
	"github.com/jwalton/kitsch-prompt/internal/perf"
	"github.com/jwalton/kitsch-prompt/internal/shellprompt"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Show the prompt",
	Run: func(cmd *cobra.Command, args []string) {
		performance := perf.New(4)

		cacheDir := filepath.Join(userConfigDir, "cache")

		jobs, _ := cmd.Flags().GetInt("jobs")
		status, _ := cmd.Flags().GetInt("status")
		terminalWidth, _ := cmd.Flags().GetInt("terminal-width")
		keymap, _ := cmd.Flags().GetString("keymap")
		shell, _ := cmd.Flags().GetString("shell")
		perf, _ := cmd.Flags().GetBool("perf")
		demo, _ := cmd.Flags().GetString("demo")

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
		performance.End("Option parsing")

		// Read configuration
		configuration, err := readConfig()
		if err != nil {
			println(gchalk.Red("Fatal error parsing configuration: ", err.Error()))
			fmt.Print("$ ")
			os.Exit(1)
		}

		styles := styling.Registry{}
		styles.AddCustomColors(configuration.Colors)

		performance.End("Config parsing")

		// Create our context.
		var context modules.Context
		if demo != "" {
			demoConfig := &modules.DemoConfig{}
			err := demoConfig.Load(demo)
			if err != nil {
				log.Error("Failed to load demo config:", err)
				os.Exit(1)
			}
			context = modules.NewDemoContext(*demoConfig, styles)
		} else {
			globals := modules.NewGlobals(shell, terminalWidth, status, jobs, cmdDuration, keymap)
			context = modules.NewContext(globals, configuration.ProjectsTypes, cacheDir, styles)
		}
		performance.End("Context setup")

		// Execute the prompt
		result := configuration.Prompt.Module.Execute(&context)

		performance.EndWithChildren("Prompt", result.ChildDurations)

		if perf {
			performance.Print()
		}

		withEscapes := shellprompt.AddZeroWidthCharacterEscapes(context.Globals.Shell, result.Text)
		fmt.Print(withEscapes)
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().String("shell", "", "The type of shell")
	promptCmd.Flags().StringP("cmd-duration", "d", "", "The execution duration of the last command, in milliseconds")
	promptCmd.Flags().StringP("keymap", "k", "", "The keymap of fish/zsh")
	promptCmd.Flags().IntP("jobs", "j", 0, "The number of currently running jobs")
	promptCmd.Flags().IntP("status", "s", 0, "The status code of the previously run command")
	promptCmd.Flags().Int("terminal-width", 0, "The width of the terminal")
	promptCmd.Flags().Bool("perf", false, "Print performance information about each module")
	promptCmd.Flags().Bool("verbose", false, "Print verbose output")
	promptCmd.Flags().String("demo", "", "If present, kitsch-prompt will run in demo mode, loading values from the specified file.")
}
