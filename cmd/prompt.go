package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/go-supportscolor"
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
		start := time.Now()
		cacheDir := filepath.Join(userConfigDir, "cache")

		jobs, _ := cmd.Flags().GetInt("jobs")
		status, _ := cmd.Flags().GetInt("status")
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

		setupDuration := time.Since(start)

		// Read configuration
		start = time.Now()
		configuration, err := readConfig()
		if err != nil {
			println(gchalk.Red("Fatal error parsing configuration: ", err.Error()))
			fmt.Print("$ ")
			os.Exit(1)
		}

		styles := styling.Registry{}
		styles.AddCustomColors(configuration.Colors)

		configurationParsingDuration := time.Since(start)

		// Create our context.
		start = time.Now()
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
			globals := modules.NewGlobals(shell, status, jobs, cmdDuration, keymap)
			context = modules.NewContext(globals, configuration.ProjectsTypes, cacheDir, styles)
		}
		contextDuration := time.Since(start)

		// Execute the prompt
		result := configuration.Prompt.Module.Execute(&context)

		if perf {
			fmt.Printf("Setup time: %v\n", setupDuration)
			fmt.Printf("Parsing configuration: %v\n", configurationParsingDuration)
			fmt.Printf("Context setup: %v\n", contextDuration)
			renderPerf(result.ChildDurations, 0)
		}

		withEscapes := shellprompt.AddZeroWidthCharacterEscapes(context.Globals.Shell, result.Text)
		fmt.Print(withEscapes)
	},
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
	promptCmd.Flags().String("demo", "", "If present, kitsch-prompt will run in demo mode, loading values from the specified file.")
}
